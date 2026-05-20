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
	_ resource.Resource              = &serviceControlPolicyResource{}
	_ resource.ResourceWithConfigure = &serviceControlPolicyResource{}
)

func NewServiceControlPolicyResource() resource.Resource {
	return &serviceControlPolicyResource{}
}

type serviceControlPolicyResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *serviceControlPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_service_control_policy"
}

func (r *serviceControlPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Organization Service Control Policy",
		Attributes: map[string]schema.Attribute{
			"policy_id": schema.StringAttribute{
				Description:         "Policy ID",
				MarkdownDescription: "Policy ID",
				Optional:            true,
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Service Control Policy Name",
				MarkdownDescription: "Service Control Policy Name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Description:         "Description",
				MarkdownDescription: "Description",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				Description:         "Policy Type (SYSTEM_MANAGED or USER_DEFINED)",
				MarkdownDescription: "Policy Type (SYSTEM_MANAGED or USER_DEFINED)",
				Optional:            true,
				Computed:            true,
			},
			"document": schema.SingleNestedAttribute{
				Description:         "Service Control Policy Document",
				MarkdownDescription: "Service Control Policy Document",
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
									Optional:            true,
								},
								"condition": schema.MapAttribute{
									Description:         "Condition",
									MarkdownDescription: "Condition",
									ElementType:         types.StringType,
									Optional:            true,
								},
								"effect": schema.StringAttribute{
									Description:         "Effect (Allow/Deny)",
									MarkdownDescription: "Effect (Allow/Deny)",
									Required:            true,
								},
								"not_action": schema.ListAttribute{
									Description:         "Not Action",
									MarkdownDescription: "Not Action",
									ElementType:         types.StringType,
									Optional:            true,
								},
								"principal": schema.StringAttribute{
									Description:         "Principal (e.g., \"*\")",
									MarkdownDescription: "Principal (e.g., \"*\")",
									Optional:            true,
								},
								"resource": schema.ListAttribute{
									Description:         "Resource",
									MarkdownDescription: "Resource",
									ElementType:         types.StringType,
									Optional:            true,
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
						Optional:            true,
					},
				},
			},
			"category": schema.StringAttribute{
				Computed:            true,
				Description:         "Category",
				MarkdownDescription: "Category",
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
			"service_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Service Name",
				MarkdownDescription: "Service Name",
			},
			"source": schema.StringAttribute{
				Computed:            true,
				Description:         "Source",
				MarkdownDescription: "Source",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "State",
				MarkdownDescription: "State",
			},
		},
	}
}

