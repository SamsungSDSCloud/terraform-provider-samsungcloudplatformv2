package cloudcontrol

import (
	"context"
	"fmt"
	"strings"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v6/samsungcloudplatform/client/cloudcontrol"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const guardrailBindingTargetType = "OU"

var (
	_ resource.Resource                = &guardrailBindingResource{}
	_ resource.ResourceWithConfigure   = &guardrailBindingResource{}
	_ resource.ResourceWithImportState = &guardrailBindingResource{}
)

func NewGuardrailBindingResource() resource.Resource {
	return &guardrailBindingResource{}
}

type guardrailBindingResource struct {
	config  *scpsdk.Configuration
	client  *cloudcontrol.Client
	clients *client.SCPClient
}

func (r *guardrailBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloudcontrol_guardrail_binding"
}

func (r *guardrailBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = inst.Client.CloudControl
	r.clients = inst.Client
}

func (r *guardrailBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Binds (enables) one or more guardrails to one or more organization units. \n" +
			"This is a batch operation over the cross product of guardrail_ids x unit_ids - the API " +
			"has no update endpoint, so any change to guardrail_ids, unit_ids, or landing_zone_id " +
			"replaces this resource (disable old combination, enable new one).",
		Attributes: map[string]schema.Attribute{
			"landing_zone_id": schema.StringAttribute{
				Optional: true,
				Description: "Landing Zone ID that contains the organization units the guardrails are bound to. \n" +
					"  - example : '2c8a138f8d78e1fc29a449dbfa8681' \n",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"guardrail_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Guardrail IDs to bind. \n" +
					"  - example : ['f98e76d54c32b10a9z8y7x6w5v4u3'] \n",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"unit_ids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Organization Unit IDs to bind the guardrails to. \n" +
					"  - example : ['ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5'] \n",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"success_ids": schema.ListNestedAttribute{
				Computed: true,
				Description: "Guardrail/unit pairs that were successfully bound. \n" +
					"  - example : '[{guardrail_id: f98e76d54c32b10a9z8y7x6w5v4u3, unit_id: ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"guardrail_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the guardrail. \n" +
								"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
						},
						"unit_id": schema.StringAttribute{
							Computed: true,
							Description: "Organization Unit ID. \n" +
								"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
						},
					},
				},
			},
			"failed_ids": schema.ListNestedAttribute{
				Computed: true,
				Description: "Guardrail/unit pairs that failed to bind. \n" +
					"  - example : '[{guardrail_id: f98e76d54c32b10a9z8y7x6w5v4u3, unit_id: ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5, ...}]' \n",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"guardrail_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the guardrail. \n" +
								"  - example : 'f98e76d54c32b10a9z8y7x6w5v4u3' \n",
						},
						"unit_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier of the organization unit the guardrail was to be bound to. \n" +
								"  - example : 'ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5' \n",
						},
						"error_code": schema.StringAttribute{
							Computed: true,
							Description: "Error code returned when the binding failed. \n" +
								"  - example : 'CloudControl.AlreadyEnabledGuardrails' \n",
						},
						"failed_caused": schema.StringAttribute{
							Computed: true,
							Description: "Failure reason. \n" +
								"  - example : 'guardrail already enabled' \n",
						},
						"response": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Raw response payload for the failure. \n" +
								"  - example : '{}' \n",
						},
					},
				},
			},
		},
	}
}

func toStringSlice(ctx context.Context, l types.List) ([]string, diag.Diagnostics) {
	out := make([]string, 0)
	diags := l.ElementsAs(ctx, &out, false)
	return out, diags
}

func buildGuardrailBindingResponseMap(response map[string]interface{}) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	m := make(map[string]attr.Value, len(response))
	for k, val := range response {
		m[k] = types.StringValue(fmt.Sprintf("%v", val))
	}

	v, d := types.MapValue(types.StringType, m)
	diags.Append(d...)
	if diags.HasError() {
		return types.MapNull(types.StringType), diags
	}
	return v, diags
}

