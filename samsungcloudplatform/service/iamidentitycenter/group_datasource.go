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
	_ datasource.DataSource              = &iamIdentityCenterGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterGroupDataSource{}
)

func NewIamIdentityCenterGroupDataSource() datasource.DataSource {
	return &iamIdentityCenterGroupDataSource{}
}

type iamIdentityCenterGroupDataSource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_group"
}

func (r *iamIdentityCenterGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves details of an IAM Identity Center Group.",
		MarkdownDescription: "Retrieves details of an IAM Identity Center Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Group ID\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Group Name\n  - example: Admin_Group",
				MarkdownDescription: "Group Name\n  - example: Admin_Group",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"description": schema.StringAttribute{
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

func (r *iamIdentityCenterGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type groupDataSourceModel struct {
	InstanceId  types.String                `tfsdk:"instance_id"`
	Name        types.String                `tfsdk:"name"`
	Id          types.String                `tfsdk:"id"`
	Description types.String                `tfsdk:"description"`
	Group       *groupDataSourceDetailModel `tfsdk:"group"`
}

type groupDataSourceDetailModel struct {
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

func (r *iamIdentityCenterGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()
	instanceId := state.InstanceId.ValueString()
	result, httpResp, err := r.clients.IamIdentityCenter.GetGroup(ctx, groupId, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.Diagnostics.AddError(
				"Not Found",
				"Group not found",
			)
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
		state.Name = types.StringValue(result.Group.Name)
		state.Description = types.StringPointerValue(result.Group.Description.Get())

		state.Group = &groupDataSourceDetailModel{}
		state.Group.Id = types.StringValue(result.Group.Id)
		state.Group.Name = types.StringValue(result.Group.Name)
		state.Group.Description = types.StringPointerValue(result.Group.Description.Get())
		state.Group.InstanceId = types.StringValue(result.Group.InstanceId)
		state.Group.CreatedAt = types.StringValue(result.Group.CreatedAt.Format("2006-01-02T15:04:05Z"))
		state.Group.CreatedBy = types.StringValue(result.Group.CreatedBy)
		state.Group.CreatorName = types.StringValue(result.Group.CreatorName)
		state.Group.ModifiedAt = types.StringValue(result.Group.ModifiedAt.Format("2006-01-02T15:04:05Z"))
		state.Group.ModifiedBy = types.StringValue(result.Group.ModifiedBy)
		state.Group.ModifierName = types.StringValue(result.Group.ModifierName)
		if result.Group.UserCount.Get() != nil {
			state.Group.UserCount = types.Int64Value(int64(*result.Group.UserCount.Get()))
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
