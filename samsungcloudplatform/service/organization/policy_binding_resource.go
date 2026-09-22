package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	sdkorganization "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/organization/1.3"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &policyBindingResource{}
	_ resource.ResourceWithConfigure   = &policyBindingResource{}
	_ resource.ResourceWithImportState = &policyBindingResource{}
)

func NewPolicyBindingResource() resource.Resource {
	return &policyBindingResource{}
}

type policyBindingResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *policyBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_policy_binding"
}

func (r *policyBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attaches a set of control policies to a single target within the organization (entity type is fixed to TARGET). Use this resource to bind one or more policies to a target such as an account, organization unit, or the root.",
		Attributes: map[string]schema.Attribute{
			"entity": schema.StringAttribute{
				Description: "Entity type for the binding. \n" +
					"  - allowed_values : ['POLICY', 'TARGET'] \n" +
					"  - example : 'POLICY' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization that owns the target. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_ids": schema.SetAttribute{
				Description: "Set of target IDs to bind policies to. \n" +
					"  - example : ['ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5'] \n",
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"policy_ids": schema.SetAttribute{
				Description: "Set of control policy IDs to attach to the target. \n" +
					"  - example : ['f98e76d54c32b10a9z8y7x6w5v4u3'] \n",
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"success_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of IDs that were successfully bound. \n" +
					"  - example : ['ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5'] \n",
			},
			"failed_ids": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of IDs that failed to bind. \n" +
					"  - example : '[{failed_id: ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, error_code: Conflict, failed_caused: Account 482d111a302547fbbd11ee0141ea23bb is not found}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"failed_id": schema.StringAttribute{
							Computed: true,
							Description: "Failed target or policy ID. \n" +
								"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
						},
						"error_code": schema.StringAttribute{
							Computed: true,
							Description: "Error code returned when the operation fails. \n" +
								"  - example : 'Conflict' \n",
						},
						"failed_caused": schema.StringAttribute{
							Computed: true,
							Description: "Cause of failure. \n" +
								"  - example : 'Account 482d111a302547fbbd11ee0141ea23bb is not found' \n",
						},
						"response": schema.StringAttribute{
							Computed: true,
							Description: "Response body from the failed request. \n" +
								"  - example : '{}' \n",
						},
					},
				},
			},
		},
	}
}

