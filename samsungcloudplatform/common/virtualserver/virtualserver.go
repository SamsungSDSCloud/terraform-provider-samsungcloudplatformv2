package virtualserver

import (
	"context"
	"encoding/json"
	"fmt"
	scpsdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatformv2/v3/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func UnsetNilFields(v interface{}) {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)

		if field.Kind() == reflect.Struct {
			valueField := field.FieldByName("value")
			if valueField.IsValid() && (valueField.IsNil() || valueField.Elem().String() == "") {
				unsetMethod := field.Addr().MethodByName("Unset")
				if unsetMethod.IsValid() {
					unsetMethod.Call(nil)
				}
			}
		}
	}
}

func ToNullableStringValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func ToNullableInt32Value(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
}

func SetResourceIdentifier(stateIdentifier types.String, availableIdentifiers []types.String, diags *diag.Diagnostics) types.String {
	// SetResourceIdentifier validates or auto-selects resource identifiers
	//
	// When stateIdentifier exists:
	//    - Checks existence in availableIdentifiers
	//    - Returns valid identifier or adds error to diagnostics
	// When stateIdentifier is null:
	//    - Returns first available identifier (requires non-empty slice)

	result := types.String{}
	exist := false

	if !stateIdentifier.IsNull() {
		for _, v := range availableIdentifiers {
			if v == stateIdentifier {
				result = stateIdentifier
				exist = true
				break
			}
		}
		if !exist {
			diags.AddError(
				"Error Invalid Resource Identifier",
				"Could not find Resource "+stateIdentifier.ValueString())
		}
	} else {
		result = availableIdentifiers[0]
	}

	return result
}

func AsyncRequestPollingWithState[T any](ctx context.Context, resourceId string, maxAttempts int, internal time.Duration,
	stateField string, DesiredState string, errorState string, getFunc func(string) (T, error)) (T, error) {
	var zero T
	ticker := time.NewTicker(internal)
	defer ticker.Stop()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resource, err := getFunc(resourceId)
		if err != nil {
			return zero, fmt.Errorf("attempt %d/%d failed: %w",
				attempt, maxAttempts, err)
		}

		currentState, err := getField(resource, stateField)
		if err != nil {
			return zero, fmt.Errorf("field extraction failed: %w", err)
		}

		if errorState != "" && currentState == errorState {
			return zero, fmt.Errorf("resource state:" + errorState)
		}

		switch currentState := currentState.(type) {
		case string:
			if currentState == DesiredState {
				return resource, nil
			}
		case bool:
			if strconv.FormatBool(currentState) == DesiredState {
				return resource, nil
			}
		case int:
			if strconv.Itoa(currentState) == DesiredState {
				return resource, nil
			}
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

func getField(obj interface{}, fieldName string) (interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("nil input")
	}

	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct: %v", val.Kind())
	}

	fields := strings.Split(fieldName, ".")

	currentVal := val

	for _, field := range fields {
		fieldVal := currentVal.FieldByName(field)
		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("field '%s' not found", fieldName)
		}

		if !fieldVal.CanInterface() {
			return nil, fmt.Errorf("cannot access private field '%s'", fieldName)
		}

		currentVal = fieldVal
	}

	return currentVal.Interface(), nil
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

func GetStatusCodeFromError(err error) float64 {
	var status float64
	var data map[string]interface{}
	body := err.(*scpsdk.GenericOpenAPIError).Body()
	err = json.Unmarshal(body, &data)

	errors := data["errors"]
	for _, err := range errors.([]interface{}) {
		errorMap := err.(map[string]interface{})
		status = errorMap["status"].(float64)
	}

	return status
}
