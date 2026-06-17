package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/common/importstate"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	scpsdkiam "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/iam/1.4"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.Resource                = &iamGroupMemberResource{}
	_ resource.ResourceWithConfigure   = &iamGroupMemberResource{}
	_ resource.ResourceWithImportState = &iamGroupMemberResource{}
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
		Description: "List of members belonging to the group.",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional: true,
				Description: "Unique identifier of the group to add the member to.\n" +
					"  - example : 'grp-1234567890abcdef'",
			},
			"user_id": schema.StringAttribute{
				Optional: true,
				Description: "Unique identifier of the user to add as a member.\n" +
					"  - example : 'usr-1234567890abcdef'",
			},
			"group_member": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Detailed information about the group member.\n" +
					"  - example : '{created_at: 2024-05-17T00:23:17Z, created_by: ef50cdc207f05f6fb8f20219f229ed1f, creator_email: -, user_id: f39c460fade34fecb05ede8f904b24b7, user_name: -, ...}'",
				PlanModifiers: []planmodifier.Object{
					groupMemberModifier{},
				},
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "생성 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						MarkdownDescription: "생성 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "생성자\n" +
							"  - example : 'user@example.com'",
						MarkdownDescription: "생성자\n" +
							"  - example : 'user@example.com'",
					},
					"creator_created_at": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "생성 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						MarkdownDescription: "생성 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"creator_email": schema.StringAttribute{
						Computed: true,
						Description: "생성자 Email\n" +
							"  - example : 'admin@example.com'",
						MarkdownDescription: "생성자 Email\n" +
							"  - example : 'admin@example.com'",
					},
					"creator_last_login_at": schema.StringAttribute{
						Computed: true,
						Description: "생성자 마지막 로그인 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						MarkdownDescription: "생성자 마지막 로그인 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "생성자 성, 이름\n" +
							"  - example : 'Admin User'",
						MarkdownDescription: "생성자 성, 이름\n" +
							"  - example : 'Admin User'",
						Default: stringdefault.StaticString("-"),
					},
					"groups": schema.ListNestedAttribute{
						Computed: true,
						Description: "List of groups the user belongs to.\n" +
							"  - example : '[{id: grp-1234567890abcdef, name: MyGroup}]'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
									Description: "Group ID\n" +
										"  - example : 'grp-1234567890abcdef'",
									MarkdownDescription: "Group ID\n" +
										"  - example : 'grp-1234567890abcdef'",
								},
								"name": schema.StringAttribute{
									Computed: true,
									Description: "Group Name\n" +
										"  - example : 'MyGroup'",
									MarkdownDescription: "Group Name\n" +
										"  - example : 'MyGroup'",
								},
							},
						},
					},
					"user_created_at": schema.StringAttribute{
						Computed: true,
						Optional: true,
						Description: "생성 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						MarkdownDescription: "생성 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"user_email": schema.StringAttribute{
						Computed: true,
						Description: "Email address of the user who is a member of the group.\n" +
							"  - example : 'user@example.com'",
						MarkdownDescription: "Email address of the user who is a member of the group.\n" +
							"  - example : 'user@example.com'",
					},
					"user_id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the user who is a member of the group.\n" +
							"  - example : 'usr-1234567890abcdef'",
						MarkdownDescription: "Unique identifier of the user who is a member of the group.\n" +
							"  - example : 'usr-1234567890abcdef'",
					},
					"user_last_login_at": schema.StringAttribute{
						Computed: true,
						Description: "User 마지막 로그인 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
						MarkdownDescription: "User 마지막 로그인 일시\n" +
							"  - example : '2024-01-01T00:00:00Z'",
					},
					"user_name": schema.StringAttribute{
						Computed: true,
						Description: "User 성, 이름\n" +
							"  - example : 'John Doe'",
						MarkdownDescription: "User 성, 이름\n" +
							"  - example : 'John Doe'",
						Default: stringdefault.StaticString("-"),
					},
				},
			},
		},
	}
}

func (r *iamGroupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importstate.ImportState(ctx, req, resp,
		path.Root("group_id"),
		path.Root("user_id"),
	)
}

type groupMemberModifier struct{}

func (m groupMemberModifier) Description(ctx context.Context) string {
	return "Keeps the existing group_member state if group_id and user_id remain unchanged."
}

func (m groupMemberModifier) MarkdownDescription(ctx context.Context) string {
	return "Keeps the existing `group_member` state if `group_id` and `user_id` remain unchanged."
}

