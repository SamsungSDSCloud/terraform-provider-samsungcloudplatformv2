package iamidentitycenter

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	iamidentitycenterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
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
	_ resource.Resource                = &iamIdentityCenterUserResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterUserResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterUserResource{}
)

func NewIamIdentityCenterUserResource() resource.Resource {
	return &iamIdentityCenterUserResource{}
}

type iamIdentityCenterUserResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_user"
}

func (r *iamIdentityCenterUserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Identity Center User.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"user_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
				Description:         "User Login ID\n  - example: johndoe\n  - maxLength: 128",
				MarkdownDescription: "User Login ID\n  - example: johndoe\n  - maxLength: 128",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
				Description:         "Real Name\n  - example: John Doe\n  - maxLength: 256",
				MarkdownDescription: "Real Name\n  - example: John Doe\n  - maxLength: 256",
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Description:         "Password",
				MarkdownDescription: "Password",
			},
			"email": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User Email\n  - example: john.doe@example.com",
				MarkdownDescription: "User Email\n  - example: john.doe@example.com",
			},
			"phone_number": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Phone Number\n  - example: 010-1234-5678",
				MarkdownDescription: "Phone Number\n  - example: 010-1234-5678",
			},
			"temporary_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Temporary Password Status",
				MarkdownDescription: "Temporary Password Status",
			},
			"business_unit": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Business Unit\n  - example: Business Unit A",
				MarkdownDescription: "Business Unit\n  - example: Business Unit A",
			},
			"department": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Department\n  - example: Department X",
				MarkdownDescription: "Department\n  - example: Department X",
			},
			"manager": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Manager\n  - example: Alice Smith",
				MarkdownDescription: "Manager\n  - example: Alice Smith",
			},
			"employee_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Employee ID\n  - example: emp-12345",
				MarkdownDescription: "Employee ID\n  - example: emp-12345",
			},
			"nation_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Nationality ID\n  - example: +82",
				MarkdownDescription: "Nationality ID\n  - example: +82",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User Description\n  - example: Description of John Doe",
				MarkdownDescription: "User Description\n  - example: Description of John Doe",
			},
			"user_uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "User ID",
				MarkdownDescription: "User ID",
			},
		},
	}
}

func (r *iamIdentityCenterUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.Instance, got: %T", req.ProviderData),
		)
		return
	}

	r.clients = inst.Client
}

