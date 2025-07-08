package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventPolicyDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventPolicyDataSources{}
)

func NewCloudMonitoringEventPolicyDataSources() datasource.DataSource {
	return &cloudMonitoringEventPolicyDataSources{}
}

type cloudMonitoringEventPolicyDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventPolicyDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event_policies" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventPolicyDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of event policy.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Optional:    true,
			},
			common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
				Description: "ProductResourceId",
				Optional:    true,
			},
			common.ToSnakeCase("MetricKey"): schema.StringAttribute{
				Description: "MetricKey",
				Optional:    true,
			},
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page",
				Optional:    true,
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size",
				Optional:    true,
			},
			common.ToSnakeCase("EventPolicySummaries"): schema.ListNestedAttribute{
				Description: "Event Policy List.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("CreateById"): schema.StringAttribute{
							Description: "CreateById",
							Computed:    true,
						},
						common.ToSnakeCase("UpdateById"): schema.StringAttribute{
							Description: "UpdateById",
							Computed:    true,
						},
						//common.ToSnakeCase("UpdateBy"): schema.SingleNestedAttribute{
						//	Description: "UpdateBy",
						//	Computed:    true,
						//},
						common.ToSnakeCase("ModifiedDt"): schema.StringAttribute{
							Description: "ModifiedDt",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Description: "ModifiedBy",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedByName"): schema.StringAttribute{
							Description: "CreatedByName",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Description: "CreatedBy",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedByName"): schema.StringAttribute{
							Description: "ModifiedByName",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedById"): schema.StringAttribute{
							Description: "CreatedById",
							Computed:    true,
						},
						//common.ToSnakeCase("CreateBy"): schema.SingleNestedAttribute{
						//	Description: "CreateBy",
						//	Computed:    true,
						//},
						common.ToSnakeCase("UpdatedById"): schema.StringAttribute{
							Description: "UpdatedById",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedDt"): schema.StringAttribute{
							Description: "CreatedDt",
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
						common.ToSnakeCase("ProductSq"): schema.Int64Attribute{
							Description: "ProductSq",
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
						common.ToSnakeCase("ProductTargetType"): schema.StringAttribute{
							Description: "ProductTargetType",
							Computed:    true,
						},
						common.ToSnakeCase("ProductTargetTypeEn"): schema.StringAttribute{
							Description: "ProductTargetTypeEn",
							Computed:    true,
						},
						common.ToSnakeCase("IsLogMetric"): schema.StringAttribute{
							Description: "IsLogMetric",
							Computed:    true,
						},
						common.ToSnakeCase("ObjectName"): schema.StringAttribute{
							Description: "ObjectName",
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
						common.ToSnakeCase("CheckAsg"): schema.BoolAttribute{
							Description: "CheckAsg",
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
						common.ToSnakeCase("EventThreshold"): schema.SingleNestedAttribute{
							Description: "EventThreshold",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ThresholdType"): schema.StringAttribute{
									Description: "ThresholdType",
									Computed:    true,
								},
								common.ToSnakeCase("MetricFunction"): schema.StringAttribute{
									Description: "MetricFunction",
									Computed:    true,
								},
								common.ToSnakeCase("SingleThreshold"): schema.SingleNestedAttribute{
									Description: "SingleThreshold",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										common.ToSnakeCase("ComparisonOperator"): schema.StringAttribute{
											Description: "ComparisonOperator",
											Computed:    true,
										},
										common.ToSnakeCase("Value"): schema.Float64Attribute{
											Description: "Value",
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudMonitoringEventPolicyDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *cloudMonitoringEventPolicyDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventPolicyDataSourceIds

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	eventPolicy, err := d.client.GetProductEventPolicyList(state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Servers : "+detail,
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, eventPolicySummariesElement := range eventPolicy.Contents {

		eventPolicyStatistics := cloudmonitoring.EventPolicyStatistics{
			EventPolicyStatisticsType:   types.StringValue(eventPolicySummariesElement.GetEventPolicyStatistics().EventPolicyStatisticsType),
			EventPolicyStatisticsPeriod: types.Int64Value(eventPolicySummariesElement.GetEventPolicyStatistics().EventPolicyStatisticsPeriod),
		}

		singleThreshold := cloudmonitoring.SingleThreshold{

			ComparisonOperator: types.StringValue(eventPolicySummariesElement.EventThreshold.GetSingleThreshold().ComparisonOperator),
			Value:              types.Float64Value(eventPolicySummariesElement.EventThreshold.GetSingleThreshold().Value),
		}

		//rangeThreshold := cloudmonitoring.RangeThreshold{
		//	MaxComparisonOperator types.String  `tfsdk:"max_comparison_operator"`
		//	MinComparisonOperator types.String  `tfsdk:"min_comparison_operator"`
		//	MaxValue              types.Float64 `tfsdk:"max_value"`
		//	MinValue              types.Float64 `tfsdk:"min_value"`
		//}

		eventThreshold := cloudmonitoring.EventThreshold{
			ThresholdType: types.StringValue(eventPolicySummariesElement.GetEventThreshold().ThresholdType),
			//MetricFunction: types.StringValue(eventPolicyElement.GetEventThreshold().GetMetricFunction()),
			SingleThreshold: singleThreshold,
			//RangeThreshold:  rangeThreshold,
		}

		eventPoliciesSummariesState := cloudmonitoring.EventPolicySummary{
			EventPolicyId:         types.Int64Value(eventPolicySummariesElement.GetEventPolicyId()),
			ProductResourceId:     types.StringValue(eventPolicySummariesElement.ProductResourceId),
			ProductSq:             types.Int64Value(eventPolicySummariesElement.ProductSq),
			ProductName:           types.StringValue(eventPolicySummariesElement.GetProductName()),
			MetricKey:             types.StringValue(eventPolicySummariesElement.MetricKey),
			MetricName:            types.StringValue(eventPolicySummariesElement.GetMetricName()),
			MetricDescription:     types.StringValue(eventPolicySummariesElement.GetMetricDescription()),
			MetricDescriptionEn:   types.StringValue(eventPolicySummariesElement.GetMetricDescriptionEn()),
			ProductTargetType:     types.StringValue(eventPolicySummariesElement.GetProductTargetType()),
			ProductTargetTypeEn:   types.StringValue(eventPolicySummariesElement.GetProductTargetTypeEn()),
			IsLogMetric:           types.StringValue(eventPolicySummariesElement.GetIsLogMetric()),
			ObjectName:            types.StringValue(eventPolicySummariesElement.GetObjectName()),
			EventLevel:            types.StringValue(eventPolicySummariesElement.EventLevel),
			FtCount:               types.Int64Value(eventPolicySummariesElement.FtCount),
			EventMessagePrefix:    types.StringValue(eventPolicySummariesElement.GetEventMessagePrefix()),
			CheckAsg:              types.BoolValue(eventPolicySummariesElement.GetCheckAsg()),
			EventOccurTimeZone:    types.StringValue(eventPolicySummariesElement.GetEventOccurTimeZone()),
			EventThreshold:        eventThreshold,
			EventPolicyStatistics: eventPolicyStatistics,
		}

		state.EventPolicySummaries = append(state.EventPolicySummaries, eventPoliciesSummariesState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
