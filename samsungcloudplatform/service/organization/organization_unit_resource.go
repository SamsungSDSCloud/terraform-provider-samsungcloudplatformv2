package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v5/samsungcloudplatform/client/organization"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v5/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &organizationUnitResource{}
	_ resource.ResourceWithConfigure   = &organizationUnitResource{}
	_ resource.ResourceWithImportState = &organizationUnitResource{}
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
		Description: "Manages organizational units for grouping and structuring accounts within the organization",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Organization Unit Name. \n" +
					"  - example : 'score-organization-unit' \n",
				Required: true,
			},
			"description": schema.StringAttribute{
				Description: "Organization Unit Description. \n" +
					"  - example : 'Score Organization Unit' \n",
				Optional: true,
			},
			"organization_id": schema.StringAttribute{
				Description: "Unique identifier of the organization. \n" +
					"  - example : 'o-x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5' \n",
				Optional: true,
			},
			"parent_unit_id": schema.StringAttribute{
				Description: "Root or Parent Organization Unit ID. \n" +
					"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
				Optional: true,
				Computed: true,
			},
			"policy_ids": schema.ListAttribute{
				Description: "Control Policy List ID. \n" +
					"  - example : ['f98e76d54c32b10a9z8y7x6w5v4u3'] \n",
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
									Description: "Control Policy ID. \n" +
										"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
								},
								"policy_name": schema.StringAttribute{
									Computed: true,
									Description: "Control Policy Name. \n" +
										"  - example : 'test-policy-name' \n",
								},
							},
						},
						Computed: true,
						Description: "Service Control Policies. \n" +
							"  - example : '[{policy_id: f98e76d54c32b10a9z8y7x6w5v4u3, policy_name: test-policy-name}]' \n",
					},
					"created_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the organization unit was created. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"created_by": schema.StringAttribute{
						Computed: true,
						Description: "User who created the organization unit. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"creator_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the organization unit creator. \n" +
							"  - example : 'John Doe na' \n",
					},
					"depth": schema.Int64Attribute{
						Computed: true,
						Description: "Hierarchy (0~5). \n" +
							"  - example : 1 \n",
					},
					"description": schema.StringAttribute{
						Computed: true,
						Description: "Organization Unit Description. \n" +
							"  - example : 'Score Organization Unit' \n",
					},
					"id": schema.StringAttribute{
						Computed: true,
						Description: "Organization Unit ID. \n" +
							"  - example : 'ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
						Description: "Timestamp when the organization unit was modified. \n" +
							"  - example : '2025-01-01T00:00:00.000Z' \n",
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
						Description: "User who modified the organization unit. \n" +
							"  - example : 'c23fb561c689455993874fa5d5ed4a2f' \n",
					},
					"modifier_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the organization unit modifier. \n" +
							"  - example : 'Alice' \n",
					},
					"name": schema.StringAttribute{
						Computed: true,
						Description: "Organization Unit Name. \n" +
							"  - example : 'score-organization-unit' \n",
					},
					"parent_unit_id": schema.StringAttribute{
						Computed: true,
						Description: "Root or Parent Organization Unit ID. \n" +
							"  - example : 'ou-fc8c29a138d78e24bf1fa86812fc8b' \n",
					},
					"service_name": schema.StringAttribute{
						Computed: true,
						Description: "Name of the service to which the policy applies. \n" +
							"  - example : 'Organization' \n",
					},
					"srn": schema.StringAttribute{
						Computed: true,
						Description: "Samsung Resource Name (SRN) uniquely identifying this resource. \n" +
							"  - example : 'srn:dev2::1b8c29a138d78e24bf1fa86812fcaa:kr-west1::organizations/ou/ou-c29a138f8f1d78e24dbfa8681fc2fc8' \n",
					},
					"type": schema.StringAttribute{
						Computed: true,
						Description: "Type (Root or Organization Unit). \n" +
							"  - example : 'OU' \n",
					},
				},
				Computed: true,
				Description: "Organization Unit Info. \n" +
					"  - example : '{id: ou-c29a138f8f1d78e24dbfa8681fc2fc8, name: score-organization-unit, ...}' \n",
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

	unitId := unitInfo.Id
	err = waitForOrganizationUnitReady(ctx, r.client, unitId, orgId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error waiting for organization unit creation",
			"Error waiting for organization unit to be ready: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

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
		if !strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddError(
				"Error reading Organization Unit",
				"Could not read Organization Unit: "+err.Error(),
			)
			return
		}
		resp.State.RemoveResource(ctx)
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

func (r *organizationUnitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) < 1 {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Expected format: organization_unit_id or organization_unit_id:organization_id",
		)
		return
	}

	orgUnitId := parts[0]
	orgId := ""
	if len(parts) > 2 {
		orgId = parts[1]
	}

	orgUnitValue, diags := types.ObjectValue(organization.OrganizationUnitValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"control_policies": types.ListNull(types.ObjectType{AttrTypes: organization.ControlPoliciesValue{}.AttributeTypes(ctx)}),
		"created_at":       types.StringValue(""),
		"created_by":       types.StringValue(""),
		"creator_name":     types.StringValue(""),
		"depth":            types.Int64Value(0),
		"description":      types.StringValue(""),
		"id":               types.StringValue(orgUnitId),
		"modified_at":      types.StringValue(""),
		"modified_by":      types.StringValue(""),
		"modifier_name":    types.StringValue(""),
		"name":             types.StringValue(""),
		"parent_unit_id":   types.StringValue(""),
		"service_name":     types.StringValue(""),
		"srn":              types.StringValue(""),
		"type":             types.StringValue(""),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := organization.OrganizationUnitResource{
		OrganizationId:   types.StringValue(orgId),
		OrganizationUnit: orgUnitValue,
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func waitForOrganizationUnitReady(ctx context.Context, orgClient *organization.Client, unitId string, orgId string) error {
	return client.WaitForStatus(ctx, nil, []string{}, []string{"READY"}, func() (interface{}, string, error) {
		unitResp, err := orgClient.GetOrganizationUnit(ctx, unitId, orgId)
		if err != nil {
			return nil, "", err
		}
		return unitResp, "READY", nil
	}, -1, -1, -1, -1)
}
