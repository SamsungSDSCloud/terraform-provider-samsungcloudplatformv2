package backup

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
	"strings"
	"unicode"
)

func SetResourceIdentifier(stateIdentifier types.String, availableIdentifiers []types.String, diags *diag.Diagnostics) types.String {
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
				"Error Invalid Backup Resource Identifier",
				"Could not find Backup Resource "+stateIdentifier.ValueString())
		}
	} else {
		result = availableIdentifiers[0]
	}

	return result
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
