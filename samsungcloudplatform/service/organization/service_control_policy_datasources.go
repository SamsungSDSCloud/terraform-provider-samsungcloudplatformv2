package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v4/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v4/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &serviceControlPolicyListDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceControlPolicyListDataSource{}
)

func NewServiceControlPoliciesDataSource() datasource.DataSource {
	return &serviceControlPolicyListDataSource{}
}

type serviceControlPolicyListDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *serviceControlPolicyListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_service_control_policies"
}

func (d *serviceControlPolicyListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serviceControlPolicyListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List Organization Service Control Policies",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
				Computed: false,
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: false,
				Description: "Service Control Policy Name. \n" +
					"  - example : 'MyPolicy' \n",
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: false,
				Description: "Service Control Policy Type. \n" +
					"  - example : 'MANAGED' \n",
			},
			"policies": schema.ListNestedAttribute{
				Computed: true,
				Description: "Service Control Policy List. \n" +
					"  - example : '[{policy_id: 138c2fc8c29a449dbfa8681f8f1d78e2, name: MyPolicy, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_id": schema.StringAttribute{
							Computed: true,
							Description: "Service Control Policy ID. \n" +
								"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization. \n" +
								"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
						},
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Service Control Policy Name. \n" +
								"  - example : 'MyPolicy' \n",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Policy Description. \n" +
								"  - example : 'This is an example policy.' \n",
						},
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Service Control Policy Type. \n" +
								"  - example : 'MANAGED' \n",
						},
						"category": schema.StringAttribute{
							Computed: true,
							Description: "Policy Category. \n" +
								"  - example : 'SCP' \n",
						},
						"state": schema.StringAttribute{
							Computed: true,
							Description: "service control policy state. \n" +
								"  - example : 'ACTIVE' \n",
						},
						"source": schema.StringAttribute{
							Computed: true,
							Description: "Policy Creation Subject. \n" +
								"  - example : 'ORGANIZATION' \n",
						},
						"document": schema.SingleNestedAttribute{
							Computed: true,
							Description: "Service Control Policy Document. \n" +
								"  - example : '{version: 2024-07-01, statement: [...]}' \n",
							Attributes: map[string]schema.Attribute{
								"statement": schema.ListNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.ListAttribute{
												ElementType: types.StringType,
												Computed:    true,
												Description: "List of actions permitted or denied by the policy statement. \n" +
													"  - example : ['s3:PutObject', 's3:GetObject'] \n",
											},
											"condition": schema.MapAttribute{
												ElementType: types.StringType,
												Computed:    true,
												Description: "Policy Condition. \n" +
													"  - example : '{StringEquals: {aws:RequestedRegion: us-east-1}}' \n",
											},
											"effect": schema.StringAttribute{
												Computed: true,
												Description: "Policy Effect. \n" +
													"  - example : 'Allow' \n" +
													"  - allowed_values : ['Allow', 'Deny'] \n",
											},
											"not_action": schema.ListAttribute{
												ElementType: types.StringType,
												Computed:    true,
												Description: "Policy Exclusion Action. \n" +
													"  - example : ['s3:DeleteObject'] \n",
											},
											"principal": schema.StringAttribute{
												Computed: true,
												Description: "Principal entity to which the policy statement applies. \n" +
													"  - example : '*' \n",
											},
											"resource": schema.ListAttribute{
												ElementType: types.StringType,
												Computed:    true,
												Description: "List of resources to which the policy statement applies. \n" +
													"  - example : ['*'] \n",
											},
											"sid": schema.StringAttribute{
												Computed: true,
												Description: "Syntax ID. \n" +
													"  - example : 'statement1' \n",
											},
										},
									},
									Computed: true,
									Description: "Policy Syntax. \n" +
										"  - example : '[{effect: Allow, action: [*], resource: [*], sid: statement1}]' \n",
								},
								"version": schema.StringAttribute{
									Computed: true,
									Description: "Policy Version. \n" +
										"  - example : '2012-10-17' \n",
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the service control policy unit created. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"created_by": schema.StringAttribute{
							Computed: true,
							Description: "User who created the service control policy. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"creator_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the service control policy creator. \n" +
								"  - example : 'John Doe na' \n",
						},
						"modified_at": schema.StringAttribute{
							Computed: true,
							Description: "Timestamp when the service control policy was modified. \n" +
								"  - example : '2025-01-01T00:00:00.000Z' \n",
						},
						"modified_by": schema.StringAttribute{
							Computed: true,
							Description: "User who modified the service control policy. \n" +
								"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
						},
						"modifier_name": schema.StringAttribute{
							Computed: true,
							Description: "Name of the service control policy modifier. \n" +
								"  - example : 'Alice' \n",
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
				Description: "Total service control policy count. \n" +
					"  - example : 10 \n",
			},
			"page": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page number. \n" +
					"  - example : 0 \n",
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Page size. \n" +
					"  - example : 20 \n",
			},
			"sort": schema.StringAttribute{
				Optional: true,
				Description: "Sort criteria. \n" +
					"  - example : 'created_at:desc' \n",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort Criteria from response. \n" +
					"  - example : ['created_at:desc'] \n",
			},
			"exclude_target_id": schema.StringAttribute{
				Optional: true,
				Computed: false,
				Description: "Target ID to Exclude. \n" +
					"  - example : '112c2fc8c29a449dbfa8681f8f1d7442' \n",
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: false,
				Description: "Service Control Policy ID. \n" +
					"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
			},
		},
	}
}

