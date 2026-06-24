package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/vpcv1d2"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &VPCSubnetVipNatIpResource{}
	_ resource.ResourceWithConfigure   = &VPCSubnetVipNatIpResource{}
	_ resource.ResourceWithImportState = &VPCSubnetVipNatIpResource{}
)

// NewVPCSubnetVipNatIpResource is a helper function to simplify the provider implementation.
func NewVPCSubnetVipNatIpResource() resource.Resource {
	return &VPCSubnetVipNatIpResource{}
}

// VPCSubnetVipNatIpResource is the resource implementation.
type VPCSubnetVipNatIpResource struct {
	_config *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VPCSubnetVipNatIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet_vip_nat_ip"
}

// Schema defines the schema for the resource.
func (r *VPCSubnetVipNatIpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource of Subnet VIP's NAT IP",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the subnet vip nat ip belongs to.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("VipId"): schema.StringAttribute{
				Description: "The unique identifier of the subnet vip.\n" +
					"  - example : 0466a9448d9a4411a86055939e451c8f",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("PublicipId"): schema.StringAttribute{
				Description: "The identifier of the public IP address.\n" +
					"  - example : 12f56e27070248a6a240a497e43fbe18",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("NatType"): schema.StringAttribute{
				Description: "The type of the NAT.\n" +
					"  - example : PUBLIC",
				Required: true,
			},

			// Output
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the static nat ip.\n" +
					"  - example : 0009e49548154745948e9722adefbf40",
				Computed: true,
			},
			common.ToSnakeCase("State"): schema.StringAttribute{
				Description: "The current lifecycle state of the static nat.\n" +
					"  - example : ACTIVE",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VPCSubnetVipNatIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *VPCSubnetVipNatIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d2.SubnetVipNatIpResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := r.client.CreateSubnetVipNatIp(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to create VPC Subnet VIP NAT IP",
			fmt.Sprintf("An error occurred while creating VPC Subnet VIP NAT IP: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map API response to object
	plan.Id = types.StringValue(apiResponse.Id)
	plan.State = types.StringValue(apiResponse.State)

	waitForState := []string{"ACTIVE", "ERROR"}
	natIPState, err := waitForVpcSubnetNatIpStatus(ctx, r.client, plan.SubnetId.ValueString(), plan.VipId.ValueString(), plan.PublicipId.ValueString(), []string{}, waitForState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VPC Subnet VIP NAT IP",
			"Error waiting for VPC Subnet VIP NAT IP to become active: "+err.Error(),
		)
		return
	}

	plan.State = types.StringValue(natIPState)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VPCSubnetVipNatIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcv1d2.SubnetVipNatIpResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ShowSubnetVip(ctx, state.SubnetId.ValueString(), state.VipId.ValueString())
	if err != nil {
		// Subnet VIP was deleted => remove Subnet VIP NAT IP too
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error reading Subnet VIP",
			"Could not read Subnet VIP Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
	if data == nil {
		resp.Diagnostics.AddError(
			"Error reading data",
			"An error occurred while reading data. Empty response",
		)
		return
	}

	// Map API response to object
	if data.SubnetVip.StaticNat.IsSet() {
		staticNat := data.SubnetVip.StaticNat.Get()
		if staticNat != nil && staticNat.PublicipId == state.PublicipId.ValueString() {
			state.Id = types.StringValue(staticNat.Id)
			state.State = types.StringValue(staticNat.State)
			state.PublicipId = types.StringValue(staticNat.PublicipId)
			// state.NatType is not exist in api response
		} else {
			// Subnet VIP NAT IP was changed without us knowing so we are not managed this VIP NAT IP resource anymore
			resp.State.RemoveResource(ctx)
			return
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VPCSubnetVipNatIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcv1d2.SubnetVipNatIpResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSubnetVipNatIp(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting VPC subnet VIP NAT IP",
			"Could not delete VPC subnet VIP NAT IP with Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	_, err = waitForVpcSubnetNatIpStatus(ctx, r.client, state.SubnetId.ValueString(), state.VipId.ValueString(), state.PublicipId.ValueString(), []string{}, []string{"DELETED"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VPC Subnet VIP NAT IP",
			"Error waiting for VPC Subnet VIP NAT IP to be deleted: "+err.Error(),
		)
		return
	}
}

func (r *VPCSubnetVipNatIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"VPC Subnet VIP NAT IP resource do not support update operations. The resource will not be updated.",
	)
}

func waitForVpcSubnetNatIpStatus(ctx context.Context, vpcClient *vpcv1d2.Client, subnetId string, vipId string, publicIpId string, pendingStates []string, targetStates []string) (string, error) {
	var state string
	return state, client.WaitForStatus(ctx, nil, pendingStates, targetStates, func() (interface{}, string, error) {
		info, err := vpcClient.ShowSubnetVip(ctx, subnetId, vipId)
		if err != nil {
			return nil, "", err
		}
		if info.SubnetVip.StaticNat.IsSet() {
			staticNat := info.SubnetVip.StaticNat.Get()
			if staticNat != nil && staticNat.PublicipId == publicIpId {
				state = string(info.SubnetVip.StaticNat.Get().State)
				return info, state, nil
			}
		}
		return info, "DELETED", nil
	}, -1, -1, -1, -1)
}

// ImportState imports an existing resource into Terraform state.
func (r *VPCSubnetVipNatIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected ID format: subnetId/vipId/publicipId
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: subnetId/vipId/publicipId, got: %q", req.ID),
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("subnet_id"), types.StringValue(parts[0]))
	resp.State.SetAttribute(ctx, path.Root("vip_id"), types.StringValue(parts[1]))
	resp.State.SetAttribute(ctx, path.Root("publicip_id"), types.StringValue(parts[2]))
}
