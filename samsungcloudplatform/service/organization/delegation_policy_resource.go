package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &delegationPolicyResource{}
	_ resource.ResourceWithConfigure   = &delegationPolicyResource{}
	_ resource.ResourceWithImportState = &delegationPolicyResource{}
)

func NewDelegationPolicyResource() resource.Resource {
	return &delegationPolicyResource{}
}

type delegationPolicyResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *delegationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_delegation_policy"
}

func (r *delegationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages delegation policies for controlling account permissions within the organization",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
				Computed: true,
			},
			"document": schema.SingleNestedAttribute{
				Description: "Delegation Policy Document. \n" +
					"  - example : '{version: 2024-07-01, statement: [...]}' \n",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						Description: "Policy Statement. \n" +
							"  - example : '[{effect: Allow, action: [scp:iam:ListUsers], resource: [*], ...}]' \n",
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									Description: "List of actions permitted or denied by the policy statement. \n" +
										"  - example : ['scp:iam:ListUsers', 'scp:iam:GetUser'] \n",
									ElementType: types.StringType,
									Required:    true,
								},
								"effect": schema.StringAttribute{
									Description: "Policy Effect. \n" +
										"  - example : 'Allow' \n" +
										"  - allowed_values : ['Allow', 'Deny'] \n",
									Required: true,
								},
								"principal": schema.SingleNestedAttribute{
									Description: "Principal (object with scp list). Use this for SRN-based principals like {scp: [\"srn:dev2::...:::iam:user/...\"]}. \n" +
										"  - example : {'scp': ['srn:dev2::...:::iam:user/...']} \n",
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"scp": schema.ListAttribute{
											Description: "SCP Principal SRN list. \n" +
												"  - example : ['srn:dev2::...:::iam:user/...'] \n",
											ElementType: types.StringType,
											Optional:    true,
										},
									},
								},
								"resource": schema.ListAttribute{
									Description: "List of resources to which the policy statement applies. \n" +
										"  - example : ['scp:compute:*'] \n",
									ElementType: types.StringType,
									Required:    true,
								},
								"sid": schema.StringAttribute{
									Description: "Statement ID. \n" +
										"  - example : 'Statement1' \n",
									Optional: true,
								},
							},
						},
					},
					"version": schema.StringAttribute{
						Description: "Policy Version. \n" +
							"  - example : '2024-07-01' \n",
						Required: true,
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Computed: true,
				Description: "Delegation Policy Info. \n" +
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
													Description: "SCP Principal SRN list. \n" +
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
								Description: "Policy Syntax. \n" +
									"  - example : '[{effect: Allow, action: [scp:iam:ListUsers], resource: [*], ...}]' \n",
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

func (r *delegationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	inst, ok := req.ProviderData.(client.Instance)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Instance, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = inst.Client.Organization
	r.clients = inst.Client
}

func (r *delegationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.DelegationPolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateDelegationPolicy(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Delegation Policy",
			"Could not create Delegation Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policyValue, diags := r.buildDelegationPolicyValue(ctx, &data.Policy)
	resp.Diagnostics.Append(diags...)

	plan.Policy = policyValue

	orgId := plan.OrganizationId.ValueString()
	err = waitForDelegationPolicyReady(ctx, r.client, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for delegation policy creation",
			"Error waiting for delegation policy to be ready: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *delegationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.DelegationPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()
	if orgId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Delegation Policy",
			"Organization ID is empty",
		)
		return
	}

	data, err := r.client.GetDelegationPolicy(ctx, orgId)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	policyValue, diags := r.buildDelegationPolicyValue(ctx, &data.Policy)
	resp.Diagnostics.Append(diags...)

	state.OrganizationId = types.StringValue(data.Policy.OrganizationId)
	state.Policy = policyValue

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *delegationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organization.DelegationPolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.UpdateDelegationPolicy(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Delegation Policy",
			"Could not update Delegation Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	policyValue, diags := r.buildDelegationPolicyValue(ctx, &data.Policy)
	resp.Diagnostics.Append(diags...)

	plan.Policy = policyValue

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *delegationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.DelegationPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()
	if orgId == "" {
		resp.Diagnostics.AddError(
			"Unable to Delete Delegation Policy",
			"Organization ID is empty",
		)
		return
	}

	_, err := r.client.DeleteDelegationPolicy(ctx, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Delegation Policy",
			"Could not delete Delegation Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *delegationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	orgId := req.ID
	if orgId == "" {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"The import ID cannot be empty",
		)
		return
	}

	state := organization.DelegationPolicyResource{
		OrganizationId: types.StringValue(orgId),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *delegationPolicyResource) buildDelegationPolicyValue(ctx context.Context, policy *sdkorganization.DelegationPolicy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build document value
	documentValue, docDiags := r.buildDocumentValue(ctx, &policy.Document)
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

func (r *delegationPolicyResource) buildDocumentValue(ctx context.Context, doc *sdkorganization.DelegationPolicyDocument) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build statements
	var statements []attr.Value
	for _, stmt := range doc.Statement {
		statementValue, stmtDiags := r.buildStatementValue(ctx, &stmt)
		diags.Append(stmtDiags...)
		statements = append(statements, statementValue)
	}
	statementsList, statementsListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)
	diags.Append(statementsListDiags...)
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

func (r *delegationPolicyResource) buildStatementValue(ctx context.Context, stmt *sdkorganization.DelegationPolicyStatement) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build action list
	var actions []attr.Value
	for _, a := range stmt.Action {
		actions = append(actions, types.StringValue(a))
	}
	actionsList, actionsDiags := types.ListValue(types.StringType, actions)
	diags.Append(actionsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	// Build resource list
	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	resourcesList, resourcesDiags := types.ListValue(types.StringType, resources)
	diags.Append(resourcesDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.DelegationPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	// Build principal - supports only object (scp list)
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

func waitForDelegationPolicyReady(ctx context.Context, orgClient *organization.Client, orgId string) error {
	return client.WaitForStatus(ctx, nil, []string{}, []string{"READY"}, func() (interface{}, string, error) {
		policyResp, err := orgClient.GetDelegationPolicy(ctx, orgId)
		if err != nil {
			return nil, "", err
		}
		return policyResp, "READY", nil
	}, -1, -1, -1, -1)
}
