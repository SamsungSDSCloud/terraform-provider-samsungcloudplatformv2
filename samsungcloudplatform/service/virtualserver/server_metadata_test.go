package virtualserver

// SCPEQC-54892: metadata 를 지정하고 apply 하면
// "Provider produced inconsistent result after apply" 로 실패하던 건.
//
// API 가 응답 metadata 에 HA_Enabled 같은 키를 스스로 덧붙이는데, metadata 가
// Optional 이라 config 로 확정한 값과 응답이 달라지면서 terraform 의
// "확정된 값은 apply 후에도 같아야 한다" 계약이 깨졌다.
//
// metadata 를 Computed 전용으로 바꿔서 비교 대상인 config 값 자체를 없앴다.
// 그래서 플랫폼이 무슨 키를 덧붙이든 apply 가 깨지지 않는다.
//
// mock API 하네스(mockAPI / newMockedResource / showResponse / importState)는
// server_import_test.go 의 것을 그대로 쓴다.

import (
	"context"
	"testing"

	scpvirtualserver "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v6/library/virtualserver/1.5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	metadataInjectedKey = "HA_Enabled"
	metadataUserKey     = "os_type"
	metadataUserValue   = "linux"
)

// serverSchema: 리소스가 선언한 스키마를 그대로 가져온다.
func serverSchema(t *testing.T) schema.Schema {
	t.Helper()

	resp := &resource.SchemaResponse{}
	(&virtualServerServerResource{}).Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("스키마를 만들 수 없다: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// metadataOnlyMock: metadata 매핑만 보기 위해 SG 를 비운 최소 구성.
func metadataOnlyMock() mockAPI {
	return mockAPI{
		interfaces:     []scpvirtualserver.InterfaceResponseV1Dot2{iface(testPort1, testSubnet1, testFixedIp1, true, false)},
		securityGroups: []scpvirtualserver.SecurityGroupResponse{},
		serverVolumes:  []scpvirtualserver.ServersVolumeResponse{attached(testBootVolumeId, deviceOsDisk, true)},
		volumes:        []scpvirtualserver.VolumeShowResponseV1Dot4{volume(testBootVolumeId, 104, volumeTypeSSD, nil, nil)},
		defaultVolType: volumeTypeSSD,
	}
}

// mapServerMetadata: 주어진 API 응답 metadata 로 매핑을 돌려 state 의 metadata 를 얻는다.
func mapServerMetadata(t *testing.T, apiMetadata map[string]interface{}) types.Map {
	t.Helper()

	r := newMockedResource(t, metadataOnlyMock())

	resp := showResponse()
	resp.Metadata = apiMetadata

	mapped, err := r.MapGetResponseToState(context.Background(), resp, importState(), types.MapNull(types.StringType))
	if err != nil {
		t.Fatalf("metadata 매핑 실패: %v", err)
	}
	return mapped.Metadata
}

func assertMetadata(t *testing.T, got types.Map, want map[string]string) {
	t.Helper()

	elements := got.Elements()
	if len(elements) != len(want) {
		t.Fatalf("metadata 키 개수가 다르다. want %d(%v), got %d(%v)", len(want), want, len(elements), elements)
	}
	for k, wantValue := range want {
		v, ok := elements[k]
		if !ok {
			t.Fatalf("metadata 에 %q 가 없다. got %v", k, elements)
		}
		s, ok := v.(types.String)
		if !ok {
			t.Fatalf("metadata[%q] 가 문자열이 아니다: %T", k, v)
		}
		if s.ValueString() != wantValue {
			t.Fatalf("metadata[%q] want %q, got %q", k, wantValue, s.ValueString())
		}
	}
}

// 이 결함의 핵심 수정. metadata 가 config 로 지정 가능하면(Optional)
// 플랫폼이 덧붙인 키 때문에 apply 가 다시 깨진다.
func TestMetadata_SchemaIsComputedOnly(t *testing.T) {
	attribute, ok := serverSchema(t).Attributes["metadata"].(schema.MapAttribute)
	if !ok {
		t.Fatalf("metadata 속성을 찾을 수 없다")
	}

	if attribute.Optional {
		t.Errorf("metadata 가 Optional 이면 config 로 지정할 수 있게 되어 " +
			"플랫폼이 덧붙인 키 때문에 apply 가 inconsistent result 로 깨진다")
	}
	if attribute.Required {
		t.Errorf("metadata 는 사용자가 지정하는 값이 아니다")
	}
	if !attribute.Computed {
		t.Errorf("metadata 는 Computed 여야 API 가 준 값을 그대로 담을 수 있다")
	}
}

// 플랫폼이 덧붙인 키는 그대로 state 에 담긴다. Computed 전용이므로
// 비교할 config 값이 없어서 이래도 apply 가 깨지지 않는다.
func TestMetadata_PlatformInjectedKeyIsReflected(t *testing.T) {
	got := mapServerMetadata(t, map[string]interface{}{
		metadataUserKey:     metadataUserValue,
		metadataInjectedKey: "True",
	})

	assertMetadata(t, got, map[string]string{
		metadataUserKey:     metadataUserValue,
		metadataInjectedKey: "True",
	})
}

// 덧붙는 값이 문자열이 아니어도(bool 등) panic 없이 문자열로 담겨야 한다.
func TestMetadata_NonStringValueIsStringified(t *testing.T) {
	got := mapServerMetadata(t, map[string]interface{}{
		metadataInjectedKey: true,
	})

	assertMetadata(t, got, map[string]string{metadataInjectedKey: "true"})
}

// metadata 가 비어 있어도 null 이 아닌 빈 맵이어야 한다.
// Computed 속성이 null 로 남으면 plan 이 unknown 을 계속 물고 간다.
func TestMetadata_EmptyResponseIsEmptyMap(t *testing.T) {
	got := mapServerMetadata(t, map[string]interface{}{})

	if got.IsNull() || got.IsUnknown() {
		t.Fatalf("빈 metadata 는 빈 맵이어야 한다. got %v", got)
	}
	assertMetadata(t, got, map[string]string{})
}
