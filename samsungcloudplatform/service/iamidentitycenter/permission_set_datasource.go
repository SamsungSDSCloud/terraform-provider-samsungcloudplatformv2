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
	_ datasource.DataSource              = &iamIdentityCenterPermissionSetDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterPermissionSetDataSource{}
)

func NewIamIdentityCenterPermissionSetDataSource() datasource.DataSource {
	return &iamIdentityCenterPermissionSetDataSource{}
}

type iamIdentityCenterPermissionSetDataSource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_set"
}

func (r *iamIdentityCenterPermissionSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves details of an IAM Identity Center Permission Set.",
		MarkdownDescription: "Retrieves details of an IAM Identity Center Permission Set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
				MarkdownDescription: "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
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
			"account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Account ID\n  - example: 123456789012",
				MarkdownDescription: "Account ID\n  - example: 123456789012",
			},
			"resource_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource Name\n  - example: ExamplePermissionSet",
				MarkdownDescription: "Resource Name\n  - example: ExamplePermissionSet",
			},
			"resource_type": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource Type\n  - example: PermissionSet",
				MarkdownDescription: "Resource Type\n  - example: PermissionSet",
			},
			"resource_type_display_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Resource Type Display Name\n  - example: Permission Set",
				MarkdownDescription: "Resource Type Display Name\n  - example: Permission Set",
			},
			"service": schema.StringAttribute{
				Computed:            true,
				Description:         "Service\n  - example: IAM Identity Center",
				MarkdownDescription: "Service\n  - example: IAM Identity Center",
			},
			"service_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Service Name\n  - example: scp-iam-identity-center",
				MarkdownDescription: "Service Name\n  - example: scp-iam-identity-center",
			},
		},
	}
}

type permissionSetDataSourceModel struct {
	InstanceId              types.String `tfsdk:"instance_id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	SessionDuration         types.Int64  `tfsdk:"session_duration"`
	CustomPolicies          types.List   `tfsdk:"custom_policies"`
	ManagedPolicyIds        types.List   `tfsdk:"managed_policy_ids"`
	Id                      types.String `tfsdk:"id"`
	Srn                     types.String `tfsdk:"srn"`
	State                   types.String `tfsdk:"state"`
	CreatedAt               types.String `tfsdk:"created_at"`
	CreatedBy               types.String `tfsdk:"created_by"`
	CreatorName             types.String `tfsdk:"creator_name"`
	ModifiedAt              types.String `tfsdk:"modified_at"`
	ModifiedBy              types.String `tfsdk:"modified_by"`
	ModifierName            types.String `tfsdk:"modifier_name"`
	AccountId               types.String `tfsdk:"account_id"`
	ResourceName            types.String `tfsdk:"resource_name"`
	ResourceType            types.String `tfsdk:"resource_type"`
	ResourceTypeDisplayName types.String `tfsdk:"resource_type_display_name"`
	Service                 types.String `tfsdk:"service"`
	ServiceName             types.String `tfsdk:"service_name"`
}

func (r *iamIdentityCenterPermissionSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *iamIdentityCenterPermissionSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state permissionSetDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	permissionSetId := state.Id.ValueString()

	if permissionSetId == "" {
		resp.Diagnostics.AddError(
			"Permission Set ID is required",
			"Please provide the permission set ID",
		)
		return
	}

	result, httpResp, err := r.clients.IamIdentityCenter.GetPermissionSet(ctx, permissionSetId, instanceId)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.Diagnostics.AddError(
				"Permission Set Not Found",
				fmt.Sprintf("Permission Set with ID %s not found", permissionSetId),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set",
			err.Error(),
		)
		return
	}

	if result != nil && result.Id != "" {
		state.Id = types.StringValue(result.Id)
		state.Name = types.StringValue(result.Name)
		state.Description = types.StringPointerValue(result.Description.Get())
		state.SessionDuration = types.Int64Value(int64(result.SessionDuration))
		state.Srn = types.StringValue(result.Srn)
		if result.State != nil {
			state.State = types.StringValue(string(*result.State))
		}
		state.CreatedAt = types.StringValue(result.CreatedAt.Format("2006-01-02T15:04:05Z"))
		state.CreatedBy = types.StringValue(result.CreatedBy)
		state.CreatorName = types.StringValue(result.CreatorName)
		state.ModifiedAt = types.StringValue(result.ModifiedAt.Format("2006-01-02T15:04:05Z"))
		state.ModifiedBy = types.StringValue(result.ModifiedBy)
		state.ModifierName = types.StringValue(result.ModifierName)
		state.AccountId = types.StringPointerValue(result.AccountId.Get())
		state.ResourceName = types.StringPointerValue(result.ResourceName.Get())
		state.ResourceType = types.StringPointerValue(result.ResourceType.Get())
		state.ResourceTypeDisplayName = types.StringPointerValue(result.ResourceTypeDisplayName.Get())
		state.Service = types.StringPointerValue(result.Service.Get())
		state.ServiceName = types.StringPointerValue(result.ServiceName.Get())

		if len(result.CustomPolicies) > 0 {
			state.CustomPolicies, diags = types.ListValueFrom(ctx, types.StringType, result.CustomPolicies)
			resp.Diagnostics.Append(diags...)
		}

		if len(result.ManagedPolicies) > 0 {
			var managedPolicyIds []string
			for _, mp := range result.ManagedPolicies {
				mpMap := mp
				id, ok := mpMap["id"].(string)
				if !ok {
					continue
				}
				managedPolicyIds = append(managedPolicyIds, id)
			}
			state.ManagedPolicyIds, diags = types.ListValueFrom(ctx, types.StringType, managedPolicyIds)
			resp.Diagnostics.Append(diags...)
		}

		if state.CustomPolicies.IsNull() || state.CustomPolicies.IsUnknown() {
			state.CustomPolicies = types.ListNull(types.StringType)
		}
		if state.ManagedPolicyIds.IsNull() || state.ManagedPolicyIds.IsUnknown() {
			state.ManagedPolicyIds = types.ListNull(types.StringType)
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
