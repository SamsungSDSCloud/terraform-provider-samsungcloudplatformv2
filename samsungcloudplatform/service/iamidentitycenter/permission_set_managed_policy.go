package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	iamidentitycenterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &iamIdentityCenterPermissionSetManagedPolicyResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterPermissionSetManagedPolicyResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterPermissionSetManagedPolicyResource{}
)

func NewIamIdentityCenterPermissionSetManagedPolicyResource() resource.Resource {
	return &iamIdentityCenterPermissionSetManagedPolicyResource{}
}

type iamIdentityCenterPermissionSetManagedPolicyResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_set_managed_policy"
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a single AWS managed policy attached to an IAM Identity Center Permission Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Composed resource ID (<permission_set_id>,<instance_id>,<managed_policy_id>)",
				MarkdownDescription: "Composed resource ID (`<permission_set_id>,<instance_id>,<managed_policy_id>`)",
			},
			"permission_set_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
				MarkdownDescription: "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"managed_policy_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "IAM Managed Policy ID (the AWS managed policy identifier)\n  - example: 37f2e31ff86b415698d7e8eeafab445d",
				MarkdownDescription: "IAM Managed Policy ID (the AWS managed policy identifier)\n  - example: 37f2e31ff86b415698d7e8eeafab445d",
			},
			"managed_policy_name": schema.StringAttribute{
				Required:            true,
				Description:         "IAM Managed Policy Display Name\n  - example: AmazonS3ReadOnlyAccess",
				MarkdownDescription: "IAM Managed Policy Display Name\n  - example: AmazonS3ReadOnlyAccess",
			},
		},
	}
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clients = inst.Client
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type permissionSetManagedPolicyResourceModel struct {
	PermissionSetId   types.String `tfsdk:"permission_set_id"`
	InstanceId        types.String `tfsdk:"instance_id"`
	ManagedPolicyId   types.String `tfsdk:"managed_policy_id"`
	ManagedPolicyName types.String `tfsdk:"managed_policy_name"`
	Id                types.String `tfsdk:"id"`
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionSetManagedPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := plan.PermissionSetId.ValueString()
	instanceId := plan.InstanceId.ValueString()
	managedPolicyId := plan.ManagedPolicyId.ValueString()

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryManaged)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Idempotent: the backend PUT is additive (union) and never replaces.
	if findManagedPolicy(policies, managedPolicyId) == nil {
		// The managed PUT object requires BOTH the IAM managed policy id and its
		// display name; the backend raises an error if "name" is missing.
		_, err := r.clients.IamIdentityCenter.SetPolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesSetRequest{
			InstanceId: instanceId,
			ManagedPolicies: []map[string]interface{}{
				{"id": managedPolicyId, "name": plan.ManagedPolicyName.ValueString()},
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create IAM Identity Center Permission Set Managed Policy",
				err.Error(),
			)
			return
		}
	}

	plan.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, managedPolicyId))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionSetManagedPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set Managed Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryManaged)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := findManagedPolicy(policies, identifier)
	if p == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.PermissionSetId = types.StringValue(permissionSetId)
	state.InstanceId = types.StringValue(instanceId)
	state.ManagedPolicyId = types.StringValue(identifier)
	// managed_policy_name is supplied by the user; on import (when it is not in
	// state) fall back to the display name the backend stores in contents.
	if state.ManagedPolicyName.IsNull() || state.ManagedPolicyName.IsUnknown() {
		if name, ok := managedPolicyDisplayName(*p); ok && name != "" {
			state.ManagedPolicyName = types.StringValue(name)
		}
	}
	state.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, identifier))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state permissionSetManagedPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan permissionSetManagedPolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM Identity Center Permission Set Managed Policy",
			err.Error(),
		)
		return
	}

	// Only the display name is mutable. Re-PUT the single managed policy with
	// the new name (the backend PUT is additive, so other policies are untouched).
	_, err = r.clients.IamIdentityCenter.SetPolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesSetRequest{
		InstanceId: instanceId,
		ManagedPolicies: []map[string]interface{}{
			{"id": identifier, "name": plan.ManagedPolicyName.ValueString()},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM Identity Center Permission Set Managed Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryManaged)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := findManagedPolicy(policies, identifier)
	if p == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ManagedPolicyName = plan.ManagedPolicyName
	state.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, identifier))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetManagedPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionSetManagedPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Permission Set Managed Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryManaged)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete uses the IDC-internal policy.Id, never the IAM managed policy id
	// (which lives in contents.id).
	if p := findManagedPolicy(policies, identifier); p != nil {
		if err := r.clients.IamIdentityCenter.DeletePolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesDeleteRequest{
			InstanceId: instanceId,
			PolicyIds:  []string{p.Id},
		}); err != nil {
			resp.Diagnostics.AddError(
				"Failed to delete IAM Identity Center Permission Set Managed Policy",
				err.Error(),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
