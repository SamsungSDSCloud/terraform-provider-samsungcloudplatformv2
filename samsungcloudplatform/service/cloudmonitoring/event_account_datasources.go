package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventAccountDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventAccountDataSources{}
)

func NewCloudMonitoringEventAccountDataSources() datasource.DataSource {
	return &cloudMonitoringEventAccountDataSources{}
}

type cloudMonitoringEventAccountDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventAccountDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_account_events" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventAccountDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of event.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("EventState"): schema.StringAttribute{
				Description: "EventState",
				Optional:    true,
			},
			common.ToSnakeCase("QueryStartDt"): schema.StringAttribute{
				Description: "queryStartDt",
				Optional:    true,
			},
			common.ToSnakeCase("QueryEndDt"): schema.StringAttribute{
				Description: "QueryEndDt",
				Optional:    true,
			},
			common.ToSnakeCase("Events"): schema.ListNestedAttribute{
				Description: "A list of event.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("DurationSecond"): schema.Int64Attribute{
							Description: "DurationSecond",
							Computed:    true,
						},
						common.ToSnakeCase("EndDt"): schema.StringAttribute{
							Description: "EndDt",
							Computed:    true,
						},
						common.ToSnakeCase("EventId"): schema.StringAttribute{
							Description: "EventId",
							Computed:    true,
						},
						common.ToSnakeCase("EventLevel"): schema.StringAttribute{
							Description: "EventLevel",
							Computed:    true,
						},
						common.ToSnakeCase("EventMessage"): schema.StringAttribute{
							Description: "EventMessage",
							Computed:    true,
						},
						common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
							Description: "EventPolicyId",
							Computed:    true,
						},
						common.ToSnakeCase("EventState"): schema.StringAttribute{
							Description: "EventState",
							Computed:    true,
						},
						common.ToSnakeCase("MetricKey"): schema.StringAttribute{
							Description: "MetricKey",
							Computed:    true,
						},
						common.ToSnakeCase("MetricName"): schema.StringAttribute{
							Description: "MetricName",
							Computed:    true,
						},
						common.ToSnakeCase("EventMessage"): schema.StringAttribute{
							Description: "EventMessage",
							Computed:    true,
						},
						common.ToSnakeCase("ObjectDisplayName"): schema.StringAttribute{
							Description: "ObjectDisplayName",
							Computed:    true,
						},
						common.ToSnakeCase("ObjectName"): schema.StringAttribute{
							Description: "ObjectName",
							Computed:    true,
						},
						common.ToSnakeCase("ObjectType"): schema.StringAttribute{
							Description: "ObjectType",
							Computed:    true,
						},
						common.ToSnakeCase("ObjectTypeName"): schema.StringAttribute{
							Description: "ObjectTypeName",
							Computed:    true,
						},
						common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
							Description: "ProductResourceId",
							Computed:    true,
						},
						common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
							Description: "ProductTypeCode",
							Computed:    true,
						},
						common.ToSnakeCase("StartDt"): schema.StringAttribute{
							Description: "StartDt",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{ // 필터는 Block 으로 정의한다.
			"filter": filter.DataSourceSchema(), // 필터 스키마는 공통으로 제공되는 함수를 이용하여 정의한다.
		},
	}
}

func (d *cloudMonitoringEventAccountDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *cloudMonitoringEventAccountDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventAccountDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	event, err := d.client.GetAccountEventList(state.EventState, state.QueryStartDt, state.QueryEndDt)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, eventElement := range event.Contents {
		eventState := cloudmonitoring.Event{
			DurationSecond:    types.Int64Value(eventElement.DurationSecond),
			EndDt:             types.StringValue(eventElement.GetEndDt().String()),
			EventId:           types.StringValue(eventElement.EventId),
			EventLevel:        types.StringValue(eventElement.EventLevel),
			EventMessage:      types.StringValue(eventElement.EventMessage),
			EventPolicyId:     types.Int64Value(eventElement.EventPolicyId),
			EventState:        types.StringValue(eventElement.EventState),
			MetricKey:         types.StringValue(eventElement.MetricKey),
			MetricName:        types.StringValue(eventElement.MetricName),
			ObjectDisplayName: types.StringValue(eventElement.GetObjectDisplayName()),
			ObjectName:        types.StringValue(eventElement.GetObjectName()),
			ObjectType:        types.StringValue(eventElement.GetObjectType()),
			ObjectTypeName:    types.StringValue(eventElement.GetObjectTypeName()),
			ProductResourceId: types.StringValue(eventElement.ProductResourceId),
			StartDt:           types.StringValue(eventElement.GetStartDt().String()),
		}

		state.Events = append(state.Events, eventState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
