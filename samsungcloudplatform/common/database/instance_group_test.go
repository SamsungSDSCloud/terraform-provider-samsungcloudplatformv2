package database

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ctx = context.Background()

// --- compareSlices ---

func TestCompareSlices_Equal(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{3, 2, 1}
	if !compareSlices(a, b) {
		t.Fatal("expected true")
	}
}

func TestCompareSlices_DifferentLength(t *testing.T) {
	a := []int{1, 2}
	b := []int{1, 2, 3}
	if compareSlices(a, b) {
		t.Fatal("expected false")
	}
}

func TestCompareSlices_DifferentValues(t *testing.T) {
	a := []int{1, 2}
	b := []int{1, 3}
	if compareSlices(a, b) {
		t.Fatal("expected false")
	}
}

func TestCompareSlices_Empty(t *testing.T) {
	if !compareSlices([]int{}, []int{}) {
		t.Fatal("expected true")
	}
}

func TestCompareSlices_Duplicates(t *testing.T) {
	a := []string{"a", "a", "b"}
	b := []string{"a", "b", "a"}
	if !compareSlices(a, b) {
		t.Fatal("expected true")
	}
}

func TestCompareSlices_DifferentDuplicates(t *testing.T) {
	a := []string{"a", "a", "b"}
	b := []string{"a", "b", "b"}
	if compareSlices(a, b) {
		t.Fatal("expected false")
	}
}

// --- CompareBlockStorages ---

func TestCompareBlockStorages_Equal(t *testing.T) {
	a := []BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}}
	b := []BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}}
	if !CompareBlockStorages(a, b) {
		t.Fatal("expected true")
	}
}

func TestCompareBlockStorages_Different(t *testing.T) {
	a := []BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}}
	b := []BSKey{{RoleType: "DATA", SizeGb: 100, VolumeType: "SSD"}}
	if CompareBlockStorages(a, b) {
		t.Fatal("expected false")
	}
}

func TestCompareBlockStorages_DifferentOrder(t *testing.T) {
	a := []BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}, {RoleType: "OS", SizeGb: 100, VolumeType: "HDD"}}
	b := []BSKey{{RoleType: "OS", SizeGb: 100, VolumeType: "HDD"}, {RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}}
	if !CompareBlockStorages(a, b) {
		t.Fatal("expected true")
	}
}

// --- CompareInstances ---

func TestCompareInstances_Equal(t *testing.T) {
	a := []string{"MASTER", "SLAVE"}
	b := []string{"SLAVE", "MASTER"}
	if !CompareInstances(a, b) {
		t.Fatal("expected true")
	}
}

func TestCompareInstances_Different(t *testing.T) {
	a := []string{"MASTER"}
	b := []string{"SLAVE"}
	if CompareInstances(a, b) {
		t.Fatal("expected false")
	}
}

// --- CompareInstanceGroupKeys ---

func TestCompareInstanceGroupKeys_AllMatch(t *testing.T) {
	if !CompareInstanceGroupKeys("MASTER", "t1", "MASTER", "t1",
		[]string{"MASTER"}, []string{"MASTER"},
		[]BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}},
		[]BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}}) {
		t.Fatal("expected true")
	}
}

func TestCompareInstanceGroupKeys_RoleTypeMismatch(t *testing.T) {
	if CompareInstanceGroupKeys("MASTER", "t1", "SLAVE", "t1",
		[]string{"MASTER"}, []string{"MASTER"},
		nil, nil) {
		t.Fatal("expected false")
	}
}

func TestCompareInstanceGroupKeys_ServerTypeMismatch(t *testing.T) {
	if CompareInstanceGroupKeys("MASTER", "t1", "MASTER", "t2",
		[]string{"MASTER"}, []string{"MASTER"},
		nil, nil) {
		t.Fatal("expected false")
	}
}

func TestCompareInstanceGroupKeys_InstancesMismatch(t *testing.T) {
	if CompareInstanceGroupKeys("MASTER", "t1", "MASTER", "t1",
		[]string{"MASTER"}, []string{"SLAVE"},
		nil, nil) {
		t.Fatal("expected false")
	}
}

func TestCompareInstanceGroupKeys_BSKeysMismatch(t *testing.T) {
	if CompareInstanceGroupKeys("MASTER", "t1", "MASTER", "t1",
		nil, nil,
		[]BSKey{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}},
		[]BSKey{{RoleType: "DATA", SizeGb: 100, VolumeType: "SSD"}}) {
		t.Fatal("expected false")
	}
}