func (r *serviceControlPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serviceControlPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.ServiceControlPolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateServiceControlPolicy(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Service Control Policy",
			"Could not create Service Control Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	documentValue, diags := r.buildDocumentValueFromResponse(ctx, &data.Policy.Document, &plan.Document)
	resp.Diagnostics.Append(diags...)

	plan.PolicyId = types.StringValue(data.Policy.Id)
	plan.Document = documentValue
	plan.Category = types.StringValue(data.Policy.Category)
	plan.CreatedAt = types.StringValue(data.Policy.CreatedAt.Format("2006-01-02T15:04:05.000Z"))
	plan.CreatedBy = types.StringValue(data.Policy.CreatedBy)
	plan.CreatorName = types.StringValue(data.Policy.GetCreatorName())

	desc := data.Policy.Description.Get()
	if desc != nil {
		plan.Description = types.StringValue(*desc)
	}

	plan.ModifiedAt = types.StringValue(data.Policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z"))
	plan.ModifiedBy = types.StringValue(data.Policy.ModifiedBy)
	plan.ModifierName = types.StringValue(data.Policy.GetModifierName())
	plan.ServiceName = types.StringValue(data.Policy.ServiceName)
	plan.Source = types.StringValue(data.Policy.Source)
	plan.State = types.StringValue(data.Policy.State)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceControlPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.ServiceControlPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyId.ValueString()
	orgId := state.OrganizationId.ValueString()

	if policyId == "" {
		resp.Diagnostics.AddError(
			"Unable to Read Service Control Policy",
			"Policy ID is empty",
		)
		return
	}

	data, err := r.client.GetServiceControlPolicy(ctx, policyId, orgId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Service Control Policy",
			err.Error(),
		)
		return
	}

	documentValue, diags := r.buildDocumentValueFromResponse(ctx, &data.Policy.Document, &state.Document)
	resp.Diagnostics.Append(diags...)

	state.PolicyId = types.StringValue(data.Policy.Id)
	state.OrganizationId = types.StringValue(data.Policy.OrganizationId)
	state.Document = documentValue
	state.Category = types.StringValue(data.Policy.Category)
	state.CreatedAt = types.StringValue(data.Policy.CreatedAt.Format("2006-01-02T15:04:05.000Z"))
	state.CreatedBy = types.StringValue(data.Policy.CreatedBy)
	state.CreatorName = types.StringValue(data.Policy.GetCreatorName())

	desc := data.Policy.Description.Get()
	if desc != nil {
		state.Description = types.StringValue(*desc)
	}

	state.ModifiedAt = types.StringValue(data.Policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z"))
	state.ModifiedBy = types.StringValue(data.Policy.ModifiedBy)
	state.ModifierName = types.StringValue(data.Policy.GetModifierName())
	state.ServiceName = types.StringValue(data.Policy.ServiceName)
	state.Source = types.StringValue(data.Policy.Source)
	state.State = types.StringValue(data.Policy.State)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceControlPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organization.ServiceControlPolicyResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := plan.PolicyId.ValueString()
	orgId := plan.OrganizationId.ValueString()
	if policyId == "" || orgId == "" {
		var state organization.ServiceControlPolicyResource
		stateDiags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(stateDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if policyId == "" {
			policyId = state.PolicyId.ValueString()
		}
		if orgId == "" {
			orgId = state.OrganizationId.ValueString()
		}
	}

	if policyId == "" || orgId == "" {
		resp.Diagnostics.AddError(
			"Unable to Update Service Control Policy",
			"Policy ID or Organization ID is empty",
		)
		return
	}

	_, err := r.client.UpdateServiceControlPolicy(ctx, policyId, orgId, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Service Control Policy",
			"Could not update Service Control Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	readData, err := r.client.GetServiceControlPolicy(ctx, policyId, orgId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Service Control Policy after update",
			err.Error(),
		)
		return
	}

	plan.PolicyId = types.StringValue(readData.Policy.Id)
	plan.Name = types.StringValue(readData.Policy.Name)
	plan.Category = types.StringValue(readData.Policy.Category)
	plan.CreatedAt = types.StringValue(readData.Policy.CreatedAt.Format("2006-01-02T15:04:05.000Z"))
	plan.CreatedBy = types.StringValue(readData.Policy.CreatedBy)
	plan.CreatorName = types.StringValue(readData.Policy.GetCreatorName())
	plan.ModifiedAt = types.StringValue(readData.Policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z"))
	plan.ModifiedBy = types.StringValue(readData.Policy.ModifiedBy)
	plan.ModifierName = types.StringValue(readData.Policy.GetModifierName())
	plan.ServiceName = types.StringValue(readData.Policy.ServiceName)
	plan.Source = types.StringValue(readData.Policy.Source)
	plan.State = types.StringValue(readData.Policy.State)

	desc := readData.Policy.Description.Get()
	if desc != nil {
		plan.Description = types.StringValue(*desc)
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *serviceControlPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.ServiceControlPolicyResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyId := state.PolicyId.ValueString()
	orgId := state.OrganizationId.ValueString()

	if policyId == "" {
		resp.Diagnostics.AddError(
			"Unable to Delete Service Control Policy",
			"Policy ID is empty",
		)
		return
	}

	_, err := r.client.DeleteServiceControlPolicies(ctx, orgId, []string{policyId})
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Service Control Policy",
			"Could not delete Service Control Policy, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}

func (r *serviceControlPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: policy_id or policy_id:organization_id
	id := req.ID
	orgId := ""

	// Parse the import ID
	for i := len(id) - 1; i >= 0; i-- {
		if id[i] == ':' {
			orgId = id[i+1:]
			id = id[:i]
			break
		}
	}

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"The import ID cannot be empty",
		)
		return
	}

	state := organization.ServiceControlPolicyResource{
		PolicyId:       types.StringValue(id),
		OrganizationId: types.StringValue(orgId),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *serviceControlPolicyResource) buildServiceControlPolicyValue(ctx context.Context, policy *sdkorganization.ServiceControlPolicy) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build document value
	documentValue, docDiags := r.buildDocumentValue(ctx, &policy.Document)
	diags.Append(docDiags...)

	desc := policy.Description.Get()
	var descStr string
	if desc != nil {
		descStr = *desc
	}

	policyValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"category":        types.StringValue(policy.Category),
		"created_at":      types.StringValue(policy.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":      types.StringValue(policy.CreatedBy),
		"creator_name":    types.StringValue(policy.GetCreatorName()),
		"description":     types.StringValue(descStr),
		"document":        documentValue,
		"id":              types.StringValue(policy.Id),
		"modified_at":     types.StringValue(policy.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":     types.StringValue(policy.ModifiedBy),
		"modifier_name":   types.StringValue(policy.GetModifierName()),
		"name":            types.StringValue(policy.Name),
		"organization_id": types.StringValue(policy.OrganizationId),
		"source":          types.StringValue(policy.Source),
		"state":           types.StringValue(policy.State),
		"type":            types.StringValue(string(policy.Type)),
	})
	diags.Append(valueDiags...)

	return policyValue, diags
}

func (r *serviceControlPolicyResource) buildDocumentValue(ctx context.Context, doc *sdkorganization.ServiceControlPolicyDocument) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build statements
	var statements []attr.Value
	for _, stmt := range doc.Statement {
		statementValue, stmtDiags := r.buildStatementValue(ctx, &stmt)
		diags.Append(stmtDiags...)
		statements = append(statements, statementValue)
	}
	statementsList := types.ListValueMust(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)

	documentValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   types.StringValue(doc.GetVersion()),
	})
	diags.Append(valueDiags...)

	return documentValue, diags
}

