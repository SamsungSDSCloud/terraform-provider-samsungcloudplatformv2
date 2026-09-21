package iamidentitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	iamidentitycenterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &iamIdentityCenterAccountAssignmentResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterAccountAssignmentResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterAccountAssignmentResource{}
)

// NewIamIdentityCenterAccountAssignmentResource instantiates the IAM Identity Center account assignment resource.
func NewIamIdentityCenterAccountAssignmentResource() resource.Resource {
	return &iamIdentityCenterAccountAssignmentResource{}
}

type iamIdentityCenterAccountAssignmentResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterAccountAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_account_assignment"
}

func (r *iamIdentityCenterAccountAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a single IAM Identity Center Account Assignment, " +
			"binding one principal to one permission set in one target account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account Assignment ID",
				MarkdownDescription: "Account Assignment ID",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:         "IAM Identity Center Instance ID\n  - example: 6tjgl36qyf5d",
				MarkdownDescription: "IAM Identity Center Instance ID\n  - example: 6tjgl36qyf5d",
			},
			"target_account_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:         "Target Account ID\n  - example: 1458ae4b3f2d45dfa54c5645505d56af",
				MarkdownDescription: "Target Account ID\n  - example: 1458ae4b3f2d45dfa54c5645505d56af",
			},
			"principal_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:         "Principal ID (user or group)\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Principal ID (user or group)\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"principal_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("USER", "GROUP"),
				},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:         "Principal type: USER or GROUP",
				MarkdownDescription: "Principal type: USER or GROUP",
			},
			"permission_set_id": schema.StringAttribute{
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:         "Permission Set ID\n  - example: a9e1d249b69a4abc9c15ba53ff8468a8",
				MarkdownDescription: "Permission Set ID\n  - example: a9e1d249b69a4abc9c15ba53ff8468a8",
			},
			"account_assignment_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Backend Account Assignment ID",
				MarkdownDescription: "Backend Account Assignment ID",
			},
			"principal_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Principal name",
				MarkdownDescription: "Principal name",
			},
			"target_account_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Target account name",
				MarkdownDescription: "Target account name",
			},
			"target_account_email": schema.StringAttribute{
				Computed:            true,
				Description:         "Target account email",
				MarkdownDescription: "Target account email",
			},
			"role_srn": schema.StringAttribute{
				Computed:            true,
				Description:         "Role SRN",
				MarkdownDescription: "Role SRN",
			},
			"permission_set_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Permission set name",
				MarkdownDescription: "Permission set name",
			},
			"permission_set_srn": schema.StringAttribute{
				Computed:            true,
				Description:         "Permission set SRN",
				MarkdownDescription: "Permission set SRN",
			},
		},
	}
}

func (r *iamIdentityCenterAccountAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.clients = inst.Client
}

func (r *iamIdentityCenterAccountAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "|")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: <instance_id>|<target_account_id>|<account_assignment_id>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("target_account_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_assignment_id"), parts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
}

func (r *iamIdentityCenterAccountAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accountAssignmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	principalType := sdk.PRINCIPALTYPEENUM_USER
	if plan.PrincipalType.ValueString() == "GROUP" {
		principalType = sdk.PRINCIPALTYPEENUM_GROUP
	}

	createReq := iamidentitycenterClient.AccountAssignmentCreateRequest{
		InstanceId: plan.InstanceId.ValueString(),
		Accounts: []sdk.AccountV1Dot6{
			{TargetAccountId: plan.TargetAccountId.ValueString()},
		},
		PermissionSets: []sdk.PermissionSetInfoV1Dot6{
			{PermissionSetId: plan.PermissionSetId.ValueString()},
		},
		Principals: []sdk.Principal{
			{PrincipalId: plan.PrincipalId.ValueString(), PrincipalType: principalType},
		},
	}

	result, err := r.clients.IamIdentityCenter.CreateAccountAssignment(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create IAM Identity Center Account Assignment",
			err.Error(),
		)
		return
	}

	if result == nil || result.AccountAssignments == nil || len(result.AccountAssignments) == 0 {
		resp.Diagnostics.AddError(
			"Failed to create IAM Identity Center Account Assignment",
			"Create succeeded but no assignment was returned",
		)
		return
	}

	created := result.AccountAssignments[0]
	plan.Id = types.StringValue(created.Id)
	plan.AccountAssignmentId = types.StringValue(created.Id)

	// Create 응답에는 표시용 필드가 없어 Read로 보완
	if !r.refreshData(ctx, &plan) {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *iamIdentityCenterAccountAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accountAssignmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.refreshData(ctx, &state) {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *iamIdentityCenterAccountAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accountAssignmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.refreshData(ctx, &plan) {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *iamIdentityCenterAccountAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accountAssignmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteReq := iamidentitycenterClient.AccountAssignmentDeleteRequest{
		InstanceId:      state.InstanceId.ValueString(),
		TargetAccountId: state.TargetAccountId.ValueString(),
	}

	err := r.clients.IamIdentityCenter.DeleteAccountAssignment(ctx, state.AccountAssignmentId.ValueString(), deleteReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Account Assignment",
			err.Error(),
		)
		return
	}
}

// refreshData lists the target account's assignments and fills state from the
// matching account_assignment_id. It returns false when the assignment no longer exists.
func (r *iamIdentityCenterAccountAssignmentResource) refreshData(ctx context.Context, state *accountAssignmentModel) bool {
	result, err := r.clients.IamIdentityCenter.ListAllAccountAssignments(ctx, state.InstanceId.ValueString(), state.TargetAccountId.ValueString())
	if err != nil {
		// 드리프트 처리 대상: 조회 실패 시 리소스 제거로 처리
		return false
	}

	for _, a := range result.AccountAssignments {
		for _, ps := range a.PermissionSets {
			if ps.AccountAssignmentId != state.AccountAssignmentId.ValueString() {
				continue
			}

			state.Id = types.StringValue(ps.AccountAssignmentId)
			state.AccountAssignmentId = types.StringValue(ps.AccountAssignmentId)
			state.PrincipalName = types.StringValue(a.PrincipalName)
			state.TargetAccountName = types.StringValue(a.TargetAccountName)
			state.TargetAccountEmail = types.StringValue(a.TargetAccountEmail)
			state.RoleSrn = types.StringPointerValue(a.RoleSrn.Get())
			state.PermissionSetName = types.StringValue(ps.PermissionSetName)
			state.PermissionSetSrn = types.StringValue(ps.PermissionSetSrn)
			return true
		}
	}

	return false
}

type accountAssignmentModel struct {
	Id                  types.String `tfsdk:"id"`
	InstanceId          types.String `tfsdk:"instance_id"`
	TargetAccountId     types.String `tfsdk:"target_account_id"`
	PrincipalId         types.String `tfsdk:"principal_id"`
	PrincipalType       types.String `tfsdk:"principal_type"`
	PermissionSetId     types.String `tfsdk:"permission_set_id"`
	AccountAssignmentId types.String `tfsdk:"account_assignment_id"`
	PrincipalName       types.String `tfsdk:"principal_name"`
	TargetAccountName   types.String `tfsdk:"target_account_name"`
	TargetAccountEmail  types.String `tfsdk:"target_account_email"`
	RoleSrn             types.String `tfsdk:"role_srn"`
	PermissionSetName   types.String `tfsdk:"permission_set_name"`
	PermissionSetSrn    types.String `tfsdk:"permission_set_srn"`
}
