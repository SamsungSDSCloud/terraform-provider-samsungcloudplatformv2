package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	vpc "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/vpcv1d3"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/common"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &VpcCidrResource{}
	_ resource.ResourceWithConfigure   = &VpcCidrResource{}
	_ resource.ResourceWithImportState = &VpcCidrResource{}
)

// NewVpcCidrResource is a helper function to simplify the provider implementation.
func NewVpcCidrResource() resource.Resource {
	return &VpcCidrResource{}
}

// VpcCidrResource is the resource implementation.
type VpcCidrResource struct {
	_config *scpsdk.Configuration
	client  *vpc.Client
	clients *client.SCPClient
}

// Metadata returns the resource type name.
func (r *VpcCidrResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_cidr"
}

// Schema defines the schema for the resource.
func (r *VpcCidrResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CIDR blocks that can be used in a VPC",
		Attributes: map[string]schema.Attribute{
			// Input
			common.ToSnakeCase("VpcId"): schema.StringAttribute{
				Description: "The identifier of the VPC that the resource belongs to.\n" +
					"  - example : '023c57b14f11483689338d085e061492'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			common.ToSnakeCase("Cidr"): schema.StringAttribute{
				Description: "The IP address range of the vpc in CIDR notation.\n" +
					"  - example : '192.168.0.0/24'",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Output
			common.ToSnakeCase("Id"): schema.StringAttribute{
				Description: "The unique identifier of the vpc cidr.\n" +
					"  - example : '192.168.0.0/24'",
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VpcCidrResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.clients = inst.Client
}

// Create creates the resource and sets the initial Terraform state.
func (r *VpcCidrResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpc.VpcCidrResource

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.AddVpcCidr(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to add VPC CIDR",
			fmt.Sprintf("An error occurred while adding VPC CIDR: %s. Details: %s", err.Error(), detail),
		)
		return
	}
	if !data.HasCidr() {
		resp.Diagnostics.AddError(
			"Failed to add VPC CIDR",
			"An error occurred while adding VPC CIDR. Empty response.",
		)
		return
	}

	plan.Id = types.StringValue(data.GetCidr())

	// Set state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *VpcCidrResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpc.VpcCidrResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get VPC detail which includes all CIDRs
	vpcData, err := r.client.GetVpc(ctx, state.VpcId.ValueString())
	if err != nil {
		if strings.HasPrefix(err.Error(), "404") {
			// VPC itself was deleted
			resp.State.RemoveResource(ctx)
			return
		}
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Failed to read VPC CIDR",
			fmt.Sprintf("An error occurred while reading VPC CIDR: %s. Details: %s", err.Error(), detail),
		)
		return
	}

	// Check if the CIDR still exists in the VPC's CIDR list
	cidrFound := false
	for _, cidr := range vpcData.Vpc.Cidrs {
		if cidr.Cidr == state.Cidr.ValueString() {
			cidrFound = true
			break
		}
	}

	if !cidrFound {
		// CIDR no longer exists on the VPC, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	// CIDR still exists, refresh state
	state.Id = types.StringValue(state.Cidr.ValueString())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *VpcCidrResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not supported",
		"VPC CIDR do not support Update operations. The resource will not be updated.",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VpcCidrResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"VPC CIDR cannot be removed",
		"The VPC CIDR API does not support deletion. The CIDR block remains attached to the VPC; "+
			"only the Terraform state entry has been removed. "+
			"Remove the CIDR manually if it is no longer needed.",
	)
	// Nothing else to do — the framework removes the state automatically
}

// ImportState imports an existing VPC CIDR into Terraform state.
func (r *VpcCidrResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: vpcId,cidr, got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vpc_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidr"), parts[1])...)
}
