package database

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
	"strings"
	"time"
	"unicode"
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
			return zero, fmt.Errorf(string("resource state:" + errorState))
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

func ConvertListtoStringSlice(list types.List) ([]string, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	result := make([]string, len(list.Elements()))
	for i, elem := range list.Elements() {
		strVal, ok := elem.(types.String)
		if !ok {
			return nil, fmt.Errorf("unexpected elementtype: %T", elem)
		}
		result[i] = strVal.ValueString()
	}
	return result, nil
}

func CompareIPAddresses(state []string, plan []string) ([]string, []string) {

	toSet := func(items []string) map[string]bool {
		set := make(map[string]bool)
		for _, item := range items {
			set[item] = true
		}
		return set
	}

	stateIPSet := toSet(state)
	planIPSet := toSet(plan)

	diff := func(sourceSet map[string]bool, targetSet map[string]bool) []string {
		var result []string
		for key := range sourceSet {
			if !targetSet[key] {
				result = append(result, key)
			}
		}
		return result
	}

	addedIPs := diff(planIPSet, stateIPSet)
	removedIPs := diff(stateIPSet, planIPSet)

	return addedIPs, removedIPs
}