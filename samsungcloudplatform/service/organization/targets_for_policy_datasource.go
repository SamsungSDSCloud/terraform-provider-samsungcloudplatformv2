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
)

var (
	_ datasource.DataSource              = &targetsForPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &targetsForPolicyDataSource{}
)

func NewTargetsForPolicyDataSource() datasource.DataSource {
	return &targetsForPolicyDataSource{}
}

type targetsForPolicyDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *targetsForPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_targets_for_policy"
}

func (d *targetsForPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *targetsForPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List the targets (accounts, organization units, or root) to which a given control policy is assigned.",
		Attributes: map[string]schema.Attribute{
			"policy_id": schema.StringAttribute{
				Description: "Unique identifier of the control policy whose assigned targets are listed. \n" +
					"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
				Required: true,
			},
			"target_type": schema.StringAttribute{
				Description: "Type of the target to list. \n" +
					"  - example : 'OU' \n" +
					"  - allowed_values : ['ROOT', 'ACCOUNT', 'OU'] \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ROOT", "ACCOUNT", "OU"),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization that owns the policy. \n" +
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
				Description: "Filter by target name. \n" +
					"  - example : 'example-target' \n",
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
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total number of targets assigned to the policy. \n" +
					"  - example : 3 \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort criteria from response. \n" +
					"  - example : ['created_at:desc'] \n",
			},
			"targets": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of targets assigned to the policy. \n" +
					"  - example : '[{id: 0a36e0746dbf4908acf0357829701381, target_name: score-account, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Description: "Target ID the policy is assigned to. \n" +
								"  - example : '0a36e0746dbf4908acf0357829701381' \n",
						},
						"target_name": schema.StringAttribute{
							Computed: true,
							Description: "Target name the policy is assigned to. \n" +
								"  - example : 'example-target' \n",
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the binding was created. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the binding. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the binding was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the binding. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"control_policies": schema.ListNestedAttribute{
							Computed: true,
							Description: "List of control policies connected to the target. \n" +
								"  - example : '[{policy_id: f98e76d54c32b10a9z8y7x6w5v4u3, policy_name: test-policy-name}]' \n",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"policy_id": schema.StringAttribute{
										Computed: true,
										Description: "Control policy ID. \n" +
											"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
									},
									"policy_name": schema.StringAttribute{
										Computed: true,
										Description: "Control policy name. \n" +
											"  - example : 'test-policy-name' \n",
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

func (d *targetsForPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.TargetsForPolicyDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyId.ValueString()
	targetType := state.TargetType.ValueString()
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

	data, err := d.client.ListTargetsForPolicy(ctx, policyId, targetType, organizationId, policyCategory, name, size, page, sort)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Unable to Read Targets for Policy",
			"Could not list targets for policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	targets := make([]attr.Value, 0, len(data.Targets))
	for i := range data.Targets {
		targetValue, targetDiags := d.buildTargetsForPolicyValue(ctx, &data.Targets[i])
		resp.Diagnostics.Append(targetDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		targets = append(targets, targetValue)
	}

	targetsList, targetsListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.TargetsForPolicyValue{}.AttributeTypes(ctx)},
		targets,
	)
	resp.Diagnostics.Append(targetsListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sortList, sortListDiags := types.ListValueFrom(ctx, types.StringType, data.Sort)
	resp.Diagnostics.Append(sortListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Targets = targetsList
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

func (d *targetsForPolicyDataSource) buildTargetsForPolicyValue(ctx context.Context, target *sdkorganization.TargetsForPolicySummary) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	controlPolicies := make([]attr.Value, 0, len(target.ControlPolicies))
	for i := range target.ControlPolicies {
		policyValue, policyDiags := types.ObjectValue(organization.ConnectedPolicyValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"policy_id":   types.StringValue(target.ControlPolicies[i].PolicyId),
			"policy_name": types.StringValue(target.ControlPolicies[i].PolicyName),
		})
		diags.Append(policyDiags...)
		if diags.HasError() {
			return types.ObjectNull(organization.TargetsForPolicyValue{}.AttributeTypes(ctx)), diags
		}
		controlPolicies = append(controlPolicies, policyValue)
	}

	controlPoliciesList, listDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.ConnectedPolicyValue{}.AttributeTypes(ctx)},
		controlPolicies,
	)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.TargetsForPolicyValue{}.AttributeTypes(ctx)), diags
	}

	targetValue, valueDiags := types.ObjectValue(organization.TargetsForPolicyValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"id":               types.StringValue(target.Id),
		"target_name":      types.StringValue(nullableStringValue(target.TargetName)),
		"created_at":       types.StringValue(target.CreatedAt.Format(TimeFormat)),
		"created_by":       types.StringValue(target.CreatedBy),
		"modified_at":      types.StringValue(target.ModifiedAt.Format(TimeFormat)),
		"modified_by":      types.StringValue(target.ModifiedBy),
		"control_policies": controlPoliciesList,
	})
	diags.Append(valueDiags...)

	return targetValue, diags
}
