package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/organization/1.3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &delegationPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &delegationPolicyDataSource{}
)

func NewDelegationPolicyDataSource() datasource.DataSource {
	return &delegationPolicyDataSource{}
}

type delegationPolicyDataSource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (d *delegationPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_delegation_policy"
}

func (d *delegationPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *delegationPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Show Organization Delegation Policy",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
				Computed: true,
			},
			"document": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Delegation Policy Document. \n" +
					"  - example : '{version: 2024-07-01, statement: [...]}' \n",
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "List of actions permitted or denied by the policy statement. \n" +
										"  - example : ['scp:iam:ListUsers'] \n",
								},
								"effect": schema.StringAttribute{
									Computed: true,
									Description: "Policy Effect. \n" +
										"  - example : 'Allow' \n",
								},
								"principal": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"scp": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Description: "Principal entity to which the policy statement applies. \n" +
												"  - example : ['srn:dev2::...:::iam:user/...'] \n",
										},
									},
									Description: "Principal entity to which the policy statement applies. \n" +
										"  - example : '{scp: [srn:dev2::...:::iam:user/...]}' \n",
								},
								"resource": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "List of resources to which the policy statement applies. \n" +
										"  - example : ['scp:compute:*'] \n",
								},
								"sid": schema.StringAttribute{
									Computed: true,
									Description: "Statement ID. \n" +
										"  - example : 'Statement1' \n",
								},
							},
						},
						Computed: true,
						Description: "Policy Statement. \n" +
							"  - example : '[{effect: Allow, action: [organization:ListAccounts], resource: [*], ...}]' \n",
					},
					"version": schema.StringAttribute{
						Computed: true,
						Description: "Policy Version. \n" +
							"  - example : '2024-07-01' \n",
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Delegation Policy. \n" +
					"  - example : '{organization_id: o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5, created_at: 2025-01-01T00:00:00.000Z, ...}' \n",
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the delegation policy was created. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "User who created the delegation policy. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"document": schema.SingleNestedAttribute{
						Computed: true,
						Description: "Delegation Policy Document. \n" +
							"  - example : '{version: 2024-07-01, statement: [...]}' \n",
						Attributes: map[string]schema.Attribute{
							"statement": schema.ListNestedAttribute{
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"action": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Description: "List of actions permitted or denied by the policy statement. \n" +
												"  - example : ['scp:iam:ListUsers'] \n",
										},
										"effect": schema.StringAttribute{
											Computed: true,
											Description: "Policy Effect. \n" +
												"  - example : 'Allow' \n",
										},
										"principal": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"scp": schema.ListAttribute{
													ElementType: types.StringType,
													Computed:    true,
													Description: "SCP Principal. \n" +
														"  - example : ['srn:dev2::...:::iam:user/...'] \n",
												},
											},
											Description: "Principal entity to which the policy statement applies. \n" +
												"  - example : '{scp: [srn:dev2::...:::iam:user/...]}' \n",
										},
										"resource": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
											Description: "List of resources to which the policy statement applies. \n" +
												"  - example : ['scp:compute:*'] \n",
										},
										"sid": schema.StringAttribute{
											Computed: true,
											Description: "Statement ID. \n" +
												"  - example : 'Statement1' \n",
										},
									},
								},
								Computed: true,
								Description: "Policy Statement. \n" +
									"  - example : '[{effect: Allow, action: [organization:ListAccounts], resource: [*], ...}]' \n",
							},
							"version": schema.StringAttribute{
								Computed: true,
								Description: "Policy Version. \n" +
									"  - example : '2024-07-01' \n",
							},
						},
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the delegation policy was modified. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "User who modified the delegation policy. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"organization_id": schema.StringAttribute{
						Computed: true,
						Description: "Unique identifier of the organization. \n" +
							"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
					},
				},
			},
		},
	}
}

func (d *delegationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state organization.DelegationPolicyDataSource
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()

	data, err := d.client.GetDelegationPolicy(ctx, orgId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Read Delegation Policy",
			err.Error(),
		)
		return
	}

	// Build document value for the top-level document attribute
	documentValue, docDiags := d.buildDocumentValue(ctx, &data.Policy.Document)
	resp.Diagnostics.Append(docDiags...)

	// Build policy value for the computed policy attribute
	policyValue, policyDiags := d.buildPolicyValue(ctx, &data.Policy)
	resp.Diagnostics.Append(policyDiags...)

	state.Document = documentValue
	state.Policy = policyValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *delegationPolicyDataSource) buildPolicyValue(ctx context.Context, policy *sdkorganization.DelegationPolicy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	documentValue, docDiags := d.buildDocumentValue(ctx, &policy.Document)
	diags.Append(docDiags...)

	policyValue, valueDiags := types.ObjectValue(organization.DelegationPolicyValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"created_at":      types.StringValue(policy.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":      types.StringValue(policy.CreatedBy),
		"document":        documentValue,
		"modified_at":     types.StringValue(policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":     types.StringValue(policy.ModifiedBy),
		"organization_id": types.StringValue(policy.OrganizationId),
	})
	diags.Append(valueDiags...)

	return policyValue, diags
}

func (d *delegationPolicyDataSource) buildDocumentValue(ctx context.Context, doc *sdkorganization.DelegationPolicyDocument) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var statements []attr.Value
	for _, stmt := range doc.Statement {
		statementValue, stmtDiags := d.buildStatementValue(ctx, &stmt)
		diags.Append(stmtDiags...)
		statements = append(statements, statementValue)
	}
	statementsList, statementsDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)
	diags.Append(statementsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.DelegationPolicyDocumentValue{}.AttributeTypes(ctx)), diags
	}

	documentValue, valueDiags := types.ObjectValue(organization.DelegationPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   types.StringValue(doc.GetVersion()),
	})
	diags.Append(valueDiags...)

	return documentValue, diags
}

func (d *delegationPolicyDataSource) buildStatementValue(ctx context.Context, stmt *sdkorganization.DelegationPolicyStatement) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var actions []attr.Value
	for _, a := range stmt.Action {
		actions = append(actions, types.StringValue(a))
	}
	actionsList, actionsDiags := types.ListValue(types.StringType, actions)
	diags.Append(actionsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	resourcesList, resourcesDiags := types.ListValue(types.StringType, resources)
	diags.Append(resourcesDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	var principalValue types.Object
	principal := stmt.GetPrincipal()
	if principal.MapmapOfStringarrayOfString != nil {
		var scps []attr.Value
		if scpList, ok := (*principal.MapmapOfStringarrayOfString)["scp"]; ok {
			for _, s := range scpList {
				scps = append(scps, types.StringValue(s))
			}
		}
		scpsList, scpsDiags := types.ListValue(types.StringType, scps)
		diags.Append(scpsDiags...)
		if diags.HasError() {
			return types.ObjectNull(organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)), diags
		}

		var err diag.Diagnostics
		principalValue, err = types.ObjectValue(organization.DelegationPolicyPrincipalValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"scp": scpsList,
		})
		diags.Append(err...)
	} else {
		principalValue = types.ObjectNull(organization.DelegationPolicyPrincipalValue{}.AttributeTypes(ctx))
	}

	statementValue, valueDiags := types.ObjectValue(organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"action":    actionsList,
		"effect":    types.StringValue(stmt.GetEffect()),
		"principal": principalValue,
		"resource":  resourcesList,
		"sid":       types.StringValue(stmt.GetSid()),
	})
	diags.Append(valueDiags...)

	return statementValue, diags
}
