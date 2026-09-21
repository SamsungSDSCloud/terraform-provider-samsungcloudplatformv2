package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/firewall"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	scpfirewall "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/firewall/1.1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &tgwFirewallResource{}
	_ resource.ResourceWithConfigure = &tgwFirewallResource{}
)

// NewVpcTgwFirewallResource is a helper function to simplify the provider implementation.
func NewVpcTgwFirewallResource() resource.Resource {
	return &tgwFirewallResource{}
}

type tgwFirewallResource struct {
	config          *scpsdk.Configuration
	client          *vpcv1d3.Client
	client_firewall *firewall.Client
	clients         *client.SCPClient
}

// Metadata returns the data source type name.
func (r *tgwFirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_firewall"
}

// Schema defines the schema for the data source.
func (r *tgwFirewallResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Transit Gateway Firewall",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "The identifier of the transit gateway that the firewall belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},
			common.ToSnakeCase("ProductType"): schema.StringAttribute{
				Description: "The type of the firewall service.\n" +
					"  - enum: TGW_IGW | TGW_GGW | TGW_DGW | TGW_BM\n" +
					"  - example: TGW_IGW",
				Required: true,
			},
			common.ToSnakeCase("UplinkActiveZone"): schema.StringAttribute{
				Description: "Uplink Active Zone.\n" +
					"  - example: kr-west1-a",
				Optional: true,
			},
			common.ToSnakeCase("UplinkStandbyZone"): schema.StringAttribute{
				Description: "Uplink Standby Zone.\n" +
					"  - example: kr-west1-b",
				Optional: true,
			},

			// Output
			common.ToSnakeCase("TransitGatewayFirewall"): schema.SingleNestedAttribute{
				Description: "Transit Gateway Firewall details",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"state": schema.StringAttribute{
						Description: "The current lifecycle state of the transit gateway.\n" +
							"  - enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					"uplink_active_zone": schema.StringAttribute{
						Description: "Uplink Active Zone.",
						Computed:    true,
					},
					"uplink_standby_zone": schema.StringAttribute{
						Description: "Uplink Standby Zone.",
						Computed:    true,
					},
					"uplink_zone_state": schema.StringAttribute{
						Description: "The state of the uplink zone.\n" +
							"  - enum: ATTACHING | ACTIVE | DETACHING | DELETED | INACTIVE | ERROR | EDITING",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *tgwFirewallResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.VpcV1Dot3
	r.client_firewall = inst.Client.Firewall
	r.clients = inst.Client
}

func (r *tgwFirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d3.TransitGatewayFirewallResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := r.client.CreateTransitGatewayFirewall(ctx, plan.TransitGatewayId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating CreateTransitGatewayFirewall",
			"Could not create CreateTransitGatewayFirewall, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// Map API response to object
	tgwFirewall := vpcv1d3.TransitGatewayFirewall{
		UplinkActiveZone:  nullableStringValue(data.GetUplinkActiveZoneOk()),
		UplinkStandbyZone: nullableStringValue(data.GetUplinkStandbyZoneOk()),
		UplinkZoneState:   types.StringValue(string(data.GetUplinkZoneState())),
		State:             types.StringValue(string(data.GetState())),
	}

	tgwFirewallObjectValue, d := types.ObjectValueFrom(ctx, tgwFirewall.AttributeTypes(), tgwFirewall)
	resp.Diagnostics.Append(d...)
	plan.TransitGatewayFirewall = tgwFirewallObjectValue

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tgwFirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Does not have detail API
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tgwFirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Implemented",
		"TGW Firewall update function is not yet implemented.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tgwFirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpcv1d3.TransitGatewayFirewallResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var transitGatewayFirewall vpcv1d3.TransitGatewayFirewall
	errR := state.TransitGatewayFirewall.As(ctx, &transitGatewayFirewall, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse TGW firewall",
			fmt.Sprintf("An error occurred while parsing TGW firewall: %s", errR),
		)
		return
	}

	tgwInfo, err := r.clients.VpcV1Dot3.GetTransitGatewayInfo(ctx, state.TransitGatewayId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall created",
			fmt.Sprintf("An error occurred whileGetTransitGatewayInfo  info: %s", err),
		)
		return
	}

	firewallIdsStr, ok := tgwInfo.TransitGateway.GetFirewallIdsOk()
	if !ok || firewallIdsStr == nil || *firewallIdsStr == "" {
		resp.Diagnostics.AddError(
			"Firewall IDs not found",
			"Transit gateway has no firewall IDs associated.",
		)
		return
	}

	firewallIds := strings.Split(*firewallIdsStr, ",")
	var firewallId string
	for _, id := range firewallIds {
		firewallDetail, err := r.clients.FirewallV1d1.ShowFirewall(ctx, id)
		if err == nil {
			if firewallDetail.Firewall.ProductType == scpfirewall.FirewallProductType(state.ProductType.ValueString()) {
				firewallId = id
				break
			}
		}
	}

	if len(firewallId) == 0 {
		errMsg := fmt.Sprintf("no firewall with product_type %q found among firewall IDs: %s",
			state.ProductType.ValueString(), *firewallIdsStr)
		resp.Diagnostics.AddError(
			"Firewall not found",
			errMsg,
		)
		return
	}

	// Delete existing tgw vpc connection
	_, err = r.client.DeleteTransitGatewayFirewall(ctx, state.TransitGatewayId.ValueString(), firewallId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Delete TransitGatewayFirewall",
			"Could not delete TransitGatewayFirewall, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func nullableStringValue(val *string, isSet bool) basetypes.StringValue {
	if isSet && val != nil {
		return types.StringValue(*val)
	}
	return types.StringNull()
}
