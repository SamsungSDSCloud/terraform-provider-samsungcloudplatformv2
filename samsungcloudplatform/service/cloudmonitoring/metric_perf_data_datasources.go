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

// ListMetricPerfData
// [POST] /v1/cloudmonitorings/product/v2/metric-data

type cloudMonitoringMetricPerfDataDataSources struct {
	config  *scpsdk.Configuration
	client  *cloudmonitoring.Client
	clients *client.SCPClient
}

var (
	_ datasource.DataSource              = &cloudMonitoringMetricPerfDataDataSources{}
	_ datasource.DataSourceWithConfigure = &cloudMonitoringMetricPerfDataDataSources{}
)

func NewCloudMonitoringMetricPerfDataDataSources() datasource.DataSource {
	return &cloudMonitoringMetricPerfDataDataSources{}
}

func (d *cloudMonitoringMetricPerfDataDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudmonitoring_metricperfdatas"
}

func (d *cloudMonitoringMetricPerfDataDataSources) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Schema of MetricPerfDataDataSources.",
		Attributes: map[string]schema.Attribute{
			// request
			common.ToSnakeCase("XResourceType"): schema.StringAttribute{
				Description: "X-ResourceType",
				Required:    true,
			},
			common.ToSnakeCase("IgnoreInvalid"): schema.StringAttribute{
				Description: "IgnoreInvalid",
				Optional:    true,
			},
			common.ToSnakeCase("MetricKey"): schema.StringAttribute{
				Description: "MetricKey",
				Required:    true,
			},
			common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
				Description: "ProductResourceId",
				Required:    true,
			},
			common.ToSnakeCase("ObjectType"): schema.StringAttribute{
				Description: "ObjectType",
				Optional:    true,
			},
			common.ToSnakeCase("ObjectList"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "ObjectList",
				Optional:    true,
			},
			common.ToSnakeCase("StatisticsPeriod"): schema.Int32Attribute{
				Description: "StatisticsPeriod",
				Optional:    true,
			},
			common.ToSnakeCase("StatisticsTypeList"): schema.ListAttribute{
				ElementType: types.StringType,
				Description: "StatisticsTypeList",
				Optional:    true,
			},
			// 계층구조 입력 불가
			//common.ToSnakeCase("MetricDataConditions"): schema.ListNestedAttribute{
			//	Description: "MetricDataConditions",
			//	Computed:    true,
			//	NestedObject: schema.NestedAttributeObject{
			//		Attributes: map[string]schema.Attribute{
			//			common.ToSnakeCase("MetricKey"): schema.StringAttribute{
			//				Description: "MetricKey",
			//				Required:    true,
			//			},
			//			common.ToSnakeCase("ProductResourceInfos"): schema.ListNestedAttribute{
			//				Description: "ProductResourceInfos",
			//				Computed:    true,
			//				NestedObject: schema.NestedAttributeObject{
			//					Attributes: map[string]schema.Attribute{
			//						common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
			//							Description: "ProductResourceId",
			//							Required:    true,
			//						},
			//					},
			//				},
			//			},
			//			common.ToSnakeCase("StatisticsPeriod"): schema.Int32Attribute{
			//				Description: "StatisticsPeriod",
			//				Optional:    true,
			//			},
			//			common.ToSnakeCase("StatisticsTypeList"): schema.ListAttribute{
			//				ElementType: types.StringType,
			//				Description: "StatisticsTypeList",
			//				Optional:    true,
			//			},
			//		},
			//	},
			//},
			common.ToSnakeCase("QueryEndDt"): schema.StringAttribute{
				Description: "QueryEndDt",
				Required:    true,
			},
			common.ToSnakeCase("QueryStartDt"): schema.StringAttribute{
				Description: "QueryStartDt",
				Required:    true,
			},

			// response
			common.ToSnakeCase("MetricPerfDatas"): schema.ListNestedAttribute{
				Description: "MetricPerfDatas",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("MetricKey"): schema.StringAttribute{
							Description: "MetricKey",
							Optional:    true,
						},
						common.ToSnakeCase("MetricName"): schema.StringAttribute{
							Description: "MetricName",
							Optional:    true,
						},
						common.ToSnakeCase("MetricType"): schema.StringAttribute{
							Description: "MetricType",
							Optional:    true,
						},
						common.ToSnakeCase("MetricUnit"): schema.StringAttribute{
							Description: "MetricUnit",
							Optional:    true,
						},
						common.ToSnakeCase("ObjectDisplayName"): schema.StringAttribute{
							Description: "ObjectDisplayName",
							Optional:    true,
						},
						common.ToSnakeCase("ObjectName"): schema.StringAttribute{
							Description: "ObjectName",
							Optional:    true,
						},
						common.ToSnakeCase("ObjectType"): schema.StringAttribute{
							Description: "ObjectType",
							Optional:    true,
						},
						common.ToSnakeCase("PerfData"): schema.ListNestedAttribute{
							Description: "PerfData",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("ts"): schema.Int64Attribute{
										Description: "ts",
										Optional:    true,
									},
									common.ToSnakeCase("value"): schema.StringAttribute{
										Description: "value",
										Optional:    true,
									},
								},
							},
						},
						common.ToSnakeCase("ProductName"): schema.StringAttribute{
							Description: "ProductName",
							Optional:    true,
						},
						common.ToSnakeCase("ProductResourceId"): schema.StringAttribute{
							Description: "ProductResourceId",
							Optional:    true,
						},
						common.ToSnakeCase("StatisticsPeriod"): schema.Int32Attribute{
							Description: "StatisticsPeriod",
							Required:    true,
						},
						common.ToSnakeCase("StatisticsType"): schema.StringAttribute{
							Description: "StatisticsType",
							Required:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": filter.DataSourceSchema(),
		},
	}
}

func (d *cloudMonitoringMetricPerfDataDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudMonitoringMetricPerfDataDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cloudmonitoring.MetricPerfDataDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.GetMetricPerfDataList(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Servers",
			err.Error(),
		)
		return
	}

	for _, element := range res.Contents {
		resState := cloudmonitoring.MetricPerfData{
			MetricKey:         types.StringValue(element.GetMetricKey()),
			MetricName:        types.StringValue(element.GetMetricName()),
			MetricType:        types.StringValue(element.GetMetricType()),
			MetricUnit:        types.StringValue(element.GetMetricUnit()),
			ObjectDisplayName: types.StringValue(element.GetObjectDisplayName()),
			ObjectName:        types.StringValue(element.GetObjectName()),
			ObjectType:        types.StringValue(element.GetObjectType()),
			ProductName:       types.StringValue(element.GetProductName()),
			ProductResourceId: types.StringValue(element.GetProductResourceId()),
			StatisticsPeriod:  types.Int32Value(element.StatisticsPeriod),
			StatisticsType:    types.StringValue(element.StatisticsType),
		}

		perfDataResponse := element.GetPerfData()
		if perfDataResponse != nil {
			perfDataEntry := make([]cloudmonitoring.MetricPerfDataItem, 0)
			for _, perfData := range perfDataResponse {
				perfDataEntry = append(perfDataEntry, cloudmonitoring.MetricPerfDataItem{
					Ts:    perfData["ts"].(float64),
					Value: perfData["value"].(string),
				})
			}
			resState.PerfData = perfDataEntry
		}

		state.MetricPerfDatas = append(state.MetricPerfDatas, resState)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
