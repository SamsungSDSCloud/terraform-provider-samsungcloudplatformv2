package database

import "fmt"

// BlockStorageResize describes an in-place size change of an existing block storage.
type BlockStorageResize struct {
	Id     string
	SizeGb int32
}

// BlockStorageUpdatePlan is the identity-matched classification of a desired block
// storage list against the current (state) list. Entries are matched by their
// computed Id, which InstanceGroupsPlanModifier copies from prior state onto every
// block storage that still exists (whether unchanged or merely resized). This makes
// the update independent of list position, so inserting or reordering entries no
// longer misattributes a change to an unrelated storage sitting at the same index.
type BlockStorageUpdatePlan struct {
	Resizes []BlockStorageResize // existing storages whose size_gb changed
	Adds    []BlockStorageGroup  // storages with no prior identity (new)
	Removed []BlockStorageGroup  // storages dropped from the desired list
}

type blockStorageRoleVolumeKey struct {
	roleType   string
	volumeType string
}

func PlanBlockStorageUpdate(currentBS, desiredBS []BlockStorageGroup) (BlockStorageUpdatePlan, error) {
	currentById := make(map[string]BlockStorageGroup, len(currentBS))
	for _, c := range currentBS {
		currentById[c.Id.ValueString()] = c
	}

	var plan BlockStorageUpdatePlan
	seen := make(map[string]bool, len(currentBS))
	addsByRoleVolume := make(map[blockStorageRoleVolumeKey]bool)
	resizesByRoleVolume := make(map[blockStorageRoleVolumeKey]bool)
	for _, d := range desiredBS {
		roleVolumeKey := blockStorageRoleVolumeKey{
			roleType:   d.RoleType.ValueString(),
			volumeType: d.VolumeType.ValueString(),
		}
		id := d.Id.ValueString()
		cur, ok := currentById[id]
		if id == "" || !ok {
			plan.Adds = append(plan.Adds, d)
			addsByRoleVolume[roleVolumeKey] = true
			continue
		}
		seen[id] = true
		if d.RoleType.ValueString() != cur.RoleType.ValueString() ||
			d.VolumeType.ValueString() != cur.VolumeType.ValueString() {
			return plan, fmt.Errorf("immutable fields cannot be modified: RoleType, VolumeType")
		}
		if d.SizeGb.ValueInt32() != cur.SizeGb.ValueInt32() {
			plan.Resizes = append(plan.Resizes, BlockStorageResize{Id: id, SizeGb: d.SizeGb.ValueInt32()})
			resizesByRoleVolume[roleVolumeKey] = true
		}
	}

	for roleVolumeKey := range resizesByRoleVolume {
		if addsByRoleVolume[roleVolumeKey] {
			return plan, fmt.Errorf(
				"ambiguous block storage update for role_type=%q, volume_type=%q: resizing an existing storage and adding a new storage with the same role_type/volume_type in one apply is not supported; apply the resize first, then add the new storage",
				roleVolumeKey.roleType,
				roleVolumeKey.volumeType,
			)
		}
	}

	for _, c := range currentBS {
		if !seen[c.Id.ValueString()] {
			plan.Removed = append(plan.Removed, c)
		}
	}
	return plan, nil
}
