package vpn

import (
	"context"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/client"
	scpvpn "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/library/vpn/1.0"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"time"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpnVpnTunnelResource{}
	_ resource.ResourceWithConfigure = &vpnVpnTunnelResource{}
)

// NewVpnVpnTunnelResource is a helper function to simplify the provider implementation.
func NewVpnVpnTunnelResource() resource.Resource {
	return &vpnVpnTunnelResource{}
}

// vpnVpnTunnelResource is the data source implementation.
type vpnVpnTunnelResource struct {
	config *scpsdk.Configuration
	client *vpn.Client
}

// Metadata returns the data source type name.
func (r *vpnVpnTunnelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpn_vpn_tunnel"
}

// Schema defines the schema for the data source.
func (r *vpnVpnTunnelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Vpn tunnel",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
			"id": schema.StringAttribute{
				Description: "Identifier of the resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			"vpn_gateway_id": schema.StringAttribute{
				Description: "VpnGatewayId",
				Required:    true,
			},
			"phase1": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"diffie_hellman_groups": schema.ListAttribute{
						Description: "DiffieHellmanGroups",
						Required:    true,
						ElementType: types.Int32Type,
					},
					"encryptions": schema.ListAttribute{
						Description: "Encryptions",
						Required:    true,
						ElementType: types.StringType,
					},
					"dpd_retry_interval": schema.Int32Attribute{
						Description: "DpdRetryInterval",
						Required:    true,
					},
					"ike_version": schema.Int32Attribute{
						Description: "IkeVersion",
						Required:    true,
					},
					"life_time": schema.Int32Attribute{
						Description: "LifeTime",
						Required:    true,
					},
					"peer_gateway_ip": schema.StringAttribute{
						Description: "PeerGatewayIp",
						Required:    true,
					},
					"pre_shared_key": schema.StringAttribute{
						Description: "PreSharedKey",
						Required:    true,
					},
				},
			},
			"phase2": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"diffie_hellman_groups": schema.ListAttribute{
						Description: "DiffieHellmanGroups",
						Required:    true,
						ElementType: types.Int32Type,
					},
					"encryptions": schema.ListAttribute{
						Description: "Encryptions",
						Required:    true,
						ElementType: types.StringType,
					},
					"life_time": schema.Int32Attribute{
						Description: "LifeTime",
						Required:    true,
					},
					"perfect_forward_secrecy": schema.StringAttribute{
						Description: "PerfectForwardSecrecy",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("ENABLE", "DISABLE"),
						},
					},
					"remote_subnet": schema.StringAttribute{
						Description: "RemoteSubnet",
						Required:    true,
					},
				},
			},
			common.ToSnakeCase("VpnTunnel"): schema.SingleNestedAttribute{
				Description: "Vpn tunnel",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					common.ToSnakeCase("AccountId"): schema.StringAttribute{
						Description: "AccountId",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
						Description: "CreatedAt",
						Computed:    true,
					},
					common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
						Description: "CreatedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Description"): schema.StringAttribute{
						Description: "Description",
						Computed:    true,
					},
					common.ToSnakeCase("Id"): schema.StringAttribute{
						Description: "Id",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
						Description: "ModifiedAt",
						Computed:    true,
					},
					common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
						Description: "ModifiedBy",
						Computed:    true,
					},
					common.ToSnakeCase("Name"): schema.StringAttribute{
						Description: "Name",
						Computed:    true,
					},
					common.ToSnakeCase("State"): schema.StringAttribute{
						Description: "State",
						Computed:    true,
					},
					common.ToSnakeCase("VpcId"): schema.StringAttribute{
						Description: "VpcId",
						Computed:    true,
					},
					common.ToSnakeCase("VpcName"): schema.StringAttribute{
						Description: "VpcName",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
						Description: "VpnGatewayId",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayIpAddress"): schema.StringAttribute{
						Description: "VpnGatewayIpAddress",
						Computed:    true,
					},
					common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
						Description: "VpnGatewayName",
						Computed:    true,
					},
					common.ToSnakeCase("Phase1"): schema.SingleNestedAttribute{
						Description: "Phase1",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "DiffieHellmanGroups",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"dpd_retry_interval": schema.Int32Attribute{
								Description: "DpdRetryInterval",
								Computed:    true,
							},
							"encryptions": schema.ListAttribute{
								Description: "Encryptions",
								Computed:    true,
								ElementType: types.StringType,
							},
							"ike_version": schema.Int32Attribute{
								Description: "IkeVersion",
								Computed:    true,
							},
							"life_time": schema.Int32Attribute{
								Description: "LifeTime",
								Computed:    true,
							},
							"peer_gateway_ip": schema.StringAttribute{
								Description: "PeerGatewayIp",
								Computed:    true,
							},
							"pre_shared_key": schema.StringAttribute{
								Description: "PreSharedKey",
								Computed:    true,
							},
						},
					},
					common.ToSnakeCase("Phase2"): schema.SingleNestedAttribute{
						Description: "Phase2",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"diffie_hellman_groups": schema.ListAttribute{
								Description: "DiffieHellmanGroups",
								Computed:    true,
								ElementType: types.Int32Type,
							},
							"encryptions": schema.ListAttribute{
								Description: "Encryptions",
								Computed:    true,
								ElementType: types.StringType,
							},
							"life_time": schema.Int32Attribute{
								Description: "LifeTime",
								Computed:    true,
							},
							"perfect_forward_secrecy": schema.StringAttribute{
								Description: "PerfectForwardSecrecy",
								Computed:    true,
							},
							"remote_subnet": schema.StringAttribute{
								Description: "RemoteSubnet",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *vpnVpnTunnelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.Vpn
}

// Create creates the resource and sets the initial Terraform state.
func (r *vpnVpnTunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vpn.VpnTunnelResource
	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpn tunnel
	data, err := r.client.CreateVpnTunnel(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpn tunnel",
			"Could not create vpn tunnel, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vpnTunnel := data.VpnTunnel
	plan.Id = types.StringValue(vpnTunnel.Id)
	diags = resp.State.Set(ctx, plan)

	err = waitForVpnTunnelStatus(ctx, r.client, vpnTunnel.Id, []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vpn tunnel",
			"Error waiting for vpn tunnel to become active: "+err.Error(),
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

// Read refreshes the Terraform state with the latest data.
func (r *vpnVpnTunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vpn.VpnTunnelResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetVpnTunnel(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpn tunnel",
			"Could not read vpn tunnel ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	vt := data.VpnTunnel
	vtModel := vpn.VpnTunnel{
		AccountId:           types.StringValue(vt.AccountId),
		CreatedAt:           types.StringValue(vt.CreatedAt.Format(time.RFC3339)),
		CreatedBy:           types.StringValue(vt.CreatedBy),
		Description:         types.StringPointerValue(vt.Description.Get()),
		Id:                  types.StringValue(vt.Id),
		Phase1:              mapPhase1Detail(vt.Phase1),
		Phase2:              mapPhase2Detail(vt.Phase2),
		ModifiedAt:          types.StringValue(vt.ModifiedAt.Format(time.RFC3339)),
		ModifiedBy:          types.StringValue(vt.ModifiedBy),
		Name:                types.StringValue(vt.Name),
		State:               types.StringValue(string(vt.State)),
		VpcId:               types.StringValue(vt.VpcId),
		VpcName:             types.StringValue(vt.VpcName),
		VpnGatewayId:        types.StringValue(vt.VpnGatewayId),
		VpnGatewayIpAddress: types.StringValue(vt.VpnGatewayIpAddress),
		VpnGatewayName:      types.StringValue(vt.VpnGatewayName),
	}

	vtObjectValue, diags := types.ObjectValueFrom(ctx, vtModel.AttributeTypes(), vtModel)
	state.VpnTunnel = vtObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vpnVpnTunnelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, changedPlan, state vpn.VpnTunnelResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	var nullString types.String
	changedPlan = plan

	if state.Phase1.PeerGatewayIp.Equal(plan.Phase1.PeerGatewayIp) {
		changedPlan.Phase1.PeerGatewayIp = nullString
	}

	if state.Phase2.RemoteSubnet.Equal(plan.Phase2.RemoteSubnet) {
		changedPlan.Phase2.RemoteSubnet = nullString
	}

	// Update existing order
	_, err := r.client.UpdateVpnTunnel(ctx, state.Id.ValueString(), changedPlan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Updating vpn tunnel",
			"Could not update vpn tunnel, unexpected error: "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpnTunnelStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"ACTIVE"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating vpn tunnel",
			"Error waiting for vpn tunnel to become active: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *vpnVpnTunnelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vpn.VpnTunnelResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpn tunnel
	err := r.client.DeleteVpnTunnel(ctx, state.Id.ValueString())
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Deleting vpn tunnel",
			"Could not delete vpn tunnel, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	err = waitForVpnTunnelStatus(ctx, r.client, state.Id.ValueString(), []string{}, []string{"DELETED"})
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError(
			"Error deleting vpn tunnel",
			"Error waiting for vpn tunnel to become deleted: "+err.Error(),
		)
		return
	}
}

func mapPhase1Detail(phase1 scpvpn.VpnPhase1Detail) vpn.VpnPhase1Detail {
	return vpn.VpnPhase1Detail{
		DiffieHellmanGroups: convertToTypesInt32Slice(phase1.DiffieHellmanGroups),
		Encryptions:         convertToTypesStringSlice(phase1.Encryptions),
		DpdRetryInterval:    types.Int32Value(phase1.DpdRetryInterval),
		IkeVersion:          types.Int32Value(phase1.IkeVersion),
		LifeTime:            types.Int32Value(phase1.LifeTime),
		PeerGatewayIp:       types.StringValue(phase1.PeerGatewayIp),
		//PreSharedKey:        types.StringValue("**********"),
	}
}

func mapPhase2Detail(phase2 scpvpn.VpnPhase2Detail) vpn.VpnPhase2Detail {
	return vpn.VpnPhase2Detail{
		DiffieHellmanGroups:   convertToTypesInt32Slice(phase2.DiffieHellmanGroups),
		Encryptions:           convertToTypesStringSlice(phase2.Encryptions),
		LifeTime:              types.Int32Value(phase2.LifeTime),
		PerfectForwardSecrecy: types.StringValue(phase2.PerfectForwardSecrecy),
		RemoteSubnet:          types.StringValue(phase2.RemoteSubnet),
	}
}

func convertToTypesInt32Slice(intSlice []int32) []types.Int32 {
	result := make([]types.Int32, len(intSlice))
	for i, v := range intSlice {
		result[i] = types.Int32Value(v)
	}
	return result
}

func convertToTypesStringSlice(strSlice []string) []types.String {
	result := make([]types.String, len(strSlice))
	for i, v := range strSlice {
		result[i] = types.StringValue(v)
	}
	return result
}

func waitForVpnTunnelStatus(ctx context.Context, vpnClient *vpn.Client, id string, pendingStates []string, targetStates []string) error {
	return client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpnClient.GetVpnTunnel(ctx, id)
		if err != nil {
			return nil, "", err
		}
		return info, string(info.VpnTunnel.State), nil
	})
}
