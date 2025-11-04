package cloudmonitoring

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/cloudmonitoring"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &cloudMonitoringEventPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringEventPolicyDataSource{}
)

func NewCloudMonitoringEventPolicyDataSource() datasource.DataSource {
	return &cloudMonitoringEventPolicyDataSource{}
}

type cloudMonitoringEventPolicyDataSource struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *cloudMonitoringEventPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_event_policy" // service 의 metadata 를 {{ provider명 }}_{{ 서비스명 }}_{{ 단수형 리소스명 }} 형태로 추가한다.
}

func (d *cloudMonitoringEventPolicyDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) { // 아직 정의하지 않은 Schema 메서드를 추가한다.
	resp.Schema = schema.Schema{
		Description: "list of event.",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("xResourceType"): schema.StringAttribute{
				Description: "xResourceType",
				Optional:    true,
			},
			common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
				Description: "EventPolicyId",
				Optional:    true,
			},
			common.ToSnakeCase("EventPolicyDetail"): schema.SingleNestedAttribute{
				Description: "EventPolicyDetail",
				Computed:    true,
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
					//common.ToSnakeCase("UpdatedById"): schema.StringAttribute{
					//	Description: "UpdatedById",
					//	Computed:    true,
					//},
					common.ToSnakeCase("CreatedDt"): schema.StringAttribute{
						Description: "CreatedDt",
						Computed:    true,
					},
					common.ToSnakeCase("EventPolicyId"): schema.Int64Attribute{
						Description: "EventPolicyId",
						Computed:    true,
					},
					common.ToSnakeCase("ProductSummary"): schema.SingleNestedAttribute{
						Description: "ProductSummary",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("ProductSq"): schema.Int64Attribute{
								Description: "ProductSq",
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
							common.ToSnakeCase("ProductTypeCode"): schema.StringAttribute{
								Description: "ProductTypeCode",
								Computed:    true,
							},
							common.ToSnakeCase("ProductTypeName"): schema.StringAttribute{
								Description: "ProductTypeName",
								Computed:    true,
							},
							common.ToSnakeCase("ProductIpAddress"): schema.StringAttribute{
								Description: "ProductIpAddress",
								Computed:    true,
							},
							common.ToSnakeCase("ProductState"): schema.StringAttribute{
								Description: "ProductState",
								Computed:    true,
							},
							common.ToSnakeCase("AgentState"): schema.StringAttribute{
								Description: "AgentState",
								Computed:    true,
							},
							common.ToSnakeCase("ProductSubName"): schema.StringAttribute{
								Description: "ProductSubName",
								Computed:    true,
							},
							common.ToSnakeCase("ProductSubType"): schema.StringAttribute{
								Description: "ProductSubType",
								Computed:    true,
							},
							common.ToSnakeCase("LbName"): schema.StringAttribute{
								Description: "LbName",
								Computed:    true,
							},
							common.ToSnakeCase("VpcName"): schema.StringAttribute{
								Description: "VpcName",
								Computed:    true,
							},
							common.ToSnakeCase("LbSize"): schema.StringAttribute{
								Description: "LbSize",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("MetricSummary"): schema.SingleNestedAttribute{
						Description: "MetricSummary",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							common.ToSnakeCase("MetricKey"): schema.StringAttribute{
								Description: "MetricKey",
								Computed:    true,
							},
							common.ToSnakeCase("MetricName"): schema.StringAttribute{
								Description: "MetricName",
								Computed:    true,
							},
							common.ToSnakeCase("MetricSetKey"): schema.StringAttribute{
								Description: "MetricSetKey",
								Computed:    true,
							},
							common.ToSnakeCase("MetricSetName"): schema.StringAttribute{
								Description: "MetricSetName",
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
							common.ToSnakeCase("MetricType"): schema.StringAttribute{
								Description: "MetricType",
								Computed:    true,
							},
							common.ToSnakeCase("MetricUnit"): schema.StringAttribute{
								Description: "MetricUnit",
								Computed:    true,
							},
							common.ToSnakeCase("IsObjectExist"): schema.StringAttribute{
								Description: "IsObjectExist",
								Computed:    true,
							},
							common.ToSnakeCase("IsLogMetric"): schema.StringAttribute{
								Description: "IsLogMetric",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("ProductSq"): schema.Int64Attribute{
						Description: "ProductSq",
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
					}, common.ToSnakeCase("ObjectDisplayName"): schema.StringAttribute{
						Description: "ObjectDisplayName",
						Computed:    true,
					},
					common.ToSnakeCase("EventLevel"): schema.StringAttribute{
						Description: "EventLevel",
						Computed:    true,
					}, common.ToSnakeCase("FtCount"): schema.Int64Attribute{
						Description: "FtCount",
						Computed:    true,
					}, common.ToSnakeCase("EventMessagePrefix"): schema.StringAttribute{
						Description: "EventMessagePrefix",
						Computed:    true,
					}, common.ToSnakeCase("ObjectType"): schema.StringAttribute{
						Description: "ObjectType",
						Computed:    true,
					},
					common.ToSnakeCase("ObjectTypeName"): schema.StringAttribute{
						Description: "ObjectTypeName",
						Computed:    true,
					},
					common.ToSnakeCase("ProductInfoAttrs"): schema.StringAttribute{
						Description: "ProductInfoAttrs",
						Computed:    true,
					},
					//common.ToSnakeCase("ProductInfoAttrs"): schema.ListNestedAttribute{
					//	Description: "ProductInfoAttrs",
					//	Computed:    true,
					//	NestedObject: schema.NestedAttributeObject{
					//		Attributes: map[string]schema.Attribute{
					//			common.ToSnakeCase("attrName"): schema.StringAttribute{
					//				Description: "attrName",
					//				Computed:    true,
					//			},
					//			common.ToSnakeCase("attrValue"): schema.StringAttribute{
					//				Description: "attrValue",
					//				Computed:    true,
					//			},
					//		},
					//	},
					//},
					common.ToSnakeCase("DisableObject"): schema.StringAttribute{
						Description: "DisableObject",
						Computed:    true,
					},
					common.ToSnakeCase("UserNames"): schema.StringAttribute{
						Description: "UserNames",
						Computed:    true,
					},
					common.ToSnakeCase("UserNameStr"): schema.StringAttribute{
						Description: "UserNameStr",
						Computed:    true,
					},
					common.ToSnakeCase("DisableYn"): schema.StringAttribute{
						Description: "DisableYn",
						Computed:    true,
					},
					common.ToSnakeCase("AttrListStr"): schema.StringAttribute{
						Description: "AttrListStr",
						Computed:    true,
					},
					common.ToSnakeCase("AsgYn"): schema.StringAttribute{
						Description: "AsgYn",
						Computed:    true,
					},
					common.ToSnakeCase("StartDt"): schema.StringAttribute{
						Description: "StartDt",
						Computed:    true,
					},
					common.ToSnakeCase("DisplayEventRule"): schema.StringAttribute{
						Description: "DisplayEventRule",
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
					//common.ToSnakeCase("EventPolicyStatistics"): schema.SingleNestedAttribute{
					//	Description: "EventPolicyStatistics",
					//	Computed:    true,
					//	Attributes: map[string]schema.Attribute{
					//		common.ToSnakeCase("EventPolicyStatisticsType"): schema.StringAttribute{
					//			Description: "EventPolicyStatisticsType",
					//			Computed:    true,
					//		},
					//		common.ToSnakeCase("EventPolicyStatisticsPeriod"): schema.Int64Attribute{
					//			Description: "EventPolicyStatisticsPeriod",
					//			Computed:    true,
					//		},
					//	},
					//},
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
	}
}

func (d *cloudMonitoringEventPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *cloudMonitoringEventPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) { // 아직 정의하지 않은 Read 메서드를 추가한다.
	var state cloudmonitoring.EventPolicyDataSource

	diags := req.Config.Get(ctx, &state) // datasource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetEventPolicyDetail(state.ResourceType, state.EventPolicyId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Servers : "+detail,
			err.Error(),
		)
		return
	}

	var productSummaryElement = data.GetProductSummary() //	ProductSummary
	productSummary := cloudmonitoring.ProductSummary{
		ProductSq:         types.Int64Value(productSummaryElement.GetProductSq()),
		ProductResourceId: types.StringValue(productSummaryElement.ProductResourceId),
		ProductName:       types.StringValue(productSummaryElement.ProductName),
		ProductTypeCode:   types.StringValue(productSummaryElement.ProductTypeCode),
		ProductTypeName:   types.StringValue(productSummaryElement.ProductTypeName),
		ProductIpAddress:  types.StringValue(productSummaryElement.GetProductIpAddress()),
		ProductState:      types.StringValue(productSummaryElement.ProductState),
		AgentState:        types.StringValue(productSummaryElement.AgentState),
		ProductSubName:    types.StringValue(productSummaryElement.GetProductSubName()),
		ProductSubType:    types.StringValue(productSummaryElement.GetProductSubType()),
		LbName:            types.StringValue(productSummaryElement.GetLbName()),
		VpcName:           types.StringValue(productSummaryElement.GetVpcName()),
		LbSize:            types.StringValue(productSummaryElement.GetLbSize()),
	}

	// MetricSummary
	var metricSummaryElement = data.GetMetricSummary()
	metricSummary := cloudmonitoring.MetricSummary{
		MetricKey:           types.StringValue(metricSummaryElement.MetricKey),
		MetricName:          types.StringValue(metricSummaryElement.MetricName),
		MetricSetKey:        types.StringValue(metricSummaryElement.MetricSetKey),
		MetricSetName:       types.StringValue(metricSummaryElement.MetricSetName),
		MetricDescription:   types.StringValue(metricSummaryElement.GetMetricDescription()),
		MetricDescriptionEn: types.StringValue(metricSummaryElement.GetMetricDescriptionEn()),
		ProductTargetType:   types.StringValue(metricSummaryElement.GetProductTargetType()),
		ProductTargetTypeEn: types.StringValue(metricSummaryElement.GetProductTargetTypeEn()),
		MetricType:          types.StringValue(metricSummaryElement.MetricType),
		MetricUnit:          types.StringValue(metricSummaryElement.MetricUnit),
		IsObjectExist:       types.StringValue(metricSummaryElement.IsObjectExist),
		IsLogMetric:         types.StringValue(metricSummaryElement.GetIsLogMetric()),
	}

	////	EventPolicyStatistics
	//eventPolicyStatistics := cloudmonitoring.EventPolicyStatistics{
	//	EventPolicyStatisticsType:   types.StringValue(data.GetEventPolicyStatistics().EventPolicyStatisticsType),
	//	EventPolicyStatisticsPeriod: types.Int64Value(data.GetEventPolicyStatistics().EventPolicyStatisticsPeriod),
	//}

	singleThreshold := cloudmonitoring.SingleThreshold{
		ComparisonOperator: types.StringValue(data.EventThreshold.GetSingleThreshold().ComparisonOperator),
		Value:              types.Float64Value(data.EventThreshold.GetSingleThreshold().Value),
	}

	//}
	//rangeThreshold := cloudmonitoring.RangeThreshold{
	//	MaxComparisonOperator: types.StringValue(data.EventThreshold.GetRangeThreshold().MaxComparisonOperator),
	//	MinComparisonOperator: types.StringValue(data.EventThreshold.GetRangeThreshold().MinComparisonOperator),
	//	MaxValue:              types.Float64Value(data.EventThreshold.RangeThreshold.MaxValue),
	//	MinValue:              types.Float64Value(data.EventThreshold.RangeThreshold.MinValue),
	//}

	//	EventThreshold
	eventThresholdElement := data.GetEventThreshold()
	eventThreshold := cloudmonitoring.EventThreshold{
		ThresholdType:   types.StringValue(eventThresholdElement.ThresholdType),
		MetricFunction:  types.StringValue(eventThresholdElement.GetMetricFunction()),
		SingleThreshold: singleThreshold,
		//RangeThreshold:  rangeThreshold,
	}

	eventPolicyDetailModel := cloudmonitoring.EventPolicyDetail{
		//CreatedByName: types.StringValue(data.GetCreatedByName()),
		//CreateBy:      createby,
		//UpdateById: types.StringValue(data.GetUpdateById()),
		//UpdateBy:            createby,
		EventPolicyId:       types.Int64Value(data.EventPolicyId),
		ProductSummary:      productSummary,
		MetricSummary:       metricSummary,
		ProductSq:           types.Int64Value(data.GetProductSq()),
		ProductResourceId:   types.StringValue(data.GetProductResourceId()),
		MetricKey:           types.StringValue(data.MetricKey),
		MetricDescription:   types.StringValue(data.GetMetricDescription()),
		MetricDescriptionEn: types.StringValue(data.GetMetricDescriptionEn()),
		ProductTargetType:   types.StringValue(data.GetProductTargetType()),
		ProductTargetTypeEn: types.StringValue(data.GetProductTargetTypeEn()),
		IsLogMetric:         types.StringValue(data.GetIsLogMetric()),
		ObjectName:          types.StringValue(data.GetObjectName()),
		ObjectDisplayName:   types.StringValue(data.GetObjectDisplayName()),
		EventLevel:          types.StringValue(data.EventLevel),
		FtCount:             types.Int64Value(data.FtCount),
		EventMessagePrefix:  types.StringValue(data.GetEventMessagePrefix()),
		ObjectType:          types.StringValue(data.GetObjectType()),
		ObjectTypeName:      types.StringValue(data.GetObjectTypeName()),
		//ProductInfoAttrs:    types.StringValue(data.GetProductInfoAttrs(),
		DisableObject:      types.StringValue(data.GetDisableObject()),
		UserNames:          types.StringValue(data.GetUserNames()),
		UserNameStr:        types.StringValue(data.GetUserNameStr()),
		DisableYn:          types.StringValue(data.GetDisableYn()),
		AttrListStr:        types.StringValue(data.GetAttrListStr()),
		AsgYn:              types.StringValue(data.GetAsgYn()),
		StartDt:            types.StringValue(data.GetStartDt().String()),
		DisplayEventRule:   types.StringValue(data.GetDisplayEventRule()),
		CheckAsg:           types.BoolValue(data.GetCheckAsg()),
		EventOccurTimeZone: types.StringValue(data.GetEventOccurTimeZone()),
		EventThreshold:     eventThreshold,
		//EventPolicyStatistics: eventPolicyStatistics,
	}

	eventPolicyDetailObjectValue, diags := types.ObjectValueFrom(ctx, eventPolicyDetailModel.AttributeTypes(), eventPolicyDetailModel)

	state.EventPolicyDetail = eventPolicyDetailObjectValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
