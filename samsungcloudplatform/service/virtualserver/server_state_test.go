package virtualserver

// SCPEQC-54893: 스키마에는 state 로 ACTIVE / SHUTOFF 가 모두 명시돼 있지만
// API 는 생성 시 ACTIVE 만 허용한다. 종전에는 Create() 안에서만 검사해서
// plan 은 멀쩡히 통과하고 apply 에서야 실패했다.
//
// ModifyPlan 이 prior state 가 없는 경우(= create)를 잡아 plan 단계에서 거부해야 한다.

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// objectWithState: state 만 지정하고 나머지 속성은 전부 null 인 값.
func objectWithState(t *testing.T, s schema.Schema, state tftypes.Value) tftypes.Value {
	t.Helper()

	objType, ok := s.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("서버 스키마가 object 가 아니다")
	}

	attributes := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attrType, nil)
	}
	attributes["state"] = state

	return tftypes.NewValue(objType, attributes)
}

// modifyPlanOnCreate: prior state 가 없는(=create) 상황으로 ModifyPlan 을 돌린다.
func modifyPlanOnCreate(t *testing.T, planState tftypes.Value) resource.ModifyPlanResponse {
	t.Helper()

	ctx := context.Background()
	s := serverSchema(t)

	plan := tfsdk.Plan{Schema: s, Raw: objectWithState(t, s, planState)}

	req := resource.ModifyPlanRequest{
		Plan: plan,
		// Raw 를 비워두면 IsNull() 이 true 라서 create 로 취급된다.
		State: tfsdk.State{Schema: s},
	}
	resp := resource.ModifyPlanResponse{Plan: plan}

	(&virtualServerServerResource{}).ModifyPlan(ctx, req, &resp)
	return resp
}

// 생성 시 SHUTOFF 는 plan 단계에서 거부돼야 한다. apply 까지 가면 안 된다.
func TestModifyPlan_CreateWithShutoffIsRejected(t *testing.T) {
	resp := modifyPlanOnCreate(t, tftypes.NewValue(tftypes.String, ServerStateShutoff))

	if !resp.Diagnostics.HasError() {
		t.Fatalf("state=%s 로 생성하는 plan 이 통과했다. plan 단계에서 막아야 한다", ServerStateShutoff)
	}

	if !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), ServerStateActive) {
		t.Errorf("에러 메시지가 %s 만 허용된다는 걸 알려주지 않는다: %s",
			ServerStateActive, resp.Diagnostics.Errors()[0].Detail())
	}
}

// ACTIVE 는 정상 통과해야 한다.
func TestModifyPlan_CreateWithActiveIsAllowed(t *testing.T) {
	resp := modifyPlanOnCreate(t, tftypes.NewValue(tftypes.String, ServerStateActive))

	if resp.Diagnostics.HasError() {
		t.Fatalf("state=%s 로 생성하는 plan 이 막혔다: %v", ServerStateActive, resp.Diagnostics)
	}
}

// state 를 config 에 쓰지 않은 경우(null)도 통과해야 한다.
func TestModifyPlan_CreateWithoutStateIsAllowed(t *testing.T) {
	resp := modifyPlanOnCreate(t, tftypes.NewValue(tftypes.String, nil))

	if resp.Diagnostics.HasError() {
		t.Fatalf("state 를 지정하지 않은 plan 이 막혔다: %v", resp.Diagnostics)
	}
}

// 아직 값이 정해지지 않은 경우(unknown)는 판단할 수 없으므로 통과시켜야 한다.
func TestModifyPlan_CreateWithUnknownStateIsAllowed(t *testing.T) {
	resp := modifyPlanOnCreate(t, tftypes.NewValue(tftypes.String, tftypes.UnknownValue))

	if resp.Diagnostics.HasError() {
		t.Fatalf("state 가 unknown 인 plan 이 막혔다: %v", resp.Diagnostics)
	}
}