func (r *serviceControlPolicyResource) buildDocumentValueFromResponse(ctx context.Context, doc *sdkorganization.ServiceControlPolicyDocument, inputDoc *types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	if doc == nil {
		return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
	}

	inputStatements := make([]types.Object, 0)
	if inputDoc != nil && !inputDoc.IsNull() && !inputDoc.IsUnknown() {
		if attrs := inputDoc.Attributes(); attrs != nil {
			if stmtAttr, ok := attrs["statement"]; ok {
				if stmtList, ok := stmtAttr.(types.List); ok {
					for _, elem := range stmtList.Elements() {
						if stmtObj, ok := elem.(types.Object); ok {
							inputStatements = append(inputStatements, stmtObj)
						}
					}
				}
			}
		}
	}

	var statements []attr.Value
	for i, stmt := range doc.Statement {
		var inputStmt *types.Object
		if i < len(inputStatements) {
			inputStmt = &inputStatements[i]
		}
		statementValue, stmtDiags := r.buildStatementValueFromResponse(ctx, &stmt, inputStmt)
		diags.Append(stmtDiags...)
		statements = append(statements, statementValue)
	}

	versionVal := types.StringNull()
	if doc.GetVersion() != "" {
		versionVal = types.StringValue(doc.GetVersion())
	} else if inputDoc != nil && !inputDoc.IsNull() && !inputDoc.IsUnknown() {
		if attrs := inputDoc.Attributes(); attrs != nil {
			if vAttr, ok := attrs["version"]; ok && !vAttr.IsNull() && !vAttr.IsUnknown() {
				versionVal = vAttr.(types.String)
			}
		}
	}

	statementsList := types.ListValueMust(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)

	documentValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   versionVal,
	})
	diags.Append(valueDiags...)

	return documentValue, diags
}

