package database

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type instanceGroupsPlanModifier struct{}

func InstanceGroupsPlanModifier() planmodifier.List {
	return instanceGroupsPlanModifier{}
}

func (m instanceGroupsPlanModifier) Description(_ context.Context) string {
	return "Copies computed instance-group/block-storage/instance fields from prior state onto identity-matched elements to avoid spurious plan diffs."
}

func (m instanceGroupsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m instanceGroupsPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	var planGroups []InstanceGroup
	if diags := req.PlanValue.ElementsAs(ctx, &planGroups, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	var stateGroups []InstanceGroup
	if diags := req.StateValue.ElementsAs(ctx, &stateGroups, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	stateByRole := make(map[string]InstanceGroup, len(stateGroups))
	for _, sg := range stateGroups {
		if _, exists := stateByRole[sg.RoleType.ValueString()]; !exists {
			stateByRole[sg.RoleType.ValueString()] = sg
		}
	}

	for i := range planGroups {
		stateGroup, ok := stateByRole[planGroups[i].RoleType.ValueString()]
		if !ok {
			continue
		}
		planGroups[i].Id = stateGroup.Id

		newBS, ok := ReconcileBlockStorages(ctx, planGroups[i].BlockStorageGroups, stateGroup.BlockStorageGroups, &resp.Diagnostics)
		if !ok {
			return
		}
		planGroups[i].BlockStorageGroups = newBS

		newIt, ok := ReconcileInstances(ctx, planGroups[i].Instances, stateGroup.Instances, &resp.Diagnostics)
		if !ok {
			return
		}
		planGroups[i].Instances = newIt
	}

	lst, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InstanceGroup{}.AttributeTypes()}, planGroups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = lst
}

func firstUnused(used []bool, n int, match func(j int) bool) int {
	for j := 0; j < n; j++ {
		if !used[j] && match(j) {
			return j
		}
	}
	return -1
}

func ReconcileBlockStorages(ctx context.Context, planList, stateList types.List, diagsOut *diag.Diagnostics) (types.List, bool) {
	if planList.IsNull() || planList.IsUnknown() {
		return planList, true
	}

	var planBS []BlockStorageGroup
	if diags := planList.ElementsAs(ctx, &planBS, false); diags.HasError() {
		diagsOut.Append(diags...)
		return planList, false
	}

	var stateBS []BlockStorageGroup
	if !stateList.IsNull() && !stateList.IsUnknown() {
		if diags := stateList.ElementsAs(ctx, &stateBS, false); diags.HasError() {
			diagsOut.Append(diags...)
			return planList, false
		}
	}

	used := make([]bool, len(stateBS))
	matched := make([]bool, len(planBS))

	adopt := func(i, j int) {
		planBS[i].Id = stateBS[j].Id
		planBS[i].Name = stateBS[j].Name
		used[j] = true
		matched[i] = true
	}

	for i := range planBS {
		j := firstUnused(used, len(stateBS), func(j int) bool {
			return planBS[i].RoleType.ValueString() == stateBS[j].RoleType.ValueString() &&
				planBS[i].VolumeType.ValueString() == stateBS[j].VolumeType.ValueString() &&
				planBS[i].SizeGb.ValueInt32() == stateBS[j].SizeGb.ValueInt32()
		})
		if j >= 0 {
			adopt(i, j)
		}
	}

	for i := range planBS {
		if matched[i] {
			continue
		}
		j := firstUnused(used, len(stateBS), func(j int) bool {
			return planBS[i].RoleType.ValueString() == stateBS[j].RoleType.ValueString() &&
				planBS[i].VolumeType.ValueString() == stateBS[j].VolumeType.ValueString()
		})
		if j >= 0 {
			adopt(i, j)
		}
	}

	lst, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BlockStorageGroup{}.AttributeTypes()}, planBS)
	if diags.HasError() {
		diagsOut.Append(diags...)
		return planList, false
	}
	return lst, true
}

func ReconcileInstances(ctx context.Context, planList, stateList types.List, diagsOut *diag.Diagnostics) (types.List, bool) {
	if planList.IsNull() || planList.IsUnknown() {
		return planList, true
	}

	var planIt []Instance
	if diags := planList.ElementsAs(ctx, &planIt, false); diags.HasError() {
		diagsOut.Append(diags...)
		return planList, false
	}

	var stateIt []Instance
	if !stateList.IsNull() && !stateList.IsUnknown() {
		if diags := stateList.ElementsAs(ctx, &stateIt, false); diags.HasError() {
			diagsOut.Append(diags...)
			return planList, false
		}
	}

	matched := reconcileInstanceMatches(planIt, stateIt)
	fillUnmatchedInstances(planIt, matched)

	lst, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Instance{}.AttributeTypes()}, planIt)
	if diags.HasError() {
		diagsOut.Append(diags...)
		return planList, false
	}
	return lst, true
}

func restoreInstance(dst *Instance, src Instance) {
	dst.Name = src.Name
	if dst.ServiceIpAddress.IsUnknown() {
		dst.ServiceIpAddress = src.ServiceIpAddress
	}
	if dst.PublicIpId.IsUnknown() {
		dst.PublicIpId = src.PublicIpId
	}
}

func reconcileInstanceMatches(planIt, stateIt []Instance) []bool {
	used := make([]bool, len(stateIt))
	matched := make([]bool, len(planIt))

	adopt := func(i, j int) {
		restoreInstance(&planIt[i], stateIt[j])
		used[j] = true
		matched[i] = true
	}

	for i := range planIt {
		if planIt[i].ServiceIpAddress.IsNull() || planIt[i].ServiceIpAddress.IsUnknown() {
			continue
		}
		j := firstUnused(used, len(stateIt), func(j int) bool {
			return planIt[i].RoleType.ValueString() == stateIt[j].RoleType.ValueString() &&
				planIt[i].ServiceIpAddress.ValueString() == stateIt[j].ServiceIpAddress.ValueString()
		})
		if j >= 0 {
			adopt(i, j)
		}
	}

	for i := range planIt {
		if matched[i] {
			continue
		}
		j := firstUnused(used, len(stateIt), func(j int) bool {
			return planIt[i].RoleType.ValueString() == stateIt[j].RoleType.ValueString()
		})
		if j >= 0 {
			adopt(i, j)
		}
	}

	return matched
}

func fillUnmatchedInstances(planIt []Instance, matched []bool) {
	for i := range planIt {
		if matched[i] {
			continue
		}
		if planIt[i].Name.IsNull() {
			planIt[i].Name = types.StringUnknown()
		}
		if planIt[i].ServiceIpAddress.IsNull() {
			planIt[i].ServiceIpAddress = types.StringUnknown()
		}
		if planIt[i].PublicIpId.IsNull() {
			planIt[i].PublicIpId = types.StringUnknown()
		}
	}
}
