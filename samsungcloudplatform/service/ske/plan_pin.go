package ske

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// pinNestedObject rebuilds a computed nested object for an update plan.
// Attributes in keep take their prior state value; all others are unknown and
// filled by apply. Only keep attributes the update cannot change: a kept value
// that differs after apply fails with "inconsistent result after apply".
func pinNestedObject(ctx context.Context, state types.Object, keep map[string]bool) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	attrTypes := state.AttributeTypes(ctx)
	stateAttrs := state.Attributes()
	attrs := make(map[string]attr.Value, len(attrTypes))
	for name, t := range attrTypes {
		if keep[name] {
			attrs[name] = stateAttrs[name]
			continue
		}
		v, err := t.ValueFromTerraform(ctx, tftypes.NewValue(t.TerraformType(ctx), tftypes.UnknownValue))
		if err != nil {
			diags.AddError("Error Modifying Plan", "Could not build unknown value for "+name+": "+err.Error())
			return types.ObjectUnknown(attrTypes), diags
		}
		attrs[name] = v
	}
	obj, d := types.ObjectValue(attrTypes, attrs)
	diags.Append(d...)
	return obj, diags
}

// useStateIfUnconfigured keeps the prior state value of an Optional+Computed
// attribute that is not set in the configuration, instead of marking it unknown.
func useStateIfUnconfigured[T attr.Value](plan *T, state T, config T) {
	if (*plan).IsUnknown() && config.IsNull() && !state.IsUnknown() {
		*plan = state
	}
}

// rootAttributesUnchanged reports whether every root attribute except the
// computed nested object has the same value in plan and state.
func rootAttributesUnchanged(plan, state tftypes.Value, nested string) bool {
	var planAttrs, stateAttrs map[string]tftypes.Value
	if plan.As(&planAttrs) != nil || state.As(&stateAttrs) != nil {
		return false
	}
	for name, v := range planAttrs {
		if name != nested && !v.Equal(stateAttrs[name]) {
			return false
		}
	}
	return true
}

// canPinNestedObject reports whether the plan's nested object is unknown and
// the state has a known value to pin from.
func canPinNestedObject(plan types.Object, state types.Object) bool {
	return plan.IsUnknown() && !state.IsNull() && !state.IsUnknown()
}