func (r *guardrailBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudcontrol.GuardrailBindingResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	guardrailIds, gDiags := toStringSlice(ctx, plan.GuardrailIds)
	resp.Diagnostics.Append(gDiags...)
	unitIds, uDiags := toStringSlice(ctx, plan.UnitIds)
	resp.Diagnostics.Append(uDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	landingZoneId := plan.LandingZoneId.ValueString()

	data, err := r.client.EnableGuardrailBindings(ctx, guardrailIds, unitIds, landingZoneId)
	if err != nil {
		detail := client.GetDetailFromError(err)
		resp.Diagnostics.AddError(
			"Error enabling Guardrail Bindings",
			"Could not enable Guardrail Bindings, unexpected error: "+err.Error()+"\nReason: "+detail,
		)
		return
	}

	successItems := make([]cloudcontrol.GuardrailBindingSuccessItem, 0)
	if data != nil {
		for _, s := range data.SuccessIds {
			successItems = append(successItems, cloudcontrol.GuardrailBindingSuccessItem{
				GuardrailId: types.StringValue(s.GuardrailId),
				UnitId:      types.StringValue(s.UnitId),
			})
		}
	}
	successList, successListDiags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: cloudcontrol.GuardrailBindingSuccessItem{}.AttributeTypes(ctx)},
		successItems,
	)
	resp.Diagnostics.Append(successListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.SuccessIds = successList

	failedItems := make([]cloudcontrol.GuardrailBindingFailedItem, 0)
	var failedMsg strings.Builder
	if data != nil {
		for _, f := range data.FailedIds {
			respMap, mapDiags := buildGuardrailBindingResponseMap(f.Response)
			resp.Diagnostics.Append(mapDiags...)
			if resp.Diagnostics.HasError() {
				return
			}

			failedItems = append(failedItems, cloudcontrol.GuardrailBindingFailedItem{
				GuardrailId:  types.StringValue(f.GuardrailId),
				UnitId:       types.StringValue(f.UnitId),
				ErrorCode:    types.StringValue(f.ErrorCode),
				FailedCaused: types.StringValue(f.FailedCaused),
				Response:     respMap,
			})
			failedMsg.WriteString(fmt.Sprintf("\n- guardrail=%s unit=%s: %s (%s)", f.GuardrailId, f.UnitId, f.FailedCaused, f.ErrorCode))
		}
	}
	failedList, failedListDiags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: cloudcontrol.GuardrailBindingFailedItem{}.AttributeTypes(ctx)},
		failedItems,
	)
	resp.Diagnostics.Append(failedListDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.FailedIds = failedList

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if failedMsg.Len() > 0 {
		resp.Diagnostics.AddWarning(
			"Some guardrail bindings failed",
			"The following guardrail/unit combinations failed to bind:"+failedMsg.String()+
				"\n\nSee the failed_ids attribute for details.",
		)
	}
}

