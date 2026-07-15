package database

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type computedNullToUnknownPlanModifier struct{}

func ComputedNullToUnknown() planmodifier.String {
	return computedNullToUnknownPlanModifier{}
}

func (m computedNullToUnknownPlanModifier) Description(_ context.Context) string {
	return "Marks a computed attribute unknown for newly added collection elements so server-assigned values do not cause inconsistent-result errors."
}

func (m computedNullToUnknownPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m computedNullToUnknownPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// User explicitly set a value: keep it.
	if !req.ConfigValue.IsNull() {
		return
	}
	// Existing element (prior state at this path): keep the state-derived value.
	if !req.StateValue.IsNull() {
		return
	}
	// Already unknown: nothing to do.
	if req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = types.StringUnknown()
}