func (r *policyBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.PolicyBindingResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	targetIds := stringSetToSlice(plan.TargetIds)
	policyIds := stringSetToSlice(plan.PolicyIds)
	entity := sdkorganization.AssignmentEntityType(plan.Entity.ValueString())
	result, err := r.client.AttachPolicyBindings(ctx, entity, policyIds, targetIds, organizationId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError("Error creating Policy Binding", err.Error()+"\nReason: "+detail)
		return
	}

	successIds, diags := types.ListValueFrom(ctx, types.StringType, result.SuccessIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	failedIdValues := make([]attr.Value, 0, len(result.FailedIds))
	for _, f := range result.FailedIds {
		obj, objDiags := types.ObjectValue(organization.FailedBindingValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"failed_id":     types.StringValue(f.FailedId),
			"error_code":    types.StringValue(f.ErrorCode),
			"failed_caused": types.StringValue(f.FailedCaused),
			"response":      types.StringValue(fmt.Sprintf("%v", f.Response)),
		})
		resp.Diagnostics.Append(objDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		failedIdValues = append(failedIdValues, obj)
	}

	failedIds, diags := types.ListValue(
		types.ObjectType{AttrTypes: organization.FailedBindingValue{}.AttributeTypes(ctx)},
		failedIdValues,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(result.FailedIds) > 0 {
		for _, f := range result.FailedIds {
			resp.Diagnostics.AddWarning(
				"Policy Binding partially failed",
				fmt.Sprintf("id=%s error=%s cause=%s", f.FailedId, f.ErrorCode, f.FailedCaused),
			)
		}
	}

	plan.SuccessIds = successIds
	plan.FailedIds = failedIds
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.PolicyBindingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := state.OrganizationId.ValueString()
	policyIds := stringSetToSlice(state.PolicyIds)
	targetIds := stringSetToSlice(state.TargetIds)

	for _, targetId := range targetIds {
		const pageSize = 100
		allBound := false

		for page := 0; page < 100; page++ {
			data, err := r.client.GetPoliciesForTarget(ctx, targetId, orgId, "", "", int32(page), pageSize, "")
			if err != nil {
				detail := client.GetDetailFromError(err)
				resp.Diagnostics.AddError("Error reading Policy Binding", err.Error()+"\nReason: "+detail)
				return
			}
			if data == nil || len(data.Policies) == 0 {
				break
			}
			if allPoliciesBound(data.Policies, policyIds) {
				allBound = true
				break
			}
			if len(data.Policies) < pageSize {
				break
			}
		}

		if !allBound {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *policyBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// empty left.
}

func (r *policyBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.PolicyBindingResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	policyIds := stringSetToSlice(state.PolicyIds)
	targetIds := stringSetToSlice(state.TargetIds)
	entity := sdkorganization.AssignmentEntityType(state.Entity.ValueString())

	result, err := r.client.RemovePolicyBindings(ctx, entity, policyIds, targetIds, organizationId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError("Error deleting Policy Binding", err.Error()+"\nReason: "+detail)
		return
	}
	if result != nil && len(result.FailedIds) > 0 {
		for _, f := range result.FailedIds {
			resp.Diagnostics.AddError(
				"Policy Binding removal partially failed",
				fmt.Sprintf("id=%s error=%s cause=%s", f.FailedId, f.ErrorCode, f.FailedCaused),
			)
		}
	}
}

func (r *policyBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	const format = "organization_id:entity:policy_id1,policy_id2,...:target_id1,target_id2,..."

	parts := strings.SplitN(req.ID, ":", 4)
	if len(parts) != 4 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: "+format+"\nGot: "+req.ID)
		return
	}

	orgId := strings.TrimSpace(parts[0])
	entity := strings.TrimSpace(parts[1])
	if orgId == "" || entity == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"organization_id and entity must not be empty.\nExpected format: "+format)
		return
	}

	policyIds := splitTrimFilter(parts[2])
	if len(policyIds) == 0 {
		resp.Diagnostics.AddError("Invalid import ID",
			"At least one policy_id is required.\nExpected format: "+format)
		return
	}

	targetIds := splitTrimFilter(parts[3])
	if len(targetIds) == 0 {
		resp.Diagnostics.AddError("Invalid import ID",
			"At least one target_id is required.\nExpected format: "+format)
		return
	}

	policyIdsValue, diags := types.SetValueFrom(ctx, types.StringType, policyIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	targetIdsValue, diags := types.SetValueFrom(ctx, types.StringType, targetIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, organization.PolicyBindingResource{
		OrganizationId: types.StringValue(orgId),
		Entity:         types.StringValue(entity),
		PolicyIds:      policyIdsValue,
		TargetIds:      targetIdsValue,
		SuccessIds:     types.ListValueMust(types.StringType, []attr.Value{}),
		FailedIds: types.ListValueMust(
			types.ObjectType{AttrTypes: organization.FailedBindingValue{}.AttributeTypes(ctx)},
			[]attr.Value{},
		),
	})...)
}

func nullableStringValue(ns sdkorganization.NullableString) string {
	if ns.IsSet() && ns.Get() != nil {
		return *ns.Get()
	}
	return ""
}

func allPoliciesBound(policies []sdkorganization.PoliciesForTargetSummary, wanted []string) bool {
	present := make(map[string]bool, len(policies))
	for _, p := range policies {
		present[p.Id] = true
	}
	for _, id := range wanted {
		if !present[id] {
			return false
		}
	}
	return true
}

func stringSetToSlice(set types.Set) []string {
	elems := set.Elements()
	result := make([]string, 0, len(elems))
	for _, e := range elems {
		if s, ok := e.(types.String); ok {
			result = append(result, s.ValueString())
		}
	}
	return result
}

func splitTrimFilter(s string) []string {
	raw := strings.Split(s, ",")
	result := make([]string, 0, len(raw))
	for _, r := range raw {
		if v := strings.TrimSpace(r); v != "" {
			result = append(result, v)
		}
	}
	return result
}
