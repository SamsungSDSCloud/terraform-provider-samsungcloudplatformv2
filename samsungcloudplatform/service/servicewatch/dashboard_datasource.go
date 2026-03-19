package servicewatch

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	servicewatch2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/servicewatch/1.2"
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
				Description: "Dashboard ID",
				Optional:    true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Dashboard name",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Dashboard type",
				Optional:    true,
			},
			common.ToSnakeCase("FavoriteEnabled"): schema.BoolAttribute{
				Description: "Whether it is a favorite dashboard",
				Optional:    true,
			},
			common.ToSnakeCase("Srn"): schema.StringAttribute{
				Description: "Service resource name",
				Optional:    true,
			},
			common.ToSnakeCase("ShareType"): schema.StringAttribute{
				Description: "Sharing type",
				Optional:    true,
			},
			common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
				Description: "Created date time",
				Computed:    true,
			},
			common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
				Description: "Creator ID",
				Computed:    true,
			},
			common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{

				Description: "Modified date time",
				Computed:    true,
			},
			common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
				Description: "Modifier ID",
				Computed:    true,
			},
			common.ToSnakeCase("Widgets"): schema.ListNestedAttribute{
				Description: "List of widgets",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Widget ID",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Widget type",
							Computed:    true,
						},
						common.ToSnakeCase("Width"): schema.Int32Attribute{
							Description: "Widget width",
							Computed:    true,
						},
						common.ToSnakeCase("Height"): schema.Int32Attribute{
							Description: "Widget height",
							Computed:    true,
						},
						common.ToSnakeCase("Order"): schema.Int32Attribute{
							Description: "Widget's order in the dashboard",
							Computed:    true,
						},
						common.ToSnakeCase("Properties"): schema.SingleNestedAttribute{
							Description: "Widget's detailed properties",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("Title"): schema.StringAttribute{
									Description: "Widget title",
									Computed:    true,
								},
								common.ToSnakeCase("Period"): schema.Int32Attribute{
									Description: "Query period (seconds)",
									Computed:    true,
								},
								common.ToSnakeCase("Stacked"): schema.BoolAttribute{
									Description: "Whether the graph is stacked",
									Computed:    true,
								},
								common.ToSnakeCase("StatisticType"): schema.StringAttribute{
									Description: "Statistical function",
									Computed:    true,
								},
								common.ToSnakeCase("View"): schema.StringAttribute{
									Description: "View type",
									Computed:    true,
								},
								common.ToSnakeCase("Metrics"): schema.ListNestedAttribute{
									Description: "List of metrics included in the widget",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											common.ToSnakeCase("Name"): schema.StringAttribute{
												Description: "Metric name",
												Computed:    true,
											},
											common.ToSnakeCase("NamespaceName"): schema.StringAttribute{
												Description: "Namespace name",
												Computed:    true,
											},
											common.ToSnakeCase("DisplayName"): schema.StringAttribute{
												Description: "Display name (label) of the metric",
												Computed:    true,
											},
											common.ToSnakeCase("Color"): schema.StringAttribute{
												Description: "Metric line color",
												Computed:    true,
											},
											common.ToSnakeCase("Dimensions"): schema.ListNestedAttribute{
												Description: "List of dimensions",
												Optional:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														common.ToSnakeCase("Key"): schema.StringAttribute{
															Description: "Dimension key",
															Computed:    true,
														},
														common.ToSnakeCase("Value"): schema.StringAttribute{
															Description: "Dimension value",
															Computed:    true,
														},
													},
												},
											},
											common.ToSnakeCase("Period"): schema.Int32Attribute{
												Description: "Query period (seconds)",
												Computed:    true,
											},
											common.ToSnakeCase("StatisticType"): schema.StringAttribute{
												Description: "Statistical function",
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
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
			"Error Reading Dashboard",
			"Could not read Dashboard ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
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
	state.CreatedAt = types.StringValue(dashboard.GetCreatedAt().Format("2006-01-02 15:04:05"))
	state.ModifiedAt = types.StringValue(dashboard.GetModifiedAt().Format("2006-01-02 15:04:05"))
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
