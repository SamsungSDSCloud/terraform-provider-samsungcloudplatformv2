package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &delegationPolicyResource{}
	_ resource.ResourceWithConfigure = &delegationPolicyResource{}
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
		Description: "Organization Delegation Policy",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
			},
			"document": schema.SingleNestedAttribute{
				Description:         "Delegation Policy Document",
				MarkdownDescription: "Delegation Policy Document",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						Description:         "Policy Statement",
						MarkdownDescription: "Policy Statement",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									Description:         "Action",
									MarkdownDescription: "Action",
									ElementType:         types.StringType,
									Required:            true,
								},
								"effect": schema.StringAttribute{
									Description:         "Effect (Allow/Deny)",
									MarkdownDescription: "Effect (Allow/Deny)",
									Required:            true,
								},
								"principal": schema.SingleNestedAttribute{
									Description:         "Principal (object with scp list)",
									MarkdownDescription: "Principal (object with scp list). Use this for SRN-based principals like {scp: [\"srn:dev2::...:::iam:user/...\"]}",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"scp": schema.ListAttribute{
											Description:         "SCP Principal SRN list",
											MarkdownDescription: "SCP Principal SRN list",
											ElementType:         types.StringType,
											Optional:            true,
										},
									},
								},
								"resource": schema.ListAttribute{
									Description:         "Resource",
									MarkdownDescription: "Resource",
									ElementType:         types.StringType,
									Required:            true,
								},
								"sid": schema.StringAttribute{
									Description:         "Statement ID",
									MarkdownDescription: "Statement ID",
									Optional:            true,
								},
							},
						},
					},
					"version": schema.StringAttribute{
						Description:         "Policy Version",
						MarkdownDescription: "Policy Version",
						Required:            true,
					},
				},
			},
			"policy": schema.SingleNestedAttribute{
				Computed: true,
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
											ElementType: types.StringType,
											Computed:    true,
										},
										"effect": schema.StringAttribute{
											Computed: true,
										},
										"principal": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"scp": schema.ListAttribute{
													ElementType: types.StringType,
													Computed:    true,
												},
											},
										},
										"resource": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
										"sid": schema.StringAttribute{
											Computed: true,
										},
									},
								},
								Computed: true,
							},
							"version": schema.StringAttribute{
								Computed: true,
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
		resp.Diagnostics.AddError(
			"Unable to Read Delegation Policy",
			err.Error(),
		)
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

func (r *delegationPolicyResource) buildStatementValue(ctx context.Context, stmt *sdkorganization.DelegationPolicyStatement) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build action list
	var actions []attr.Value
	for _, a := range stmt.Action {
		actions = append(actions, types.StringValue(a))
	}
	actionsList := types.ListValueMust(types.StringType, actions)

	// Build resource list
	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	resourcesList := types.ListValueMust(types.StringType, resources)

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