// --- MatchByDefOrder ---

func TestMatchByDefOrder_AllMatch(t *testing.T) {
	defs := []int{1, 2, 3}
	resps := []int{3, 1, 2}
	match := func(d, r int) bool { return d == r }
	result := MatchByDefOrder(defs, resps, match)
	if len(result) != 3 || result[0] != 1 || result[1] != 2 || result[2] != 0 {
		t.Fatalf("unexpected: %v", result)
	}
}

func TestMatchByDefOrder_NoMatch(t *testing.T) {
	defs := []int{1, 2}
	resps := []int{3, 4}
	match := func(d, r int) bool { return d == r }
	result := MatchByDefOrder(defs, resps, match)
	if len(result) != 0 {
		t.Fatalf("expected empty, got: %v", result)
	}
}

func TestMatchByDefOrder_DuplicatesConsumedOnce(t *testing.T) {
	defs := []string{"a", "a"}
	resps := []string{"a"}
	match := func(d, r string) bool { return d == r }
	result := MatchByDefOrder(defs, resps, match)
	if len(result) != 1 {
		t.Fatalf("expected 1 match, got: %d", len(result))
	}
}

// --- ToAnySlice ---

func TestToAnySlice(t *testing.T) {
	src := []int{1, 2, 3}
	out := ToAnySlice(src)
	if len(out) != 3 {
		t.Fatal("unexpected length")
	}
	for i, v := range out {
		if v != any(src[i]) {
			t.Fatal("unexpected value")
		}
	}
}

func TestToAnySlice_Empty(t *testing.T) {
	out := ToAnySlice([]int{})
	if len(out) != 0 {
		t.Fatal("expected empty")
	}
}

// --- matchBlockStorage ---

func TestMatchBlockStorage_Match(t *testing.T) {
	def := BlockStorageGroup{
		RoleType:   types.StringValue("DATA"),
		VolumeType: types.StringValue("SSD"),
		SizeGb:     types.Int32Value(50),
		Name:       types.StringValue("mybs"),
	}
	resp := BlockStorageGroupResponse{
		RoleType:   "DATA",
		VolumeType: "SSD",
		SizeGb:     50,
		Name:       "mybs",
	}
	if !matchBlockStorage(def, resp) {
		t.Fatal("expected true")
	}
}

func TestMatchBlockStorage_NameNull(t *testing.T) {
	def := BlockStorageGroup{
		RoleType:   types.StringValue("DATA"),
		VolumeType: types.StringValue("SSD"),
		SizeGb:     types.Int32Value(50),
		Name:       types.StringNull(),
	}
	resp := BlockStorageGroupResponse{
		RoleType:   "DATA",
		VolumeType: "SSD",
		SizeGb:     50,
		Name:       "anything",
	}
	if !matchBlockStorage(def, resp) {
		t.Fatal("expected true — null name should ignore response name")
	}
}

func TestMatchBlockStorage_NameMismatch(t *testing.T) {
	def := BlockStorageGroup{
		RoleType:   types.StringValue("DATA"),
		VolumeType: types.StringValue("SSD"),
		SizeGb:     types.Int32Value(50),
		Name:       types.StringValue("a"),
	}
	resp := BlockStorageGroupResponse{
		RoleType:   "DATA",
		VolumeType: "SSD",
		SizeGb:     50,
		Name:       "b",
	}
	if matchBlockStorage(def, resp) {
		t.Fatal("expected false")
	}
}

func TestMatchBlockStorage_SizeMismatch(t *testing.T) {
	def := BlockStorageGroup{
		RoleType:   types.StringValue("DATA"),
		VolumeType: types.StringValue("SSD"),
		SizeGb:     types.Int32Value(100),
		Name:       types.StringNull(),
	}
	resp := BlockStorageGroupResponse{
		RoleType:   "DATA",
		VolumeType: "SSD",
		SizeGb:     50,
		Name:       "x",
	}
	if matchBlockStorage(def, resp) {
		t.Fatal("expected false")
	}
}

// --- matchInstance ---

func TestMatchInstance_Match(t *testing.T) {
	def := Instance{
		RoleType: types.StringValue("MASTER"),
		Name:     types.StringValue("inst1"),
	}
	resp := InstanceResponse{
		RoleType: "MASTER",
		Name:     "inst1",
	}
	if !matchInstance(def, resp) {
		t.Fatal("expected true")
	}
}

