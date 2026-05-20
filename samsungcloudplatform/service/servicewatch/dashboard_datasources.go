package servicewatch

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceWatchDashboardDataSources{}
	_ datasource.DataSourceWithConfigure = &serviceWatchDashboardDataSources{}
)

// NewServiceWatchDashboardDataSources is a helper function to simplify the provider implementation.
func NewServiceWatchDashboardDataSources() datasource.DataSource {
	return &serviceWatchDashboardDataSources{}
}

// serviceWatchDashboardDataSources is the data source implementation.
type serviceWatchDashboardDataSources struct {
	config  *scpsdk.Configuration
	client  *servicewatch.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *serviceWatchDashboardDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicewatch_dashboards"
}

// Schema defines the schema for the data source.
func (d *serviceWatchDashboardDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Dashboard Data Sources",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Dashboard name",
				Optional:    true,
			},
			common.ToSnakeCase("NameLike"): schema.StringAttribute{
				Description: "Wildcard search for dashboard names",
				Optional:    true,
			},
			common.ToSnakeCase("FavoriteEnabled"): schema.BoolAttribute{
				Description: "Whether it is a favorite dashboard",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Dashboard type",
				Optional:    true,
			},
			common.ToSnakeCase("ServiceCode"): schema.StringAttribute{
				Description: "Associated service code",
				Optional:    true,
			},
			common.ToSnakeCase("Dashboards"): schema.ListNestedAttribute{
				Description: "Dashboards",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "Dashboard ID",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "Dashboard name",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Dashboard type",
							Computed:    true,
						},
						common.ToSnakeCase("FavoriteEnabled"): schema.BoolAttribute{
							Description: "Whether it is a favorite dashboard",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "Created date time",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "Modified date time",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *serviceWatchDashboardDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *serviceWatchDashboardDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicewatch.DashboardDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetDashboardList(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Dashboard",
			err.Error(),
		)
	}
	if data == nil {
		state.Dashboards = types.ListValueMust(types.ObjectType{AttrTypes: servicewatch.Dashboard{}.AttributeTypes()}, []attr.Value{})

	} else {
		var dashboardList []servicewatch.Dashboard
		for _, dashboard := range data.Dashboards {
			dashboardState := servicewatch.Dashboard{
				Id:              types.StringValue(dashboard.GetId()),
				Name:            types.StringValue(dashboard.GetName()),
				Type:            types.StringValue(dashboard.GetType()),
				FavoriteEnabled: types.BoolValue(dashboard.GetFavoriteEnabled()),
				CreatedAt:       types.StringValue(dashboard.GetCreatedAt().Format(TimeFormatDisplay)),
				ModifiedAt:      types.StringValue(dashboard.GetModifiedAt().Format(TimeFormatDisplay)),
			}
			dashboardList = append(dashboardList, dashboardState)
		}
		state.Dashboards, diags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: servicewatch.Dashboard{}.AttributeTypes()}, dashboardList)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
