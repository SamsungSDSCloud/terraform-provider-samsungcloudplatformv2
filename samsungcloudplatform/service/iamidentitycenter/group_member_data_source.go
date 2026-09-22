package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamIdentityCenterGroupMemberDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterGroupMemberDataSource{}
)

func NewIamIdentityCenterGroupMemberDataSource() datasource.DataSource {
	return &iamIdentityCenterGroupMemberDataSource{}
}

type iamIdentityCenterGroupMemberDataSource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterGroupMemberDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_group_member"
}

func (r *iamIdentityCenterGroupMemberDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages IAM Identity Center Group Members (Data Source).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group Member ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Group Member ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				Description:         "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"user_uuids": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "List of User IDs in the group\n  - example: [138c2fc8c29a449dbfa8681f8f1d78e2]",
				MarkdownDescription: "List of User IDs in the group\n  - example: [138c2fc8c29a449dbfa8681f8f1d78e2]",
			},
		},
	}
}

func (r *iamIdentityCenterGroupMemberDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type groupMemberDataSourceModel struct {
	GroupId    types.String `tfsdk:"group_id"`
	InstanceId types.String `tfsdk:"instance_id"`
	UserUuids  types.List   `tfsdk:"user_uuids"`
	Id         types.String `tfsdk:"id"`
}

func (r *iamIdentityCenterGroupMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupMemberDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.clients.IamIdentityCenter.ListGroupUsers(
		ctx,
		state.GroupId.ValueString(),
		state.InstanceId.ValueString(),
		"",
		0,
		0,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Group Members",
			err.Error(),
		)
		return
	}

	var userUuids []string
	if result != nil {
		for _, user := range result.Users {
			userUuids = append(userUuids, user.Id)
		}
	}

	state.UserUuids, diags = types.ListValueFrom(ctx, types.StringType, userUuids)
	resp.Diagnostics.Append(diags...)

	state.Id = types.StringValue(state.GroupId.ValueString() + ":" + state.InstanceId.ValueString())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