func (r *iamIdentityCenterUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.clients.IamIdentityCenter

	createReq := sdk.UserCreateRequest{
		InstanceId: plan.InstanceId.ValueString(),
		UserId:     plan.UserId.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	if !plan.Email.IsNull() && !plan.Email.IsUnknown() {
		createReq.Email = *sdk.NewNullableString(plan.Email.ValueStringPointer())
	}
	if !plan.PhoneNumber.IsNull() && !plan.PhoneNumber.IsUnknown() {
		createReq.PhoneNumber = *sdk.NewNullableString(plan.PhoneNumber.ValueStringPointer())
	}
	if !plan.BusinessUnit.IsNull() && !plan.BusinessUnit.IsUnknown() {
		createReq.BusinessUnit = *sdk.NewNullableString(plan.BusinessUnit.ValueStringPointer())
	}
	if !plan.Department.IsNull() && !plan.Department.IsUnknown() {
		createReq.Department = *sdk.NewNullableString(plan.Department.ValueStringPointer())
	}
	if !plan.Manager.IsNull() && !plan.Manager.IsUnknown() {
		createReq.Manager = *sdk.NewNullableString(plan.Manager.ValueStringPointer())
	}
	if !plan.EmployeeId.IsNull() && !plan.EmployeeId.IsUnknown() {
		createReq.EmployeeId = *sdk.NewNullableString(plan.EmployeeId.ValueStringPointer())
	}
	if !plan.NationId.IsNull() && !plan.NationId.IsUnknown() {
		createReq.NationId = *sdk.NewNullableString(plan.NationId.ValueStringPointer())
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createReq.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}
	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		encodedPassword := base64.StdEncoding.EncodeToString([]byte(plan.Password.ValueString()))
		createReq.Password = *sdk.NewNullableString(&encodedPassword)
	}

	result, err := client.CreateUser(ctx, createReq)
	if err != nil {
		var openAPIErr *scpsdk.GenericOpenAPIError
		if errors.As(err, &openAPIErr) {
			resp.Diagnostics.AddError(
				"Creating IAM Identity Center User",
				fmt.Sprintf("400 Bad Request: %s", string(openAPIErr.ResponseBody)),
			)
		} else {
			resp.Diagnostics.AddError(
				"Creating IAM Identity Center User",
				err.Error(),
			)
		}
		return
	}

	plan.Id = types.StringValue(result.User.Id)
	plan.UserUuid = types.StringValue(result.User.Id)

	if plan.Password.IsUnknown() {
		plan.Password = types.StringNull()
	}

	userResult, _, err := client.GetUser(ctx, plan.InstanceId.ValueString(), result.User.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to re-fetch IAM Identity Center User",
			err.Error(),
		)
		return
	}
	applyUserToModel(&plan, &userResult.User)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.clients.IamIdentityCenter

	var userId string
	var instanceId string
	if !state.Id.IsNull() && !state.Id.IsUnknown() {
		userId = state.Id.ValueString()
	} else if !state.UserUuid.IsNull() && !state.UserUuid.IsUnknown() {
		userId = state.UserUuid.ValueString()
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		instanceId = state.InstanceId.ValueString()
	}

	result, httpResp, err := client.GetUser(ctx, instanceId, userId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Reading IAM Identity Center User",
			err.Error(),
		)
		return
	}

	applyUserToModel(&state, &result.User)

	if state.Password.IsUnknown() {
		state.Password = types.StringNull()
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func resolveUserUpdateUserId(plan, state *UserModel) string {
	if !plan.Id.IsNull() && !plan.Id.IsUnknown() {
		return plan.Id.ValueString()
	}
	if !plan.UserUuid.IsNull() && !plan.UserUuid.IsUnknown() {
		return plan.UserUuid.ValueString()
	}
	if !state.Id.IsNull() && !state.Id.IsUnknown() {
		return state.Id.ValueString()
	}
	if !state.UserUuid.IsNull() && !state.UserUuid.IsUnknown() {
		return state.UserUuid.ValueString()
	}
	return ""
}

func resolveUserUpdateInstanceId(plan, state *UserModel) string {
	if !plan.InstanceId.IsNull() && !plan.InstanceId.IsUnknown() {
		return plan.InstanceId.ValueString()
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		return state.InstanceId.ValueString()
	}
	return ""
}

func resolveUserPassword(planPwd, statePwd types.String) types.String {
	if planPwd.IsUnknown() {
		if !statePwd.IsNull() && !statePwd.IsUnknown() {
			return statePwd
		}
		return types.StringNull()
	}
	return planPwd
}

func populateUserFromGetResult(ctx context.Context, plan *UserModel, client *iamidentitycenterClient.Client, instanceId, userId string) {
	userResult, _, err := client.GetUser(ctx, instanceId, userId)
	if err != nil {
		return
	}

	applyUserToModel(plan, &userResult.User)
}

func (r *iamIdentityCenterUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserModel
	var state UserModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.clients.IamIdentityCenter

	userId := resolveUserUpdateUserId(&plan, &state)
	instanceId := resolveUserUpdateInstanceId(&plan, &state)

	updateReq := *sdk.NewUserSetRequest(instanceId)

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updateReq.Name = *sdk.NewNullableString(plan.Name.ValueStringPointer())
	}
	if !plan.Email.IsNull() && !plan.Email.IsUnknown() {
		updateReq.Email = *sdk.NewNullableString(plan.Email.ValueStringPointer())
	}
	if !plan.PhoneNumber.IsNull() && !plan.PhoneNumber.IsUnknown() {
		updateReq.PhoneNumber = *sdk.NewNullableString(plan.PhoneNumber.ValueStringPointer())
	}
	if !plan.BusinessUnit.IsNull() && !plan.BusinessUnit.IsUnknown() {
		updateReq.BusinessUnit = *sdk.NewNullableString(plan.BusinessUnit.ValueStringPointer())
	}
	if !plan.Department.IsNull() && !plan.Department.IsUnknown() {
		updateReq.Department = *sdk.NewNullableString(plan.Department.ValueStringPointer())
	}
	if !plan.Manager.IsNull() && !plan.Manager.IsUnknown() {
		updateReq.Manager = *sdk.NewNullableString(plan.Manager.ValueStringPointer())
	}
	if !plan.EmployeeId.IsNull() && !plan.EmployeeId.IsUnknown() {
		updateReq.EmployeeId = *sdk.NewNullableString(plan.EmployeeId.ValueStringPointer())
	}
	if !plan.NationId.IsNull() && !plan.NationId.IsUnknown() {
		updateReq.NationId = *sdk.NewNullableString(plan.NationId.ValueStringPointer())
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateReq.Description = *sdk.NewNullableString(plan.Description.ValueStringPointer())
	}

	result, err := client.UpdateUser(ctx, userId, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Updating IAM Identity Center User",
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(result.User.Id)
	plan.UserUuid = types.StringValue(result.User.Id)

	plan.Password = resolveUserPassword(plan.Password, state.Password)

	populateUserFromGetResult(ctx, &plan, client, instanceId, userId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.clients.IamIdentityCenter

	var userId string
	var instanceId string
	if !state.Id.IsNull() && !state.Id.IsUnknown() {
		userId = state.Id.ValueString()
	} else if !state.UserUuid.IsNull() && !state.UserUuid.IsUnknown() {
		userId = state.UserUuid.ValueString()
	}
	if !state.InstanceId.IsNull() && !state.InstanceId.IsUnknown() {
		instanceId = state.InstanceId.ValueString()
	}

	if userId == "" || instanceId == "" {
		resp.Diagnostics.AddError(
			"Deleting IAM Identity Center User",
			"User ID or Instance ID is missing in state",
		)
		return
	}

	_, httpResp, err := client.DeleteUser(ctx, userId, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Deleting IAM Identity Center User",
			err.Error(),
		)
		return
	}
}

func (r *iamIdentityCenterUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), parts[0])...)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func stringValuePtr(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}

