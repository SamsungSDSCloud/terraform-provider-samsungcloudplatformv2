package servicewatch

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/servicewatch"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
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
				Description: "The name of dashboard.\n" +
					" - example : Production-Web-Servers\n" +
					" - minLength: 3\n" +
					" - maxLength: 512\n",
				Optional:    true,
			},
			common.ToSnakeCase("NameLike"): schema.StringAttribute{
				Description: "Wildcard search for dashboard names.\n" +
					" - example : Production\n",
				Optional:    true,
			},
			common.ToSnakeCase("FavoriteEnabled"): schema.BoolAttribute{
				Description: "Whether it is a favorite dashboard.\n" +
					" - example : true\n",
				Optional:    true,
			},
			common.ToSnakeCase("Type"): schema.StringAttribute{
				Description: "Dashboard type.\n" +
					" - example : Custom\n",
				Optional:    true,
			},
			common.ToSnakeCase("ServiceCode"): schema.StringAttribute{
				Description: "Associated service code.\n" +
					" - example : scp-compute\n",
				Optional:    true,
			},
			common.ToSnakeCase("Dashboards"): schema.ListNestedAttribute{
				Description: "List of dashboards.\n" +
					" - example : [{\"id\": \"b48e730a70e74f6aa3d2555000b5c22b\", \"name\": \"Production-Web-Servers\"}]\n",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Description: "The unique identifier of the dashboard.\n" +
								" - example : b48e730a70e74f6aa3d2555000b5c22b\n",
							Computed:    true,
						},
						common.ToSnakeCase("Name"): schema.StringAttribute{
							Description: "The name of the dashboard.\n" +
								" - example : Production-Web-Servers\n" +
								" - minLength: 3\n" +
								" - maxLength: 512\n",
							Computed:    true,
						},
						common.ToSnakeCase("Type"): schema.StringAttribute{
							Description: "Dashboard type.\n" +
								" - example : Custom\n",
							Computed:    true,
						},
						common.ToSnakeCase("FavoriteEnabled"): schema.BoolAttribute{
							Description: "Whether it is a favorite dashboard.\n" +
								" - example : true\n",
							Computed:    true,
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was created, in ISO 8601 format.\n" +
								" - example : 2024-05-17T00:23:17Z\n",
							Computed:    true,
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Description: "The timestamp when the resource was last modified, in ISO 8601 format.\n" +
								" - example : 2024-05-17T00:23:17Z\n",
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