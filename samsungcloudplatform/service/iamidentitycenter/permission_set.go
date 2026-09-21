package iamidentitycenter

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &iamIdentityCenterPermissionSetResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterPermissionSetResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterPermissionSetResource{}
)

func NewIamIdentityCenterPermissionSetResource() resource.Resource {
	return &iamIdentityCenterPermissionSetResource{}
}

type iamIdentityCenterPermissionSetResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_set"
}

func (r *iamIdentityCenterPermissionSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Identity Center Permission Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
				MarkdownDescription: "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Description:         "Permission Set Name\n  - example: Admin_Permission_Set\n  - maxLength: 128",
				MarkdownDescription: "Permission Set Name\n  - example: Admin_Permission_Set\n  - maxLength: 128",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Permission Set Description\n  - example: Permission set for administrators",
				MarkdownDescription: "Permission Set Description\n  - example: Permission set for administrators",
			},
			"session_duration": schema.Int64Attribute{
				Required:            true,
				Description:         "Maximum Session Duration (seconds)\n  - example: 3600\n  - minimum: 900\n  - maximum: 43200",
				MarkdownDescription: "Maximum Session Duration (seconds)\n  - example: 3600\n  - minimum: 900\n  - maximum: 43200",
			},
			"srn": schema.StringAttribute{
				Computed:            true,
				Description:         "Permission Set SRN\n  - example: srn:<offering>::<project-id>:<region>::permissionSet/${InstanceId}/${permissionId}",
				MarkdownDescription: "Permission Set SRN\n  - example: srn:<offering>::<project-id>:<region>::permissionSet/${InstanceId}/${permissionId}",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "Permission Set State\n  - example: ACTIVE\n  - enum: [\"ACTIVE\", \"INACTIVE\"]",
				MarkdownDescription: "Permission Set State\n  - example: ACTIVE\n  - enum: [\"ACTIVE\", \"INACTIVE\"]",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Created At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Created At\n  - example: 2024-05-17T00:23:17Z",
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
				MarkdownDescription: "Created By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
			},
			"creator_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Creator Name\n  - example: John Doe",
				MarkdownDescription: "Creator Name\n  - example: John Doe",
			},
			"modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified At\n  - example: 2024-05-17T00:23:17Z",
				MarkdownDescription: "Modified At\n  - example: 2024-05-17T00:23:17Z",
			},
			"modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
				MarkdownDescription: "Modified By\n  - example: 90dddfc2b1e04edba54ba2b41539a9ac",
			},
			"modifier_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Modifier Name\n  - example: Smith",
				MarkdownDescription: "Modifier Name\n  - example: Smith",
			},
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: 123456789012",
				MarkdownDescription: "Account ID\n  - example: 123456789012",
			},
			"resource_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource Name",
				MarkdownDescription: "Resource Name",
			},
			"resource_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource Type",
				MarkdownDescription: "Resource Type",
			},
			"resource_type_display_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource Type Display Name",
				MarkdownDescription: "Resource Type Display Name",
			},
			"service": schema.StringAttribute{
				Computed:            true,
				Description:         "Service",
				MarkdownDescription: "Service",
			},
			"service_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Service Name",
				MarkdownDescription: "Service Name",
			},
		},
	}
}

func (r *iamIdentityCenterPermissionSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamIdentityCenterPermissionSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type permissionSetResourceModel struct {
	InstanceId              types.String `tfsdk:"instance_id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	SessionDuration         types.Int64  `tfsdk:"session_duration"`
	Id                      types.String `tfsdk:"id"`
	Srn                     types.String `tfsdk:"srn"`
	State                   types.String `tfsdk:"state"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	CreatorName             types.String `tfsdk:"creator_name"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifierName            types.String `tfsdk:"modifier_name"`
	AccountId               types.String `tfsdk:"account_id"`
	ResourceName            types.String `tfsdk:"resource_name"`
	ResourceType            types.String `tfsdk:"resource_type"`
	ResourceTypeDisplayName types.String `tfsdk:"resource_type_display_name"`
	Service                 types.String `tfsdk:"service"`
	ServiceName             types.String `tfsdk:"service_name"`
}

func applyPermissionSetToModel(model *permissionSetResourceModel, result *sdk.ShowPermissionSetAddResourceV1Dot2) {
	model.Id = types.StringValue(result.Id)
	model.Name = types.StringValue(result.Name)
	model.Description = types.StringPointerValue(result.Description.Get())
	model.SessionDuration = types.Int64Value(int64(result.SessionDuration))
	model.Srn = types.StringValue(result.Srn)
	if result.State != nil {
		model.State = types.StringValue(string(*result.State))
	}
	model.CreatedAt = types.StringValue(result.CreatedAt.Format("2006-01-02T15:04:05Z"))
	model.CreatedBy = types.StringValue(result.CreatedBy)
	model.CreatorName = types.StringValue(result.CreatorName)
	model.ModifiedAt = types.StringValue(result.ModifiedAt.Format("2006-01-02T15:04:05Z"))
	model.ModifiedBy = types.StringValue(result.ModifiedBy)
	model.ModifierName = types.StringValue(result.ModifierName)
	model.AccountId = types.StringPointerValue(result.AccountId.Get())
	model.ResourceName = types.StringPointerValue(result.ResourceName.Get())
	model.ResourceType = types.StringPointerValue(result.ResourceType.Get())
	model.ResourceTypeDisplayName = types.StringPointerValue(result.ResourceTypeDisplayName.Get())
	model.Service = types.StringPointerValue(result.Service.Get())
	model.ServiceName = types.StringPointerValue(result.ServiceName.Get())
}

// populatePermissionSetAfterCreate maps the create/get response metadata onto
// the plan. It re-fetches the permission set so all computed fields are known.
func populatePermissionSetAfterCreate(ctx context.Context, plan permissionSetResourceModel, result *sdk.PermissionSetDetailsV1Dot2, client *iamidentitycenterClient.Client) permissionSetResourceModel {
	if result == nil || result.Id == "" {
		// Create response was empty — set all computed fields to known null to avoid "unknown value" error
		plan.Srn = types.StringNull()
		plan.State = types.StringNull()
		plan.CreatedAt = types.StringNull()
		plan.CreatedBy = types.StringNull()
		plan.CreatorName = types.StringNull()
		plan.ModifiedAt = types.StringNull()
		plan.ModifiedBy = types.StringNull()
		plan.ModifierName = types.StringNull()
		plan.AccountId = types.StringNull()
		plan.ResourceName = types.StringNull()
		plan.ResourceType = types.StringNull()
		plan.ResourceTypeDisplayName = types.StringNull()
		plan.Service = types.StringNull()
		plan.ServiceName = types.StringNull()
		return plan
	}
	plan.Id = types.StringValue(result.Id)

	permissionSetResult, _, err := client.GetPermissionSet(ctx, result.Id, plan.InstanceId.ValueString())
	if err != nil || permissionSetResult == nil || permissionSetResult.Id == "" {
		// GetPermissionSet failed (eventual consistency) — fall back to create response fields
		plan.Name = types.StringValue(result.Name)
		if result.Description.IsSet() {
			plan.Description = types.StringPointerValue(result.Description.Get())
		}
		plan.SessionDuration = types.Int64Value(int64(result.SessionDuration))
		plan.Srn = types.StringValue(result.Srn)
		if result.State != nil {
			plan.State = types.StringValue(string(*result.State))
		} else {
			plan.State = types.StringNull()
		}
		plan.CreatedAt = types.StringValue(result.CreatedAt.Format("2006-01-02T15:04:05Z"))
		plan.CreatedBy = types.StringValue(result.CreatedBy)
		plan.CreatorName = types.StringValue(result.CreatorName)
		plan.ModifiedAt = types.StringValue(result.ModifiedAt.Format("2006-01-02T15:04:05Z"))
		plan.ModifiedBy = types.StringValue(result.ModifiedBy)
		plan.ModifierName = types.StringValue(result.ModifierName)
		// Fields not available in create response — set to null to avoid "unknown value" error
		plan.AccountId = types.StringNull()
		plan.ResourceName = types.StringNull()
		plan.ResourceType = types.StringNull()
		plan.ResourceTypeDisplayName = types.StringNull()
		plan.Service = types.StringNull()
		plan.ServiceName = types.StringNull()
		return plan
	}

	applyPermissionSetToModel(&plan, permissionSetResult)

	return plan
}

func (r *iamIdentityCenterPermissionSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan permissionSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := iamidentitycenterClient.PermissionSetCreateRequest{
		InstanceId:      plan.InstanceId.ValueString(),
		Name:            plan.Name.ValueString(),
		SessionDuration: int32(plan.SessionDuration.ValueInt64()),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() && plan.Description.ValueString() != "" {
		createReq.Description = plan.Description.ValueString()
	}

	result, err := r.clients.IamIdentityCenter.CreatePermissionSet(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create IAM Identity Center Permission Set",
			err.Error(),
		)
		return
	}

	plan = populatePermissionSetAfterCreate(ctx, plan, result, r.clients.IamIdentityCenter)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state permissionSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := state.Id.ValueString()
	instanceId := state.InstanceId.ValueString()
	result, httpResp, err := r.clients.IamIdentityCenter.GetPermissionSet(ctx, permissionSetId, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set",
			err.Error(),
		)
		return
	}

	if result != nil && result.Id != "" {
		applyPermissionSetToModel(&state, result)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state permissionSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan permissionSetResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := state.Id.ValueString()
	instanceId := plan.InstanceId.ValueString()

	updateReq := iamidentitycenterClient.PermissionSetUpdateRequest{
		InstanceId: instanceId,
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateReq.Description = plan.Description.ValueString()
	}

	if !plan.SessionDuration.IsNull() && !plan.SessionDuration.IsUnknown() {
		updateReq.SessionDuration = int32(plan.SessionDuration.ValueInt64())
	}

	_, err := r.clients.IamIdentityCenter.UpdatePermissionSet(ctx, permissionSetId, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM Identity Center Permission Set",
			err.Error(),
		)
		return
	}

	permissionSetResult, _, err := r.clients.IamIdentityCenter.GetPermissionSet(ctx, permissionSetId, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set after update",
			err.Error(),
		)
		return
	}

	if permissionSetResult != nil && permissionSetResult.Id != "" {
		applyPermissionSetToModel(&state, permissionSetResult)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterPermissionSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state permissionSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := state.Id.ValueString()
	instanceId := state.InstanceId.ValueString()
	err := r.clients.IamIdentityCenter.DeletePermissionSet(ctx, permissionSetId, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Permission Set",
			err.Error(),
		)
		return
	}
}
