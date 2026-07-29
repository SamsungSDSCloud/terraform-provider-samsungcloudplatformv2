package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/library/organization/1.2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const TimeFormat = "2006-01-02T15:04:05.000Z"

var (
	_ resource.Resource                = &serviceControlPolicyResource{}
	_ resource.ResourceWithConfigure   = &serviceControlPolicyResource{}
	_ resource.ResourceWithImportState = &serviceControlPolicyResource{}
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
		Description: "Manages service control policies that define permission boundaries for accounts within the organization",
		Attributes: map[string]schema.Attribute{
			"policy_id": schema.StringAttribute{
				Description: "Service Control Policy ID. \n" +
					"  - example : '138c2fc8c29a449dbfa8681f8f1d78e2' \n",
				Optional: true,
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "Service Control Policy Name. \n" +
					"  - example : 'MyPolicy' \n",
				Required: true,
			},
			"description": schema.StringAttribute{
				Description: "Policy Description. \n" +
					"  - example : 'This is an example policy.' \n",
				Optional: true,
			},
			"type": schema.StringAttribute{
				Description: "Service Control Policy Type. \n" +
					"  - example : 'MANAGED' \n",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"document": schema.SingleNestedAttribute{
				Description: "Service Control Policy Document. \n" +
					"  - example : '{version: 2024-07-01, statement: [...]}' \n",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"statement": schema.ListNestedAttribute{
						Description: "Policy Syntax. \n" +
							"  - example : '[{effect: Allow, action: [*], resource: [*], sid: statement1}]' \n",
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action": schema.ListAttribute{
									Description: "List of actions permitted or denied by the policy statement. \n" +
										"  - example : ['s3:PutObject', 's3:GetObject'] \n",
									ElementType: types.StringType,
									Optional:    true,
								},
								"condition": schema.MapAttribute{
									Description: "Policy Condition. \n" +
										"  - example : '{StringEquals: {aws:RequestedRegion: us-east-1}}' \n",
									ElementType: types.StringType,
									Optional:    true,
								},
								"effect": schema.StringAttribute{
									Description: "Policy Effect. \n" +
										"  - example : 'Allow' \n" +
										"  - allowed_values : ['Allow', 'Deny'] \n",
									Required: true,
								},
								"not_action": schema.ListAttribute{
									Description: "Policy Exclusion Action. \n" +
										"  - example : '['s3:DeleteObject']' \n",
									ElementType: types.StringType,
									Optional:    true,
								},
								"principal": schema.StringAttribute{
									Description: "Principal entity to which the policy statement applies. \n" +
										"  - example : '*' \n",
									Optional: true,
								},
								"resource": schema.ListAttribute{
									Description: "List of resources to which the policy statement applies. \n" +
										"  - example : ['*'] \n",
									ElementType: types.StringType,
									Optional:    true,
								},
								"sid": schema.StringAttribute{
									Description: "Syntax ID. \n" +
										"  - example : 'statement1' \n",
									Optional: true,
								},
							},
						},
					},
					"version": schema.StringAttribute{
						Description: "Policy Version. \n" +
							"  - example : '2012-10-17' \n",
						Optional: true,
					},
				},
			},
			"category": schema.StringAttribute{
				Computed: true,
				Description: "Policy Category. \n" +
					"  - example : 'SCP' \n",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				Description: "Timestamp when the service control policy was created. \n" +
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
			"service_name": schema.StringAttribute{
				Computed: true,
				Description: "Name of the service to which the policy applies. \n" +
					"  - example : 'Organization' \n",
			},
			"source": schema.StringAttribute{
				Computed: true,
				Description: "Policy Creation Subject. \n" +
					"  - example : 'ORGANIZATION' \n",
			},
			"state": schema.StringAttribute{
				Computed: true,
				Description: "service control policy state. \n" +
					"  - example : 'ACTIVE' \n",
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
	plan.CreatedAt = types.StringValue(data.Policy.CreatedAt.Format(TimeFormat))
	plan.CreatedBy = types.StringValue(data.Policy.CreatedBy)
	plan.CreatorName = types.StringValue(data.Policy.GetCreatorName())

	desc := data.Policy.Description.Get()
	if desc != nil {
		plan.Description = types.StringValue(*desc)
	}

	plan.ModifiedAt = types.StringValue(data.Policy.ModifiedAt.Format(TimeFormat))
	plan.ModifiedBy = types.StringValue(data.Policy.ModifiedBy)
	plan.ModifierName = types.StringValue(data.Policy.GetModifierName())
	plan.ServiceName = types.StringValue(data.Policy.ServiceName)
	plan.Source = types.StringValue(data.Policy.Source)
	plan.State = types.StringValue(data.Policy.State)

	policyId := data.Policy.Id
	orgId := plan.OrganizationId.ValueString()
	err = waitForServiceControlPolicyReady(ctx, r.client, policyId, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for service control policy creation",
			"Error waiting for service control policy to be ready: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

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
		resp.State.RemoveResource(ctx)
		return
	}

	documentValue, diags := r.buildDocumentValueFromResponse(ctx, &data.Policy.Document, &state.Document)
	resp.Diagnostics.Append(diags...)

	state.PolicyId = types.StringValue(data.Policy.Id)
	state.OrganizationId = types.StringValue(data.Policy.OrganizationId)
	state.Document = documentValue
	state.Category = types.StringValue(data.Policy.Category)
	state.CreatedAt = types.StringValue(data.Policy.CreatedAt.Format(TimeFormat))
	state.CreatedBy = types.StringValue(data.Policy.CreatedBy)
	state.CreatorName = types.StringValue(data.Policy.GetCreatorName())

	desc := data.Policy.Description.Get()
	if desc != nil {
		state.Description = types.StringValue(*desc)
	}

	state.ModifiedAt = types.StringValue(data.Policy.ModifiedAt.Format(TimeFormat))
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
	plan.CreatedAt = types.StringValue(readData.Policy.CreatedAt.Format(TimeFormat))
	plan.CreatedBy = types.StringValue(readData.Policy.CreatedBy)
	plan.CreatorName = types.StringValue(readData.Policy.GetCreatorName())
	plan.ModifiedAt = types.StringValue(readData.Policy.ModifiedAt.Format(TimeFormat))
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
		"created_at":      types.StringValue(policy.CreatedAt.Format(TimeFormat)),
		"created_by":      types.StringValue(policy.CreatedBy),
		"creator_name":    types.StringValue(policy.GetCreatorName()),
		"description":     types.StringValue(descStr),
		"document":        documentValue,
		"id":              types.StringValue(policy.Id),
		"modified_at":     types.StringValue(policy.ModifiedAt.Format(TimeFormat)),
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
	statementsList, statementsListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)
	diags.Append(statementsListDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
	}

	documentValue, valueDiags := types.ObjectValue(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"statement": statementsList,
		"version":   types.StringValue(doc.GetVersion()),
	})
	diags.Append(valueDiags...)

	return documentValue, diags
}