func (d *serviceControlPolicyListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.ServiceControlPolicyListOutput
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ListServiceControlPolicies(ctx, organization.ServiceControlPolicyListDataSourceRequest{
		OrganizationId:  state.OrganizationId,
		Name:            state.Name,
		Type:            state.Type,
		ExcludeTargetId: state.ExcludeTargetId,
		Page:            state.Page,
		Size:            state.Size,
		Id:              state.Id,
		Sort:            state.Sort,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Service Control Policies",
			err.Error(),
		)
		return
	}

	policies := make([]attr.Value, 0)
	for _, policy := range data.Policies {
		policyValue, policyDiags := d.buildPolicyListValue(ctx, policy)
		resp.Diagnostics.Append(policyDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		policies = append(policies, policyValue)
	}

	policiesList, policiesListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyListValue{}.AttributeTypes(ctx)},
		policies,
	)
	resp.Diagnostics.Append(policiesListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sortValues := make([]attr.Value, 0)
	if data.Sort != nil {
		for _, s := range data.Sort {
			sortValues = append(sortValues, types.StringValue(s))
		}
	}
	sortList, sortListDiags := types.ListValue(types.StringType, sortValues)
	resp.Diagnostics.Append(sortListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Policies = policiesList
	state.TotalCount = types.Int64Value(int64(data.Count))
	state.Page = types.Int64Value(int64(data.Page))
	state.Size = types.Int64Value(int64(data.Size))
	state.SortResult = sortList

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *serviceControlPolicyListDataSource) buildPolicyListValue(ctx context.Context, policy sdkorganization.ServiceControlPolicySummary) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var desc string
	if policy.Description.IsSet() && policy.Description.Get() != nil {
		desc = *policy.Description.Get()
	}

	documentValue, documentDiags := d.buildDocumentValue(ctx, &policy.Document)
	diags.Append(documentDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyListValue{}.AttributeTypes(ctx)), diags
	}

	policyValue, _ := types.ObjectValue(organization.ServiceControlPolicyListValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"policy_id":       types.StringValue(policy.Id),
		"organization_id": types.StringValue(policy.OrganizationId),
		"name":            types.StringValue(policy.Name),
		"description":     types.StringValue(desc),
		"type":            types.StringValue(string(policy.Type)),
		"category":        types.StringValue(policy.Category),
		"state":           types.StringValue(policy.State),
		"source":          types.StringValue(policy.Source),
		"document":        documentValue,
		"created_at":      types.StringValue(policy.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":      types.StringValue(policy.CreatedBy),
		"creator_name":    types.StringValue(policy.GetCreatorName()),
		"modified_at":     types.StringValue(policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":     types.StringValue(policy.ModifiedBy),
		"modifier_name":   types.StringValue(policy.GetModifierName()),
	})

	return policyValue, diags
}