func (r *guardrailBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudcontrol.GuardrailBindingResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var successItems []cloudcontrol.GuardrailBindingSuccessItem
	diags = state.SuccessIds.ElementsAs(ctx, &successItems, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(successItems) == 0 {
		diags = resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		return
	}

	landingZoneId := state.LandingZoneId.ValueString()

	boundUnitsByGuardrail := map[string]map[string]bool{}
	verificationSucceeded := true

	for _, pair := range successItems {
		guardrailId := pair.GuardrailId.ValueString()

		if _, ok := boundUnitsByGuardrail[guardrailId]; !ok {
			set, err := buildBoundUnitsSet(ctx, r.client, guardrailId, landingZoneId)
			if err != nil {
				resp.Diagnostics.AddWarning(
					"Could not verify guardrail bindings",
					"Failed to list targets for guardrail "+guardrailId+": "+err.Error()+
						"\nExisting state has been kept as-is and may be out of date.",
				)
				verificationSucceeded = false
				break
			}
			boundUnitsByGuardrail[guardrailId] = set
		}
	}

	if !verificationSucceeded {
		diags = resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		return
	}

	stillBound := make([]cloudcontrol.GuardrailBindingSuccessItem, 0)
	for _, pair := range successItems {
		guardrailId := pair.GuardrailId.ValueString()
		unitId := pair.UnitId.ValueString()
		if boundUnitsByGuardrail[guardrailId][unitId] {
			stillBound = append(stillBound, pair)
		}
	}

	if len(stillBound) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	stillBoundList, listDiags := types.ListValueFrom(
		ctx,
		types.ObjectType{AttrTypes: cloudcontrol.GuardrailBindingSuccessItem{}.AttributeTypes(ctx)},
		stillBound,
	)
	resp.Diagnostics.Append(listDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.SuccessIds = stillBoundList

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func buildBoundUnitsSet(ctx context.Context, client *cloudcontrol.Client, guardrailId, landingZoneId string) (map[string]bool, error) {
	targets, err := client.ListTargetsForGuardrail(ctx, guardrailId, guardrailBindingTargetType, landingZoneId, "", "", 0, 0)
	if err != nil {
		return nil, err
	}

	set := map[string]bool{}
	if targets != nil {
		if targets.Targets.ArrayOfOrganizationUnitsForGuardrail != nil {
			for _, t := range *targets.Targets.ArrayOfOrganizationUnitsForGuardrail {
				set[t.Id] = true
			}
		}
		if targets.Targets.ArrayOfAccountsForGuardrail != nil {
			for _, t := range *targets.Targets.ArrayOfAccountsForGuardrail {
				set[t.Id] = true
			}
		}
	}
	return set, nil
}

func (r *guardrailBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Unsupported Operation",
		"Guardrail Binding does not support in-place update. This should be unreachable since all "+
			"identifying attributes require replacement - if you see this error, please report it.",
	)
}

func (r *guardrailBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cloudcontrol.GuardrailBindingResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var successItems []cloudcontrol.GuardrailBindingSuccessItem
	diags = state.SuccessIds.ElementsAs(ctx, &successItems, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(successItems) == 0 {
		return
	}

	unitsByGuardrail := map[string][]string{}
	seen := map[string]bool{}
	for _, pair := range successItems {
		gid := pair.GuardrailId.ValueString()
		uid := pair.UnitId.ValueString()
		if key := gid + "|" + uid; !seen[key] {
			seen[key] = true
			unitsByGuardrail[gid] = append(unitsByGuardrail[gid], uid)
		}
	}

	landingZoneId := state.LandingZoneId.ValueString()

	var realFailures strings.Builder
	var callErrors strings.Builder

	for gid, units := range unitsByGuardrail {
		data, err := r.client.DisableGuardrailBindings(ctx, []string{gid}, units, landingZoneId)
		if err != nil {
			detail := client.GetDetailFromError(err)
			callErrors.WriteString(fmt.Sprintf("\n- guardrail=%s: %s (%s)", gid, err.Error(), detail))
			continue
		}

		if data != nil {
			for _, f := range data.FailedIds {
				realFailures.WriteString(fmt.Sprintf("\n- guardrail=%s unit=%s: %s (%s)",
					f.GuardrailId, f.UnitId, f.FailedCaused, f.ErrorCode))
			}
		}
	}

	if callErrors.Len() > 0 {
		resp.Diagnostics.AddError(
			"Error disabling Guardrail Bindings",
			"Could not disable Guardrail Bindings for the following guardrails:"+callErrors.String()+
				"\n\nThe resource has been kept in state so the operation can be retried. "+
				"Bindings for other guardrails may already have been removed.",
		)
		return
	}

	if realFailures.Len() > 0 {
		resp.Diagnostics.AddWarning(
			"Some guardrail bindings failed to disable",
			"The following guardrail/unit combinations failed to unbind:"+realFailures.String()+
				"\n\nThe resource has been removed from Terraform state regardless. "+
				"Verify the actual binding state if needed.",
		)
	}
}

func splitAndTrimIds(s string) []string {
	out := make([]string, 0)
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (r *guardrailBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected format '<landing_zone_id>/<guardrail_ids>/<unit_ids>' with comma-separated IDs, got: "+req.ID+
				"\n\nExample: 2c8a138f8d78e1fc29a449dbfa8681/f98e76d54c32b10a9z8y7x6w5v4u3/ou-1a2b3c4d5e6f7g8h9i0j1k2l3m4n5",
		)
		return
	}

	landingZoneId := strings.TrimSpace(parts[0])
	guardrailIds := splitAndTrimIds(parts[1])
	unitIds := splitAndTrimIds(parts[2])

	if len(guardrailIds) == 0 || len(unitIds) == 0 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Both the guardrail ID list and the unit ID list must contain at least one non-empty ID. Got: "+req.ID,
		)
		return
	}

	successItems := make([]cloudcontrol.GuardrailBindingSuccessItem, 0, len(guardrailIds)*len(unitIds))
	for _, g := range guardrailIds {
		for _, u := range unitIds {
			successItems = append(successItems, cloudcontrol.GuardrailBindingSuccessItem{
				GuardrailId: types.StringValue(g),
				UnitId:      types.StringValue(u),
			})
		}
	}

	successList, diags := types.ListValueFrom(ctx,
		types.ObjectType{AttrTypes: cloudcontrol.GuardrailBindingSuccessItem{}.AttributeTypes(ctx)},
		successItems)
	resp.Diagnostics.Append(diags...)

	guardrailList, diags := types.ListValueFrom(ctx, types.StringType, guardrailIds)
	resp.Diagnostics.Append(diags...)

	unitList, diags := types.ListValueFrom(ctx, types.StringType, unitIds)
	resp.Diagnostics.Append(diags...)

	failedList, diags := types.ListValueFrom(ctx,
		types.ObjectType{AttrTypes: cloudcontrol.GuardrailBindingFailedItem{}.AttributeTypes(ctx)},
		[]cloudcontrol.GuardrailBindingFailedItem{})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	lz := types.StringNull()
	if landingZoneId != "" {
		lz = types.StringValue(landingZoneId)
	}

	state := cloudcontrol.GuardrailBindingResource{
		LandingZoneId: lz,
		GuardrailIds:  guardrailList,
		UnitIds:       unitList,
		SuccessIds:    successList,
		FailedIds:     failedList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
