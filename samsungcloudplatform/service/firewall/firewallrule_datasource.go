package firewall

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewallv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &firewallFirewallRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallFirewallRuleDataSource{}
)

// NewFirewallFirewallRuleDataSource is a helper function to simplify the provider implementation.
func NewFirewallFirewallRuleDataSource() datasource.DataSource {
	return &firewallFirewallRuleDataSource{}
}

// firewallFirewallRuleDataSource is the data source implementation.
type firewallFirewallRuleDataSource struct {
	config  *scpsdk.Configuration
	client  *firewall.Client
	clients *client.SCPClient
}

// Metadata returns the data source type name.
func (d *firewallFirewallRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_firewall_rule"
}

// Schema defines the schema for the data source.
func (d *firewallFirewallRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Firewall rule",
		Attributes: map[string]schema.Attribute{
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the resource.\n" +
					"  - example: 0e2b4ece64944d7d8a72983e945b867b",
				Required: true,
			},
			common.ToSnakeCase("FirewallRule"): schema.SingleNestedAttribute{
				Description: "Firewall Rule.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the resource.\n" +
							"  - example: 0e2b4ece64944d7d8a72983e945b867b",
						Computed: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The identifier of the firewall associated with the resource.\n" +
							"  - example: 68db67f78abd405da98a6056a8ee42af",
						Computed: true,
					},
					common.ToSnakeCase("Sequence"): schema.Int32Attribute{
						Description: "The order in which the rule is evaluated.\n" +
							"  - example: 100",
						Computed: true,
					},
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "The source IP addresses the rule applies to.\n" +
							"  - example: [10.10.10.0/24, 10.10.11.0/24]",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
						Description: "The destination address the rule applies to.\n" +
							"  - example: [192.168.0.0/16, 192.169.0.0/16]",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("Service"): schema.ListNestedAttribute{
						Description: "The service ports the rule applies to.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								common.ToSnakeCase("ServiceType"): schema.StringAttribute{
									Description: "The type of the service.\n" +
										"  - example: TCP",
									Computed: true,
								},
								common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
									Description: "The value of the service.\n" +
										"  - example: 80",
									Computed: true,
								},
							},
						},
					},
					common.ToSnakeCase("Action"): schema.StringAttribute{
						Description: "The action applied to traffic that matches the rule.\n" +
							"  - example: ALLOW",
						Computed: true,
					},
					common.ToSnakeCase("Direction"): schema.StringAttribute{
						Description: "The direction of the traffic the rule applies to.\n" +
							"  - example: INBOUND",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "A brief explanation or note about this resource.\n" +
							"  - example: Firewall rule for web tier",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current state of the resource.\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("Status"): schema.StringAttribute{
						Description: "The current status of the resource.\n" +
							"  - example: ENABLE",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
							"  - example: 2025-01-15T10:30:00Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user ID that created the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
							"  - example: 2025-06-01T14:22:00Z",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user ID that modified the resource.\n" +
							"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *firewallFirewallRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.FirewallV1d2
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *firewallFirewallRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallRuleDataSource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataResp, err := d.clients.FirewallV1d2.GetFirewallRule(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewall rule",
			err.Error(),
		)
		return
	}

	firewallRuleModel := createFirewallRuleModelv1D1(dataResp)

	firewallRuleObjectValue, objDiags := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
	resp.Diagnostics.Append(objDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FirewallRule = firewallRuleObjectValue

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
