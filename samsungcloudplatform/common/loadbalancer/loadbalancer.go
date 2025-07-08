package loadbalancer

import "github.com/hashicorp/terraform-plugin-framework/types"

func ToNullableInt32Value(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}

func ToNullableBoolValue(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}
