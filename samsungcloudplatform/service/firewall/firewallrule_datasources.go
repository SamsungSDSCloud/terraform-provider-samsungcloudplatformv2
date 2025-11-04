package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
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
	_ datasource.DataSource              = &firewallFirewallRuleDataSources{}
	_ datasource.DataSourceWithConfigure = &firewallFirewallRuleDataSources{}
)

// NewFirewallFirewallRuleDataSources is a helper function to simplify the provider implementation.
func NewFirewallFirewallRuleDataSources() datasource.DataSource {
	return &firewallFirewallRuleDataSources{}
}

// firewallFirewallRuleDataSources is the data source implementation.
type firewallFirewallRuleDataSources struct {
	config  *scpsdk.Configuration
	client  *firewall.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *firewallFirewallRuleDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewall_rules"
}

// Schema defines the schema for the data source.
func (d *firewallFirewallRuleDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of firewall rule",
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
			common.ToSnakeCase("FirewallId"): schema.StringAttribute{
				Description: "Firewall Id \n" +
					"  - example : 68db67f78abd405da98a6056a8ee42af",
				Required: true,
			},
			common.ToSnakeCase("SrcIp"): schema.StringAttribute{
				Description: "Source IP \n" +
					"  - example : 10.10.10.10",
				Optional: true,
			},
			common.ToSnakeCase("DstIp"): schema.StringAttribute{
				Description: "Destination IP \n" +
					"  - example : 10.10.10.10",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "Description \n" +
					"  - example : firewallDescription",
				Optional: true,
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				Description: "State \n" +
					"  - example : CREATING | ACTIVE | DELETING | EDITING | ERROR",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("CREATING", "ACTIVE", "DELETING", "EDITING", "ERROR"),
					),
				},
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "Status \n" +
					"  - example : ENABLE | DISABLE",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ENABLE", "DISABLE"),
				},
			},
			common.ToSnakeCase("FetchAll"): schema.BoolAttribute{
				Description: "Fetch All \n" +
					"  - example : True | False",
				Optional: true,
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
func (d *firewallFirewallRuleDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *firewallFirewallRuleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallRuleDataSourceIds

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := GetFirewallRuleList(d.clients, state.Page, state.Size, state.Sort, state.FirewallId,
		state.SrcIp, state.DstIp, state.Description, state.State, state.Status, state.FetchAll)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewall rules.",
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

func GetFirewallRuleList(clients *client.SCPClient, page types.Int32, size types.Int32, sort types.String, firewallId types.String,
	srcIp types.String, dstIp types.String, description types.String, state types.List, status types.String, fetchAll types.Bool) ([]types.String, error) {

	data, err := clients.Firewall.GetFirewallRuleList(page, size, sort, firewallId, srcIp, dstIp, description, state, status, fetchAll)
	if err != nil {
		return nil, err
	}

	contents := data.FirewallRules

	var ids []types.String

	// Map response body to model
	for _, content := range contents {
		ids = append(ids, types.StringValue(content.Id))
	}

	return ids, nil
}
