package servicewatch

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	servicewatch2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/servicewatch/1.2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceWatchDashboardDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceWatchDashboardDataSource{}
)

// NewServiceWatchDashboardDataSource is a helper function to simplify the provider implementation.
func NewServiceWatchDashboardDataSource() datasource.DataSource {
	return &serviceWatchDashboardDataSource{}
}

// serviceWatchDashboardDataSource is the data source implementation.
type serviceWatchDashboardDataSource struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchDashboardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_dashboard"
}

// Schema defines the schema for the data source.
func (d *serviceWatchDashboardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Dashboard Data Source",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "Dashboard ID.\n" +
					" - example : b48e730a70e74f6aa3d2555000b5c22b\n",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Dashboard name.\n" +
					" - example : Production-Web-Servers\n",
				Optional: true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Dashboard type.\n" +
					" - example : Custom\n",
				Optional: true,
			},
			common.ToSnakeCase("FavoriteEnabled"): schema.BoolAttribute{
				Description: "Whether it is a favorite dashboard.\n" +
					" - example : true\n",
				Optional: true,
			},
			common.ToSnakeCase("Srn"): schema.StringAttribute{
				Description: "Service resource name.\n" +
					" - example : srn:dev2::account-id:region::scp-servicewatch:dashboard/b48e730a70e74f6aa3d2555000b5c22b\n",
				Optional:    true,
			},
			common.ToSnakeCase("ShareType"): schema.StringAttribute{
				Description: "Sharing type.\n" +
					" - example : Private\n",
				Optional: true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
					" - example : 2024-05-17T00:23:17Z\n",
				Computed: true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "The user id that created the resource.\n" +
					" - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
				Computed: true,
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{

				Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
					" - example : 2024-05-17T00:23:17Z\n",
				Computed: true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description: "The user id that last modified the resource.\n" +
					" - example : 90dddfc2b1e04edba54ba2b41539a9ac\n",
				Computed: true,
			},
			common.ToSnakeCase("Widgets"): schema.ListNestedAttribute{
				Description: "List of widgets.\n" +
					" - example : [{\"id\": \"75da70a1a4fb486ab0282cf90693ec3c\", \"type\": \"metric\"}]\n",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Widget ID.\n" +
								" - example : 75da70a1a4fb486ab0282cf90693ec3c\n",
							Computed: true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Widget type.\n" +
								" - example : metric\n",
							Computed: true,
						},
						common.ToSnakeCase("Width"): schema.Int32Attribute{
							Description: "Widget width.\n" +
								" - example : 1\n",
							Computed: true,
						},
						common.ToSnakeCase("Height"): schema.Int32Attribute{
							Description: "Widget height.\n" +
								" - example : 1\n",
							Computed: true,
						},
						common.ToSnakeCase("Order"): schema.Int32Attribute{
							Description: "Widget's order in the dashboard.\n" +
								" - example : 1\n",
							Computed: true,
						},
						common.ToSnakeCase("Properties"): schema.SingleNestedAttribute{
							Description: "Widget's detailed properties.\n" +
								" - example : {\"title\": \"CPU Utilization\", \"view\": \"line\"}\n",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Title"): schema.StringAttribute{
									Description: "Widget title.\n" +
										" - example : Virtual Server | CPU Utilization\n",
									Computed: true,
								},
								common.ToSnakeCase("Period"): schema.Int32Attribute{
									Description: "Query period (seconds).\n" +
										" - example : 300\n",
									Computed: true,
								},
								common.ToSnakeCase("Stacked"): schema.BoolAttribute{
									Description: "Whether the graph is stacked.\n" +
										" - example : false\n",
									Computed: true,
								},
								common.ToSnakeCase("StatisticType"): schema.StringAttribute{
									Description: "Statistical function.\n" +
										" - example : AVG\n",
									Computed: true,
								},
								common.ToSnakeCase("View"): schema.StringAttribute{
									Description: "View type.\n" +
										" - example : line\n",
									Computed: true,
								},
								common.ToSnakeCase("Metrics"): schema.ListNestedAttribute{
									Description: "List of metrics included in the widget.\n" +
										" - example : [{\"name\": \"cpu_utilization\", \"namespace_name\": \"SCP/Compute\"}]\n",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("Name"): schema.StringAttribute{
												Description: "Metric name.\n" +
													" - example : cpu_utilization\n",
												Computed:    true,
											},
											common.ToSnakeCase("NamespaceName"): schema.StringAttribute{
												Description: "The name of the namespace.\n" +
													" - example : SCP/Compute\n",
												Computed:    true,
											},
											common.ToSnakeCase("DisplayName"): schema.StringAttribute{
												Description: "Display name (label) of the metric.\n" +
													" - example : CPU Utilization (%)\n",
												Computed:    true,
											},
											common.ToSnakeCase("Color"): schema.StringAttribute{
													Description: "Metric line color.\n" +
														" - example : #2ecc71\n",
												Computed:    true,
											},
                                            common.ToSnakeCase("Dimensions"): schema.ListNestedAttribute{
                                                Description: "List of dimensions.\n" +
                                                    " - example : [{\"key\": \"instance_id\", \"value\": \"i-1234567890abcdef0\"}]\n",
                                                Optional:    true,
                                                NestedObject: schema.NestedAttributeObject{
                                                    Attributes: map[string]schema.Attribute{
                                                        common.ToSnakeCase("Key"): schema.StringAttribute{
                                                            Description: "Dimension key.\n" +
                                                                " - example : instance_id\n",
                                                            Computed:    true,
                                                        },
                                                        common.ToSnakeCase("Value"): schema.StringAttribute{
                                                            Description: "Dimension value.\n" +
                                                                " - example : i-1234567890abcdef0\n",
                                                            Computed:    true,
                                                        },
                                                    },
                                                },
											},
											common.ToSnakeCase("Period"): schema.Int32Attribute{
														Description: "Query period (seconds).\n" +
															" - example : 300\n",
														Computed:    true,
											},
											common.ToSnakeCase("StatisticType"): schema.StringAttribute{
													Description: "Statistical function.\n" +
														" - example : AVG\n",
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
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceWatchDashboardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			ErrUnexpectedConfigure,
			fmt.Sprintf(ErrUnexpectedConfigureFmt, req.ProviderData),
		)

		return
	}

	d.client = inst.Client.ServiceWatch
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *serviceWatchDashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.DashboardDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dashboard, err := d.client.GetDashboard(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			ErrReadDashboard,
			fmt.Sprintf(ErrReadDashboardFmt, state.Id.ValueString(), err.Error(), detail),
		)
		return
	}

	// convert widget model
	var widgets types.List
	widgetResponses, ok := dashboard.GetWidgetsOk()
	if ok && widgetResponses != nil {
		widgets, diags = convertWidget(ctx, widgetResponses)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// convert dashboard basic model
	state.Id = types.StringValue(dashboard.GetId())
	state.Name = types.StringValue(dashboard.GetName())
	state.Type = types.StringValue(dashboard.GetType())
	state.FavoriteEnabled = types.BoolValue(dashboard.GetFavoriteEnabled())
	state.Srn = types.StringValue(dashboard.GetSrn())
	state.ShareType = types.StringValue(dashboard.GetShareType())
	state.CreatedAt = types.StringValue(dashboard.GetCreatedAt().Format(TimeFormatDisplay))
	state.ModifiedAt = types.StringValue(dashboard.GetModifiedAt().Format(TimeFormatDisplay))
	state.CreatedBy = types.StringValue(dashboard.GetCreatedBy())
	state.ModifiedBy = types.StringValue(dashboard.GetModifiedBy())
	state.Widgets = widgets

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "ObjectValueFrom failed", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}
}

