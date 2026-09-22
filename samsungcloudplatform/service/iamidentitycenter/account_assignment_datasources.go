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
	_ datasource.DataSource              = &iamIdentityCenterAccountAssignmentsDataSources{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterAccountAssignmentsDataSources{}
)

func NewIamIdentityCenterAccountAssignmentsDataSources() datasource.DataSource {
	return &iamIdentityCenterAccountAssignmentsDataSources{}
}

type iamIdentityCenterAccountAssignmentsDataSources struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterAccountAssignmentsDataSources) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_account_assignments"
}

func (r *iamIdentityCenterAccountAssignmentsDataSources) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists IAM Identity Center Account Assignments.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"target_account_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Filter by target account ID.\n  - example: 3265ab469f0d406d83073da3e11e7a6c",
				MarkdownDescription: "Filter by target account ID.\n  - example: 3265ab469f0d406d83073da3e11e7a6c",
			},
			"sort": schema.StringAttribute{
				Optional:            true,
				Description:         "Sort the results by a specific field.\n  - example: target_account_id:asc",
				MarkdownDescription: "Sort the results by a specific field.\n  - example: target_account_id:asc",
			},
			"target_account_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by target account name.\n  - example: Admin Account",
				MarkdownDescription: "Filter by target account name.\n  - example: Admin Account",
			},
			"target_account_email": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by target account email.\n  - example: target_account@example.com",
				MarkdownDescription: "Filter by target account email.\n  - example: target_account@example.com",
			},
			"permission_set_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by permission set ID.\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
				MarkdownDescription: "Filter by permission set ID.\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
			},
			"principal_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by principal name.\n  - example: GROUP0407",
				MarkdownDescription: "Filter by principal name.\n  - example: GROUP0407",
			},
			"role_srn": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by role SRN.\n  - example: srn:...:role/Admin",
				MarkdownDescription: "Filter by role SRN.\n  - example: srn:...:role/Admin",
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
			"account_assignments": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of account assignments matching the filter criteria.\n" +
					"  - example : [{\"principal_id\": \"user-id\", \"target_account_id\": \"...\"}]",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Instance ID\n  - example: ssoins-12345",
							MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
						},
						"target_account_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Target Account ID\n  - example: 3265ab469f0d406d83073da3e11e7a6c",
							MarkdownDescription: "Target Account ID\n  - example: 3265ab469f0d406d83073da3e11e7a6c",
						},
						"target_account_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Target Account Name\n  - example: Admin Account",
							MarkdownDescription: "Target Account Name\n  - example: Admin Account",
						},
						"target_account_email": schema.StringAttribute{
							Computed:            true,
							Description:         "Target Account Email\n  - example: target_account@example.com",
							MarkdownDescription: "Target Account Email\n  - example: target_account@example.com",
						},
						"principal_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Principal ID (User or Group ID)\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
							MarkdownDescription: "Principal ID (User or Group ID)\n  - example: 138c2fc8c29a449dbfa8681f8f1d78e2",
						},
						"principal_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Principal Name\n  - example: GROUP0407",
							MarkdownDescription: "Principal Name\n  - example: GROUP0407",
						},
						"principal_type": schema.StringAttribute{
							Computed:            true,
							Description:         "Principal Type\n  - example: USER\n  - enum: [\"USER\", \"GROUP\"]",
							MarkdownDescription: "Principal Type\n  - example: USER\n  - enum: [\"USER\", \"GROUP\"]",
						},
						"role_srn": schema.StringAttribute{
							Computed:            true,
							Description:         "Role Srn\n  - example: srn:...:role/Admin",
							MarkdownDescription: "Role Srn\n  - example: srn:...:role/Admin",
						},
						"permission_sets": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"account_assignment_id": schema.StringAttribute{
										Computed: true,
									},
									"permission_set_id": schema.StringAttribute{
										Computed: true,
									},
									"permission_set_name": schema.StringAttribute{
										Computed: true,
									},
									"permission_set_srn": schema.StringAttribute{
										Computed: true,
									},
								},
							},
							Description:         "Permission Sets",
							MarkdownDescription: "Permission Sets",
						},
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterAccountAssignmentsDataSources) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type accountAssignmentsDataSourcesModel struct {
	InstanceId         types.String                       `tfsdk:"instance_id"`
	TargetAccountId    types.String                       `tfsdk:"target_account_id"`
	Size               types.Int32                        `tfsdk:"size"`
	Page               types.Int32                        `tfsdk:"page"`
	Sort               types.String                       `tfsdk:"sort"`
	TargetAccountName  types.String                       `tfsdk:"target_account_name"`
	TargetAccountEmail types.String                       `tfsdk:"target_account_email"`
	PermissionSetId    types.String                       `tfsdk:"permission_set_id"`
	PrincipalName      types.String                       `tfsdk:"principal_name"`
	RoleSrn            types.String                       `tfsdk:"role_srn"`
	AccountAssignments []accountAssignmentDataSourceModel `tfsdk:"account_assignments"`
}

func (r *iamIdentityCenterAccountAssignmentsDataSources) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountAssignmentsDataSourcesModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	targetAccountId := state.TargetAccountId.ValueString()
	size := state.Size.ValueInt32()
	page := state.Page.ValueInt32()
	sort := state.Sort.ValueString()
	targetAccountName := state.TargetAccountName.ValueString()
	targetAccountEmail := state.TargetAccountEmail.ValueString()
	permissionSetId := state.PermissionSetId.ValueString()
	principalName := state.PrincipalName.ValueString()
	roleSrn := state.RoleSrn.ValueString()

	result, err := r.clients.IamIdentityCenter.ListAccountAssignments(ctx, instanceId, targetAccountId, size, page, sort, targetAccountName, targetAccountEmail, permissionSetId, principalName, roleSrn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Account Assignments",
			err.Error(),
		)
		return
	}

	if result != nil && result.AccountAssignments != nil {
		state.AccountAssignments = make([]accountAssignmentDataSourceModel, len(result.AccountAssignments))
		for i, assignment := range result.AccountAssignments {
			state.AccountAssignments[i] = accountAssignmentDataSourceModel{
				InstanceId:         types.StringValue(instanceId),
				TargetAccountId:    types.StringValue(assignment.TargetAccountId),
				TargetAccountName:  types.StringValue(assignment.TargetAccountName),
				TargetAccountEmail: types.StringValue(assignment.TargetAccountEmail),
				PrincipalId:        types.StringValue(assignment.PrincipalId),
				PrincipalName:      types.StringValue(assignment.PrincipalName),
				PrincipalType:      types.StringValue(string(assignment.PrincipalType)),
				RoleSrn:            types.StringPointerValue(assignment.RoleSrn.Get()),
			}

			if len(assignment.PermissionSets) > 0 {
				state.AccountAssignments[i].PermissionSets = make([]accountAssignmentPermissionSetModel, len(assignment.PermissionSets))
				for j, ps := range assignment.PermissionSets {
					state.AccountAssignments[i].PermissionSets[j] = accountAssignmentPermissionSetModel{
						AccountAssignmentId: types.StringValue(ps.AccountAssignmentId),
						PermissionSetId:     types.StringValue(ps.PermissionSetId),
						PermissionSetName:   types.StringValue(ps.PermissionSetName),
						PermissionSetSrn:    types.StringValue(ps.PermissionSetSrn),
					}
				}
			}
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