func extractInputStatements(inputDoc *types.Object) []types.Object {
	result := make([]types.Object, 0)
	if inputDoc == nil || inputDoc.IsNull() || inputDoc.IsUnknown() {
		return result
	}
	attrs := inputDoc.Attributes()
	if attrs == nil {
		return result
	}
	stmtAttr, ok := attrs["statement"]
	if !ok {
		return result
	}
	stmtList, ok := stmtAttr.(types.List)
	if !ok {
		return result
	}
	for _, elem := range stmtList.Elements() {
		if stmtObj, ok := elem.(types.Object); ok {
			result = append(result, stmtObj)
		}
	}
	return result
}

func (r *serviceControlPolicyResource) buildDocumentValueFromResponse(ctx context.Context, doc *sdkorganization.ServiceControlPolicyDocument, inputDoc *types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if doc == nil {
		return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
	}

	inputStatements := extractInputStatements(inputDoc) // ← 헬퍼로 위임

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

	statementsList, statementsListDiags := types.ListValue(
		types.ObjectType{AttrTypes: organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)},
		statements,
	)
	diags.Append(statementsListDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx)), diags
	}
	documentValue, valueDiags := types.ObjectValue(
		organization.ServiceControlPolicyDocumentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"statement": statementsList,
			"version":   versionVal,
		},
	)
	diags.Append(valueDiags...)
	return documentValue, diags
}

func getInputAttr(inputStmt *types.Object, key string) (attr.Value, bool) {
	if inputStmt == nil || inputStmt.IsNull() || inputStmt.IsUnknown() {
		return nil, false
	}
	attrs := inputStmt.Attributes()
	if attrs == nil {
		return nil, false
	}
	v, ok := attrs[key]
	if !ok || v.IsNull() || v.IsUnknown() {
		return nil, false
	}
	return v, true
}

