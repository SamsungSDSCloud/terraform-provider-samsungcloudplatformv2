package iamidentitycenter

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/iam-identity-center/1.6"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &iamIdentityCenterAccountAssignmentDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterAccountAssignmentDataSource{}
)

func NewIamIdentityCenterAccountAssignmentDataSource() datasource.DataSource {
	return &iamIdentityCenterAccountAssignmentDataSource{}
}

type iamIdentityCenterAccountAssignmentDataSource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterAccountAssignmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_account_assignment"
}

func (r *iamIdentityCenterAccountAssignmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves details of an IAM Identity Center Account Assignment.",
		MarkdownDescription: "Retrieves details of an IAM Identity Center Account Assignment.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"target_account_id": schema.StringAttribute{
				Required:            true,
				Description:         "Target Account ID\n  - example: 3265ab469f0d406d83073da3e11e7a6c",
				MarkdownDescription: "Target Account ID\n  - example: 3265ab469f0d406d83073da3e11e7a6c",
			},
			"principal_id": schema.StringAttribute{
				Optional:            true,
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
	}
}

func (r *iamIdentityCenterAccountAssignmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type accountAssignmentDataSourceModel struct {
	InstanceId         types.String                          `tfsdk:"instance_id"`
	TargetAccountId    types.String                          `tfsdk:"target_account_id"`
	PrincipalId        types.String                          `tfsdk:"principal_id"`
	PrincipalName      types.String                          `tfsdk:"principal_name"`
	PrincipalType      types.String                          `tfsdk:"principal_type"`
	TargetAccountName  types.String                          `tfsdk:"target_account_name"`
	TargetAccountEmail types.String                          `tfsdk:"target_account_email"`
	RoleSrn            types.String                          `tfsdk:"role_srn"`
	PermissionSets     []accountAssignmentPermissionSetModel `tfsdk:"permission_sets"`
}

type accountAssignmentPermissionSetModel struct {
	AccountAssignmentId types.String `tfsdk:"account_assignment_id"`
	PermissionSetId     types.String `tfsdk:"permission_set_id"`
	PermissionSetName   types.String `tfsdk:"permission_set_name"`
	PermissionSetSrn    types.String `tfsdk:"permission_set_srn"`
}

func (r *iamIdentityCenterAccountAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountAssignmentDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	targetAccountId := state.TargetAccountId.ValueString()
	principalId := state.PrincipalId.ValueString()

	result, err := r.clients.IamIdentityCenter.GetAccountAssignment(ctx, instanceId, targetAccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Account Assignment",
			err.Error(),
		)
		return
	}

	if result == nil || result.AccountAssignments == nil || len(result.AccountAssignments) == 0 {
		resp.Diagnostics.AddError(
			"Not Found",
			"Account Assignment not found",
		)
		return
	}

	// Find the matching account assignment
	var foundAssignment *sdk.AccountAssignment
	if principalId != "" {
		// If principal_id is provided, find the specific assignment
		for i := range result.AccountAssignments {
			if result.AccountAssignments[i].PrincipalId == principalId {
				foundAssignment = &result.AccountAssignments[i]
				break
			}
		}
	} else {
		// Otherwise, use the first one
		foundAssignment = &result.AccountAssignments[0]
	}

	if foundAssignment == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			"Account Assignment not found",
		)
		return
	}

	// Update state with found assignment
	state.PrincipalId = types.StringValue(foundAssignment.PrincipalId)
	state.PrincipalName = types.StringValue(foundAssignment.PrincipalName)
	state.PrincipalType = types.StringValue(string(foundAssignment.PrincipalType))
	state.TargetAccountName = types.StringValue(foundAssignment.TargetAccountName)
	state.TargetAccountEmail = types.StringValue(foundAssignment.TargetAccountEmail)
	state.RoleSrn = types.StringPointerValue(foundAssignment.RoleSrn.Get())

	// Handle permission_sets
	if len(foundAssignment.PermissionSets) > 0 {
		state.PermissionSets = make([]accountAssignmentPermissionSetModel, len(foundAssignment.PermissionSets))
		for i, ps := range foundAssignment.PermissionSets {
			state.PermissionSets[i] = accountAssignmentPermissionSetModel{
				AccountAssignmentId: types.StringValue(ps.AccountAssignmentId),
				PermissionSetId:     types.StringValue(ps.PermissionSetId),
				PermissionSetName:   types.StringValue(ps.PermissionSetName),
				PermissionSetSrn:    types.StringValue(ps.PermissionSetSrn),
			}
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
