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
	_ datasource.DataSource              = &iamIdentityCenterGroupsDataSources{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterGroupsDataSources{}
)

func NewIamIdentityCenterGroupsDataSources() datasource.DataSource {
	return &iamIdentityCenterGroupsDataSources{}
}

type iamIdentityCenterGroupsDataSources struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterGroupsDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_groups"
}

func (r *iamIdentityCenterGroupsDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IAM Identity Center Groups.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Filter groups by name.\n  - example: Admin_Group",
				MarkdownDescription: "Filter groups by name.\n  - example: Admin_Group",
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort the results by a specific field.\n  - example: name:asc",
				MarkdownDescription: "Sort the results by a specific field.\n  - example: name:asc",
			},
			"excluded_user_uuid": schema.StringAttribute{
				Optional:            true,
				Description:         "Exclude groups that contain the specified user UUID.\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
				MarkdownDescription: "Exclude groups that contain the specified user UUID.\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
			},
			"excluded_account_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Exclude groups that contain the specified account ID.\n  - example: 5b26e6b693a342928d1110a52e657fcb",
				MarkdownDescription: "Exclude groups that contain the specified account ID.\n  - example: 5b26e6b693a342928d1110a52e657fcb",
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
			"groups": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of groups matching the filter criteria.\n" +
					"  - example : [{\"id\": \"group-id\", \"name\": \"Admin_Group\", ...}]",
				NestedObject: schema.NestedAttributeObject{
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
		},
	}
}

func (r *iamIdentityCenterGroupsDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type groupsDataSourcesModel struct {
	InstanceId        types.String               `tfsdk:"instance_id"`
	Name              types.String               `tfsdk:"name"`
	Size              types.Int32                `tfsdk:"size"`
	Page              types.Int32                `tfsdk:"page"`
	Sort              types.String               `tfsdk:"sort"`
	ExcludedUserUuid  types.String               `tfsdk:"excluded_user_uuid"`
	ExcludedAccountId types.String               `tfsdk:"excluded_account_id"`
	Groups            []groupDataSourceItemModel `tfsdk:"groups"`
}

type groupDataSourceItemModel struct {
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

func (r *iamIdentityCenterGroupsDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupsDataSourcesModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	name := state.Name.ValueString()
	size := state.Size.ValueInt32()
	page := state.Page.ValueInt32()
	sort := state.Sort.ValueString()
	excludedUserUuid := state.ExcludedUserUuid.ValueString()
	excludedAccountId := state.ExcludedAccountId.ValueString()

	result, err := r.clients.IamIdentityCenter.ListGroups(ctx, instanceId, name, size, page, sort, excludedUserUuid, excludedAccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Groups",
			err.Error(),
		)
		return
	}

	if result != nil && result.Groups != nil {
		state.Groups = make([]groupDataSourceItemModel, len(result.Groups))
		for i, group := range result.Groups {
			state.Groups[i] = groupDataSourceItemModel{
				Id:           types.StringValue(group.Id),
				Name:         types.StringValue(group.Name),
				Description:  types.StringPointerValue(group.Description.Get()),
				InstanceId:   types.StringValue(group.InstanceId),
				CreatedAt:    types.StringValue(group.CreatedAt.Format("2006-01-02T15:04:05Z")),
				CreatedBy:    types.StringValue(group.CreatedBy),
				CreatorName:  types.StringValue(group.CreatorName),
				ModifiedAt:   types.StringValue(group.ModifiedAt.Format("2006-01-02T15:04:05Z")),
				ModifiedBy:   types.StringValue(group.ModifiedBy),
				ModifierName: types.StringValue(group.ModifierName),
			}
			if group.UserCount.Get() != nil {
				state.Groups[i].UserCount = types.Int64Value(int64(*group.UserCount.Get()))
			}
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