func (m groupMemberModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If there is no state (creating a new resource), do nothing
	if req.State.Raw.IsNull() {
		return
	}

	var stateGroupId, planGroupId, stateUserId, planUserId string

	// Retrieve group_id and user_id from both State and Plan
	req.State.GetAttribute(ctx, path.Root("group_id"), &stateGroupId)
	req.Plan.GetAttribute(ctx, path.Root("group_id"), &planGroupId)
	req.State.GetAttribute(ctx, path.Root("user_id"), &stateUserId)
	req.Plan.GetAttribute(ctx, path.Root("user_id"), &planUserId)

	// Compare IDs
	if stateGroupId == planGroupId && stateUserId == planUserId {
		// IDs haven't changed: Preserve the existing state for group_member
		resp.PlanValue = req.StateValue
	} else {
		// IDs changed: Mark group_member as Unknown to trigger an update/fetch
		resp.PlanValue = types.ObjectUnknown(req.PlanValue.AttributeTypes(ctx))
	}
}

func (r *iamGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam.GroupMemberResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.AddGroupMember(ctx, plan.GroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add group member",
			"Could not add group member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	// No polling needed - group membership is a synchronous operation (immediate association).
	// The API returns immediately after creating the membership, so we fetch the state directly.

	data, err := r.client.GetGroupMembers(ctx, plan.GroupId.ValueString(), iam.GroupMembersDataResource{Size: basetypes.NewInt32Value(20)})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Members",
			err.Error(),
		)
		return
	}

	var groupMemberState iam.MemberV1Dot4
	for _, group := range data.GetGroupMembers() {
		if plan.UserId.ValueString() == group.UserId {
			var creatorLastLoginAt *string
			var userLastLoginAt *string

			if group.CreatorLastLoginAt.Get() != nil {
				t := group.CreatorLastLoginAt.Get().Format(time.RFC3339)
				creatorLastLoginAt = &t
			}
			if group.UserLastLoginAt.Get() != nil {
				t := group.UserLastLoginAt.Get().Format(time.RFC3339)
				userLastLoginAt = &t
			}

			var groupInfos []iam.GroupInfo
			for _, groupInfo := range group.Groups {
				groupInfos = append(groupInfos, iam.GroupInfo{
					Id:   types.StringValue(groupInfo.Id),
					Name: types.StringValue(groupInfo.Name),
				})
			}

			groupMemberState = iam.MemberV1Dot4{
				CreatedAt:          types.StringValue(group.CreatedAt.Format(time.RFC3339)),
				CreatedBy:          types.StringValue(group.CreatedBy),
				CreatorCreatedAt:   types.StringValue(group.CreatorCreatedAt.Format(time.RFC3339)),
				CreatorEmail:       types.StringPointerValue(group.CreatorEmail),
				CreatorLastLoginAt: types.StringPointerValue(creatorLastLoginAt),
				CreatorName:        types.StringPointerValue(group.CreatorName),
				Groups:             groupInfos,
				UserCreatedAt:      types.StringValue(group.UserCreatedAt.Format(time.RFC3339)),
				UserEmail:          types.StringPointerValue(group.UserEmail),
				UserId:             types.StringValue(group.UserId),
				UserLastLoginAt:    types.StringPointerValue(userLastLoginAt),
				UserName:           types.StringPointerValue(group.UserName),
			}
		}
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
	_, err = r.client.AddGroupMember(ctx, plan.GroupId.ValueString(), plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error add group member",
			"Could not add group member, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetGroupMembers(ctx, plan.GroupId.ValueString(), iam.GroupMembersDataResource{Size: basetypes.NewInt32Value(20)})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Members",
			err.Error(),
		)
		return
	}

	var groupMemberState iam.MemberV1Dot4
	for _, group := range data.GetGroupMembers() {
		if plan.UserId.ValueString() == group.UserId {
			var creatorLastLoginAt *string
			var userLastLoginAt *string

			if group.CreatorLastLoginAt.Get() != nil {
				t := group.CreatorLastLoginAt.Get().Format(time.RFC3339)
				creatorLastLoginAt = &t
			}
			if group.UserLastLoginAt.Get() != nil {
				t := group.UserLastLoginAt.Get().Format(time.RFC3339)
				userLastLoginAt = &t
			}

			var groupInfos []iam.GroupInfo
			for _, groupInfo := range group.Groups {
				groupInfos = append(groupInfos, iam.GroupInfo{
					Id:   types.StringValue(groupInfo.Id),
					Name: types.StringValue(groupInfo.Name),
				})
			}

			groupMemberState = iam.MemberV1Dot4{
				CreatedAt:          types.StringValue(group.CreatedAt.Format(time.RFC3339)),
				CreatedBy:          types.StringValue(group.CreatedBy),
				CreatorCreatedAt:   types.StringValue(group.CreatorCreatedAt.Format(time.RFC3339)),
				CreatorEmail:       types.StringPointerValue(group.CreatorEmail),
				CreatorLastLoginAt: types.StringPointerValue(creatorLastLoginAt),
				CreatorName:        types.StringPointerValue(group.CreatorName),
				Groups:             groupInfos,
				UserCreatedAt:      types.StringValue(group.UserCreatedAt.Format(time.RFC3339)),
				UserEmail:          types.StringPointerValue(group.UserEmail),
				UserId:             types.StringValue(group.UserId),
				UserLastLoginAt:    types.StringPointerValue(userLastLoginAt),
				UserName:           types.StringPointerValue(group.UserName),
			}
		}
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
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Read Group Members",
			err.Error(),
		)
		return
	}

	var groupMemberState iam.MemberV1Dot4
	for _, group := range data.GetGroupMembers() {
		if state.UserId.ValueString() == group.UserId {
			var creatorLastLoginAt *string
			var userLastLoginAt *string

			if group.CreatorLastLoginAt.Get() != nil {
				t := group.CreatorLastLoginAt.Get().Format(time.RFC3339)
				creatorLastLoginAt = &t
			}
			if group.UserLastLoginAt.Get() != nil {
				t := group.UserLastLoginAt.Get().Format(time.RFC3339)
				userLastLoginAt = &t
			}

			var groupInfos []iam.GroupInfo
			for _, groupInfo := range group.Groups {
				groupInfos = append(groupInfos, iam.GroupInfo{
					Id:   types.StringValue(groupInfo.Id),
					Name: types.StringValue(groupInfo.Name),
				})
			}

			groupMemberState = iam.MemberV1Dot4{
				CreatedAt:          types.StringValue(group.CreatedAt.Format(time.RFC3339)),
				CreatedBy:          types.StringValue(group.CreatedBy),
				CreatorCreatedAt:   types.StringValue(group.CreatorCreatedAt.Format(time.RFC3339)),
				CreatorEmail:       types.StringPointerValue(group.CreatorEmail),
				CreatorLastLoginAt: types.StringPointerValue(creatorLastLoginAt),
				CreatorName:        types.StringPointerValue(group.CreatorName),
				Groups:             groupInfos,
				UserCreatedAt:      types.StringValue(group.UserCreatedAt.Format(time.RFC3339)),
				UserEmail:          types.StringPointerValue(group.UserEmail),
				UserId:             types.StringValue(group.UserId),
				UserLastLoginAt:    types.StringPointerValue(userLastLoginAt),
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

func getGroupMembersV1Dot4(_members []scpsdkiam.GroupMemberV1Dot4) []iam.MemberV1Dot4 {
	var members []iam.MemberV1Dot4

	for _, member := range _members {
		var creatorLastLoginAt *string
		var userLastLoginAt *string

		if member.CreatorLastLoginAt.Get() != nil {
			t := member.CreatorLastLoginAt.Get().Format(time.RFC3339)
			creatorLastLoginAt = &t
		}
		if member.UserLastLoginAt.Get() != nil {
			t := member.UserLastLoginAt.Get().Format(time.RFC3339)
			userLastLoginAt = &t
		}

		var groupInfos []iam.GroupInfo
		for _, groupInfo := range member.Groups {
			groupInfos = append(groupInfos, iam.GroupInfo{
				Id:   types.StringValue(groupInfo.Id),
				Name: types.StringValue(groupInfo.Name),
			})
		}

		memberState := iam.MemberV1Dot4{
			CreatedAt:          types.StringValue(member.CreatedAt.Format(time.RFC3339)),
			CreatedBy:          types.StringValue(member.CreatedBy),
			CreatorCreatedAt:   types.StringValue(member.CreatorCreatedAt.Format(time.RFC3339)),
			CreatorEmail:       types.StringPointerValue(member.CreatorEmail),
			CreatorLastLoginAt: types.StringPointerValue(creatorLastLoginAt),
			CreatorName:        types.StringPointerValue(member.CreatorName),
			Groups:             groupInfos,
			UserCreatedAt:      types.StringValue(member.UserCreatedAt.Format(time.RFC3339)),
			UserEmail:          types.StringPointerValue(member.UserEmail),
			UserId:             types.StringValue(member.UserId),
			UserLastLoginAt:    types.StringPointerValue(userLastLoginAt),
			UserName:           types.StringPointerValue(member.UserName),
		}

		members = append(members, memberState)
	}
	return members
}
