package iamidentitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	iamidentitycenterClient "github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/iamidentitycenter"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &iamIdentityCenterGroupResource{}
	_ resource.ResourceWithConfigure   = &iamIdentityCenterGroupResource{}
	_ resource.ResourceWithImportState = &iamIdentityCenterGroupResource{}
)

func NewIamIdentityCenterGroupResource() resource.Resource {
	return &iamIdentityCenterGroupResource{}
}

type iamIdentityCenterGroupResource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_group"
}

func (r *iamIdentityCenterGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM Identity Center Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
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
				Description:         "Group Name\n  - example: Admin_Group\n  - maxLength: 128",
				MarkdownDescription: "Group Name\n  - example: Admin_Group\n  - maxLength: 128",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group Description\n  - example: Group for administrators",
				MarkdownDescription: "Group Description\n  - example: Group for administrators",
			},
			"group": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Group Details",
				MarkdownDescription: "Group Details",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						Description:         "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
						MarkdownDescription: "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						Description:         "Group Name\n  - example: Admin_Group",
						MarkdownDescription: "Group Name\n  - example: Admin_Group",
					},
					"description": schema.StringAttribute{
						Computed:            true,
						Description:         "Group Description\n  - example: Group for administrators",
						MarkdownDescription: "Group Description\n  - example: Group for administrators",
					},
					"instance_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Instance ID\n  - example: ssoins-12345",
						MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
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
					"user_count": schema.Int64Attribute{
						Computed:            true,
						Description:         "Number of Users in the group\n  - example: 1",
						MarkdownDescription: "Number of Users in the group\n  - example: 1",
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iamIdentityCenterGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Read requires the instance, so the import ID is "{instance_id}:{group_id}".
	// A bare group id is also accepted for readability when the Caller supplies
	// instance_id another way, but the composite form is the reliable one.
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), parts[0])...)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type groupResourceModel struct {
	InstanceId  types.String `tfsdk:"instance_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Group       types.Object `tfsdk:"group"`
}

type groupDetailModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	InstanceId   types.String `tfsdk:"instance_id"`
	CreatedAt    types.String `tfsdk:"created_at"`
	CreatedBy    types.String `tfsdk:"created_by"`
	CreatorName  types.String `tfsdk:"creator_name"`
	ModifiedAt   types.String `tfsdk:"modified_at"`
	ModifiedBy   types.String `tfsdk:"modified_by"`
	ModifierName types.String `tfsdk:"modifier_name"`
	UserCount    types.Int64  `tfsdk:"user_count"`
}

func (g groupDetailModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"name":          types.StringType,
		"description":   types.StringType,
		"instance_id":   types.StringType,
		"created_at":    types.StringType,
		"created_by":    types.StringType,
		"creator_name":  types.StringType,
		"modified_at":   types.StringType,
		"modified_by":   types.StringType,
		"modifier_name": types.StringType,
		"user_count":    types.Int64Type,
	}
}

func toGroupDetailModel(detail *sdk.GroupDetail) groupDetailModel {
	groupDetail := groupDetailModel{
		Id:           types.StringValue(detail.Id),
		Name:         types.StringValue(detail.Name),
		Description:  types.StringPointerValue(detail.Description.Get()),
		InstanceId:   types.StringValue(detail.InstanceId),
		CreatedAt:    types.StringValue(detail.CreatedAt.Format("2006-01-02T15:04:05Z")),
		CreatedBy:    types.StringValue(detail.CreatedBy),
		CreatorName:  types.StringValue(detail.CreatorName),
		ModifiedAt:   types.StringValue(detail.ModifiedAt.Format("2006-01-02T15:04:05Z")),
		ModifiedBy:   types.StringValue(detail.ModifiedBy),
		ModifierName: types.StringValue(detail.ModifierName),
	}
	if detail.UserCount.Get() != nil {
		groupDetail.UserCount = types.Int64Value(int64(*detail.UserCount.Get()))
	}
	return groupDetail
}

func (r *iamIdentityCenterGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := iamidentitycenterClient.GroupCreateRequest{
		InstanceId: plan.InstanceId.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createReq.Description = plan.Description.ValueString()
	}

	result, err := r.clients.IamIdentityCenter.CreateGroup(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create IAM Identity Center Group",
			err.Error(),
		)
		return
	}

	if result != nil && result.Group.Id != "" {
		plan.Id = types.StringValue(result.Group.Id)

		plan.Name = types.StringValue(result.Group.Name)

		groupResult, _, err := r.clients.IamIdentityCenter.GetGroup(ctx, result.Group.Id, plan.InstanceId.ValueString())
		if err == nil && groupResult != nil && groupResult.Group.Id != "" {
			groupObj, diags := types.ObjectValueFrom(ctx, groupDetailModel{}.AttributeTypes(), toGroupDetailModel(&groupResult.Group))
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				plan.Group = groupObj
			}
		}
	}

	if plan.Group.IsNull() {
		plan.Group = types.ObjectNull(groupDetailModel{}.AttributeTypes())
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()
	instanceId := state.InstanceId.ValueString()
	result, httpResp, err := r.clients.IamIdentityCenter.GetGroup(ctx, groupId, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Group",
			err.Error(),
		)
		return
	}

	if result != nil && result.Group.Id != "" {
		state.Id = types.StringValue(result.Group.Id)

		groupObj, diags := types.ObjectValueFrom(ctx, groupDetailModel{}.AttributeTypes(), toGroupDetailModel(&result.Group))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Group = groupObj
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan groupResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()
	instanceId := plan.InstanceId.ValueString()

	updateReq := iamidentitycenterClient.GroupUpdateRequest{
		InstanceId: instanceId,
		Name:       plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateReq.Description = plan.Description.ValueString()
	}

	_, err := r.clients.IamIdentityCenter.UpdateGroup(ctx, groupId, updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update IAM Identity Center Group",
			err.Error(),
		)
		return
	}

	groupResult, _, err := r.clients.IamIdentityCenter.GetGroup(ctx, groupId, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Group after update",
			err.Error(),
		)
		return
	}

	if groupResult != nil && groupResult.Group.Id != "" {
		state.Id = types.StringValue(groupResult.Group.Id)
		state.Name = types.StringValue(groupResult.Group.Name)
		state.Description = types.StringPointerValue(groupResult.Group.Description.Get())

		groupObj, diags := types.ObjectValueFrom(ctx, groupDetailModel{}.AttributeTypes(), toGroupDetailModel(&groupResult.Group))
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.Group = groupObj
		}
	}

	if state.Group.IsNull() {
		state.Group = types.ObjectNull(groupDetailModel{}.AttributeTypes())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *iamIdentityCenterGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()
	instanceId := state.InstanceId.ValueString()
	_, httpResp, err := r.clients.IamIdentityCenter.DeleteGroup(ctx, groupId, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to delete IAM Identity Center Group",
			err.Error(),
		)
		return
	}
}