func convertWidget(ctx context.Context, widgetResponses []servicewatch2.WidgetDetailDtoV1Dot1) (types.List, diag.Diagnostics) {
	if widgetResponses == nil {
		return types.ListNull(types.ObjectType{AttrTypes: servicewatch.Widget{}.AttributeTypes()}), nil
	}

	var items []servicewatch.Widget
	for _, v := range widgetResponses {
		properties, diags := convertProperties(ctx, &v.Properties)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: servicewatch.Widget{}.AttributeTypes()}), diags
		}
		items = append(items, servicewatch.Widget{
			Type:       types.StringValue(v.Type),
			Width:      types.Int32Value(v.Width),
			Height:     types.Int32Value(v.Height),
			Order:      types.Int32Value(v.Order),
			Properties: properties,
		})
	}
	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: servicewatch.Widget{}.AttributeTypes()}, items)
}

func convertProperties(ctx context.Context, property *servicewatch2.PropertiesDtoV1Dot1) (types.Object, diag.Diagnostics) {
	if property == nil {
		return types.ObjectNull(servicewatch.Properties{}.AttributeTypes()), nil
	}
	metricList, diags := convertMetrics(ctx, property.GetMetrics())
	if diags.HasError() {
		return types.ObjectNull(servicewatch.Properties{}.AttributeTypes()), diags
	}
	properties := servicewatch.Properties{
		Title:         types.StringValue(property.Title),
		Stacked:       types.BoolValue(property.Stacked),
		View:          types.StringValue(property.View),
		Period:        nullableInt32Types(property.GetPeriodOk()),
		StatisticType: nullableStringTypes(property.GetStatisticTypeOk()),
		Metrics:       metricList,
	}
	return types.ObjectValueFrom(ctx, servicewatch.Properties{}.AttributeTypes(), properties)
}

func convertMetrics(ctx context.Context, metricResponses []servicewatch2.MetricDtoV1Dot1) (types.List, diag.Diagnostics) {
	if metricResponses == nil {
		return types.ListNull(types.ObjectType{AttrTypes: servicewatch.Metric{}.AttributeTypes()}), nil
	}

	var items []servicewatch.Metric
	for _, v := range metricResponses {
		dimensionList, diags := convertDimensions(ctx, v.Dimensions)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: servicewatch.Metric{}.AttributeTypes()}), diags
		}
		items = append(items, servicewatch.Metric{
			Name:          types.StringValue(v.GetName()),
			NamespaceName: types.StringValue(v.GetNamespaceName()),
			DisplayName:   types.StringValue(v.GetDisplayName()),
			Color:         types.StringValue(v.GetColor()),
			Dimensions:    dimensionList,
			Period:        types.Int32Value(v.GetPeriod()),
			StatisticType: types.StringValue(v.GetStatisticType()),
		})
	}
	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: servicewatch.Metric{}.AttributeTypes()}, items)
}

func convertDimensions(ctx context.Context, dimensionResponse []servicewatch2.DimensionDTO) (types.List, diag.Diagnostics) {
	if dimensionResponse == nil {
		return types.ListNull(types.ObjectType{AttrTypes: servicewatch.Dimension{}.AttributeTypes()}), nil
	}

	var items []servicewatch.Dimension
	for _, v := range dimensionResponse {
		items = append(items, servicewatch.Dimension{
			Key:   types.StringValue(v.GetKey()),
			Value: types.StringValue(v.GetValue()),
		})
	}
	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: servicewatch.Dimension{}.AttributeTypes()}, items)
}
