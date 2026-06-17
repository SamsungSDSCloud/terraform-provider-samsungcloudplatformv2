package iam

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/iam"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
				Optional: true,
				Description: "Group ID to filter members.\n" +
					"  - example : 'group-12345678'",
			},
			"size": schema.Int32Attribute{
				Optional: true,
				Description: "Size (between 1 and 10000)\n" +
					"  - example : 100",
				Validators: []validator.Int32{
					int32validator.Between(1, 10000),
				},
			},
			"page": schema.Int32Attribute{
				Optional: true,
				Description: "Page (between 0 and 10000)\n" +
					"  - example : 0",
				Validators: []validator.Int32{
					int32validator.Between(0, 10000),
				},
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort order for results.\n" +
					"  - example : 'created_at,desc'",
			},
			"user_name": schema.StringAttribute{
				Optional: true,
				Description: "Filter by user name.\n" +
					"  - example : 'john.doe'",
			},
			"user_email": schema.StringAttribute{
				Optional: true,
				Description: "Filter by user email.\n" +
					"  - example : 'user@example.com'",
			},
			"creator_name": schema.StringAttribute{
				Optional: true,
				Description: "Name of the user who created this group member.\n" +
					"  - example : 'John Doe'",
				MarkdownDescription: "Name of the user who created this group member.\n  - example : 'John Doe'",
			},
			"creator_email": schema.StringAttribute{
				Optional: true,
				Description: "Email address of the user who created this group member.\n" +
					"  - example : 'user@example.com'",
				MarkdownDescription: "Email address of the user who created this group member.\n  - example : 'user@example.com'",
			},
			"group_members": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group Members",
				MarkdownDescription: "Group Members",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the group member was added.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the group member was added.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who added the group member.\n" +
								"  - example : 'user@example.com'",
							MarkdownDescription: "User who added the group member.\n" +
								"  - example : 'user@example.com'",
						},
						"creator_created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the creator was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the creator was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"creator_email": schema.StringAttribute{
							Computed: true,
							Description: "Email of the user who created this group member.\n" +
								"  - example : 'user@example.com'",
							MarkdownDescription: "Email of the user who created this group member.\n" +
								"  - example : 'user@example.com'",
						},
						"creator_last_login_at": schema.StringAttribute{
							Optional: true,
							Description: "Timestamp when the creator last logged in.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the creator last logged in.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"creator_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the user who created this group member.\n" +
								"  - example : 'John Doe'",
							MarkdownDescription: "Name of the user who created this group member.\n" +
								"  - example : 'John Doe'",
						},
						"groups": schema.ListNestedAttribute{
							Computed:            true,
							Description:         "Groups",
							MarkdownDescription: "Groups",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
										Description: "Group ID.\n" +
											"  - example : 'grp-1234567890abcdef'",
									},
									"name": schema.StringAttribute{
										Computed: true,
										Description: "Group Name.\n" +
											"  - example : 'MyGroup'",
									},
								},
							},
						},
						"user_created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the user account was created.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"user_email": schema.StringAttribute{
							Computed: true,
							Description: "Email address of the user.\n" +
								"  - example : 'user@example.com'",
						},
						"user_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier for the user.\n" +
								"  - example : 'usr-1234567890abcdef'",
						},
						"user_last_login_at": schema.StringAttribute{
							Optional: true,
							Description: "Timestamp when the user last logged in.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
							MarkdownDescription: "Timestamp when the user last logged in.\n" +
								"  - example : '2024-01-01T00:00:00Z'",
						},
						"user_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the user.\n" +
								"  - example : 'Jane Doe'",
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

	state.GroupMembers = getGroupMembersV1Dot4(data.GroupMembers)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
