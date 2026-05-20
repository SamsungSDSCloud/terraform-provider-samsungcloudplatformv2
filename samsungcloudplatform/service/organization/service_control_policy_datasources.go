package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
				Required:    true,
				Computed:    false,
				Description: "Organization ID (optional filter)",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    false,
				Description: "Filter by policy name",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    false,
				Description: "Filter by policy type (SYSTEM_MANAGED or USER_DEFINED)",
			},
			"policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy ID",
							MarkdownDescription: "Policy ID",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							Description:         "Organization ID",
							MarkdownDescription: "Organization ID",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy Name",
							MarkdownDescription: "Policy Name",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "Description",
							MarkdownDescription: "Description",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy Type",
							MarkdownDescription: "Policy Type",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							Description:         "Category",
							MarkdownDescription: "Category",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							Description:         "State",
							MarkdownDescription: "State",
						},
						"source": schema.StringAttribute{
							Computed:            true,
							Description:         "Policy Creation Subject",
							MarkdownDescription: "Policy Creation Subject",
						},
						"document": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "Service Control Policy Document",
							MarkdownDescription: "Service Control Policy Document",
							Attributes: map[string]schema.Attribute{
								"statement": schema.ListNestedAttribute{
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"action": schema.ListAttribute{
												ElementType:         types.StringType,
												Computed:            true,
												Description:         "Action",
												MarkdownDescription: "Action",
											},
											"condition": schema.MapAttribute{
												ElementType:         types.StringType,
												Computed:            true,
												Description:         "Condition",
												MarkdownDescription: "Condition",
											},
											"effect": schema.StringAttribute{
												Computed:            true,
												Description:         "Effect (Allow/Deny)",
												MarkdownDescription: "Effect (Allow/Deny)",
											},
											"not_action": schema.ListAttribute{
												ElementType:         types.StringType,
												Computed:            true,
												Description:         "Not Action",
												MarkdownDescription: "Not Action",
											},
											"principal": schema.StringAttribute{
												Computed:            true,
												Description:         "Principal (e.g., \"*\")",
												MarkdownDescription: "Principal (e.g., \"*\")",
											},
											"resource": schema.ListAttribute{
												ElementType:         types.StringType,
												Computed:            true,
												Description:         "Resource",
												MarkdownDescription: "Resource",
											},
											"sid": schema.StringAttribute{
												Computed:            true,
												Description:         "Statement ID",
												MarkdownDescription: "Statement ID",
											},
										},
									},
									Computed:            true,
									Description:         "Policy Statement",
									MarkdownDescription: "Policy Statement",
								},
								"version": schema.StringAttribute{
									Computed:            true,
									Description:         "Policy Version",
									MarkdownDescription: "Policy Version",
								},
							},
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Created At",
							MarkdownDescription: "Created At",
						},
						"created_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Created By",
							MarkdownDescription: "Created By",
						},
						"creator_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Creator Name",
							MarkdownDescription: "Creator Name",
						},
						"modified_at": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified At",
							MarkdownDescription: "Modified At",
						},
						"modified_by": schema.StringAttribute{
							Computed:            true,
							Description:         "Modified By",
							MarkdownDescription: "Modified By",
						},
						"modifier_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Modifier Name",
							MarkdownDescription: "Modifier Name",
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:            true,
				Description:         "Total count",
				MarkdownDescription: "Total count",
			},
			"page": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Page number",
				MarkdownDescription: "Page number",
			},
			"size": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "Page size",
				MarkdownDescription: "Page size",
			},
			"sort": schema.StringAttribute{
				Optional:    true,
				Description: "Sort criteria (e.g., 'created_at:desc')",
			},
			"sort_result": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Sort Criteria from response",
			},
			"exclude_target_id": schema.StringAttribute{
				Optional:            true,
				Computed:            false,
				Description:         "Exclude target ID",
				MarkdownDescription: "Exclude target ID",
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            false,
				Description:         "Filter by policy ID",
				MarkdownDescription: "Filter by policy ID",
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
		policyValue := d.buildPolicyListValue(ctx, policy)
		policies = append(policies, policyValue)
	}

	policiesList := types.ListValueMust(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyListValue{}.AttributeTypes(ctx)},
		policies,
	)

	sortValues := make([]attr.Value, 0)
	if data.Sort != nil {
		for _, s := range data.Sort {
			sortValues = append(sortValues, types.StringValue(s))
		}
	}
	sortList := types.ListValueMust(types.StringType, sortValues)

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

func (d *serviceControlPolicyListDataSource) buildPolicyListValue(ctx context.Context, policy sdkorganization.ServiceControlPolicySummary) types.Object {
	var desc string
	if policy.Description.IsSet() && policy.Description.Get() != nil {
		desc = *policy.Description.Get()
	}

	documentValue := d.buildDocumentValue(ctx, &policy.Document)

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

	return policyValue
}

func (d *serviceControlPolicyListDataSource) buildDocumentValue(ctx context.Context, doc *sdkorganization.ServiceControlPolicyDocument) types.Object {
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
				conditionValue := d.buildConditionValue(ctx, stmt.Condition)

				statementValue, _ := types.ObjectValue(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx), map[string]attr.Value{
					"action":     d.buildStringList(ctx, stmt.Action),
					"condition":  conditionValue,
					"effect":     types.StringValue(effectStr),
					"not_action": d.buildStringList(ctx, stmt.NotAction),
					"principal":  principalValue,
					"resource":   d.buildStringList(ctx, stmt.Resource),
					"sid":        types.StringValue(sidStr),
				})
				statements = append(statements, statementValue)
			}
			statementsList = types.ListValueMust(
				types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
				statements,
			)
		}
		if doc.Version != nil {
			versionVal = types.StringValue(*doc.Version)
		}
	}

	documentValue, _ := types.ObjectValue(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   versionVal,
	})

	return documentValue
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

func (d *serviceControlPolicyListDataSource) buildConditionValue(ctx context.Context, condition any) types.Map {
	if condition == nil {
		return types.MapNull(types.StringType)
	}
	if c, ok := condition.(map[string]map[string][]string); ok {
		if len(c) == 0 {
			return types.MapNull(types.StringType)
		}
		result := make(map[string]attr.Value)
		for k, v := range c {
			result[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		return types.MapValueMust(types.StringType, result)
	}
	return types.MapNull(types.StringType)
}