func TestMatchInstance_NameNull(t *testing.T) {
	def := Instance{
		RoleType: types.StringValue("MASTER"),
		Name:     types.StringNull(),
	}
	resp := InstanceResponse{
		RoleType: "MASTER",
		Name:     "anything",
	}
	if !matchInstance(def, resp) {
		t.Fatal("expected true")
	}
}

func TestMatchInstance_NameMismatch(t *testing.T) {
	def := Instance{
		RoleType: types.StringValue("MASTER"),
		Name:     types.StringValue("a"),
	}
	resp := InstanceResponse{
		RoleType: "MASTER",
		Name:     "b",
	}
	if matchInstance(def, resp) {
		t.Fatal("expected false")
	}
}

func TestMatchInstance_RoleMismatch(t *testing.T) {
	def := Instance{
		RoleType: types.StringValue("MASTER"),
		Name:     types.StringNull(),
	}
	resp := InstanceResponse{
		RoleType: "SLAVE",
		Name:     "x",
	}
	if matchInstance(def, resp) {
		t.Fatal("expected false")
	}
}

// --- helper to build plan InstanceGroup ---

func newPlanIG(bs []BlockStorageGroup, inst []Instance) InstanceGroup {
	bsList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BlockStorageGroup{}.AttributeTypes()}, bs)
	itList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Instance{}.AttributeTypes()}, inst)
	return InstanceGroup{
		Id:                 types.StringNull(),
		BlockStorageGroups: bsList,
		Instances:          itList,
		RoleType:           types.StringValue("MASTER"),
		ServerTypeName:     types.StringValue("t1"),
	}
}

func newRespIG(id string, bs []BlockStorageGroupResponse, inst []InstanceResponse) InstanceGroupResponse {
	return InstanceGroupResponse{
		Id:                 id,
		BlockStorageGroups: bs,
		Instances:          inst,
		RoleType:           "MASTER",
		ServerTypeName:     "t1",
	}
}

// --- isEqualInstanceGroup ---

func TestIsEqualInstanceGroup_Match(t *testing.T) {
	def := newPlanIG(
		[]BlockStorageGroup{{RoleType: types.StringValue("DATA"), SizeGb: types.Int32Value(50), VolumeType: types.StringValue("SSD"), Name: types.StringNull()}},
		[]Instance{{RoleType: types.StringValue("MASTER"), Name: types.StringNull()}},
	)
	resp := newRespIG("1",
		[]BlockStorageGroupResponse{{RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}},
		[]InstanceResponse{{RoleType: "MASTER"}},
	)
	if !isEqualInstanceGroup(ctx, resp, def) {
		t.Fatal("expected true")
	}
}

func TestIsEqualInstanceGroup_RoleMismatch(t *testing.T) {
	def := newPlanIG(nil, []Instance{{RoleType: types.StringValue("MASTER"), Name: types.StringNull()}})
	resp := newRespIG("1", nil, []InstanceResponse{{RoleType: "SLAVE"}})
	if isEqualInstanceGroup(ctx, resp, def) {
		t.Fatal("expected false")
	}
}

func TestIsEqualInstanceGroup_ServerTypeMismatch(t *testing.T) {
	def := InstanceGroup{
		Id:                 types.StringNull(),
		BlockStorageGroups: types.ListNull(types.ObjectType{AttrTypes: BlockStorageGroup{}.AttributeTypes()}),
		Instances:          types.ListNull(types.ObjectType{AttrTypes: Instance{}.AttributeTypes()}),
		RoleType:           types.StringValue("MASTER"),
		ServerTypeName:     types.StringValue("t2"),
	}
	resp := newRespIG("1", nil, nil)
	if isEqualInstanceGroup(ctx, resp, def) {
		t.Fatal("expected false")
	}
}

// --- MapInstanceGroup ---

