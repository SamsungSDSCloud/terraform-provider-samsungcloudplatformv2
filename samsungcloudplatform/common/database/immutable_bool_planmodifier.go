package database

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// immutableBoolPlanModifier rejects changes to an immutable bool attribute at
// plan time. Enforcing immutability only in Update() surfaces the error at apply
// and leaves a permanent diff; doing it here fails `terraform plan` and never
// lets the changed value reach state.
type immutableBoolPlanModifier struct{}

// ImmutableBool returns a bool plan modifier that errors when the attribute
// changes from a known prior-state value. Place it after UseStateForUnknown()
// so a config-omitted value (resolved from state) is treated as unchanged.
func ImmutableBool() planmodifier.Bool {
	return immutableBoolPlanModifier{}
}

func (m immutableBoolPlanModifier) Description(_ context.Context) string {
	return "Rejects changes to this immutable attribute at plan time."
}

func (m immutableBoolPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m immutableBoolPlanModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// Create (no prior state), destroy (null plan), or not-yet-known plan: nothing to compare.
	if req.StateValue.IsNull() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	if !req.PlanValue.Equal(req.StateValue) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Immutable Attribute Change",
			req.Path.String()+" is immutable and cannot be changed after the cluster is created.",
		)
	}
}
