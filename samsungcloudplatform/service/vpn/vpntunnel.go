package vpn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpn"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/tag"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpvpn "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/vpn/1.1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vpnVpnTunnelResource{}
	_ resource.ResourceWithConfigure = &vpnVpnTunnelResource{}
	_ resource.ResourceWithImportState = &vpnVpnTunnelResource{}
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
		Description: "Manages VPN tunnels for encrypted site-to-site connections.",
		Attributes: map[string]schema.Attribute{
			"tags": tag.ResourceSchema(),
		"id": schema.StringAttribute{
			Description: "The unique identifier of the resource.\n" +
				"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
			Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		"description": schema.StringAttribute{
			Description: "A brief explanation or note about this resource.\n" +
				"  - example: VPN test\n" +
				"  - constraints: maxLength: 40",
			Optional: true,
			Computed: true,
			Default: stringdefault.StaticString(""),
		},
		"name": schema.StringAttribute{
			Description: "The name of the resource.\n" +
				"  - example: vpnGWProd\n" +
				"  - valid: English letters and numbers only\n" +
				"  - constraints: minLength: 1, maxLength: 20",
			Required: true,
		},
		"vpn_gateway_id": schema.StringAttribute{
			Description: "The identifier of the VPN gateway that the resource belongs to.\n" +
				"  - example: 01c543eb4b8d42a9a3502345d4025147",
			Required: true,
		},
			"phase1": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"dpd_retry_interval": schema.Int32Attribute{
						Description: "The Dead Peer Detection retry interval in seconds.\n" +
							"  - example: 60\n" +
							"  - constraints: 1-3600",
						Required: true,
					},
					"ike_version": schema.Int32Attribute{
						Description: "The IKE (Internet Key Exchange) protocol version.\n" +
							"  - example: 2\n" +
							"  - valid: 1, 2",
						Required: true,
					},
					"peer_gateway_ip": schema.StringAttribute{
						Description: "The IP address of the peer VPN gateway.\n" +
							"  - example: 123.0.0.2",
						Required: true,
					},
					"phase1_diffie_hellman_groups": schema.ListAttribute{
						Description: "The list of Diffie-Hellman groups for IKE phase 1.\n" +
							"  - example: [30, 31, 32]\n" +
							"  - valid: 32, 31, 30, 29, 28, 27, 21, 20, 19, 18, 17, 16, 15, 14, 5, 2, 1\n" +
							"  - constraints: min 1, max 3",
						Required:    true,
						ElementType: types.Int32Type,
					},
					"phase1_encryptions": schema.ListAttribute{
						Description: "The list of encryption algorithms for IKE phase 1.\n" +
							"  - example: [\"des-md5\", \"chacha20poly1305\"]\n" +
							"  - valid: null-md5, null-sha1, null-sha256, null-sha384, null-sha512, des-null, des-md5, des-sha1, des-sha256, des-sha384, des-sha512, 3des-null, 3des-md5, 3des-sha1, 3des-sha256, 3des-sha384, 3des-sha512, aes128-null, aes128-md5, aes128-sha1, aes128-sha256, aes128-sha384, aes128-sha512, aes128gcm, aes192-null, aes192-md5, aes192-sha1, aes192-sha256, aes192-sha384, aes192-sha512, aes256-null, aes256-md5, aes256-sha1, aes256-sha256, aes256-sha384, aes256-sha512, aes256gcm, chacha20poly1305, aria128-null, aria128-md5, aria128-sha1, aria128-sha256, aria128-sha384, aria128-sha512, aria192-null, aria192-md5, aria192-sha1, aria192-sha256, aria192-sha384, aria192-sha512, aria256-null, aria256-md5, aria256-sha1, aria256-sha256, aria256-sha384, aria256-sha512, seed-null, seed-md5, seed-sha1, seed-sha256, seed-sha384, seed-sha512\n" +
							"  - constraints: min 1, max 10",
						Required:    true,
						ElementType: types.StringType,
					},
					"phase1_life_time": schema.Int32Attribute{
						Description: "The lifetime of the IKE phase 1 security association in seconds.\n" +
							"  - example: 86400\n" +
							"  - constraints: 120-172800",
						Required: true,
					},
					"pre_shared_key": schema.StringAttribute{
						Description: "The pre-shared key (PSK) for IKE mutual authentication between VPN gateways.\n" +
							"  - example: PreSharedKey1234567890\n" +
							"  - constraints: 8-64 characters, 32-character alphanumeric recommended",
						Required: true,
					},
				},
			},
			"phase2": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"perfect_forward_secrecy": schema.StringAttribute{
						Description: "The Perfect Forward Secrecy setting for IKE phase 2.\n" +
							"  - example: ENABLE\n" +
							"  - valid: ENABLE, DISABLE",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ENABLE", "DISABLE"),
						},
					},
					"phase2_diffie_hellman_groups": schema.ListAttribute{
						Description: "The list of Diffie-Hellman groups for IKE phase 2.\n" +
							"  - example: [30, 31, 32]\n" +
							"  - valid: 32, 31, 30, 29, 28, 27, 21, 20, 19, 18, 17, 16, 15, 14, 5, 2, 1\n" +
							"  - constraints: min 1, max 3",
						Required:    true,
						ElementType: types.Int32Type,
					},
					"phase2_encryptions": schema.ListAttribute{
						Description: "The list of encryption algorithms for IKE phase 2.\n" +
							"  - example: [\"des-md5\", \"chacha20poly1305\"]\n" +
							"  - valid: null-md5, null-sha1, null-sha256, null-sha384, null-sha512, des-null, des-md5, des-sha1, des-sha256, des-sha384, des-sha512, 3des-null, 3des-md5, 3des-sha1, 3des-sha256, 3des-sha384, 3des-sha512, aes128-null, aes128-md5, aes128-sha1, aes128-sha256, aes128-sha384, aes128-sha512, aes128gcm, aes192-null, aes192-md5, aes192-sha1, aes192-sha256, aes192-sha384, aes192-sha512, aes256-null, aes256-md5, aes256-sha1, aes256-sha256, aes256-sha384, aes256-sha512, aes256gcm, chacha20poly1305, aria128-null, aria128-md5, aria128-sha1, aria128-sha256, aria128-sha384, aria128-sha512, aria192-null, aria192-md5, aria192-sha1, aria192-sha256, aria192-sha384, aria192-sha512, aria256-null, aria256-md5, aria256-sha1, aria256-sha256, aria256-sha384, aria256-sha512, seed-null, seed-md5, seed-sha1, seed-sha256, seed-sha384, seed-sha512\n" +
							"  - constraints: min 1, max 10",
						Required:    true,
						ElementType: types.StringType,
					},
					"phase2_life_time": schema.Int32Attribute{
						Description: "The lifetime of the IKE phase 2 security association in seconds.\n" +
							"  - example: 86400\n" +
							"  - constraints: 120-172800",
						Required: true,
					},
					"remote_subnets": schema.ListAttribute{
						Description: "The list of remote subnets for IKE phase 2.\n" +
							"  - example: [\"10.1.1.0/24\", \"10.1.2.0/24\", \"10.1.3.0/24\"]\n" +
							"  - constraints: min 1, max 20",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
		common.ToSnakeCase("VpnTunnel"): schema.SingleNestedAttribute{
			Description: "VPN Tunnel",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				common.ToSnakeCase("AccountId"): schema.StringAttribute{
					Description: "The account ID associated with the resource.\n" +
						"  - example: 297615908b8e4ec69520a99a6777add3",
					Computed:    true,
				},
				common.ToSnakeCase("CreatedAt"): schema.StringAttribute{
					Description: "The timestamp when the resource was created in ISO 8601 format.\n" +
						"  - example: 2025-01-15T10:30:00Z",
					Computed:    true,
				},
				common.ToSnakeCase("CreatedBy"): schema.StringAttribute{
					Description: "The user ID that created the resource.\n" +
						"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
					Computed:    true,
				},
				common.ToSnakeCase("Description"): schema.StringAttribute{
					Description: "A brief explanation or note about this resource.\n" +
						"  - example: VPN test",
					Computed:    true,
				},
				common.ToSnakeCase("Id"): schema.StringAttribute{
					Description: "The unique identifier of the resource.\n" +
						"  - example: 6a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d",
					Computed:    true,
				},
				common.ToSnakeCase("ModifiedAt"): schema.StringAttribute{
					Description: "The timestamp when the resource was last modified in ISO 8601 format.\n" +
						"  - example: 2025-06-01T14:22:00Z",
					Computed:    true,
				},
				common.ToSnakeCase("ModifiedBy"): schema.StringAttribute{
					Description: "The user ID that modified the resource.\n" +
						"  - example: 6a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
					Computed:    true,
				},
				common.ToSnakeCase("Name"): schema.StringAttribute{
					Description: "The name of the resource.\n" +
						"  - example: vpnGWProd",
					Computed:    true,
				},
				common.ToSnakeCase("State"): schema.StringAttribute{
					Description: "The current state of the resource.\n" +
						"  - example: ACTIVE",
					Computed:    true,
				},
				common.ToSnakeCase("Status"): schema.StringAttribute{
					Description: "The current status of the vpn tunnel.\n" +
						"  - example: UP",
					Computed:    true,
				},
				common.ToSnakeCase("VpcId"): schema.StringAttribute{
					Description: "The identifier of the VPC that the resource belongs to.\n" +
						"  - example: f32265726b694b32920aa3b111f4c715",
					Computed:    true,
				},
				common.ToSnakeCase("VpcName"): schema.StringAttribute{
					Description: "The name of the VPC that the resource belongs to.\n" +
						"  - example: vpcProd01",
					Computed:    true,
				},
				common.ToSnakeCase("VpnGatewayId"): schema.StringAttribute{
					Description: "The identifier of the VPN gateway that the resource belongs to.\n" +
						"  - example: 01c543eb4b8d42a9a3502345d4025147",
					Computed:    true,
				},
				common.ToSnakeCase("VpnGatewayIpAddress"): schema.StringAttribute{
					Description: "The IP address of the VPN gateway.\n" +
						"  - example: 10.0.0.0/24",
					Computed:    true,
				},
				common.ToSnakeCase("VpnGatewayName"): schema.StringAttribute{
					Description: "The name of the VPN gateway that the resource belongs to.\n" +
						"  - example: vpnGWProd",
					Computed:    true,
				},
				common.ToSnakeCase("Phase1"): schema.SingleNestedAttribute{
					Description: "The IKE phase 1 negotiation settings of the VPN tunnel.",
					Computed:    true,
						Attributes: map[string]schema.Attribute{
						"dpd_retry_interval": schema.Int32Attribute{
							Description: "The Dead Peer Detection retry interval in seconds.\n" +
								"  - example: 60",
							Computed:    true,
						},
						"ike_version": schema.Int32Attribute{
							Description: "The IKE (Internet Key Exchange) protocol version.\n" +
								"  - example: 2",
							Computed:    true,
						},
						"life_time": schema.Int32Attribute{
							Description: "The lifetime of the IKE phase 1 security association in seconds.\n" +
								"  - example: 86400",
							Computed:    true,
						},
						"peer_gateway_ip": schema.StringAttribute{
							Description: "The IP address of the peer VPN gateway.\n" +
								"  - example: 123.0.0.2",
							Computed:    true,
						},
						"diffie_hellman_groups": schema.ListAttribute{
							Description: "The list of Diffie-Hellman groups for IKE phase 1.\n" +
								"  - example: [30, 31, 32]",
							Computed:    true,
							ElementType: types.Int32Type,
						},
						"encryptions": schema.ListAttribute{
							Description: "The list of encryption algorithms for IKE phase 1.\n" +
								"  - example: [\"null-md5\", \"aes128gcm\", \"chacha20poly1305\"]",
							Computed:    true,
							ElementType: types.StringType,
						},
						},
					},
				common.ToSnakeCase("Phase2"): schema.SingleNestedAttribute{
					Description: "The IKE phase 2 negotiation settings of the VPN tunnel.",
					Computed:    true,
					Attributes: map[string]schema.Attribute{
						"life_time": schema.Int32Attribute{
							Description: "The lifetime of the IKE phase 2 security association in seconds.\n" +
								"  - example: 86400",
							Computed:    true,
						},
						"perfect_forward_secrecy": schema.StringAttribute{
							Description: "The Perfect Forward Secrecy setting for IKE phase 2.\n" +
								"  - example: ENABLE",
							Computed:    true,
						},
						"remote_subnets": schema.ListAttribute{
							Description: "The list of remote subnets for IKE phase 2.\n" +
								"  - example: [\"10.1.1.0/24\", \"10.1.2.0/24\", \"10.1.3.0/24\"]",
							Computed:    true,
							ElementType: types.StringType,
						},
						"diffie_hellman_groups": schema.ListAttribute{
							Description: "The list of Diffie-Hellman groups for IKE phase 2.\n" +
								"  - example: [30, 31, 32]",
							Computed:    true,
							ElementType: types.Int32Type,
						},
						"encryptions": schema.ListAttribute{
							Description: "The list of encryption algorithms for IKE phase 2.\n" +
								"  - example: [\"null-md5\", \"aes128gcm\", \"chacha20poly1305\"]",
							Computed:    true,
							ElementType: types.StringType,
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
	var plan vpn.VpnTunnel1d1Resource
	diags := req.Plan.Get(ctx, &plan) // resource 블록에 작성된 configuration data 를 읽어온다.
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new vpn tunnel
	data, err := r.client.CreateVpnTunnel1d1(ctx, plan)

	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating vpn tunnel",
			"Could not create vpn tunnel, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error creating vpn tunnel",
			"Create vpn tunnel returned empty response",
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
	var state vpn.VpnTunnel1d1Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from port
	data, err := r.client.GetVpnTunnel(ctx, state.Id.ValueString())

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error Reading vpn tunnel",
			"Could not read vpn tunnel ID "+state.Id.ValueString()+": "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Error Reading vpn tunnel",
			"Get vpn tunnel returned empty response",
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
		Status:              types.StringValue(string(vt.Status)),
		VpcId:               types.StringValue(vt.VpcId),
		VpcName:             types.StringValue(vt.VpcName),
		VpnGatewayId:        types.StringValue(vt.VpnGatewayId),
		VpnGatewayIpAddress: types.StringValue(vt.VpnGatewayIpAddress),
		VpnGatewayName:      types.StringValue(vt.VpnGatewayName),
	}

	vtObjectValue, diags := types.ObjectValueFrom(ctx, vtModel.AttributeTypes(), vtModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Description = types.StringPointerValue(vt.Description.Get())
	state.Name = types.StringValue(vt.Name)
	state.VpnGatewayId = types.StringValue(vt.VpnGatewayId)
	state.VpnTunnel = vtObjectValue

	state.Phase1 = &vpn.VpnPhase1v1d1Detail{
		DpdRetryInterval:          vtModel.Phase1.DpdRetryInterval,
		IkeVersion:                vtModel.Phase1.IkeVersion,
		PeerGatewayIp:             vtModel.Phase1.PeerGatewayIp,
		Phase1DiffieHellmanGroups: vtModel.Phase1.DiffieHellmanGroups,
		Phase1Encryptions:         vtModel.Phase1.Encryptions,
		Phase1LifeTime:            vtModel.Phase1.LifeTime,
		PreSharedKey:              getPreSharedKey(state.Phase1),
	}

	state.Phase2 = &vpn.VpnPhase2v1d1Detail{
		PerfectForwardSecrecy:     vtModel.Phase2.PerfectForwardSecrecy,
		Phase2DiffieHellmanGroups: vtModel.Phase2.DiffieHellmanGroups,
		Phase2Encryptions:         vtModel.Phase2.Encryptions,
		Phase2LifeTime:            vtModel.Phase2.LifeTime,
		RemoteSubnets:             vtModel.Phase2.RemoteSubnets,
	}

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
	var plan, changedPlan, state vpn.VpnTunnel1d1Resource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	req.State.Get(ctx, &state)

	var nullString types.String
	changedPlan = plan

	if state.Phase1 != nil && plan.Phase1 != nil && state.Phase1.PeerGatewayIp.Equal(plan.Phase1.PeerGatewayIp) {
		changedPlan.Phase1.PeerGatewayIp = nullString
	}

	// Comment this condition, since convert version v1.0 -> 1.1, RemoteSubnet -> RemoteSubnets[]
	//if state.Phase2.RemoteSubnet.Equal(plan.Phase2.RemoteSubnet) {
	//	changedPlan.Phase2.RemoteSubnet = nullString
	//}

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
	var state vpn.VpnTunnel1d1Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vpn tunnel
	err := r.client.DeleteVpnTunnel1d1(ctx, state.Id.ValueString())
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

// getPreSharedKey returns the PSK from prior state. The API never echoes it back,
// so we preserve the value to avoid "inconsistent result" errors.
func getPreSharedKey(phase1 *vpn.VpnPhase1v1d1Detail) types.String {
	if phase1 != nil {
		return phase1.PreSharedKey
	}
	return types.StringValue("")
}

func mapPhase1Detail(phase1 scpvpn.VpnPhase1DetailV1Dot1) vpn.VpnPhase1v1Dot1Detail {
	// Handle empty/zero-value Phase1
	if phase1.PeerGatewayIp == "" {
		return vpn.VpnPhase1v1Dot1Detail{}
	}
	return vpn.VpnPhase1v1Dot1Detail{
		DpdRetryInterval:    types.Int32Value(phase1.DpdRetryInterval),
		IkeVersion:          types.Int32Value(phase1.IkeVersion),
		LifeTime:            types.Int32Value(phase1.LifeTime),
		PeerGatewayIp:       types.StringValue(phase1.PeerGatewayIp),
		DiffieHellmanGroups: convertToTypesInt32Slice(phase1.DiffieHellmanGroups),
		Encryptions:         convertToTypesStringSlice(phase1.Encryptions),
	}
}

func mapPhase2Detail(phase2 scpvpn.VpnPhase2DetailV1Dot1) vpn.VpnPhase2v1Dot1Detail {
	// Handle empty/zero-value Phase2
	if phase2.PerfectForwardSecrecy == "" {
		return vpn.VpnPhase2v1Dot1Detail{}
	}
	return vpn.VpnPhase2v1Dot1Detail{
		DiffieHellmanGroups:   convertToTypesInt32Slice(phase2.DiffieHellmanGroups),
		Encryptions:           convertToTypesStringSlice(phase2.Encryptions),
		LifeTime:              types.Int32Value(phase2.LifeTime),
		PerfectForwardSecrecy: types.StringValue(phase2.PerfectForwardSecrecy),
		RemoteSubnets:         convertToTypesStringSlice(phase2.RemoteSubnets),
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
	}, -1, -1, -1, -1)
}

// ImportState imports an existing resource into Terraform state using its ID.
func (r *vpnVpnTunnelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
  resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
