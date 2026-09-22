package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &iamIdentityCenterGroupMemberResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterGroupMemberResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterGroupMemberResource{}
)

func NewIamIdentityCenterGroupMemberResource() resource.Resource {
	return &iamIdentityCenterGroupMemberResource{}
}

type iamIdentityCenterGroupMemberResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterGroupMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_group_member"
}

func (r *iamIdentityCenterGroupMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// This resource follows the canonical one-membership-per-resource model:
	// each resource manages exactly one (group, member) relationship, so there
	// is never an empty-membership state and removal is a plain destroy.
	resp.Schema = schema.Schema{
		Description: "Manages a single IAM Identity Center group membership (one user per resource).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Group Member ID (`{group_id}:{member_id}`)\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2:2fb5e7f8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Group Member ID (`{group_id}:{member_id}`)\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2:2fb5e7f8c29a449dbfa8681f8f1d78e2",
			},
			"group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"instance_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"member_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "ID of the user to add as a group member (single user per resource)\n  - example: 2fb5e7f8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "ID of the user to add as a group member (single user per resource)\n  - example: 2fb5e7f8c29a449dbfa8681f8f1d78e2",
			},
		},
	}
}

func (r *iamIdentityCenterGroupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamIdentityCenterGroupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type groupMemberResourceModel struct {
	GroupId    types.String `tfsdk:"group_id"`
	InstanceId types.String `tfsdk:"instance_id"`
	MemberId   types.String `tfsdk:"member_id"`
	Id         types.String `tfsdk:"id"`
}

func (m *groupMemberResourceModel) membershipID() string {
	return m.GroupId.ValueString() + ":" + m.MemberId.ValueString()
}

func (r *iamIdentityCenterGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupMemberResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := sdk.GroupUsersRequest{
		InstanceId: plan.InstanceId.ValueString(),
		UserUuids:  []string{plan.MemberId.ValueString()},
	}
	if _, err := r.clients.IamIdentityCenter.AddUsersToGroup(ctx, plan.GroupId.ValueString(), addReq); err != nil {
		resp.Diagnostics.AddError(
			"Failed to add user to IAM Identity Center Group",
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(plan.membershipID())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterGroupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupMemberResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Page through the group's members (page index is 0-based) and check whether
	// this membership's user is present. The specific (group, member) relationship
	// is considered gone — the only case that removes the resource from state —
	// when the user is not among the group's members.
	const pageSize = 100
	memberId := state.MemberId.ValueString()
	found := false
	page := 0
	for {
		result, err := r.clients.IamIdentityCenter.ListGroupUsers(
			ctx,
			state.GroupId.ValueString(),
			state.InstanceId.ValueString(),
			"",
			pageSize,
			int32(page),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read IAM Identity Center Group Member",
				err.Error(),
			)
			return
		}
		if result == nil {
			break
		}
		for _, user := range result.Users {
			if user.Id == memberId {
				found = true
				break
			}
		}
		if found || len(result.Users) < pageSize {
			break
		}
		page++
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(state.membershipID())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterGroupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupMemberResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// group_id, instance_id and member_id are all RequiresReplace, so a change
	// here is expressed as destroy+create rather than an in-place update. Just
	// restore the planned state defensively.
	plan.Id = types.StringValue(plan.membershipID())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterGroupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupMemberResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removeReq := sdk.GroupUsersRequest{
		InstanceId: state.InstanceId.ValueString(),
		UserUuids:  []string{state.MemberId.ValueString()},
	}
	if err := r.clients.IamIdentityCenter.RemoveUsersFromGroup(ctx, state.GroupId.ValueString(), removeReq); err != nil {
		resp.Diagnostics.AddError(
			"Failed to remove user from IAM Identity Center Group",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}
