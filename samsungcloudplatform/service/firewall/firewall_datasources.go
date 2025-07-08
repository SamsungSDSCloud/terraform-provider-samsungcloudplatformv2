package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &firewallFirewallDataSources{}
	_ datasource.DataSourceWithConfigure = &firewallFirewallDataSources{}
)

func NewFirewallFirewallDataSources() datasource.DataSource {
	return &firewallFirewallDataSources{}
}

type firewallFirewallDataSources struct {
	config  *scpsdk.Configuration
	client  *firewall.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *firewallFirewallDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewalls"
}

// Schema defines the schema for the data source.

func (d *firewallFirewallDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of firewall",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Page"): schema.Int32Attribute{
				Description: "Page \n" +
					"  - example : 0",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "Size \n" +
					"  - example : 20",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "Sort \n" +
					"  - example : created_at:desc",
				Optional: true,
			},
			common.ToSnakeCase("Name"): schema.StringAttribute{
				Description: "Firewall Name \n" +
					"  - example : subnetName",
				Optional: true,
			},
			common.ToSnakeCase("VpcName"): schema.StringAttribute{
				Description: "Vpc Name \n" +
					"  - example : vpcName",
				Optional: true,
			},
			common.ToSnakeCase("ProductType"): schema.ListAttribute{
				Description: "ProductType \n" +
					"  - example : IGW | GGW | DGW | LB | SIGW | TGW_IGW | TGW_GGW | TGW_DGW | TGW_SIGW | TGW_BM",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("IGW", "GGW", "DGW", "LB", "SIGW", "TGW_IGW", "TGW_GGW", "TGW_DGW", "TGW_SIGW", "TGW_BM"),
					),
				},
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				Description: "State \n" +
					"  - example : CREATING | ACTIVE | EDITING | DELETING | ERROR | DEPLOYING",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("CREATING", "ACTIVE", "EDITING", "DELETING", "ERROR", "DEPLOYING"),
					),
				},
			},
			common.ToSnakeCase("Ids"): schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Firewall Id List",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *firewallFirewallDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Firewall
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *firewallFirewallDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetFirewallList(d.clients, state.Page, state.Size, state.Sort, state.Name, state.VpcName, state.ProductType, state.State)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewalls.",
			err.Error(),
		)
	}

	state.Ids = ids

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetFirewallList(clients *client.SCPClient, page types.Int32, size types.Int32, sort types.String, name types.String, vpcName types.String, productType types.List, state types.List) ([]types.String, error) {

	data, err := clients.Firewall.GetFirewallList(page, size, sort, name, vpcName, productType, state)
	if err != nil {
		return nil, err
	}

	contents := data.Firewalls

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