func TestMapInstanceGroup_Reorder(t *testing.T) {
	defBs := []BlockStorageGroup{
		{RoleType: types.StringValue("OS"), SizeGb: types.Int32Value(100), VolumeType: types.StringValue("HDD"), Name: types.StringNull()},
		{RoleType: types.StringValue("DATA"), SizeGb: types.Int32Value(50), VolumeType: types.StringValue("SSD"), Name: types.StringNull()},
	}
	defIt := []Instance{
		{RoleType: types.StringValue("SLAVE"), Name: types.StringNull()},
		{RoleType: types.StringValue("MASTER"), Name: types.StringNull()},
	}
	bsResps := []BlockStorageGroupResponse{
		{Id: "bs2", Name: "data", RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"},
		{Id: "bs1", Name: "os", RoleType: "OS", SizeGb: 100, VolumeType: "HDD"},
	}
	itResps := []InstanceResponse{
		{Name: "slave1", RoleType: "SLAVE", ServiceIpAddress: "10.0.0.2"},
		{Name: "master1", RoleType: "MASTER", ServiceIpAddress: "10.0.0.1"},
	}

	p := &MapInstanceGroupParams{DefBs: defBs, DefIt: defIt, BsResps: bsResps, ItResps: itResps}
	bsList, itList := MapInstanceGroup(ctx, p)

	var bsOut []BlockStorageGroup
	bsList.ElementsAs(ctx, &bsOut, false)
	if len(bsOut) != 2 || bsOut[0].RoleType.ValueString() != "OS" || bsOut[1].RoleType.ValueString() != "DATA" {
		t.Fatalf("bs order wrong: %v", bsOut)
	}

	var itOut []Instance
	itList.ElementsAs(ctx, &itOut, false)
	if len(itOut) != 2 || itOut[0].RoleType.ValueString() != "SLAVE" || itOut[1].RoleType.ValueString() != "MASTER" {
		t.Fatalf("it order wrong: %v", itOut)
	}
}

// --- MapInstanceGroups ---

func TestMapInstanceGroups_Found(t *testing.T) {
	def := newPlanIG(
		[]BlockStorageGroup{{RoleType: types.StringValue("DATA"), SizeGb: types.Int32Value(50), VolumeType: types.StringValue("SSD"), Name: types.StringNull()}},
		[]Instance{{RoleType: types.StringValue("MASTER"), Name: types.StringNull()}},
	)
	resps := []InstanceGroupResponse{
		newRespIG("ig1",
			[]BlockStorageGroupResponse{{Id: "bs1", Name: "d", RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}},
			[]InstanceResponse{{Name: "m1", RoleType: "MASTER", ServiceIpAddress: "10.0.0.1"}},
		),
	}

	remaining, mapped := MapInstanceGroups(ctx, resps, def)
	if len(remaining) != 0 {
		t.Fatal("expected remaining to be empty")
	}
	if mapped.Id.ValueString() != "ig1" {
		t.Fatal("expected id ig1")
	}
}

func TestMapInstanceGroups_NotFound(t *testing.T) {
	def := newPlanIG(nil, []Instance{{RoleType: types.StringValue("MASTER"), Name: types.StringNull()}})
	resps := []InstanceGroupResponse{
		newRespIG("ig1", nil, []InstanceResponse{{Name: "m1", RoleType: "SLAVE"}}),
	}

	remaining, mapped := MapInstanceGroups(ctx, resps, def)
	if len(remaining) != 1 {
		t.Fatal("expected 1 remaining")
	}
	if mapped.Id.ValueString() != "" {
		t.Fatal("expected empty id")
	}
}

// --- MapInstanceGroupsList ---

func TestMapInstanceGroupsList_Single(t *testing.T) {
	def := newPlanIG(
		[]BlockStorageGroup{{RoleType: types.StringValue("DATA"), SizeGb: types.Int32Value(50), VolumeType: types.StringValue("SSD"), Name: types.StringNull()}},
		[]Instance{{RoleType: types.StringValue("MASTER"), Name: types.StringNull()}},
	)
	defList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InstanceGroup{}.AttributeTypes()}, []InstanceGroup{def})
	resps := []InstanceGroupResponse{
		newRespIG("ig1",
			[]BlockStorageGroupResponse{{Id: "bs1", Name: "d", RoleType: "DATA", SizeGb: 50, VolumeType: "SSD"}},
			[]InstanceResponse{{Name: "m1", RoleType: "MASTER", ServiceIpAddress: "10.0.0.1"}},
		),
	}

	result := MapInstanceGroupsList(ctx, defList, resps)
	var out []InstanceGroup
	result.ElementsAs(ctx, &out, false)
	if len(out) != 1 || out[0].Id.ValueString() != "ig1" {
		t.Fatalf("unexpected: %+v", out)
	}
}
