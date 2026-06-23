package database

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BSKey struct {
	RoleType   string
	SizeGb     int32
	VolumeType string
}

func MapInstanceGroupsList(
	ctx context.Context,
	planInstanceGroups types.List,
	respInstanceGroups []InstanceGroupResponse,
) types.List {
	var instanceGroups []InstanceGroup

	igVals := make([]InstanceGroup, 0, len(planInstanceGroups.Elements()))
	planInstanceGroups.ElementsAs(ctx, &igVals, false)

	for _, instanceGroup := range igVals {
		cutIGs, updatedIG := MapInstanceGroups(ctx, respInstanceGroups, instanceGroup)
		instanceGroups = append(instanceGroups, updatedIG)
		respInstanceGroups = cutIGs
	}

	lst, _ := types.ListValueFrom(ctx,
		types.ObjectType{AttrTypes: InstanceGroup{}.AttributeTypes()},
		instanceGroups)
	return lst
}

func compareSlices[T comparable](actual, expect []T) bool {
	if len(actual) != len(expect) {
		return false
	}

	count := make(map[T]int, len(actual))
	for _, v := range actual {
		count[v]++
	}
	for _, v := range expect {
		count[v]--
		if count[v] == 0 {
			delete(count, v)
		}
	}
	return len(count) == 0
}

func CompareBlockStorages(actual, expect []BSKey) bool {
	return compareSlices(actual, expect)
}

func CompareInstances(actual, expect []string) bool {
	return compareSlices(actual, expect)
}

// CompareInstanceGroupKeys compares an instance group by its extracted keys.
func CompareInstanceGroupKeys(
	actualRoleType, actualServerTypeName string,
	expectRoleType, expectServerTypeName string,
	actualItKey, expectItKey []string,
	actualBSKey, expectBSKey []BSKey,
) bool {
	equal := actualRoleType == expectRoleType
	equal = equal && actualServerTypeName == expectServerTypeName
	equal = equal && CompareInstances(actualItKey, expectItKey)
	equal = equal && CompareBlockStorages(actualBSKey, expectBSKey)
	return equal
}

// MatchByDefOrder returns indices into resps that match each def in definition order.
// Each response item is consumed once (handles duplicates).
func MatchByDefOrder[T any, R any](defs []T, resps []R, match func(T, R) bool) []int {
	result := make([]int, 0, len(defs))
	used := make([]bool, len(resps))
	for _, d := range defs {
		for i, r := range resps {
			if used[i] {
				continue
			}
			if match(d, r) {
				used[i] = true
				result = append(result, i)
				break
			}
		}
	}
	return result
}

func MapInstanceGroups(ctx context.Context, instanceGroups []InstanceGroupResponse, def InstanceGroup) ([]InstanceGroupResponse, InstanceGroup) {

	for rm, instanceGroup := range instanceGroups {

		if !isEqualInstanceGroup(ctx, instanceGroup, def) {
			continue
		}

		defBs := make([]BlockStorageGroup, len(def.BlockStorageGroups.Elements()))
		def.BlockStorageGroups.ElementsAs(ctx, &defBs, false)
		defIt := make([]Instance, len(def.Instances.Elements()))
		def.Instances.ElementsAs(ctx, &defIt, false)

		bsList, itList := MapInstanceGroup(ctx, &MapInstanceGroupParams{
			DefBs:   defBs,
			DefIt:   defIt,
			BsResps: instanceGroup.BlockStorageGroups,
			ItResps: instanceGroup.Instances,
		})

		return append(instanceGroups[:rm], instanceGroups[rm+1:]...), InstanceGroup{
			Id:                 types.StringValue(instanceGroup.Id),
			BlockStorageGroups: bsList,
			Instances:          itList,
			RoleType:           types.StringValue(string(instanceGroup.RoleType)),
			ServerTypeName:     types.StringValue(instanceGroup.ServerTypeName),
		}

	}

	return instanceGroups, def
}

type MapInstanceGroupParams struct {
	DefBs   []BlockStorageGroup
	DefIt   []Instance
	BsResps []BlockStorageGroupResponse
	ItResps []InstanceResponse
}

