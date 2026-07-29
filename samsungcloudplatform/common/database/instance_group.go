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
	// Import/최초 Read: plan이 비어있으면 응답 기반 매핑
	if len(planInstanceGroups.Elements()) == 0 {
		return MapInstanceGroupsFromResponse(ctx, respInstanceGroups)
	}

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

func MapInstanceGroupsFromResponse(
	ctx context.Context,
	respInstanceGroups []InstanceGroupResponse,
) types.List {
	instanceGroups := make([]InstanceGroup, 0, len(respInstanceGroups))

	for _, ig := range respInstanceGroups {
		bsList, itList := mapInstanceGroupFromResponse(ctx, ig)
		instanceGroups = append(instanceGroups, InstanceGroup{
			Id:                 types.StringValue(ig.Id),
			BlockStorageGroups: bsList,
			Instances:          itList,
			RoleType:           types.StringValue(ig.RoleType),
			ServerTypeName:     types.StringValue(ig.ServerTypeName),
		})
	}

	lst, _ := types.ListValueFrom(ctx,
		types.ObjectType{AttrTypes: InstanceGroup{}.AttributeTypes()},
		instanceGroups)
	return lst
}

func mapInstanceGroupFromResponse(ctx context.Context, ig InstanceGroupResponse) (types.List, types.List) {
	bsVals := make([]BlockStorageGroup, 0, len(ig.BlockStorageGroups))
	for _, bs := range ig.BlockStorageGroups {
		bsVals = append(bsVals, BlockStorageGroup{
			Id:         types.StringValue(bs.Id),
			Name:       types.StringValue(bs.Name),
			RoleType:   types.StringValue(bs.RoleType),
			SizeGb:     types.Int32Value(bs.SizeGb),
			VolumeType: types.StringValue(bs.VolumeType),
		})
	}
	bsList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BlockStorageGroup{}.AttributeTypes()}, bsVals)

	itVals := make([]Instance, 0, len(ig.Instances))
	for _, it := range ig.Instances {
		itVals = append(itVals, Instance{
			Name:             types.StringValue(it.Name),
			RoleType:         types.StringValue(it.RoleType),
			ServiceIpAddress: types.StringPointerValue(&it.ServiceIpAddress),
			PublicIpId:       types.StringPointerValue(&it.PublicIpId),
		})
	}
	itList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Instance{}.AttributeTypes()}, itVals)

	return bsList, itList
}

func CompareSlices[T comparable](actual, expect []T) bool {
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
	return CompareSlices(actual, expect)
}

func CompareInstances(actual, expect []string) bool {
	return CompareSlices(actual, expect)
}

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

		if !IsEqualInstanceGroup(ctx, instanceGroup, def) {
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
			RoleType:           types.StringValue(instanceGroup.RoleType),
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

func MapInstanceGroup(ctx context.Context, p *MapInstanceGroupParams) (types.List, types.List) {
	bsIndices := MatchByDefOrder(p.DefBs, p.BsResps, MatchBlockStorage)
	bsVals := make([]BlockStorageGroup, 0, len(bsIndices))
	for _, idx := range bsIndices {
		bs := p.BsResps[idx]
		bsVals = append(bsVals, BlockStorageGroup{
			Id:         types.StringValue(bs.Id),
			Name:       types.StringValue(bs.Name),
			RoleType:   types.StringValue(bs.RoleType),
			SizeGb:     types.Int32Value(bs.SizeGb),
			VolumeType: types.StringValue(bs.VolumeType),
		})
	}
	bsList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BlockStorageGroup{}.AttributeTypes()}, bsVals)

	itIndices := MatchByDefOrder(p.DefIt, p.ItResps, MatchInstance)
	itVals := make([]Instance, 0, len(itIndices))
	for _, idx := range itIndices {
		it := p.ItResps[idx]
		itVals = append(itVals, Instance{
			Name:             types.StringValue(it.Name),
			RoleType:         types.StringValue(it.RoleType),
			ServiceIpAddress: types.StringPointerValue(&it.ServiceIpAddress),
			PublicIpId:       types.StringPointerValue(&it.PublicIpId),
		})
	}
	itList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Instance{}.AttributeTypes()}, itVals)

	return bsList, itList
}

func ToAnySlice[T any](src []T) []any {
	out := make([]any, len(src))
	for i, v := range src {
		out[i] = v
	}
	return out
}

func MatchBlockStorage(dv BlockStorageGroup, rv BlockStorageGroupResponse) bool {
	nameExists := !dv.Name.IsNull() && dv.Name.ValueString() != ""
	nameMatch := nameExists && rv.Name == dv.Name.ValueString()
	attrMatch := rv.RoleType == dv.RoleType.ValueString() &&
		rv.VolumeType == dv.VolumeType.ValueString() &&
		rv.SizeGb == dv.SizeGb.ValueInt32()
	return attrMatch && (!nameExists || nameMatch)
}

func MatchInstance(dv Instance, rv InstanceResponse) bool {
	nameExists := !dv.Name.IsNull() && dv.Name.ValueString() != ""
	nameMatch := nameExists && rv.Name == dv.Name.ValueString()
	roleMatch := rv.RoleType == dv.RoleType.ValueString()
	return roleMatch && (!nameExists || nameMatch)
}

func IsEqualInstanceGroup(ctx context.Context, actual InstanceGroupResponse, expect InstanceGroup) bool {

	expectIt := make([]Instance, len(expect.Instances.Elements()))
	expect.Instances.ElementsAs(ctx, &expectIt, false)
	expectItKey := make([]string, len(expectIt))
	for i, it := range expectIt {
		expectItKey[i] = it.RoleType.ValueString()
	}
	actualItKey := make([]string, len(actual.Instances))
	for i, it := range actual.Instances {
		actualItKey[i] = it.RoleType
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
			RoleType:   bs.RoleType,
			VolumeType: bs.VolumeType,
			SizeGb:     bs.SizeGb,
		}
	}

	return CompareInstanceGroupKeys(
		actual.RoleType, actual.ServerTypeName,
		expect.RoleType.ValueString(), expect.ServerTypeName.ValueString(),
		actualItKey, expectItKey,
		actualBSKey, expectBSKey,
	)

}
