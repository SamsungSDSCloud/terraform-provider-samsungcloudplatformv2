package database

import (
	"context"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ClusterLifecycleController interface {
	StartCluster(ctx context.Context, clusterId string) error
	StopCluster(ctx context.Context, clusterId string) error
}

func GetChangedFields[T any](plan, state T, targetFields []string) ([]string, error) {
	targetSet := make(map[string]bool)
	for _, f := range targetFields {
		targetSet[f] = true
	}

	planValue := reflect.ValueOf(plan)
	stateValue := reflect.ValueOf(state)

	var changedFields []string
	planType := planValue.Type()

	for i := 0; i < planValue.NumField(); i++ {
		fieldName := planType.Field(i).Name

		if !targetSet[fieldName] {
			continue
		}

		planField := planValue.Field(i).Interface()
		stateField := stateValue.Field(i).Interface()

		switch p := planField.(type) {
		case basetypes.StringValue:
			s := stateField.(basetypes.StringValue)
			if (p.IsUnknown() && s.IsNull()) || (p.IsNull() && s.IsUnknown()) {
				continue
			}
		}

		if !reflect.DeepEqual(planField, stateField) {
			changedFields = append(changedFields, fieldName)
		}
	}

	return changedFields, nil
}

// OnlyBackupOptionChanged reports whether plan and state differ only in their
// BackupOption field (or not at all). It lets BackupOption — the sole updatable
// field inside init_config_option — change while every other field stays immutable,
// without hand-rolling the neutralize-and-compare per service.
//
// plan and state must be values of the same struct type that exposes a BackupOption
// field, either directly or through an embedded struct.
func OnlyBackupOptionChanged(plan, state any) bool {
	p := reflect.ValueOf(plan)
	s := reflect.ValueOf(state)
	if p.Kind() != reflect.Struct || s.Kind() != reflect.Struct || p.Type() != s.Type() {
		return false
	}

	// Work on an addressable copy of plan so BackupOption can be neutralized.
	planCopy := reflect.New(p.Type()).Elem()
	planCopy.Set(p)

	planBackup := planCopy.FieldByName("BackupOption")
	stateBackup := s.FieldByName("BackupOption")
	if !planBackup.IsValid() || !stateBackup.IsValid() || !planBackup.CanSet() {
		return false
	}
	planBackup.Set(stateBackup)

	return reflect.DeepEqual(planCopy.Interface(), state)
}

func SnakeToPascal(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])

		words[i] = string(runes)
	}
	return strings.Join(words, "")
}

// OverlapFields returns the elements of target that also appear in changed,
// preserving target's order. Use it to report exactly which restricted fields a
// change touched; IsOverlapFields only reports whether any overlap exists.
func OverlapFields(target, changed []string) []string {
	changedSet := make(map[string]struct{}, len(changed))
	for _, c := range changed {
		changedSet[c] = struct{}{}
	}

	var overlap []string
	for _, t := range target {
		if _, ok := changedSet[t]; ok {
			overlap = append(overlap, t)
		}
	}
	return overlap
}

func IsOverlapFields(fields []string, changedFields []string) bool {
	for _, field := range fields {
		for _, changed := range changedFields {
			if field == changed {
				return true
			}
		}
	}
	return false
}

// setToMap converts a types.Set to a map[string]struct{} for efficient lookups
func setToMap(ips types.Set) map[string]struct{} {
	if ips.IsNull() || ips.IsUnknown() {
		return make(map[string]struct{})
	}

	ipSet := make(map[string]struct{}, len(ips.Elements()))
	for _, ip := range ips.Elements() {
		if strVal, ok := ip.(types.String); ok {
			ipSet[strVal.ValueString()] = struct{}{}
		}
	}
	return ipSet
}

// diffSets returns elements in sourceSet that are not in targetSet
func diffSets(sourceSet, targetSet map[string]struct{}) []string {
	result := make([]string, 0, len(sourceSet))
	for key := range sourceSet {
		if _, exists := targetSet[key]; !exists {
			result = append(result, key)
		}
	}
	return result
}

// CompareIPAddresses compares two IP address sets and returns added and removed IPs
func CompareIPAddresses(state types.Set, plan types.Set) ([]string, []string) {
	stateIPSet := setToMap(state)
	planIPSet := setToMap(plan)

	addedIPs := diffSets(planIPSet, stateIPSet)
	removedIPs := diffSets(stateIPSet, planIPSet)

	return addedIPs, removedIPs
}

func GetPendingStates(currentState string) []string {
	if currentState == "RUNNING" {
		return []string{"RUNNING", "STOPPING"}
	}
	return []string{"STOPPED", "STARTING"}
}

func GetStateTransitions(controller ClusterLifecycleController) map[string]map[string]func(ctx context.Context, clusterId string) error {
	return map[string]map[string]func(ctx context.Context, clusterId string) error{
		"STOPPED": {
			"RUNNING": controller.StartCluster,
		},
		"RUNNING": {
			"STOPPED": controller.StopCluster,
		},
	}
}

// WaitForSettledState는 클러스터가 전이 중 상태(STOPPING/STARTING)이면 대응되는 종료
// 상태(STOPPED/RUNNING)에 도달할 때까지 대기하고 그 종료 상태를 반환한다. 전이 중 상태가
// 아니면 currentState를 그대로 반환한다. 실제 대기는 서비스별 구현을 waitFn으로 주입받는다.
func WaitForSettledState(
	ctx context.Context,
	currentState string,
	clusterId string,
	waitFn func(ctx context.Context, clusterId string, pendingStates, targetStates []string) error,
) (string, error) {
	switch currentState {
	case "STOPPING":
		if err := waitFn(ctx, clusterId, []string{"STOPPING"}, []string{"STOPPED"}); err != nil {
			return currentState, err
		}
		return "STOPPED", nil
	case "STARTING":
		if err := waitFn(ctx, clusterId, []string{"STARTING"}, []string{"RUNNING"}); err != nil {
			return currentState, err
		}
		return "RUNNING", nil
	default:
		return currentState, nil
	}
}
