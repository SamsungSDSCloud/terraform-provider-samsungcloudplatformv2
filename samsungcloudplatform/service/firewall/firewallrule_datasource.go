package firewall

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
			common.ToSnakeCase("FirewallId"): schema.StringAttribute{
				Description: "The identifier of the firewall associated with the resource.\n" +
					"  - example: 68db67f78abd405da98a6056a8ee42af",
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
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the resource.\n" +
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
					common.ToSnakeCase("SourceInterface"): schema.StringAttribute{
						Description: "The source interface the rule applies to.\n" +
							"  - example: L2FW-DGW2800up",
						Computed: true,
					},
					common.ToSnakeCase("SourceAddress"): schema.ListAttribute{
						Description: "The source IP addresses the rule applies to.\n" +
							"  - example: [10.10.10.0/24, 10.10.11.0/24]",
						Computed:    true,
						ElementType: types.StringType,
					},
					common.ToSnakeCase("DestinationInterface"): schema.StringAttribute{
						Description: "The destination interface the rule applies to.\n" +
							"  - example: L2FW-DGW2800dn",
						Computed: true,
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
					common.ToSnakeCase("VendorRuleId"): schema.StringAttribute{
						Description: "The firewall device's unique identifier for the rule.\n" +
							"  - example: 20",
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

	d.client = inst.Client.Firewall
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

	var defaultSize types.Int32
	var defaultPage types.Int32
	var defaultSort types.String
	var defaultSrcIp types.String
	var defaultDstIp types.String
	var defaultDescription types.String
	var defaultState types.List
	var defaultStatus types.String
	var defaultFetchAll types.Bool

	ids, err := GetFirewallRuleList(d.clients, defaultPage, defaultSize, defaultSort, state.FirewallId,
		defaultSrcIp, defaultDstIp, defaultDescription, defaultState, defaultStatus, defaultFetchAll)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read firewall rule",
			err.Error(),
		)
	}

	if len(ids) > 0 {
		exist := false
		for _, v := range ids {
			if v == state.Id {
				exist = true
				break
			}
		}

		if exist {
			data, err := d.client.GetFirewallRule(state.Id.ValueString()) // client 를 호출한다.
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError(
					"Error Reading Firewall Rule",
					"Could not read Firewall Rule ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
				)
				return
			}

			firewallRuleElement := data.FirewallRule

			sourceAddresses := make([]string, 0, len(firewallRuleElement.SourceAddress))
			for _, address := range firewallRuleElement.SourceAddress {
				sourceAddresses = append(sourceAddresses, address)
			}
			destinationAddresses := make([]string, 0, len(firewallRuleElement.DestinationAddress))
			for _, address := range firewallRuleElement.DestinationAddress {
				destinationAddresses = append(destinationAddresses, address)
			}
			services := make([]firewall.FirewallPort, 0, len(firewallRuleElement.Service))
			for _, service := range firewallRuleElement.Service {
				services = append(services, firewall.FirewallPort{
					ServiceType:  types.StringValue(string(service.ServiceType)),
					ServiceValue: types.StringValue(*service.ServiceValue),
				})
			}

			firewallRuleModel := firewall.FirewallRule{
				Id:                   types.StringValue(firewallRuleElement.Id),
				Name:                 types.StringPointerValue(firewallRuleElement.Name.Get()),
				FirewallId:           types.StringValue(firewallRuleElement.FirewallId),
				Sequence:             types.Int32Value(firewallRuleElement.Sequence),
				SourceInterface:      types.StringValue(firewallRuleElement.SourceInterface),
				SourceAddress:        sourceAddresses,
				DestinationInterface: types.StringValue(firewallRuleElement.DestinationInterface),
				DestinationAddress:   destinationAddresses,
				Service:              services,
				Action:               types.StringValue(string(firewallRuleElement.Action)),
				Direction:            types.StringValue(string(firewallRuleElement.Direction)),
				VendorRuleId:         types.StringValue(firewallRuleElement.VendorRuleId),
				Description:          types.StringPointerValue(firewallRuleElement.Description.Get()),
				State:                types.StringValue(string(firewallRuleElement.State)),
				Status:               types.StringValue(string(firewallRuleElement.Status)),
				CreatedAt:            types.StringValue(firewallRuleElement.CreatedAt.Format(time.RFC3339)),
				CreatedBy:            types.StringValue(firewallRuleElement.CreatedBy),
				ModifiedAt:           types.StringValue(firewallRuleElement.ModifiedAt.Format(time.RFC3339)),
				ModifiedBy:           types.StringValue(firewallRuleElement.ModifiedBy),
			}

			firewallRuleObjectValue, _ := types.ObjectValueFrom(ctx, firewallRuleModel.AttributeTypes(), firewallRuleModel)
			state.FirewallRule = firewallRuleObjectValue
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
