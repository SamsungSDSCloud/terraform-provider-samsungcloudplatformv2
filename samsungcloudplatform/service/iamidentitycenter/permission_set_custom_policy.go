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
	_ resource.Resource                = &iamIdentityCenterPermissionSetCustomPolicyResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterPermissionSetCustomPolicyResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterPermissionSetCustomPolicyResource{}
)

func NewIamIdentityCenterPermissionSetCustomPolicyResource() resource.Resource {
	return &iamIdentityCenterPermissionSetCustomPolicyResource{}
}

type iamIdentityCenterPermissionSetCustomPolicyResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_set_custom_policy"
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a single custom (customer managed) policy attached to an IAM Identity Center Permission Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Composed resource ID (<permission_set_id>,<instance_id>,<policy name>)",
				MarkdownDescription: "Composed resource ID (`<permission_set_id>,<instance_id>,<policy name>`)",
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
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Custom Policy Name\n  - example: MyCustomPolicy",
				MarkdownDescription: "Custom Policy Name\n  - example: MyCustomPolicy",
			},
		},
	}
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type permissionSetCustomPolicyResourceModel struct {
	PermissionSetId types.String `tfsdk:"permission_set_id"`
	InstanceId      types.String `tfsdk:"instance_id"`
	Name            types.String `tfsdk:"name"`
	Id              types.String `tfsdk:"id"`
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionSetCustomPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := plan.PermissionSetId.ValueString()
	instanceId := plan.InstanceId.ValueString()
	name := plan.Name.ValueString()

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryCustom)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Idempotent: the backend PUT is additive (union) and never replaces, so if
	// the policy already exists we no-op instead of erroring.
	if findCustomPolicy(policies, name) == nil {
		_, err := r.clients.IamIdentityCenter.SetPolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesSetRequest{
			InstanceId:     instanceId,
			CustomPolicies: []string{name},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create IAM Identity Center Permission Set Custom Policy",
				err.Error(),
			)
			return
		}
	}

	plan.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, name))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionSetCustomPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set Custom Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryCustom)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := findCustomPolicy(policies, identifier)
	if p == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.PermissionSetId = types.StringValue(permissionSetId)
	state.InstanceId = types.StringValue(instanceId)
	state.Name = types.StringValue(p.GetName())
	state.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, p.GetName()))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes are RequiresReplace, so no in-place update path exists.
	// Persist the plan unchanged.
	var plan permissionSetCustomPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetCustomPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionSetCustomPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Permission Set Custom Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryCustom)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete uses the IDC-internal policy.Id, not the policy name.
	if p := findCustomPolicy(policies, identifier); p != nil {
		if err := r.clients.IamIdentityCenter.DeletePolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesDeleteRequest{
			InstanceId: instanceId,
			PolicyIds:  []string{p.Id},
		}); err != nil {
			resp.Diagnostics.AddError(
				"Failed to delete IAM Identity Center Permission Set Custom Policy",
				err.Error(),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
