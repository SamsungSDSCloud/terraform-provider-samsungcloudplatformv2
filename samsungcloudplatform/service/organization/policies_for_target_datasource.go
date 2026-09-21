package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/organization/1.3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ datasource.DataSource              = &policiesForTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &policiesForTargetDataSource{}
)

func NewPoliciesForTargetDataSource() datasource.DataSource {
	return &policiesForTargetDataSource{}
}

type policiesForTargetDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *policiesForTargetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_policies_for_target"
}

func (d *policiesForTargetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = inst.Client.Organization
	d.clients = inst.Client
}

func (d *policiesForTargetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List the control policies assigned to a given target within the organization.",
		Attributes: map[string]schema.Attribute{
			"target_id": schema.StringAttribute{
				Description: "Unique identifier of the target (account, organization unit, or root) whose assigned policies are listed. \n" +
					"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
				Required: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization that owns the target. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
			},
			"policy_category": schema.StringAttribute{
				Description: "Filter by policy category. \n" +
					"  - example : 'SCP' \n" +
					"  - allowed_values : ['SCP', 'RCP', 'TAG', 'DGP'] \n",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("SCP", "RCP", "TAG", "DGP"),
				},
			},
			"name": schema.StringAttribute{
				Description: "Filter by control policy name. \n" +
					"  - example : 'test-policy-name' \n",
				Optional: true,
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page size. \n" +
					"  - example : 20 \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page number. \n" +
					"  - example : 0 \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria. \n" +
					"  - example : 'created_at:desc' \n",
			},
			"id": schema.StringAttribute{
				Computed: true,
				Description: "Data source identifier. \n" +
					"  - example : '0a36e0746dbf4908acf0357829701381' \n",
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total number of policies assigned to the target. \n" +
					"  - example : 3 \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria from response. \n" +
					"  - example : ['created_at:desc'] \n",
			},
			"policies": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of policies assigned to the target. \n" +
					"  - example : '[{id: 0a36e0746dbf4908acf0357829701381, policy_name: test-policy-name, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Control policy ID. \n" +
								"  - example : '0a36e0746dbf4908acf0357829701381' \n",
						},
						"policy_name": schema.StringAttribute{
							Computed: true,
							Description: "Control policy name. \n" +
								"  - example : 'test-policy-name' \n",
						},
						"policy_type": schema.StringAttribute{
							Computed: true,
							Description: "Control policy type. \n" +
								"  - example : 'SYSTEM_MANAGED' \n" +
								"  - allowed_values : ['SYSTEM_MANAGED', 'USER_DEFINED'] \n",
						},
						"source": schema.StringAttribute{
							Computed: true,
							Description: "Policy creation subject. \n" +
								"  - example : 'ORGANIZATION' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the policy binding was created. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the policy binding. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the policy binding was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the policy binding. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"link_types": schema.MapAttribute{
							Computed: true,
							Description: "Link type map describing the targets the policy is linked to. " +
								"Each key (e.g. 'DIRECTED', 'INHERITED') maps to a list of target objects. \n" +
								"  - target_id : unique identifier of the target \n" +
								"    - example : '1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n" +
								"  - target_name : name of the target \n" +
								"    - example : 'example-target' \n" +
								"  - target_type : type of the target \n" +
								"    - allowed values : ['ROOT', 'ACCOUNT', 'OU'] \n" +
								"    - example : 'OU' \n" +
								"  - example : '{DIRECTED: [{target_id: 1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, target_name: example-target, target_type: OU}], INHERITED: [{target_id: 2a2b3c4d5e6f7g8h9i0j1k2l3m4n5, target_name: root, target_type: ROOT}]}' \n",
							ElementType: types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"target_id":   basetypes.StringType{},
										"target_name": basetypes.StringType{},
										"target_type": basetypes.StringType{},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *policiesForTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.PoliciesForTargetDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	targetId := state.TargetId.ValueString()
	organizationId := state.OrganizationId.ValueString()
	policyCategory := state.PolicyCategory.ValueString()
	name := state.Name.ValueString()

	size, page := int32(0), int32(0)
	if !state.Size.IsNull() {
		size = int32(state.Size.ValueInt64())
	}
	if !state.Page.IsNull() {
		page = int32(state.Page.ValueInt64())
	}
	sort := ""
	if !state.Sort.IsNull() {
		sort = state.Sort.ValueString()
	}

	data, err := d.client.GetPoliciesForTarget(ctx, targetId, organizationId, policyCategory, name, size, page, sort)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Policies for Target",
			"Could not list policies for target, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policies := make([]attr.Value, 0, len(data.Policies))
	for i := range data.Policies {
		policyValue, policyDiags := d.buildPoliciesForTargetValue(ctx, &data.Policies[i])
		resp.Diagnostics.Append(policyDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		policies = append(policies, policyValue)
	}

	policiesList, policiesListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.PoliciesForTargetValue{}.AttributeTypes(ctx)},
		policies,
	)
	resp.Diagnostics.Append(policiesListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sortList, sortListDiags := types.ListValueFrom(ctx, types.StringType, data.Sort)
	resp.Diagnostics.Append(sortListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(targetId)
	state.Policies = policiesList
	state.TotalCount = types.Int64Value(int64(data.Count))
	state.Page = types.Int64Value(int64(data.Page))
	state.Size = types.Int64Value(int64(data.Size))
	state.SortResult = sortList

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *policiesForTargetDataSource) buildPoliciesForTargetValue(ctx context.Context, policy *sdkorganization.PoliciesForTargetSummary) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	linkTypes := make(map[string]attr.Value)
	for key, infos := range policy.LinkTypes {
		infosList := make([]attr.Value, 0, len(infos))
		for _, info := range infos {
			infoValue, infoDiags := types.ObjectValue(organization.PolicyLinkInfoValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"target_id":   types.StringValue(info.TargetId),
				"target_name": types.StringValue(nullableStringValue(info.TargetName)),
				"target_type": types.StringValue(info.TargetType),
			})
			diags.Append(infoDiags...)
			if diags.HasError() {
				return types.ObjectNull(organization.PoliciesForTargetValue{}.AttributeTypes(ctx)), diags
			}
			infosList = append(infosList, infoValue)
		}
		infosListValue, listDiags := types.ListValue(
			types.ObjectType{AttrTypes: organization.PolicyLinkInfoValue{}.AttributeTypes(ctx)},
			infosList,
		)
		diags.Append(listDiags...)
		if diags.HasError() {
			return types.ObjectNull(organization.PoliciesForTargetValue{}.AttributeTypes(ctx)), diags
		}
		linkTypes[key] = infosListValue
	}

	linkTypesValue, mapDiags := types.MapValue(
		types.ListType{ElemType: types.ObjectType{AttrTypes: organization.PolicyLinkInfoValue{}.AttributeTypes(ctx)}},
		linkTypes,
	)
	diags.Append(mapDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.PoliciesForTargetValue{}.AttributeTypes(ctx)), diags
	}

	policyValue, valueDiags := types.ObjectValue(organization.PoliciesForTargetValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"id":          types.StringValue(policy.Id),
		"policy_name": types.StringValue(policy.PolicyName),
		"policy_type": types.StringValue(string(policy.PolicyType)),
		"source":      types.StringValue(policy.Source),
		"created_at":  types.StringValue(policy.CreatedAt.Format(TimeFormat)),
		"created_by":  types.StringValue(policy.CreatedBy),
		"modified_at": types.StringValue(policy.ModifiedAt.Format(TimeFormat)),
		"modified_by": types.StringValue(policy.ModifiedBy),
		"link_types":  linkTypesValue,
	})
	diags.Append(valueDiags...)

	return policyValue, diags
}
