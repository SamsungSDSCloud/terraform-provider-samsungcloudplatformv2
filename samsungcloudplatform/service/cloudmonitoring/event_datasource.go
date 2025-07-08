package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventDataSource{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventDataSource{}
)

func NewCloudMonitoringEventDataSource() datasource.DataSource {
	return &cloudMonitoringEventDataSource{}
}

type cloudMonitoringEventDataSource struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Event Detail.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("EventId"): schema.StringAttribute{
				Description: "EventId",
				Optional:    true,
			},
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Optional:    true,
			},
			//common.ToSnakeCase("isLogMetric"): schema.StringAttribute{
			//	Description: "isLogMetric",
			//	Optional:    true,
			//},
			common.ToSnakeCase("EventLevel"): schema.StringAttribute{
				Description: "EventLevel",
				Optional:    true,
			},
			common.ToSnakeCase("EventMessage"): schema.StringAttribute{
				Description: "EventMessage",
				Computed:    true,
			},
			common.ToSnakeCase("EventState"): schema.StringAttribute{
				Description: "EventState",
				Computed:    true,
			},
			common.ToSnakeCase("StartDt"): schema.StringAttribute{
				Description: "StartDt",
				Computed:    true,
			},
			common.ToSnakeCase("EndDt"): schema.StringAttribute{
				Description: "EndDt",
				Computed:    true,
			},
			common.ToSnakeCase("DurationSecond"): schema.Int64Attribute{
				Description: "DurationSecond",
				Computed:    true,
			},
			common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
				Description: "EventPolicyId",
				Computed:    true,
			},
			common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
				Description: "ProductResourceId",
				Computed:    true,
			},
			common.ToSnakeCase("MetricKey"): schema.StringAttribute{
				Description: "MetricKey",
				Computed:    true,
			},
			common.ToSnakeCase("ObjectDisplayName"): schema.StringAttribute{
				Description: "ObjectDisplayName",
				Computed:    true,
			},
			common.ToSnakeCase("ObjectType"): schema.StringAttribute{
				Description: "ObjectType",
				Computed:    true,
			},
			common.ToSnakeCase("MetricName"): schema.StringAttribute{
				Description: "MetricName",
				Computed:    true,
			},
			common.ToSnakeCase("ObjectName"): schema.StringAttribute{
				Description: "ObjectName",
				Computed:    true,
			},
			common.ToSnakeCase("ObjectTypeName"): schema.StringAttribute{
				Description: "ObjectTypeName",
				Computed:    true,
			},
			common.ToSnakeCase("ProductIpAddress"): schema.StringAttribute{
				Description: "ProductIpAddress",
				Computed:    true,
			},
			common.ToSnakeCase("ProductName"): schema.StringAttribute{
				Description: "ProductName",
				Computed:    true,
			},
			common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
				Description: "ProductTypeCode",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

func (d *cloudMonitoringEventDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = inst.Client.CloudMonitoring
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *cloudMonitoringEventDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventDataSourceId

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	eventDetail, err := d.client.GetEvent(state.EventId, state.ResourceType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	//metricSummaryModel := eventDetail.MetricSummary{
	//	EventId:    types.StringValue(eventDetail.EventId),
	//	EventLevel: types.StringValue(eventDetail.EventLevel),
	//	EventLevel: types.StringValue(eventDetail.EventLevel),
	//}

	state.DurationSecond = types.Int64Value(eventDetail.GetDurationSecond())
	state.EndDt = types.StringValue(eventDetail.GetEndDt().String())
	state.EventId = types.StringValue(eventDetail.GetEventId())
	state.EventLevel = types.StringValue(eventDetail.GetEventLevel())
	state.EventMessage = types.StringValue(eventDetail.GetEventMessage())
	state.EventPolicyId = types.Int64Value(eventDetail.GetEventPolicyId())
	//state.EventPolicySummary
	state.EventState = types.StringValue(eventDetail.GetEventState())
	state.MetricKey = types.StringValue(eventDetail.GetMetricKey())
	state.MetricName = types.StringValue(eventDetail.GetMetricName())
	//state.metricSummary
	state.ObjectDisplayName = types.StringValue(eventDetail.GetObjectDisplayName())
	state.ObjectName = types.StringValue(eventDetail.GetObjectName())
	state.ObjectType = types.StringValue(eventDetail.GetObjectType())
	state.ObjectTypeName = types.StringValue(eventDetail.GetObjectTypeName())
	state.ProductIpAddress = types.StringValue(eventDetail.GetProductIpAddress())
	state.ProductName = types.StringValue(eventDetail.GetProductName())
	state.ProductResourceId = types.StringValue(eventDetail.GetProductResourceId())
	//state.productSummary
	state.ProductTypeCode = types.StringValue(eventDetail.GetProductTypeCode())
	state.StartDt = types.StringValue(eventDetail.GetStartDt().String())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
