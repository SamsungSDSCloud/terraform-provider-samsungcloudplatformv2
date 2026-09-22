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
	_ datasource.DataSource              = &iamIdentityCenterPermissionSetPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &iamIdentityCenterPermissionSetPoliciesDataSource{}
)

func NewIamIdentityCenterPermissionSetPoliciesDataSource() datasource.DataSource {
	return &iamIdentityCenterPermissionSetPoliciesDataSource{}
}

type iamIdentityCenterPermissionSetPoliciesDataSource struct {
	clients *client.SCPClient
}

func (r *iamIdentityCenterPermissionSetPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_identity_center_permission_set_policies"
}

func (r *iamIdentityCenterPermissionSetPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the policies attached to an IAM Identity Center Permission Set (Data Source).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Permission Set Policy reference ID `<permission_set_id>:<instance_id>`.",
				MarkdownDescription: "Permission Set Policy reference ID `<permission_set_id>:<instance_id>`.",
			},
			"permission_set_id": schema.StringAttribute{
				Required:            true,
				Description:         "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
				MarkdownDescription: "Permission Set ID\n  - example: e8d4b9f2c1a7e5d3f0b6c9a2d8e4f7b1",
			},
			"instance_id": schema.StringAttribute{
				Required:            true,
				Description:         "Instance ID\n  - example: ssoins-12345",
				MarkdownDescription: "Instance ID\n  - example: ssoins-12345",
			},
			"inline_policy": schema.StringAttribute{
				Computed:            true,
				Description:         "Inline policy document as a JSON string. Empty if no inline policy is attached.",
				MarkdownDescription: "Inline policy document as a JSON string. Empty if no inline policy is attached.",
			},
			"managed_policies": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of AWS managed policies attached to the permission set.",
				MarkdownDescription: "List of AWS managed policies attached to the permission set.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Managed policy ID.",
							MarkdownDescription: "Managed policy ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Managed policy name.",
							MarkdownDescription: "Managed policy name.",
						},
					},
				},
			},
			"custom_policies": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of customer managed policies attached to the permission set.",
				MarkdownDescription: "List of customer managed policies attached to the permission set.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							Description:         "Customer managed policy ID.",
							MarkdownDescription: "Customer managed policy ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Customer managed policy name.",
							MarkdownDescription: "Customer managed policy name.",
						},
					},
				},
			},
		},
	}
}

func (r *iamIdentityCenterPermissionSetPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

type permissionSetPolicyReference struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type permissionSetPoliciesDataSourceModel struct {
	PermissionSetId types.String                   `tfsdk:"permission_set_id"`
	InstanceId      types.String                   `tfsdk:"instance_id"`
	Id              types.String                   `tfsdk:"id"`
	InlinePolicy    types.String                   `tfsdk:"inline_policy"`
	ManagedPolicies []permissionSetPolicyReference `tfsdk:"managed_policies"`
	CustomPolicies  []permissionSetPolicyReference `tfsdk:"custom_policies"`
}

func (r *iamIdentityCenterPermissionSetPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state permissionSetPoliciesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permissionSetId := state.PermissionSetId.ValueString()
	instanceId := state.InstanceId.ValueString()

	result, err := r.clients.IamIdentityCenter.GetPolicies(ctx, permissionSetId, instanceId, "", "", 0, 0)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read IAM Identity Center Permission Set Policies",
			err.Error(),
		)
		return
	}

	var customPolicies []permissionSetPolicyReference
	var managedPolicies []permissionSetPolicyReference
	var inlinePolicy string

	if result != nil {
		for _, policy := range result.Policies {
			collectPolicyReference(policy, &customPolicies, &managedPolicies, &inlinePolicy)
		}
	}

	state.CustomPolicies = customPolicies
	state.ManagedPolicies = managedPolicies
	state.InlinePolicy = types.StringValue(inlinePolicy)
	state.Id = types.StringValue(permissionSetId + ":" + instanceId)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func collectPolicyReference(policy sdk.PolicyV1Dot2, customPolicies *[]permissionSetPolicyReference, managedPolicies *[]permissionSetPolicyReference, inlinePolicy *string) {
	switch policy.Category {
	case "CUSTOM_POLICY":
		if policy.Id != "" {
			*customPolicies = append(*customPolicies, referenceOf(policy))
		}
	case "MANAGED_POLICY":
		if policy.Id != "" {
			*managedPolicies = append(*managedPolicies, referenceOf(policy))
		}
	case "INLINE_POLICY":
		if *inlinePolicy == "" {
			*inlinePolicy = inlinePolicyDocument(policy.Contents)
		}
	}
}

func referenceOf(policy sdk.PolicyV1Dot2) permissionSetPolicyReference {
	return permissionSetPolicyReference{
		Id:   types.StringValue(policy.Id),
		Name: nullableString(policy.Name),
	}
}

func nullableString(value sdk.NullableString) types.String {
	if !value.IsSet() {
		return types.StringNull()
	}
	actual := value.Get()
	if actual == nil {
		return types.StringNull()
	}
	return types.StringValue(*actual)
}

func inlinePolicyDocument(contents sdk.NullableContents) string {
	// Reuse the inline/permisssion set resources' parser so the data source
	// reads every backend shape identically to the resources.
	policy := sdk.PolicyV1Dot2{Contents: contents}
	doc, ok := existingInlineContent(policy)
	if !ok {
		return ""
	}
	return doc
}
