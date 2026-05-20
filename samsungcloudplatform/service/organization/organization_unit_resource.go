package organization

import (
	"context"
	"fmt"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v3/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &organizationUnitResource{}
	_ resource.ResourceWithConfigure = &organizationUnitResource{}
)

func NewOrganizationUnitResource() resource.Resource {
	return &organizationUnitResource{}
}

type organizationUnitResource struct {
	config  *scpsdk.Configuration
	client  *organization.Client
	clients *client.SCPClient
}

func (r *organizationUnitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_unit"
}

func (r *organizationUnitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Organization Unit",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Organization Unit Name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Organization Unit Description",
				Optional:    true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Organization ID",
				Optional:    true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Parent Organization Unit ID (required for non-root units)",
				Optional:    true,
				Computed:    true,
			},
			"policy_ids": schema.ListAttribute{
				Description: "Policy IDs",
				Optional:    true,
				ElementType: types.StringType,
			},
			"organization_unit": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"control_policies": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"policy_id": schema.StringAttribute{
									Computed: true,
								},
								"policy_name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						Computed: true,
					},
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Computed: true,
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
					},
					"depth": schema.Int64Attribute{
						Computed: true,
					},
					"description": schema.StringAttribute{
						Computed: true,
					},
					"id": schema.StringAttribute{
						Computed: true,
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
					},
					"modifier_name": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"parent_unit_id": schema.StringAttribute{
						Computed: true,
					},
					"service_name": schema.StringAttribute{
						Computed: true,
					},
					"srn": schema.StringAttribute{
						Computed: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
				Computed: true,
			},
		},
	}
}

