package vpc

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	firewall "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/firewall"
	vpcv1d2 "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	client          *vpcv1d2.Client
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

			// Output
			common.ToSnakeCase("TransitGateway"): schema.SingleNestedAttribute{
				Description: "Transit Gateway",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "The identifier of the account that owns the transit gateway.\n" +
							"  - example : 7df8abb4912e4709b1cb237daccca7a8",
						Computed: true,
					},
					common.ToSnakeCase("Bandwidth"): schema.Int32Attribute{
						Description: "The bandwidth capacity of the connection.\n" +
							"  - example: 1",
						Computed: true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "The timestamp when the transit gateway was created in ISO 8601 format. \n" +
							"  - example : 2024-05-17T00:23:17Z",
						Computed: true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "The user id that created the transit gateway. \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Enter a brief explanation or note about this transit gateway. This help identify the purpose or usage of the resource.\n" +
							"  - example : TransitGateway Description",
						Computed: true,
					},
					common.ToSnakeCase("FirewallConnectionState"): schema.StringAttribute{
						Description: "The current lifecycle state of the firewall connection. \n" +
							"  - enum: ATTACHING | ACTIVE | DETACHING | DELETED | INACTIVE | ERROR\n" +
							"  - example: INACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("FirewallIds"): schema.StringAttribute{
						Description: "Firewall ID list\n" +
							"  - example: bbb93aca123f4bb2b2c0f206f4a86b2b",
						Computed: true,
					},
					common.ToSnakeCase("FirewallId"): schema.StringAttribute{
						Description: "The identifier of the firewall associated with the transit gateway.\n" +
							"  - example: bbb93aca123f4bb2b2c0f206f4a86b2b",
						Computed: true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "The unique identifier of the transit gateway.\n" +
							"  - example: fe860e0af0c04dcd8182b84f907f31f4",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "The timestamp when the transit gateway was last modified in ISO 8601 format.\n" +
							"  - example : 2024-05-17T00:23:17Z ",
						Computed: true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "The user id that modified the transit gateway. \n" +
							"  - example : 90dddfc2b1e04edba54ba2b41539a9ac",
						Computed: true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "The name of the transit gateway.\n" +
							"  - minLength: 3\n" +
							"  - maxLength: 20\n" +
							"  - pattern: ^[a-zA-Z0-9-]*$\n" +
							"  - example: TransitGatewayName",
						Computed: true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "The current lifecycle state of the transit gateway.\n" +
							"  - enum: CREATING | ACTIVE | DELETING | DELETED | ERROR | EDITING\n" +
							"  - example: ACTIVE",
						Computed: true,
					},
					common.ToSnakeCase("UplinkEnabled"): schema.BoolAttribute{
						Description: "Whether the uplink is enabled.\n" +
							"  - default: false\n" +
							"  - example: false",
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

	r.client = inst.Client.VpcV1Dot2
	r.client_firewall = inst.Client.Firewall
	r.clients = inst.Client
}

func (r *tgwFirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d2.TransitGatewayFireWallResource
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
	tgw := vpcv1d2.TransitGateway{
		Id:            types.StringValue(data.TransitGateway.Id),
		Name:          types.StringValue(data.TransitGateway.Name),
		AccountId:     types.StringValue(data.TransitGateway.AccountId),
		CreatedAt:     types.StringValue(data.TransitGateway.CreatedAt.Format(time.RFC3339)),
		CreatedBy:     types.StringValue(data.TransitGateway.CreatedBy),
		ModifiedAt:    types.StringValue(data.TransitGateway.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:    types.StringValue(data.TransitGateway.ModifiedBy),
		State:         types.StringValue(string(data.TransitGateway.State)),
		UplinkEnabled: types.BoolPointerValue(data.TransitGateway.UplinkEnabled),
	}
	if data.TransitGateway.Description.IsSet() {
		if val := data.TransitGateway.Description.Get(); val != nil {
			tgw.Description = types.StringValue(*val)
		}
	}
	if data.TransitGateway.FirewallIds.IsSet() {
		if val := data.TransitGateway.FirewallIds.Get(); val != nil {
			tgw.FirewallIds = types.StringValue(*val)
		}
	}
	if data.TransitGateway.Bandwidth.IsSet() {
		if val := data.TransitGateway.Bandwidth.Get(); val != nil {
			tgw.Bandwidth = types.Int32PointerValue(val)
		}
	}
	if data.TransitGateway.FirewallConnectionState.IsSet() {
		if desc := data.TransitGateway.FirewallConnectionState.Get(); desc != nil {
			tgw.FirewallConnectionState = types.StringValue(string(*desc))
		}
	}

	firewalLst, err := r.client_firewall.GetFirewallList(
		types.Int32Value(0),
		types.Int32Value(1),
		types.StringValue(""),
		types.StringValue(data.TransitGateway.Name),
		types.StringValue(""),
		types.ListValueMust(types.StringType, []attr.Value{types.StringValue(plan.ProductType.ValueString())}),
		types.ListNull(types.StringType),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get firewall created",
			fmt.Sprintf("An error occurred while getting firewall info: %s", err),
		)
		return
	}

	// Extract firewall ID from the first element in the list
	if len(firewalLst.Firewalls) > 0 {
		tgw.FirewallId = types.StringValue(firewalLst.Firewalls[0].Id)
	}

	tgwObjectValue, d := types.ObjectValueFrom(ctx, tgw.AttributeTypes(), tgw)
	resp.Diagnostics.Append(d...)
	plan.TransitGateway = tgwObjectValue

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *tgwFirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcv1d2.TransitGatewayFireWallResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var transitGateway vpcv1d2.TransitGateway
	errR := state.TransitGateway.As(ctx, &transitGateway, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse TGW firewall",
			fmt.Sprintf("An error occurred while parsing TGW firewall: %s", errR),
		)
		return
	}

	firewalLst, err := r.client_firewall.GetFirewallList(
		types.Int32Value(0),
		types.Int32Value(1),
		types.StringValue(""),
		types.StringValue(transitGateway.Name.ValueString()),
		types.StringValue(""),
		types.ListValueMust(types.StringType, []attr.Value{types.StringValue(state.ProductType.ValueString())}),
		types.ListNull(types.StringType),
	)
	if err != nil || len(firewalLst.Firewalls) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
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
	var state vpcv1d2.TransitGatewayFireWallResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var transitGateway vpcv1d2.TransitGateway
	errR := state.TransitGateway.As(ctx, &transitGateway, basetypes.ObjectAsOptions{})
	if errR != nil {
		resp.Diagnostics.AddError(
			"Failed to parse TGW firewall",
			fmt.Sprintf("An error occurred while parsing TGW firewall: %s", errR),
		)
		return
	}

	_, err := r.client.DeleteTransitGatewayFirewall(ctx, state.TransitGatewayId.ValueString(), transitGateway.FirewallId.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting TGW Firewall",
			"Could not delete TGW Firewall, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
