package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource              = &iamGroupMemberResource{}
	_ resource.ResourceWithConfigure = &iamGroupMemberResource{}
)

func NewIamGroupMemberResource() resource.Resource {
	return &iamGroupMemberResource{}
}

type iamGroupMemberResource struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (r *iamGroupMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group_member"
}

func (r *iamGroupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.client = inst.Client.Iam
	r.clients = inst.Client
}

func (r *iamGroupMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Group Member",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Group ID",
				MarkdownDescription: "Group ID",
			},
			"user_id": schema.StringAttribute{
				Optional:            true,
				Description:         "User ID",
				MarkdownDescription: "User ID",
			},
			"group_member": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Group member",
				MarkdownDescription: "Group member",
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed:            true,
						Description:         "생성 일시",
						MarkdownDescription: "생성 일시",
					},
					"created_by": schema.StringAttribute{
						Computed:            true,
						Description:         "생성자",
						MarkdownDescription: "생성자",
					},
					"creator_created_at": schema.StringAttribute{
						Optional:            true,
						Description:         "생성 일시",
						MarkdownDescription: "생성 일시",
					},
					"creator_email": schema.StringAttribute{
						Computed:            true,
						Description:         "생성자 Email",
						MarkdownDescription: "생성자 Email",
						Default:             stringdefault.StaticString("-"),
					},
					"creator_last_login_at": schema.StringAttribute{
						Optional:            true,
						Description:         "생성자 마지막 로그인 일시",
						MarkdownDescription: "생성자 마지막 로그인 일시",
					},
					"creator_name": schema.StringAttribute{
						Computed:            true,
						Description:         "생성자 성, 이름",
						MarkdownDescription: "생성자 성, 이름",
						Default:             stringdefault.StaticString("-"),
					},
					"group_names": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Description:         "Group Names",
						MarkdownDescription: "Group Names",
					},
					"user_created_at": schema.StringAttribute{
						Optional:            true,
						Description:         "생성 일시",
						MarkdownDescription: "생성 일시",
					},
					"user_email": schema.StringAttribute{
						Computed:            true,
						Description:         "User Email",
						MarkdownDescription: "User Email",
						Default:             stringdefault.StaticString("-"),
					},
					"user_id": schema.StringAttribute{
						Computed:            true,
						Description:         "User ID",
						MarkdownDescription: "User ID",
					},
					"user_last_login_at": schema.StringAttribute{
						Optional:            true,
						Description:         "User 마지막 로그인 일시",
						MarkdownDescription: "User 마지막 로그인 일시",
					},
					"user_name": schema.StringAttribute{
						Computed:            true,
						Description:         "User 성, 이름",
						MarkdownDescription: "User 성, 이름",
						Default:             stringdefault.StaticString("-"),
					},
				},
			},
		},
	}
}

func (r *iamGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.GroupMemberResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.AddGroupMember(ctx, plan.GroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add group member",
			"Could not add group member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	groupMemberState := iam.Member{
		CreatedAt:          types.StringValue(data.GroupMember.CreatedAt.Format(time.RFC3339)),
		CreatedBy:          types.StringValue(data.GroupMember.CreatedBy),
		CreatorCreatedAt:   types.StringValue(""),
		CreatorEmail:       types.StringPointerValue(data.GroupMember.CreatorEmail),
		CreatorLastLoginAt: types.StringValue(""),
		CreatorName:        types.StringPointerValue(data.GroupMember.CreatorName),
		GroupNames:         make([]types.String, 0),
		UserCreatedAt:      types.StringValue(""),
		UserEmail:          types.StringPointerValue(data.GroupMember.UserEmail),
		UserId:             types.StringValue(data.GroupMember.UserId),
		UserLastLoginAt:    types.StringValue(""),
		UserName:           types.StringPointerValue(data.GroupMember.UserName),
	}
	groupMemberObjectValue, diags := types.ObjectValueFrom(ctx, groupMemberState.AttributeTypes(), groupMemberState)
	plan.GroupMember = groupMemberObjectValue
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamGroupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam.GroupMemberResource
	var state iam.GroupMemberResource

	diags := req.Plan.Get(ctx, &plan)
	req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.UserId == state.UserId {
		resp.Diagnostics.AddWarning(
			"Could not update group member",
			"Could not update group member, unexpected error: duplicate user id",
		)
		return
	}

	// detach
	err := r.client.RemoveGroupMember(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam group member",
			"Could not delete Group Member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// attach
	data, err := r.client.AddGroupMember(ctx, plan.GroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add group member",
			"Could not add group member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	groupMemberState := iam.Member{
		CreatedAt:          types.StringValue(data.GroupMember.CreatedAt.Format(time.RFC3339)),
		CreatedBy:          types.StringValue(data.GroupMember.CreatedBy),
		CreatorCreatedAt:   types.StringValue(""),
		CreatorEmail:       types.StringPointerValue(data.GroupMember.CreatorEmail),
		CreatorLastLoginAt: types.StringValue(""),
		CreatorName:        types.StringPointerValue(data.GroupMember.CreatorName),
		GroupNames:         make([]types.String, 0),
		UserCreatedAt:      types.StringValue(""),
		UserEmail:          types.StringPointerValue(data.GroupMember.UserEmail),
		UserId:             types.StringValue(data.GroupMember.UserId),
		UserLastLoginAt:    types.StringValue(""),
		UserName:           types.StringPointerValue(data.GroupMember.UserName),
	}
	groupMemberObjectValue, diags := types.ObjectValueFrom(ctx, groupMemberState.AttributeTypes(), groupMemberState)
	plan.GroupMember = groupMemberObjectValue
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamGroupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam.GroupMemberResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveGroupMember(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting iam group member",
			"Could not delete Group Member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

}

func (r *iamGroupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam.GroupMemberResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetGroupMembers(ctx, state.GroupId.ValueString(), iam.GroupMembersDataResource{Size: basetypes.NewInt32Value(20)})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Members",
			err.Error(),
		)
		return
	}

	var groupMemberState iam.Member
	for _, group := range data.GetGroupMembers() {
		if state.UserId.ValueString() == group.UserId {
			groupMemberState = iam.Member{
				CreatedAt:          types.StringValue(group.CreatedAt.Format(time.RFC3339)),
				CreatedBy:          types.StringValue(group.CreatedBy),
				CreatorCreatedAt:   types.StringValue(""),
				CreatorEmail:       types.StringPointerValue(group.CreatorEmail),
				CreatorLastLoginAt: types.StringValue(""),
				CreatorName:        types.StringPointerValue(group.CreatorName),
				GroupNames:         make([]types.String, 0),
				UserCreatedAt:      types.StringValue(""),
				UserEmail:          types.StringPointerValue(group.UserEmail),
				UserId:             types.StringValue(group.UserId),
				UserLastLoginAt:    types.StringValue(""),
				UserName:           types.StringPointerValue(group.UserName),
			}
		}
	}

	groupMemberObjectValue, diags := types.ObjectValueFrom(ctx, groupMemberState.AttributeTypes(), groupMemberState)
	state.GroupMember = groupMemberObjectValue

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