type UserModel struct {
	Id                types.String `tfsdk:"id"`
	InstanceId        types.String `tfsdk:"instance_id"`
	UserId            types.String `tfsdk:"user_id"`
	Name              types.String `tfsdk:"name"`
	Password          types.String `tfsdk:"password"`
	Email             types.String `tfsdk:"email"`
	PhoneNumber       types.String `tfsdk:"phone_number"`
	TemporaryPassword types.Bool   `tfsdk:"temporary_password"`
	BusinessUnit      types.String `tfsdk:"business_unit"`
	Department        types.String `tfsdk:"department"`
	Manager           types.String `tfsdk:"manager"`
	EmployeeId        types.String `tfsdk:"employee_id"`
	NationId          types.String `tfsdk:"nation_id"`
	Description       types.String `tfsdk:"description"`
	UserUuid          types.String `tfsdk:"user_uuid"`
}

func applyUserToModel(model *UserModel, detail *sdk.UserDetail) {
	model.Id = types.StringValue(detail.Id)
	model.UserUuid = types.StringValue(detail.Id)
	model.UserId = types.StringValue(detail.UserId)
	model.Name = types.StringValue(detail.Name)
	model.Email = stringValuePtr(detail.Email.Get())
	model.PhoneNumber = stringValuePtr(detail.PhoneNumber.Get())
	model.BusinessUnit = stringValuePtr(detail.BusinessUnit.Get())
	model.Department = stringValuePtr(detail.Department.Get())
	model.Manager = stringValuePtr(detail.Manager.Get())
	model.EmployeeId = stringValuePtr(detail.EmployeeId.Get())
	model.NationId = stringValuePtr(detail.NationId.Get())
	model.Description = stringValuePtr(detail.Description.Get())
	model.TemporaryPassword = types.BoolValue(detail.TemporaryPassword)
}
