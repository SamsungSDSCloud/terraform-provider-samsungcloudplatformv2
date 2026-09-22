package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamIdentityCenterUsersDataSources{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterUsersDataSources{}
)

func NewIamIdentityCenterUsersDataSources() datasource.DataSource {
	return &iamIdentityCenterUsersDataSources{}
}

type iamIdentityCenterUsersDataSources struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterUsersDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_users"
}

func (r *iamIdentityCenterUsersDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IAM Identity Center Users.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"user_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Filter users by user ID.\n  - example: johndoe",
				MarkdownDescription: "Filter users by user ID.\n  - example: johndoe",
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort the results by a specific field.\n  - example: user_id:asc",
				MarkdownDescription: "Sort the results by a specific field.\n  - example: user_id:asc",
			},
			"excluded_group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Exclude users that belong to the specified group ID.\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Exclude users that belong to the specified group ID.\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"excluded_account_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Exclude users that are assigned to the specified account ID.\n  - example: 5b26e6b693a342928d1110a52e657fcb",
				MarkdownDescription: "Exclude users that are assigned to the specified account ID.\n  - example: 5b26e6b693a342928d1110a52e657fcb",
			},
			"size": schema.Int32Attribute{
				Optional: true,
				Description: "Number of results to return per page.\n" +
					"  - example : 100\n" +
					"  - min: 1, max: 10000",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional: true,
				Description: "Page number to retrieve.\n" +
					"  - example : 0\n" +
					"  - min: 0, max: 10000",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of users matching the filter criteria.\n" +
					"  - example : [{\"id\": \"user-id\", \"user_id\": \"johndoe\", \"name\": \"John Doe\"}]",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
							MarkdownDescription: "User ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
						},
						"user_id": schema.StringAttribute{
							Computed:            true,
							Description:         "User Login ID\n  - example: johndoe",
							MarkdownDescription: "User Login ID\n  - example: johndoe",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Real Name\n  - example: John Doe",
							MarkdownDescription: "Real Name\n  - example: John Doe",
						},
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterUsersDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type usersDataSourcesModel struct {
	InstanceId        types.String              `tfsdk:"instance_id"`
	UserId            types.String              `tfsdk:"user_id"`
	Size              types.Int32               `tfsdk:"size"`
	Page              types.Int32               `tfsdk:"page"`
	Sort              types.String              `tfsdk:"sort"`
	ExcludedGroupId   types.String              `tfsdk:"excluded_group_id"`
	ExcludedAccountId types.String              `tfsdk:"excluded_account_id"`
	Users             []userDataSourceItemModel `tfsdk:"users"`
}

type userDataSourceItemModel struct {
	Id     types.String `tfsdk:"id"`
	UserId types.String `tfsdk:"user_id"`
	Name   types.String `tfsdk:"name"`
}

func (r *iamIdentityCenterUsersDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourcesModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	userId := state.UserId.ValueString()
	size := state.Size.ValueInt32()
	page := state.Page.ValueInt32()
	sort := state.Sort.ValueString()
	excludedGroupId := state.ExcludedGroupId.ValueString()
	excludedAccountId := state.ExcludedAccountId.ValueString()

	result, err := r.clients.IamIdentityCenter.ListUsers(ctx, instanceId, userId, size, page, sort, excludedGroupId, excludedAccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Users",
			err.Error(),
		)
		return
	}

	if result != nil && result.Users != nil {
		state.Users = make([]userDataSourceItemModel, len(result.Users))
		for i, user := range result.Users {
			state.Users[i] = userDataSourceItemModel{
				Id:     types.StringValue(user.Id),
				UserId: types.StringValue(user.UserId),
				Name:   types.StringValue(user.Name),
			}
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
