package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func AsyncRequestPollingWithState[T any](ctx context.Context, clusterId string, maxAttempts int, internal time.Duration,
	stateField string, DesiredState string, errorState string, getFunc func(string) (T, error)) (T, error) {
	var zero T
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resource, err := getFunc(clusterId)
		if err != nil {
			//return zero, fmt.Errorf("attempt %d/%d failed: %w",
			//	attempt, maxAttempts, err)
			return zero, fmt.Errorf("%w", err)
		}

		currentState, err := getField(resource, stateField)
		if err != nil {
			return zero, fmt.Errorf("field extraction failed: %w", err)
		}

		if errorState != "" && currentState == errorState {
			return zero, fmt.Errorf("resource state: %s", errorState)
		}

		if currentState == DesiredState {
			return resource, nil
		}

		if attempt < maxAttempts {
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return zero, fmt.Errorf("polling canceled: %w", ctx.Err())
			}
		}
	}

	return zero, fmt.Errorf("max attempts reached (%d)", maxAttempts)
}

func getField(obj interface{}, fieldName string) (string, error) {
	if obj == nil {
		return "none", fmt.Errorf("nil object")
	}

	val := reflect.ValueOf(obj)

	// obj Ptr 이면 실제 값으로 변경
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// obj 구조체인지 확인
	if val.Kind() != reflect.Struct {
		return "none", fmt.Errorf("not a struct: %v", val.Kind())
	}

	fields := strings.Split(fieldName, ".")

	currentVal := val

	for _, field := range fields {
		fieldVal := currentVal.FieldByName(field)
		if !fieldVal.IsValid() {
			return "none", fmt.Errorf("field '%s' not found", fieldName)
		}

		if !fieldVal.CanInterface() {
			return "none", fmt.Errorf("cannot access private field '%s'", fieldName)
		}

		currentVal = fieldVal
	}

	return currentVal.String(), nil
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