func (r *serviceControlPolicyResource) buildStatementValueFromResponse(ctx context.Context, stmt *sdkorganization.ServiceControlPolicyStatement, inputStmt *types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var actions []attr.Value
	for _, a := range stmt.Action {
		actions = append(actions, types.StringValue(a))
	}
	if len(actions) == 0 && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if aAttr, ok := attrs["action"]; ok && !aAttr.IsNull() && !aAttr.IsUnknown() {
				if aList, ok := aAttr.(types.List); ok {
					for _, e := range aList.Elements() {
						if s, ok := e.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
							actions = append(actions, s)
						}
					}
				}
			}
		}
	}
	actionsList := types.ListValueMust(types.StringType, actions)

	var notActions []attr.Value
	for _, a := range stmt.NotAction {
		notActions = append(notActions, types.StringValue(a))
	}
	if len(notActions) == 0 && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if naAttr, ok := attrs["not_action"]; ok && !naAttr.IsNull() && !naAttr.IsUnknown() {
				if naList, ok := naAttr.(types.List); ok {
					for _, e := range naList.Elements() {
						if s, ok := e.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
							notActions = append(notActions, s)
						}
					}
				}
			}
		}
	}
	notActionsList := types.ListValueMust(types.StringType, notActions)

	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	if len(resources) == 0 && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if rAttr, ok := attrs["resource"]; ok && !rAttr.IsNull() && !rAttr.IsUnknown() {
				if rList, ok := rAttr.(types.List); ok {
					for _, e := range rList.Elements() {
						if s, ok := e.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
							resources = append(resources, s)
						}
					}
				}
			}
		}
	}
	resourcesList := types.ListValueMust(types.StringType, resources)

	conditionValue := r.buildConditionValue(ctx, stmt.Condition)
	if conditionValue.IsNull() && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if cAttr, ok := attrs["condition"]; ok && !cAttr.IsNull() && !cAttr.IsUnknown() {
				conditionValue = cAttr.(types.Map)
			}
		}
	}

	principalValue := r.buildPrincipalValue(ctx, stmt.Principal)
	if principalValue.IsNull() && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if pAttr, ok := attrs["principal"]; ok && !pAttr.IsNull() && !pAttr.IsUnknown() {
				principalValue = pAttr.(types.String)
			}
		}
	}

	effectVal := types.StringValue(stmt.GetEffect())
	if stmt.GetEffect() == "" && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if eAttr, ok := attrs["effect"]; ok && !eAttr.IsNull() && !eAttr.IsUnknown() {
				effectVal = eAttr.(types.String)
			}
		}
	}

	sidVal := types.StringValue(stmt.GetSid())
	if stmt.GetSid() == "" && inputStmt != nil && !inputStmt.IsNull() && !inputStmt.IsUnknown() {
		if attrs := inputStmt.Attributes(); attrs != nil {
			if sAttr, ok := attrs["sid"]; ok && !sAttr.IsNull() && !sAttr.IsUnknown() {
				sidVal = sAttr.(types.String)
			}
		}
	}

	statementValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"action":     actionsList,
		"condition":  conditionValue,
		"effect":     effectVal,
		"not_action": notActionsList,
		"principal":  principalValue,
		"resource":   resourcesList,
		"sid":        sidVal,
	})
	diags.Append(valueDiags...)

	return statementValue, diags
}

func (r *serviceControlPolicyResource) buildStatementValue(ctx context.Context, stmt *sdkorganization.ServiceControlPolicyStatement) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var actions []attr.Value
	for _, a := range stmt.Action {
		actions = append(actions, types.StringValue(a))
	}
	actionsList := types.ListValueMust(types.StringType, actions)

	var notActions []attr.Value
	for _, a := range stmt.NotAction {
		notActions = append(notActions, types.StringValue(a))
	}
	notActionsList := types.ListValueMust(types.StringType, notActions)

	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	resourcesList := types.ListValueMust(types.StringType, resources)

	conditionValue := r.buildConditionValue(ctx, stmt.Condition)
	principalValue := r.buildPrincipalValue(ctx, stmt.Principal)

	statementValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"action":     actionsList,
		"condition":  conditionValue,
		"effect":     types.StringValue(stmt.GetEffect()),
		"not_action": notActionsList,
		"principal":  principalValue,
		"resource":   resourcesList,
		"sid":        types.StringValue(stmt.GetSid()),
	})
	diags.Append(valueDiags...)

	return statementValue, diags
}

func (r *serviceControlPolicyResource) buildPrincipalValue(ctx context.Context, principal any) types.String {
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

func (r *serviceControlPolicyResource) buildConditionValue(ctx context.Context, condition any) types.Map {
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
