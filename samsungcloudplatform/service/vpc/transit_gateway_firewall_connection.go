package vpc

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &tgwFirewallConnectionResource{}
	_ resource.ResourceWithConfigure = &tgwFirewallConnectionResource{}
)

// NewVpcTgwFirewallConnectionResource is a helper function to simplify the provider implementation.
func NewVpcTgwFirewallConnectionResource() resource.Resource {
	return &tgwFirewallConnectionResource{}
}

type tgwFirewallConnectionResource struct {
	config     *scpsdk.Configuration
	clientv1d3 *vpcv1d3.Client
	clients    *client.SCPClient
}

// Metadata returns the data source type name.
func (r *tgwFirewallConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_transit_gateway_firewall_connection"
}

// Schema defines the schema for the data source.
func (r *tgwFirewallConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Transit Gateway Firewall Connection",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("TransitGatewayId"): schema.StringAttribute{
				Description: "The identifier of the transit gateway that the firewall connection belongs to.\n" +
					"  - example : 7df8abb4912e4709b1cb237daccca7a8",
				Required: true,
			},

			// Output
			common.ToSnakeCase("TransitGatewayFirewallConnection"): schema.SingleNestedAttribute{
				Description: "Transit Gateway Firewall Connection",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("FirewallConnectionState"): schema.StringAttribute{
						Description: "Firewall Connection State\n" +
							"  - enum: ATTACHING | ACTIVE | DETACHING | DELETED | INACTIVE | ERROR\n" +
							"  - example: INACTIVE",
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *tgwFirewallConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clientv1d3 = inst.Client.VpcV1Dot3
	r.clients = inst.Client
}

func (r *tgwFirewallConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d3.TransitGatewayFirewallConnectionResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, err := r.clientv1d3.CreateTransitGatewayFirewallConnection(ctx, plan.TransitGatewayId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Transit Gateway Firewall Connection",
			"Could not create Transit Gateway Firewall Connection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	tgwFirewallConnection := vpcv1d3.TransitGatewayFirewallConnection{}
	if ok := data.HasFirewallConnectionState(); ok {
		state := data.GetFirewallConnectionState()
		tgwFirewallConnection.FirewallConnectionState = types.StringValue(string(state))
	} else {
		tgwFirewallConnection.FirewallConnectionState = types.StringNull()
	}

	tgwFirewallConnectionObjectValue, _ := types.ObjectValueFrom(ctx, tgwFirewallConnection.AttributeTypes(), tgwFirewallConnection)
	plan.TransitGatewayFirewallConnection = tgwFirewallConnectionObjectValue

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tgwFirewallConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// does not have detail API
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tgwFirewallConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Not Implemented",
		"TGW Firewall Connection update function is not yet implemented.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tgwFirewallConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcv1d3.TransitGatewayFirewallConnectionResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.clientv1d3.DeleteTransitGatewayFirewallConnection(ctx, state.TransitGatewayId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting TGW Firewall Connection",
			"Could not delete TGW Firewall Connection, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, readReq, readResp)
	resp.State = readResp.State
}
