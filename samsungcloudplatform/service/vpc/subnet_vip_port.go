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
	_ resource.Resource                = &VPCSubnetVipPortResource{}
	_ resource.ResourceWithConfigure   = &VPCSubnetVipPortResource{}
	_ resource.ResourceWithImportState = &VPCSubnetVipPortResource{}
)

// NewVPCSubnetVipPortResource is a helper function to simplify the provider implementation.
func NewVPCSubnetVipPortResource() resource.Resource {
	return &VPCSubnetVipPortResource{}
}

// VPCSubnetVipPortResource is the resource implementation.
type VPCSubnetVipPortResource struct {
	_config *scpsdk.Configuration
	client  *vpcv1d2.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VPCSubnetVipPortResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_subnet_vip_port"
}

// Schema defines the schema for the resource.
func (r *VPCSubnetVipPortResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource of VPC Subnet VIP's Port",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("SubnetId"): schema.StringAttribute{
				Description: "The identifier of the subnet that the subnet vip port belongs to.\n" +
					"  - example : 023c57b14f11483689338d085e061492",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("VipId"): schema.StringAttribute{
				Description: "The identifier of the subnet vip.\n" +
					"  - example : 0466a9448d9a4411a86055939e451c8f",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("PortId"): schema.StringAttribute{
				Description: "The identifier of the port.\n" +
					"  - example : 35268a9f2eda4cde83b1d85c1f61f67d",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Output
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the connected port.\n" +
					"  - example : 35268a9f2eda4cde83b1d85c1f61f67d",
				Computed: true,
			},
			common.ToSnakeCase("SubnetVipId"): schema.StringAttribute{
				Description: "The unique identifier of the subnet vip.\n" +
					"  - example : 0466a9448d9a4411a86055939e451c8f",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VPCSubnetVipPortResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VPCSubnetVipPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcv1d2.SubnetVipPortResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResponse, err := r.client.CreateSubnetVipPort(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to create VPC Subnet VIP Port",
			fmt.Sprintf("An error occurred while creating VPC Subnet VIP Port: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Map API response to object
	plan.Id = types.StringValue(apiResponse.Id)
	plan.SubnetVipId = types.StringValue(apiResponse.SubnetVipId)

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VPCSubnetVipPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcv1d2.SubnetVipPortResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ShowSubnetVip(ctx, state.SubnetId.ValueString(), state.VipId.ValueString())
	if err != nil {
		// Subnet VIP was deleted => remove Subnet VIP Port too
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

	state.SubnetVipId = types.StringValue(data.SubnetVip.Id)
	// Check Connected Port to refresh resource
	if len(data.SubnetVip.ConnectedPorts) > 0 {
		isPortExist := false
		for _, port := range data.SubnetVip.ConnectedPorts {
			if port.PortId == state.PortId.ValueString() {
				state.Id = types.StringValue(port.Id)
				isPortExist = true
				break
			}
		}
		if !isPortExist {
			// Managed port Id not exist in Subnet Vip
			resp.State.RemoveResource(ctx)
			return
		}
	} else {
		// No connected port meaning the port managed was deleted too
		resp.State.RemoveResource(ctx)
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VPCSubnetVipPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcv1d2.SubnetVipPortResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSubnetVipPort(ctx, state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting VPC subnet VIP Port",
			"Could not delete VPC subnet VIP Port with Id "+state.SubnetId.ValueString()+" unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *VPCSubnetVipPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"VPC Subnet VIP Port resource do not support update operations. The resource will not be updated.",
	)
}

// ImportState imports an existing resource into Terraform state.
func (r *VPCSubnetVipPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: subnetId/vipId/portId, got: %q", req.ID),
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("subnet_id"), types.StringValue(parts[0]))
	resp.State.SetAttribute(ctx, path.Root("vip_id"), types.StringValue(parts[1]))
	resp.State.SetAttribute(ctx, path.Root("port_id"), types.StringValue(parts[2]))
}
