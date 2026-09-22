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
	_ datasource.DataSource              = &iamIdentityCenterPermissionSetsDataSources{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterPermissionSetsDataSources{}
)

func NewIamIdentityCenterPermissionSetsDataSources() datasource.DataSource {
	return &iamIdentityCenterPermissionSetsDataSources{}
}

type iamIdentityCenterPermissionSetsDataSources struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetsDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_sets"
}

func (r *iamIdentityCenterPermissionSetsDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IAM Identity Center Permission Sets.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Filter by permission set name.\n  - example: Admin_Permission_Set",
				MarkdownDescription: "Filter by permission set name.\n  - example: Admin_Permission_Set",
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort the results by a specific field.\n  - example: name:asc",
				MarkdownDescription: "Sort the results by a specific field.\n  - example: name:asc",
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
			"permission_sets": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of permission sets matching the filter criteria.\n" +
					"  - example : [{\"id\": \"ps-id\", \"name\": \"Admin_Permission_Set\", ...}]",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
							MarkdownDescription: "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Permission Set Name\n  - example: Admin_Permission_Set",
							MarkdownDescription: "Permission Set Name\n  - example: Admin_Permission_Set",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Permission Set Description\n  - example: Permission set for administrators",
							MarkdownDescription: "Permission Set Description\n  - example: Permission set for administrators",
						},
						"session_duration": schema.Int64Attribute{
							Computed:            true,
							Description:         "Maximum Session Duration (seconds)\n  - example: 3600",
							MarkdownDescription: "Maximum Session Duration (seconds)\n  - example: 3600",
						},
						"custom_policies": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "Customer Managed Policy ARNs\n  - example: [arn:aws:iam::123456789012:policy/MyCustomPolicy]",
							MarkdownDescription: "Customer Managed Policy ARNs\n  - example: [arn:aws:iam::123456789012:policy/MyCustomPolicy]",
						},
						"managed_policy_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							Description:         "Managed Policy IDs\n  - example: [37f2e31ff86b415698d7e8eeafab445d]",
							MarkdownDescription: "Managed Policy IDs\n  - example: [37f2e31ff86b415698d7e8eeafab445d]",
						},
						"srn": schema.StringAttribute{
							Computed:            true,
							Description:         "Permission Set SRN\n  - example: srn:iam:::permission-set/123456789012/us-east-1/ExamplePermissionSet",
							MarkdownDescription: "Permission Set SRN\n  - example: srn:iam:::permission-set/123456789012/us-east-1/ExamplePermissionSet",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "Permission Set State\n  - example: ACTIVE",
							MarkdownDescription: "Permission Set State\n  - example: ACTIVE",
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
							Description:         "Creator Name\n  - example: admin",
							MarkdownDescription: "Creator Name\n  - example: admin",
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
							Description:         "Modifier Name\n  - example: admin",
							MarkdownDescription: "Modifier Name\n  - example: admin",
						},
						"instance_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Instance ID\n  - example: ssoins-12345",
							MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
						},
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterPermissionSetsDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type permissionSetsDataSourcesModel struct {
	InstanceId     types.String                       `tfsdk:"instance_id"`
	Name           types.String                       `tfsdk:"name"`
	Size           types.Int32                        `tfsdk:"size"`
	Page           types.Int32                        `tfsdk:"page"`
	Sort           types.String                       `tfsdk:"sort"`
	PermissionSets []permissionSetDataSourceItemModel `tfsdk:"permission_sets"`
}

type permissionSetDataSourceItemModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	SessionDuration  types.Int64  `tfsdk:"session_duration"`
	CustomPolicies   types.List   `tfsdk:"custom_policies"`
	ManagedPolicyIds types.List   `tfsdk:"managed_policy_ids"`
	Srn              types.String `tfsdk:"srn"`
	State            types.String `tfsdk:"state"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	CreatorName      types.String `tfsdk:"creator_name"`
	ModifiedAt       types.String `tfsdk:"modified_at"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	ModifierName     types.String `tfsdk:"modifier_name"`
	InstanceId       types.String `tfsdk:"instance_id"`
}

func (r *iamIdentityCenterPermissionSetsDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state permissionSetsDataSourcesModel
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

	result, err := r.clients.IamIdentityCenter.ListPermissionSets(ctx, instanceId, name, size, page, sort)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Sets",
			err.Error(),
		)
		return
	}

	if result != nil && result.PermissionSets != nil {
		state.PermissionSets = make([]permissionSetDataSourceItemModel, len(result.PermissionSets))
		for i, ps := range result.PermissionSets {
			state.PermissionSets[i] = permissionSetDataSourceItemModel{
				Id:              types.StringValue(ps.Id),
				Name:            types.StringValue(ps.Name),
				Description:     stringValuePtr(ps.Description.Get()),
				SessionDuration: types.Int64Value(int64(ps.SessionDuration)),
				Srn:             types.StringValue(ps.Srn),
				CreatedAt:       types.StringValue(ps.CreatedAt.Format("2006-01-02T15:04:05Z")),
				CreatedBy:       types.StringValue(ps.CreatedBy),
				CreatorName:     types.StringValue(ps.CreatorName),
				ModifiedAt:      types.StringValue(ps.ModifiedAt.Format("2006-01-02T15:04:05Z")),
				ModifiedBy:      types.StringValue(ps.ModifiedBy),
				ModifierName:    types.StringValue(ps.ModifierName),
				InstanceId:      types.StringValue(ps.InstanceId),
			}

			if ps.State != nil {
				state.PermissionSets[i].State = types.StringValue(string(*ps.State))
			}

			if len(ps.CustomPolicies) > 0 {
				state.PermissionSets[i].CustomPolicies, diags = types.ListValueFrom(ctx, types.StringType, ps.CustomPolicies)
				resp.Diagnostics.Append(diags...)
			} else {
				state.PermissionSets[i].CustomPolicies = types.ListNull(types.StringType)
			}

			if len(ps.ManagedPolicies) > 0 {
				var managedPolicyIds []string
				for _, mp := range ps.ManagedPolicies {
					mpMap := mp
					id, ok := mpMap["id"].(string)
					if !ok {
						continue
					}
					managedPolicyIds = append(managedPolicyIds, id)
				}
				state.PermissionSets[i].ManagedPolicyIds, diags = types.ListValueFrom(ctx, types.StringType, managedPolicyIds)
				resp.Diagnostics.Append(diags...)
			} else {
				state.PermissionSets[i].ManagedPolicyIds = types.ListNull(types.StringType)
			}
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
