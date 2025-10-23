package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common/filter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventPolicyHistoryDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventPolicyHistoryDataSources{}
)

func NewCloudMonitoringEventPolicyHistoryDataSources() datasource.DataSource {
	return &cloudMonitoringEventPolicyHistoryDataSources{}
}

type cloudMonitoringEventPolicyHistoryDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventPolicyHistoryDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event_policy_histories" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventPolicyHistoryDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "Event Notification States.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Required:    true,
			},
			common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
				Description: "EventPolicyId",
				Required:    true,
			},
			common.ToSnakeCase("QueryStartDt"): schema.StringAttribute{
				Description: "queryStartDt",
				Optional:    true,
			},
			common.ToSnakeCase("QueryEndDt"): schema.StringAttribute{
				Description: "QueryEndDt",
				Optional:    true,
			},
			common.ToSnakeCase("EventPolicyHistories"): schema.ListNestedAttribute{
				Description: "A list of Event Policy History.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("ModifiedDt"): schema.StringAttribute{
							Description: "ModifiedDt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedByName"): schema.StringAttribute{
							Description: "ModifiedByName",
							Computed:    true,
						},
						common.ToSnakeCase("CreateById"): schema.StringAttribute{
							Description: "CreateById",
							Computed:    true,
						},
						//common.ToSnakeCase("CreateBy"): schema.SingleNestedAttribute{
						//	Description: "CreateBy",
						//	Computed:    true,
						//},
						common.ToSnakeCase("UpdateById"): schema.StringAttribute{
							Description: "UpdateById",
							Computed:    true,
						},
						//common.ToSnakeCase("UpdateBy"): schema.SingleNestedAttribute{
						//	Description: "UpdateBy",
						//	Computed:    true,
						//},
						common.ToSnakeCase("EventPolicyHistoryId"): schema.Int64Attribute{
							Description: "EventPolicyHistoryId",
							Computed:    true,
						},
						common.ToSnakeCase("EventPolicyHistoryType"): schema.StringAttribute{
							Description: "EventPolicyHistoryType",
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
						common.ToSnakeCase("ProductName"): schema.StringAttribute{
							Description: "ProductName",
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
						common.ToSnakeCase("MetricDescription"): schema.StringAttribute{
							Description: "MetricDescription",
							Computed:    true,
						},
						common.ToSnakeCase("MetricDescriptionEn"): schema.StringAttribute{
							Description: "MetricDescriptionEn",
							Computed:    true,
						},
						common.ToSnakeCase("MetricUnit"): schema.StringAttribute{
							Description: "MetricUnit",
							Computed:    true,
						},
						common.ToSnakeCase("ProductTargetType"): schema.StringAttribute{
							Description: "ProductTargetType",
							Computed:    true,
						},
						common.ToSnakeCase("ProductTargetTypeEn"): schema.StringAttribute{
							Description: "ProductTargetTypeEn",
							Computed:    true,
						},
						common.ToSnakeCase("ObjectName"): schema.StringAttribute{
							Description: "ObjectName",
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
						common.ToSnakeCase("EventLevel"): schema.StringAttribute{
							Description: "EventLevel",
							Computed:    true,
						},
						common.ToSnakeCase("FtCount"): schema.Int64Attribute{
							Description: "FtCount",
							Computed:    true,
						},
						common.ToSnakeCase("EventMessagePrefix"): schema.StringAttribute{
							Description: "EventMessagePrefix",
							Computed:    true,
						},
						common.ToSnakeCase("DisableObject"): schema.StringAttribute{
							Description: "DisableObject",
							Computed:    true,
						},
						common.ToSnakeCase("DisableYn"): schema.StringAttribute{
							Description: "DisableYn",
							Computed:    true,
						},
						common.ToSnakeCase("EventOccurTimeZone"): schema.StringAttribute{
							Description: "EventOccurTimeZone",
							Computed:    true,
						},
						common.ToSnakeCase("EventPolicyStatistics"): schema.SingleNestedAttribute{
							Description: "EventPolicyStatistics",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("EventPolicyStatisticsType"): schema.StringAttribute{
									Description: "EventPolicyStatisticsType",
									Computed:    true,
								},
								common.ToSnakeCase("EventPolicyStatisticsPeriod"): schema.Int64Attribute{
									Description: "EventPolicyStatisticsPeriod",
									Computed:    true,
								},
							},
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

func (d *cloudMonitoringEventPolicyHistoryDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *cloudMonitoringEventPolicyHistoryDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventPolicyHistoryDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetEventPolicyHistoryList(state.ResourceType, state.EventPolicyId, state.QueryStartDt, state.QueryEndDt)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, eventPolicyHistoryElement := range data.Contents {

		eventPolicyStatistics := cloudmonitoring.EventPolicyStatistics{
			EventPolicyStatisticsType:   types.StringValue(eventPolicyHistoryElement.GetEventPolicyStatistics().EventPolicyStatisticsType),
			EventPolicyStatisticsPeriod: types.Int64Value(eventPolicyHistoryElement.GetEventPolicyStatistics().EventPolicyStatisticsPeriod),
		}

		eventPolicyHistoryStates := cloudmonitoring.EventPolicyHistory{
			ModifiedDt:     types.StringValue(eventPolicyHistoryElement.GetModifiedDt().String()),
			ModifiedBy:     types.StringValue(eventPolicyHistoryElement.GetModifiedBy()),
			ModifiedByName: types.StringValue(eventPolicyHistoryElement.GetModifiedByName()),
			CreateById:     types.StringValue(eventPolicyHistoryElement.GetCreateById()),
			//CreateBy: createBy,
			UpdateById: types.StringValue(eventPolicyHistoryElement.GetUpdateById()),
			//UpdateBy: UpdateBy,
			EventPolicyHistoryId:   types.Int64Value(eventPolicyHistoryElement.GetEventPolicyHistoryId()),
			EventPolicyHistoryType: types.StringValue(eventPolicyHistoryElement.EventPolicyHistoryType),
			EventPolicyId:          types.Int64Value(eventPolicyHistoryElement.GetEventPolicyId()),
			ProductResourceId:      types.StringValue(eventPolicyHistoryElement.ProductResourceId),
			ProductName:            types.StringValue(eventPolicyHistoryElement.GetProductName()),
			MetricKey:              types.StringValue(eventPolicyHistoryElement.MetricKey),
			MetricName:             types.StringValue(eventPolicyHistoryElement.MetricName),
			MetricDescription:      types.StringValue(eventPolicyHistoryElement.GetMetricDescription()),
			MetricDescriptionEn:    types.StringValue(eventPolicyHistoryElement.GetMetricDescriptionEn()),
			MetricUnit:             types.StringValue(eventPolicyHistoryElement.GetMetricUnit()),
			ProductTargetType:      types.StringValue(eventPolicyHistoryElement.GetProductTargetType()),
			ProductTargetTypeEn:    types.StringValue(eventPolicyHistoryElement.GetProductTargetTypeEn()),
			ObjectName:             types.StringValue(eventPolicyHistoryElement.GetObjectName()),
			ObjectDisplayName:      types.StringValue(eventPolicyHistoryElement.GetObjectDisplayName()),
			ObjectType:             types.StringValue(eventPolicyHistoryElement.GetObjectType()),
			EventLevel:             types.StringValue(eventPolicyHistoryElement.EventLevel),
			FtCount:                types.Int64Value(eventPolicyHistoryElement.FtCount),
			EventMessagePrefix:     types.StringValue(eventPolicyHistoryElement.GetEventMessagePrefix()),
			DisableObject:          types.StringValue(eventPolicyHistoryElement.GetDisableObject()),
			DisableYn:              types.StringValue(eventPolicyHistoryElement.GetDisableYn()),
			EventOccurTimeZone:     types.StringValue(eventPolicyHistoryElement.GetEventOccurTimeZone()),
			//EventThreshold: EventThreshold,
			EventPolicyStatistics: eventPolicyStatistics,
		}

		state.EventPolicyHistories = append(state.EventPolicyHistories, eventPolicyHistoryStates)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