func buildStringListWithFallback(sdkValues []string, inputStmt *types.Object, key string) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	values := make([]attr.Value, 0, len(sdkValues))
	for _, v := range sdkValues {
		values = append(values, types.StringValue(v))
	}
	if len(values) > 0 {
		result, d := types.ListValue(types.StringType, values)
		diags.Append(d...)
		return result, diags
	}
	v, ok := getInputAttr(inputStmt, key)
	if !ok {
		result, d := types.ListValue(types.StringType, []attr.Value{})
		diags.Append(d...)
		return result, diags
	}
	list, ok := v.(types.List)
	if !ok {
		result, d := types.ListValue(types.StringType, []attr.Value{})
		diags.Append(d...)
		return result, diags
	}
	var resultValues []attr.Value
	for _, e := range list.Elements() {
		if s, ok := e.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
			resultValues = append(resultValues, s)
		}
	}
	resultList, d := types.ListValue(types.StringType, resultValues)
	diags.Append(d...)
	return resultList, diags
}

func (r *serviceControlPolicyResource) buildStatementValueFromResponse(ctx context.Context, stmt *sdkorganization.ServiceControlPolicyStatement, inputStmt *types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	actionsList, actionsDiags := buildStringListWithFallback(stmt.Action, inputStmt, "action")
	diags.Append(actionsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	notActionsList, notActionsDiags := buildStringListWithFallback(stmt.NotAction, inputStmt, "not_action")
	diags.Append(notActionsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	resourcesList, resourcesDiags := buildStringListWithFallback(stmt.Resource, inputStmt, "resource")
	diags.Append(resourcesDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	conditionValue, conditionDiags := r.buildConditionValue(ctx, stmt.Condition)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}
	if conditionValue.IsNull() {
		if v, ok := getInputAttr(inputStmt, "condition"); ok {
			conditionValue = v.(types.Map)
		}
	}

	principalValue := r.buildPrincipalValue(ctx, stmt.Principal)
	if principalValue.IsNull() {
		if v, ok := getInputAttr(inputStmt, "principal"); ok {
			principalValue = v.(types.String)
		}
	}

	effectVal := types.StringValue(stmt.GetEffect())
	if stmt.GetEffect() == "" {
		if v, ok := getInputAttr(inputStmt, "effect"); ok {
			effectVal = v.(types.String)
		}
	}

	sidVal := types.StringValue(stmt.GetSid())
	if stmt.GetSid() == "" {
		if v, ok := getInputAttr(inputStmt, "sid"); ok {
			sidVal = v.(types.String)
		}
	}

	statementValue, valueDiags := types.ObjectValue(
		organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"action":     actionsList,
			"condition":  conditionValue,
			"effect":     effectVal,
			"not_action": notActionsList,
			"principal":  principalValue,
			"resource":   resourcesList,
			"sid":        sidVal,
		},
	)
	diags.Append(valueDiags...)
	return statementValue, diags
}

func (r *serviceControlPolicyResource) buildStatementValue(ctx context.Context, stmt *sdkorganization.ServiceControlPolicyStatement) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	var actions []attr.Value
	for _, a := range stmt.Action {
		actions = append(actions, types.StringValue(a))
	}
	actionsList, actionsDiags := types.ListValue(types.StringType, actions)
	diags.Append(actionsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	var notActions []attr.Value
	for _, a := range stmt.NotAction {
		notActions = append(notActions, types.StringValue(a))
	}
	notActionsList, notActionsDiags := types.ListValue(types.StringType, notActions)
	diags.Append(notActionsDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	var resources []attr.Value
	for _, res := range stmt.Resource {
		resources = append(resources, types.StringValue(res))
	}
	resourcesList, resourcesDiags := types.ListValue(types.StringType, resources)
	diags.Append(resourcesDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}

	conditionValue, conditionDiags := r.buildConditionValue(ctx, stmt.Condition)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return types.ObjectNull(organization.ServiceControlPolicyStatementValue{}.AttributeTypes(ctx)), diags
	}
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

func (r *serviceControlPolicyResource) buildConditionValue(ctx context.Context, condition any) (types.Map, diag.Diagnostics) {
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

func waitForServiceControlPolicyReady(ctx context.Context, orgClient *organization.Client, policyId string, orgId string) error {
	return client.WaitForStatus(ctx, nil, []string{}, []string{"READY"}, func() (interface{}, string, error) {
		policyResp, err := orgClient.GetServiceControlPolicy(ctx, policyId, orgId)
		if err != nil {
			return nil, "", err
		}
		return policyResp, "READY", nil
	}, -1, -1, -1, -1)
}