// MapInstanceGroup maps response block storage groups and instances to terraform
// types, reordered to match definition order. Handles duplicates via consumed-index tracking.
func MapInstanceGroup(ctx context.Context, p *MapInstanceGroupParams) (types.List, types.List) {
	bsIndices := MatchByDefOrder(p.DefBs, p.BsResps, matchBlockStorage)
	bsVals := make([]BlockStorageGroup, 0, len(bsIndices))
	for _, idx := range bsIndices {
		bs := p.BsResps[idx]
		bsVals = append(bsVals, BlockStorageGroup{
			Id:         types.StringValue(bs.Id),
			Name:       types.StringValue(bs.Name),
			RoleType:   types.StringValue(string(bs.RoleType)),
			SizeGb:     types.Int32Value(bs.SizeGb),
			VolumeType: types.StringValue(string(bs.VolumeType)),
		})
	}
	bsList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BlockStorageGroup{}.AttributeTypes()}, bsVals)

	itIndices := MatchByDefOrder(p.DefIt, p.ItResps, matchInstance)
	itVals := make([]Instance, 0, len(itIndices))
	for _, idx := range itIndices {
		it := p.ItResps[idx]
		itVals = append(itVals, Instance{
			Name:             types.StringValue(it.Name),
			RoleType:         types.StringValue(string(it.RoleType)),
			ServiceIpAddress: types.StringPointerValue(&it.ServiceIpAddress),
			PublicIpId:       types.StringPointerValue(&it.PublicIpId),
		})
	}
	itList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Instance{}.AttributeTypes()}, itVals)

	return bsList, itList
}

// ToAnySlice converts a slice of any type to []any.
func ToAnySlice[T any](src []T) []any {
	out := make([]any, len(src))
	for i, v := range src {
		out[i] = v
	}
	return out
}

func matchBlockStorage(dv BlockStorageGroup, rv BlockStorageGroupResponse) bool {
	nameExists := !dv.Name.IsNull() && dv.Name.ValueString() != ""
	nameMatch := nameExists && rv.Name == dv.Name.ValueString()
	attrMatch := string(rv.RoleType) == dv.RoleType.ValueString() &&
		string(rv.VolumeType) == dv.VolumeType.ValueString() &&
		rv.SizeGb == dv.SizeGb.ValueInt32()
	return attrMatch && (!nameExists || nameMatch)
}

func matchInstance(dv Instance, rv InstanceResponse) bool {
	nameExists := !dv.Name.IsNull() && dv.Name.ValueString() != ""
	nameMatch := nameExists && rv.Name == dv.Name.ValueString()
	roleMatch := string(rv.RoleType) == dv.RoleType.ValueString()
	return roleMatch && (!nameExists || nameMatch)
}

func isEqualInstanceGroup(ctx context.Context, actual InstanceGroupResponse, expect InstanceGroup) bool {

	expectIt := make([]Instance, len(expect.Instances.Elements()))
	expect.Instances.ElementsAs(ctx, &expectIt, false)
	expectItKey := make([]string, len(expectIt))
	for i, it := range expectIt {
		expectItKey[i] = it.RoleType.ValueString()
	}
	actualItKey := make([]string, len(actual.Instances))
	for i, it := range actual.Instances {
		actualItKey[i] = string(it.RoleType)
	}

	expectBS := make([]BlockStorageGroup, len(expect.BlockStorageGroups.Elements()))
	expect.BlockStorageGroups.ElementsAs(ctx, &expectBS, false)
	expectBSKey := make([]BSKey, len(expectBS))
	for i, bs := range expectBS {
		expectBSKey[i] = BSKey{
			RoleType:   bs.RoleType.ValueString(),
			VolumeType: bs.VolumeType.ValueString(),
			SizeGb:     bs.SizeGb.ValueInt32(),
		}
	}
	actualBSKey := make([]BSKey, len(actual.BlockStorageGroups))
	for i, bs := range actual.BlockStorageGroups {
		actualBSKey[i] = BSKey{
			RoleType:   string(bs.RoleType),
			VolumeType: string(bs.VolumeType),
			SizeGb:     bs.SizeGb,
		}
	}

	return CompareInstanceGroupKeys(
		string(actual.RoleType), actual.ServerTypeName,
		expect.RoleType.ValueString(), expect.ServerTypeName.ValueString(),
		actualItKey, expectItKey,
		actualBSKey, expectBSKey,
	)

}
