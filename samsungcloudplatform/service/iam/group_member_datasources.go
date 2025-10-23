package iam

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v2/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamGroupMemberDataSources{}
	_ datasource.DataSourceWithConfigure = &iamGroupMemberDataSources{}
)

func NewIamGroupMemberDataSources() datasource.DataSource { return &iamGroupMemberDataSources{} }

type iamGroupMemberDataSources struct {
	config  *scpsdk.Configuration
	client  *iam.Client
	clients *client.SCPClient
}

func (d *iamGroupMemberDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group_members"
}

func (d *iamGroupMemberDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Iam
	d.clients = inst.Client
}

func (d *iamGroupMemberDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Group Members",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Group ID",
				MarkdownDescription: "Group ID",
			},
			"size": schema.Int32Attribute{
				Optional:            true,
				Description:         "Size (between 1 and 10000)",
				MarkdownDescription: "Size (between 1 and 10000)",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional:            true,
				Description:         "Page (between 0 and 10000)",
				MarkdownDescription: "Page (between 0 and 10000)",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort",
				MarkdownDescription: "Sort",
			},
			"user_name": schema.StringAttribute{
				Optional:            true,
				Description:         "User Name",
				MarkdownDescription: "User Name",
			},
			"user_email": schema.StringAttribute{
				Optional:            true,
				Description:         "User Email",
				MarkdownDescription: "User Email",
			},
			"creator_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Creator Name",
				MarkdownDescription: "Creator Name",
			},
			"creator_email": schema.StringAttribute{
				Optional:            true,
				Description:         "Creator Email",
				MarkdownDescription: "Creator Email",
			},
			"group_members": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group Members",
				MarkdownDescription: "Group Members",
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
							Description:         "생성 일시",
							MarkdownDescription: "생성 일시",
						},
						"creator_email": schema.StringAttribute{
							Computed:            true,
							Description:         "생성자 Email",
							MarkdownDescription: "생성자 Email",
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
						},
						"group_names": schema.ListAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							Description:         "Group names",
							MarkdownDescription: "Group names",
						},
						"user_created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "생성 일시",
							MarkdownDescription: "생성 일시",
						},
						"user_email": schema.StringAttribute{
							Computed:            true,
							Description:         "User Email",
							MarkdownDescription: "User Email",
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
						},
					},
				},
			},
		},
	}
}

func (d *iamGroupMemberDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam.GroupMembersDataResource

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.GetGroupMembers(ctx, state.GroupId.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Group Members",
			err.Error(),
		)
		return
	}

	state.GroupMembers = getGroupMembers(data.GroupMembers)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