func (r *organizationUnitResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *organizationUnitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organization.OrganizationUnitResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgId := plan.OrganizationId.ValueString()
	if orgId == "" {
		orgId = "default"
	}
	plan.OrganizationId = types.StringValue(orgId)

	data, err := r.client.CreateOrganizationUnit(ctx, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error creating Organization Unit",
			"Could not create Organization Unit, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	unitInfo := data.OrganizationUnit

	plan.Name = types.StringValue(unitInfo.Name)
	plan.Description = types.StringValue(unitInfo.GetDescription())
	if unitInfo.ParentUnitId.IsSet() {
		plan.ParentUnitId = types.StringValue(unitInfo.GetParentUnitId())
	}

	controlPolicies := make([]attr.Value, 0, len(unitInfo.GetControlPolicies()))

	for _, policy := range unitInfo.GetControlPolicies() {
		policyValue, d := types.ObjectValue(
			organization.ControlPoliciesValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"policy_id":   types.StringValue(policy.PolicyId),
				"policy_name": types.StringValue(policy.PolicyName),
			},
		)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		controlPolicies = append(controlPolicies, policyValue)
	}

	controlPoliciesList, d := types.ListValue(
		types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)},
		controlPolicies,
	)

	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitValue, diags := types.ObjectValue(organization.OrganizationUnitValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"control_policies": controlPoliciesList,
		"created_at":       types.StringValue(unitInfo.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":       types.StringValue(unitInfo.CreatedBy),
		"creator_name":     types.StringValue(unitInfo.GetCreatorName()),
		"depth":            types.Int64Value(int64(unitInfo.Depth)),
		"description":      types.StringValue(unitInfo.GetDescription()),
		"id":               types.StringValue(unitInfo.Id),
		"modified_at":      types.StringValue(unitInfo.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":      types.StringValue(unitInfo.ModifiedBy),
		"modifier_name":    types.StringValue(unitInfo.GetModifierName()),
		"name":             types.StringValue(unitInfo.Name),
		"parent_unit_id":   types.StringValue(unitInfo.GetParentUnitId()),
		"service_name":     types.StringValue(unitInfo.ServiceName),
		"srn":              types.StringValue(unitInfo.GetSrn()),
		"type":             types.StringValue(unitInfo.Type),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.OrganizationUnit = orgUnitValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *organizationUnitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organization.OrganizationUnitResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationUnit.IsNull() || state.OrganizationUnit.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Unit",
			"Organization unit is null or unknown",
		)
		return
	}

	orgUnitAttrs := state.OrganizationUnit.Attributes()
	orgUnitIdVal := orgUnitAttrs["id"]
	if orgUnitIdVal == nil || orgUnitIdVal.IsNull() || orgUnitIdVal.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Unit",
			"Organization unit ID not found in state",
		)
		return
	}
	orgUnitId := orgUnitIdVal.(types.String).ValueString()

	orgId := state.OrganizationId.ValueString()
	if orgId == "" {
		orgId = "default"
	}

	data, err := r.client.GetOrganizationUnit(ctx, orgUnitId, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Organization Unit",
			err.Error(),
		)
		return
	}

	unitInfo := data.OrganizationUnit

	state.Name = types.StringValue(unitInfo.Name)
	state.Description = types.StringValue(unitInfo.GetDescription())
	if unitInfo.ParentUnitId.IsSet() {
		state.ParentUnitId = types.StringValue(unitInfo.GetParentUnitId())
	}

	orgUnitValue, diags := types.ObjectValue(organization.OrganizationUnitValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"control_policies": types.ListNull(types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)}),
		"created_at":       types.StringValue(unitInfo.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":       types.StringValue(unitInfo.CreatedBy),
		"creator_name":     types.StringValue(unitInfo.GetCreatorName()),
		"depth":            types.Int64Value(int64(unitInfo.Depth)),
		"description":      types.StringValue(unitInfo.GetDescription()),
		"id":               types.StringValue(unitInfo.Id),
		"modified_at":      types.StringValue(unitInfo.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":      types.StringValue(unitInfo.ModifiedBy),
		"modifier_name":    types.StringValue(unitInfo.GetModifierName()),
		"name":             types.StringValue(unitInfo.Name),
		"parent_unit_id":   types.StringValue(unitInfo.GetParentUnitId()),
		"service_name":     types.StringValue(unitInfo.ServiceName),
		"srn":              types.StringValue(unitInfo.GetSrn()),
		"type":             types.StringValue(unitInfo.Type),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.OrganizationUnit = orgUnitValue

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *organizationUnitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state organization.OrganizationUnitResource
	var plan organization.OrganizationUnitResource

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitId := state.OrganizationUnit.Attributes()["id"].(types.String).ValueString()
	orgId := plan.OrganizationId.ValueString()
	if orgId == "" {
		orgId = "default"
	}

	_, err := r.client.UpdateOrganizationUnit(ctx, orgUnitId, plan)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error updating Organization Unit",
			"Could not update Organization Unit, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	data, err := r.client.GetOrganizationUnit(ctx, orgUnitId, orgId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unable to Read Organization Unit",
			"Could not read Organization Unit ID "+orgUnitId+": "+err.Error(),
		)
		return
	}

	unitInfo := data.OrganizationUnit

	plan.Name = types.StringValue(unitInfo.Name)
	plan.Description = types.StringValue(unitInfo.GetDescription())
	if unitInfo.ParentUnitId.IsSet() {
		plan.ParentUnitId = types.StringValue(unitInfo.GetParentUnitId())
	}

	orgUnitValue, diags := types.ObjectValue(organization.OrganizationUnitValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"control_policies": types.ListNull(types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)}),
		"created_at":       types.StringValue(unitInfo.CreatedAt.Format("2006-01-02T15:04:05.000Z")),
		"created_by":       types.StringValue(unitInfo.CreatedBy),
		"creator_name":     types.StringValue(unitInfo.GetCreatorName()),
		"depth":            types.Int64Value(int64(unitInfo.Depth)),
		"description":      types.StringValue(unitInfo.GetDescription()),
		"id":               types.StringValue(unitInfo.Id),
		"modified_at":      types.StringValue(unitInfo.ModifiedAt.Format("2006-01-02T15:04:05.000Z")),
		"modified_by":      types.StringValue(unitInfo.ModifiedBy),
		"modifier_name":    types.StringValue(unitInfo.GetModifierName()),
		"name":             types.StringValue(unitInfo.Name),
		"parent_unit_id":   types.StringValue(unitInfo.GetParentUnitId()),
		"service_name":     types.StringValue(unitInfo.ServiceName),
		"srn":              types.StringValue(unitInfo.GetSrn()),
		"type":             types.StringValue(unitInfo.Type),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.OrganizationUnit = orgUnitValue

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *organizationUnitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organization.OrganizationUnitResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUnitId := state.OrganizationUnit.Attributes()["id"].(types.String).ValueString()
	orgId := state.OrganizationId.ValueString()
	if orgId == "" {
		orgId = "default"
	}
	_, err := r.client.DeleteOrganizationUnit(ctx, orgUnitId, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error deleting Organization Unit",
			"Could not delete Organization Unit, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}
}