func (d *serviceControlPolicyListDataSource) buildDocumentValue(ctx context.Context, doc *sdkorganization.ServiceControlPolicyDocument) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	statementsList := types.ListNull(types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)})
	versionVal := types.StringNull()

	if doc != nil {
		if doc.Statement != nil {
			statements := make([]attr.Value, 0, len(doc.Statement))
			for _, stmt := range doc.Statement {
				var effectStr string
				if stmt.Effect != nil {
					effectStr = *stmt.Effect
				}
				var sidStr string
				if stmt.Sid != nil {
					sidStr = *stmt.Sid
				}

				principalValue := d.buildPrincipalValue(ctx, stmt.Principal)
				conditionValue, conditionDiags := d.buildConditionValue(ctx, stmt.Condition)
				diags.Append(conditionDiags...)
				if diags.HasError() {
					return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
				}

				statementValue, statementDiags := types.ObjectValue(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx), map[string]attr.Value{
					"action":     d.buildStringList(ctx, stmt.Action),
					"condition":  conditionValue,
					"effect":     types.StringValue(effectStr),
					"not_action": d.buildStringList(ctx, stmt.NotAction),
					"principal":  principalValue,
					"resource":   d.buildStringList(ctx, stmt.Resource),
					"sid":        types.StringValue(sidStr),
				})
				diags.Append(statementDiags...)
				if diags.HasError() {
					return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
				}
				statements = append(statements, statementValue)
			}
			var statementsListDiags diag.Diagnostics
			statementsList, statementsListDiags = types.ListValue(
				types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
				statements,
			)
			diags.Append(statementsListDiags...)
			if diags.HasError() {
				return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
			}
		}
		if doc.Version != nil {
			versionVal = types.StringValue(*doc.Version)
		}
	}

	documentValue, documentDiags := types.ObjectValue(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   versionVal,
	})
	diags.Append(documentDiags...)
	return documentValue, diags
}

func (d *serviceControlPolicyListDataSource) buildStringList(ctx context.Context, values []string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}
	list := make([]attr.Value, 0, len(values))
	for _, v := range values {
		list = append(list, types.StringValue(v))
	}
	return types.ListValueMust(types.StringType, list)
}

func (d *serviceControlPolicyListDataSource) buildPrincipalValue(ctx context.Context, principal any) types.String {
	if principal == nil {
		return types.StringNull()
	}
	if p, ok := principal.(string); ok {
		return types.StringValue(p)
	}
	if p, ok := principal.(map[string]any); ok {
		if len(p) == 0 {
			return types.StringNull()
		}
		if val, exists := p["scp"]; exists {
			if arr, ok := val.([]string); ok && len(arr) > 0 {
				return types.StringValue(arr[0])
			}
		}
	}
	return types.StringNull()
}

func (d *serviceControlPolicyListDataSource) buildConditionValue(ctx context.Context, condition any) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	if condition == nil {
		return types.MapNull(types.StringType), diags
	}
	if c, ok := condition.(map[string]map[string][]string); ok {
		if len(c) == 0 {
			return types.MapNull(types.StringType), diags
		}
		result := make(map[string]attr.Value)
		for k, v := range c {
			result[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		mapVal, mapDiags := types.MapValue(types.StringType, result)
		diags.Append(mapDiags...)
		return mapVal, diags
	}
	return types.MapNull(types.StringType), diags
}
