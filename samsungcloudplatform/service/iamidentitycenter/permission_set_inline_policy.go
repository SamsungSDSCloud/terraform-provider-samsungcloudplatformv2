package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	iamidentitycenterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &iamIdentityCenterPermissionSetInlinePolicyResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterPermissionSetInlinePolicyResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterPermissionSetInlinePolicyResource{}
)

func NewIamIdentityCenterPermissionSetInlinePolicyResource() resource.Resource {
	return &iamIdentityCenterPermissionSetInlinePolicyResource{}
}

type iamIdentityCenterPermissionSetInlinePolicyResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_set_inline_policy"
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a single inline policy attached to an IAM Identity Center Permission Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Composed resource ID (<permission_set_id>,<instance_id>,<policy content hash>)",
				MarkdownDescription: "Composed resource ID (`<permission_set_id>,<instance_id>,<policy content hash>`)",
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
			"policy_document": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Inline Policy JSON Document\n  - example: {\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"s3:*\",\"Resource\":\"*\"}]}",
				MarkdownDescription: "Inline Policy JSON Document\n  - example: {\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"s3:*\",\"Resource\":\"*\"}]}",
			},
		},
	}
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type permissionSetInlinePolicyResourceModel struct {
	PermissionSetId types.String `tfsdk:"permission_set_id"`
	InstanceId      types.String `tfsdk:"instance_id"`
	PolicyDocument  types.String `tfsdk:"policy_document"`
	Id              types.String `tfsdk:"id"`
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionSetInlinePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := plan.PermissionSetId.ValueString()
	instanceId := plan.InstanceId.ValueString()

	normalizedDoc, err := normalizePolicyDocument(plan.PolicyDocument.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create IAM Identity Center Permission Set Inline Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryInline)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Idempotent: the backend PUT is additive (union) and never replaces.
	if findInlinePolicy(policies, normalizedDoc) == nil {
		_, err := r.clients.IamIdentityCenter.SetPolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesSetRequest{
			InstanceId:     instanceId,
			InlinePolicies: buildInlinePoliciesMap([]string{plan.PolicyDocument.ValueString()}),
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create IAM Identity Center Permission Set Inline Policy",
				err.Error(),
			)
			return
		}
	}

	plan.PolicyDocument = types.StringValue(normalizedDoc)
	plan.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, inlinePolicyIdentifier(normalizedDoc)))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionSetInlinePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set Inline Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryInline)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Match the policy whose normalized content hashes to the identifier chunk
	// of the composed id (supports normal reads and passthrough import).
	found := false
	for i := range policies {
		if string(policies[i].Category) != permissionSetPolicyCategoryInline {
			continue
		}
		raw, ok := existingInlineContent(policies[i])
		if !ok {
			continue
		}
		norm, err := normalizePolicyDocument(raw)
		if err != nil {
			continue
		}
		if inlinePolicyIdentifier(norm) == identifier {
			state.PolicyDocument = types.StringValue(norm)
			found = true
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	state.PermissionSetId = types.StringValue(permissionSetId)
	state.InstanceId = types.StringValue(instanceId)
	state.Id = types.StringValue(composePermissionSetPolicyID(permissionSetId, instanceId, identifier))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes are RequiresReplace, so no in-place update path exists.
	// Persist the plan unchanged.
	var plan permissionSetInlinePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetInlinePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionSetInlinePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId, instanceId, identifier, err := parsePermissionSetPolicyID(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Permission Set Inline Policy",
			err.Error(),
		)
		return
	}

	policies, diags := fetchPermissionSetPolicies(ctx, r.clients, permissionSetId, instanceId, permissionSetPolicyCategoryInline)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete uses the IDC-internal policy.Id; locate it via content-hash match.
	var target *sdk.PolicyV1Dot2
	for i := range policies {
		if string(policies[i].Category) != permissionSetPolicyCategoryInline {
			continue
		}
		raw, ok := existingInlineContent(policies[i])
		if !ok {
			continue
		}
		norm, err := normalizePolicyDocument(raw)
		if err != nil {
			continue
		}
		if inlinePolicyIdentifier(norm) == identifier {
			target = &policies[i]
			break
		}
	}

	if target != nil {
		if err := r.clients.IamIdentityCenter.DeletePolicies(ctx, permissionSetId, iamidentitycenterClient.PoliciesDeleteRequest{
			InstanceId: instanceId,
			PolicyIds:  []string{target.Id},
		}); err != nil {
			resp.Diagnostics.AddError(
				"Failed to delete IAM Identity Center Permission Set Inline Policy",
				err.Error(),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}
