package loadbalancer

import (
	"time"

	util "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.3"
	utilv1d4 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/loadbalancer/1.4"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func ToNullableStringValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// NullableTime -> types.String 변환
func ToNullableTimeString(nt util.NullableTime) types.String {
	// 1. IsSet() 메서드로 값이 설정되었는지 확인
	if nt.IsSet() {
		// 2. Get() 메서드로 *time.Time 값 가져오기
		value := nt.Get()
		if value != nil {
			return types.StringValue(value.Format(time.RFC3339))
		}
	}
	return types.StringNull()
}

func ToNullableTimeStringV1d4(nt utilv1d4.NullableTime) types.String {
	// 1. IsSet() 메서드로 값이 설정되었는지 확인
	if nt.IsSet() {
		// 2. Get() 메서드로 *time.Time 값 가져오기
		value := nt.Get()
		if value != nil {
			return types.StringValue(value.Format(time.RFC3339))
		}
	}
	return types.StringNull()
}
