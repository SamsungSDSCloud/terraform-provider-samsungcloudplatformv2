package database

// InstanceUpdatePlan is the identity-matched classification of a desired instance list
// against the current (state) list. Instances have no id, so identity is carried by the
// computed Name that InstanceGroupsPlanModifier restores from prior state onto every
// instance that still exists. An instance with no Name is one with no prior identity — a
// new instance to add. Matching by Name rather than list position is what lets a new
// instance be inserted anywhere in the list without the shifted existing instances being
// seen as changed (or their service IPs being sent to the wrong new instance).
type InstanceUpdatePlan struct {
	Adds    []Instance // instances with no prior identity (Name unset)
	Removed []Instance // instances dropped from the desired list
}

// PlanInstanceUpdate classifies desiredIt against currentIt by identity (Name):
//   - a desired instance whose Name is null/unknown/empty has no prior identity and is a
//     new addition;
//   - a current instance whose Name never appears in desiredIt is a removal.
//
// It mirrors PlanBlockStorageUpdate. Callers that only support scaling up (the DB
// AddInstances API) use Adds and may ignore Removed.
func PlanInstanceUpdate(currentIt, desiredIt []Instance) InstanceUpdatePlan {
	var plan InstanceUpdatePlan
	seen := make(map[string]bool, len(desiredIt))
	for _, d := range desiredIt {
		name := d.Name.ValueString()
		if d.Name.IsNull() || d.Name.IsUnknown() || name == "" {
			plan.Adds = append(plan.Adds, d)
			continue
		}
		seen[name] = true
	}
	for _, c := range currentIt {
		if !seen[c.Name.ValueString()] {
			plan.Removed = append(plan.Removed, c)
		}
	}
	return plan
}
