package firewall

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewallv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
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
				Description: "The page number for pagination.\n" +
					"  - example: 1\n" +
					"  - constraints: min: 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Size"): schema.Int32Attribute{
				Description: "The number of items per page.\n" +
					"  - example: 20\n" +
					"  - constraints: min: 1",
				Optional: true,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			common.ToSnakeCase("Sort"): schema.StringAttribute{
				Description: "The sorting criteria.\n" +
					"  - example: created_at:desc\n" +
					"  - valid: field_name:asc or field_name:desc",
				Optional: true,
			},
			common.ToSnakeCase("FirewallId"): schema.StringAttribute{
				Description: "The identifier of the firewall associated with the resource.\n" +
					"  - example: 68db67f78abd405da98a6056a8ee42af",
				Required: true,
			},
			common.ToSnakeCase("SrcIp"): schema.StringAttribute{
				Description: "Source IP.\n" +
					"  - example: 10.10.10.10",
				Optional: true,
			},
			common.ToSnakeCase("DstIp"): schema.StringAttribute{
				Description: "Destination IP.\n" +
					"  - example: 10.10.10.10",
				Optional: true,
			},
			common.ToSnakeCase("Description"): schema.StringAttribute{
				Description: "A brief explanation or note about this resource.\n" +
					"  - example: Firewall rule for web tier\n" +
					"  - constraints: maxLength: 100",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			common.ToSnakeCase("State"): schema.ListAttribute{
				Description: "The current state of the resource.\n" +
					"  - example: ACTIVE\n" +
					"  - valid: CREATING, ACTIVE, DELETING, EDITING, ERROR",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.OneOf("CREATING", "ACTIVE", "DELETING", "EDITING", "DELETED", "ERROR"),
					),
				},
			},
			common.ToSnakeCase("Status"): schema.StringAttribute{
				Description: "The current status of the resource.\n" +
					"  - example: ENABLE\n" +
					"  - valid: ENABLE, DISABLE",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ENABLE", "DISABLE"),
				},
			},
			// Output
			common.ToSnakeCase("TotalCount"): schema.Int32Attribute{
				Computed: true,
				Description: "The total number of firewall rules.\n" +
					"  - example: 20",
			},
			common.ToSnakeCase("FirewallRules"): schema.ListNestedAttribute{
				Computed: true,
				Description: "A list of firewall rules.\n" +
					"  - example: see firewall_rules block below",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						common.ToSnakeCase("Id"): schema.StringAttribute{
							Computed: true,
							Description: "The unique identifier of the firewall rule.\n" +
								"  - example: 7df8abb4912e4709b1cb237daccca7a8",
						},
						common.ToSnakeCase("FirewallId"): schema.StringAttribute{
							Computed: true,
							Description: "The identifier of the firewall associated with the rule.\n" +
								"  - example: 68db67f78abd405da98a6056a8ee42af",
						},
						common.ToSnakeCase("Action"): schema.StringAttribute{
							Computed: true,
							Description: "The action to take when a packet matches this rule.\n" +
								"  - example: ALLOW\n" +
								"  - valid: ALLOW, DENY",
						},
						common.ToSnakeCase("Direction"): schema.StringAttribute{
							Computed: true,
							Description: "The direction of traffic this rule applies to.\n" +
								"  - example: INBOUND\n" +
								"  - valid: INBOUND, OUTBOUND",
						},
						common.ToSnakeCase("Sequence"): schema.Int32Attribute{
							Computed: true,
							Description: "The order in which this rule is evaluated.\n" +
								"  - example: 1",
						},
						common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The source IP addresses or CIDR ranges.\n" +
								"  - example: [\"10.0.0.0/24\"]",
						},
						common.ToSnakeCase("DestinationAddress"): schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "The destination IP addresses or CIDR ranges.\n" +
								"  - example: [\"192.168.0.0/16\"]",
						},
						common.ToSnakeCase("Service"): schema.ListNestedAttribute{
							Computed: true,
							Description: "The service/port configuration for this rule.\n" +
								"  - example: see service block below",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									common.ToSnakeCase("ServiceType"): schema.StringAttribute{
										Computed: true,
										Description: "The type of the service.\n" +
											"  - example: TCP",
									},
									common.ToSnakeCase("ServiceValue"): schema.StringAttribute{
										Computed: true,
										Description: "The value of the service, typically a port number or range.\n" +
											"  - example: 80",
									},
								},
							},
						},
						common.ToSnakeCase("Description"): schema.StringAttribute{
							Computed: true,
							Description: "A brief explanation or note about this firewall rule.\n" +
								"  - example: Allow HTTP traffic",
						},
						common.ToSnakeCase("State"): schema.StringAttribute{
							Computed: true,
							Description: "The current lifecycle state of the firewall rule.\n" +
								"  - example: ACTIVE\n" +
								"  - valid: CREATING, ACTIVE, DELETING, EDITING, ERROR",
						},
						common.ToSnakeCase("Status"): schema.StringAttribute{
							Computed: true,
							Description: "The current status of the firewall rule.\n" +
								"  - example: ENABLE\n" +
								"  - valid: ENABLE, DISABLE",
						},
						common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
							Computed: true,
							Description: "The timestamp when the firewall rule was created.\n" +
								"  - example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
							Computed: true,
							Description: "The user ID that created the firewall rule.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
						common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
							Computed: true,
							Description: "The timestamp when the firewall rule was last modified.\n" +
								"  - example: 2024-05-17T00:23:17Z",
						},
						common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
							Computed: true,
							Description: "The user ID that last modified the firewall rule.\n" +
								"  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
						},
					},
				},
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

	d.client = inst.Client.FirewallV1d2
	d.clients = inst.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *firewallFirewallRuleDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state firewall.FirewallRuleDataSources

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataResp, err := d.clients.FirewallV1d2.GetFirewallRuleList(state.Page, state.Size, state.Sort, state.FirewallId,
		state.SrcIp, state.DstIp, state.Description, state.State, state.Status)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewall rules.",
			err.Error(),
		)
		return
	}

	firewallRules := make([]firewall.FirewallRule, 0, len(dataResp.FirewallRules))

	for _, fwRes := range dataResp.FirewallRules {
		sourceAddresses := make([]string, 0, len(fwRes.SourceAddress))
		for _, addr := range fwRes.SourceAddress {
			sourceAddresses = append(sourceAddresses, addr)
		}

		destinationAddresses := make([]string, 0, len(fwRes.DestinationAddress))
		for _, addr := range fwRes.DestinationAddress {
			destinationAddresses = append(destinationAddresses, addr)
		}

		services := make([]firewall.FirewallPort, 0, len(fwRes.Service))
		for _, service := range fwRes.Service {
			services = append(services, firewall.FirewallPort{
				ServiceType:  types.StringValue(string(service.ServiceType)),
				ServiceValue: types.StringPointerValue(service.ServiceValue),
			})
		}

		firewallRuleModel := firewall.FirewallRule{
			Id:                 types.StringValue(fwRes.Id),
			FirewallId:         types.StringValue(fwRes.FirewallId),
			Action:             types.StringValue(string(fwRes.Action)),
			Direction:          types.StringValue(string(fwRes.Direction)),
			Sequence:           types.Int32Value(fwRes.Sequence),
			SourceAddress:      sourceAddresses,
			DestinationAddress: destinationAddresses,
			Service:            services,
			Description:        types.StringPointerValue(fwRes.Description.Get()),
			State:              types.StringValue(string(fwRes.State)),
			Status:             types.StringValue(string(fwRes.Status)),
			CreatedAt:          types.StringValue(fwRes.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(fwRes.CreatedBy),
			ModifiedAt:         types.StringValue(fwRes.ModifiedAt.Format(time.RFC3339)),
			ModifiedBy:         types.StringValue(fwRes.ModifiedBy),
		}

		firewallRules = append(firewallRules, firewallRuleModel)
	}

	state.FirewallRules = firewallRules
	state.TotalCount = types.Int32Value(dataResp.Count)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
