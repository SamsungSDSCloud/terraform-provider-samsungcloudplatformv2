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
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
			},
			"document": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Delegation Policy Document",
				MarkdownDescription: "Delegation Policy Document",
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
								"effect": schema.StringAttribute{
									Computed:            true,
									Description:         "Effect (Allow/Deny)",
									MarkdownDescription: "Effect (Allow/Deny)",
								},
								"principal": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"scp": schema.ListAttribute{
											ElementType:         types.StringType,
											Computed:            true,
											Description:         "SCP Principal",
											MarkdownDescription: "SCP Principal",
										},
									},
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
			"policy": schema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Delegation Policy",
				MarkdownDescription: "Delegation Policy",
				Attributes: map[string]schema.Attribute{
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
					"document": schema.SingleNestedAttribute{
						Computed:            true,
						Description:         "Delegation Policy Document",
						MarkdownDescription: "Delegation Policy Document",
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
										"effect": schema.StringAttribute{
											Computed:            true,
											Description:         "Effect (Allow/Deny)",
											MarkdownDescription: "Effect (Allow/Deny)",
										},
										"principal": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"scp": schema.ListAttribute{
													ElementType:         types.StringType,
													Computed:            true,
													Description:         "SCP Principal",
													MarkdownDescription: "SCP Principal",
												},
											},
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
					"organization_id": schema.StringAttribute{
						Computed:            true,
						Description:         "Organization ID",
						MarkdownDescription: "Organization ID",
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
	statementsList := types.ListValueMust(
		types.ObjectType{AttrTypes: organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)

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
	actionsList := types.ListValueMust(types.StringType, actions)

	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	resourcesList := types.ListValueMust(types.StringType, resources)

	var principalValue types.Object
	principal := stmt.GetPrincipal()
	if principal.MapmapOfStringarrayOfString != nil {
		var scps []attr.Value
		if scpList, ok := (*principal.MapmapOfStringarrayOfString)["scp"]; ok {
			for _, s := range scpList {
				scps = append(scps, types.StringValue(s))
			}
		}
		scpsList := types.ListValueMust(types.StringType, scps)

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
