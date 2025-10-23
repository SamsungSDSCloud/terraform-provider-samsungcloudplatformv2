package filter

import (
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/v2/samsungcloudplatform/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Filter struct {
	Name     types.String   `tfsdk:"name"`
	Values   []types.String `tfsdk:"values"`
	UseRegex types.Bool     `tfsdk:"use_regex"`
}

func DataSourceSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Filter",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "Filtering target name",
					Required:    true,
				},
				"values": schema.ListAttribute{
					ElementType: types.StringType,
					Description: "Filtering values. Each matching value is appended. (OR rule)",
					Required:    true,
				},
				"use_regex": schema.BoolAttribute{
					Description: "Enable regex match for values",
					Required:    true,
				},
			},
		},
	}
}

func ApplyFilter(items []map[string]interface{}, filters []Filter) []map[string]interface{} {
	if len(filters) == 0 {
		return items
	}

	for _, f := range filters {
		targetName := "struct." + f.Name.ValueString()

		var elements []string
		var err error
		if elements, err = getElements(items, targetName); err != nil {
			// Fallback to default
			elements = []string{targetName}
		}

		useRegex := f.UseRegex.ValueBool()

		// Create a string equality check strategy based on this filters "regex" flag
		stringsEqual := func(propertyVal string, filterVal string) bool {
			if useRegex {
				re, err := regexp.Compile(filterVal)
				if err != nil {
					log.Printf(`[WARN] Invalid regular expression "%s" for "%s" filter\n`, filterVal, targetName)
					return false
				}
				return re.MatchString(propertyVal)
			}

			return filterVal == propertyVal
		}

		var result []map[string]interface{}
		for _, item := range items {
			targetVal, targetValOk := getValueFromPath(item, elements)
			if targetValOk && orComparator(targetVal, convertTypesStringToInterface(f.Values), stringsEqual) {
				result = append(result, item)
			}
		}
		items = result
	}
	return items
}

func convertTypesStringToInterface(values []types.String) []interface{} {
	var result []interface{}
	for _, v := range values {
		result = append(result, v.ValueString())
	}
	return result
}

func getElements(schemaMaps []map[string]interface{}, targetName string) ([]string, error) {
	if len(schemaMaps) == 0 {
		return nil, fmt.Errorf("input schemaMaps is nil")
	}

	tokenized := strings.Split(targetName, ".")

	if len(tokenized) == 0 {
		return nil, fmt.Errorf("invalid target name : %s", targetName)
	}

	var elements []string
	currentSchemaMap := schemaMaps
	for _, t := range tokenized {
		// Iterate over schema maps to find the field
		var found bool
		for _, schemaMap := range currentSchemaMap {
			if fieldSchema, ok := schemaMap[t]; ok {
				elements = append(elements, t)

				// Check for nested items
				if nestedSchema, isMap := fieldSchema.(map[string]interface{}); isMap {
					// Update currentSchemaMap to handle nested structure
					currentSchemaMap = append(currentSchemaMap, nestedSchema)
					found = true
					break
				} else {
					// If no nested structure, end loop
					found = true
					break
				}
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid schema found for filter name %s", targetName)
		}
	}

	if len(elements) == 0 {
		return nil, fmt.Errorf("everything is filtered out")
	}

	return elements, nil
}
func getValueFromPath(item map[string]interface{}, path []string) (interface{}, bool) {
	workingMap := item
	for _, pathElement := range path {
		if value, ok := workingMap[pathElement]; ok {
			// If the value is a map, drill down further
			if nestedMap, ok := value.(map[string]interface{}); ok {
				workingMap = nestedMap
			} else {
				return value, true
			}
		} else {
			return nil, false
		}
	}
	return nil, false
}

func orComparator(target interface{}, filters []interface{}, stringsEqual func(string, string) bool) bool {
	for _, fVal := range filters {
		switch targetVal := target.(type) {
		case bool:
			fBool, err := strconv.ParseBool(fVal.(string))
			if err != nil {
				log.Println("[WARN] Invalid boolean value for filtering")
				return false
			}
			if targetVal == fBool {
				return true
			}
		case int, int32, int64:
			fInt, err := strconv.ParseInt(fVal.(string), 10, 64)
			if err != nil {
				log.Println("[WARN] Invalid integer value for filtering")
				return false
			}
			if targetVal == fInt {
				return true
			}
		case float64:
			fFloat, err := strconv.ParseFloat(fVal.(string), 64)
			if err != nil {
				log.Println("[WARN] Invalid float value for filtering")
				return false
			}
			if targetVal == fFloat {
				return true
			}
		case string:
			if stringsEqual(targetVal, fVal.(string)) {
				return true
			}
		}
	}
	return false
}

type WrappedStruct struct {
	Index  int
	Struct interface{}
}

func WrapStructs(input interface{}) ([]WrappedStruct, error) {
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}

	var result []WrappedStruct

	for i := 0; i < val.Len(); i++ {
		structValue := val.Index(i)
		if structValue.Kind() != reflect.Struct {
			return nil, fmt.Errorf("each element in the slice must be a struct")
		}

		wrapped := WrappedStruct{
			Index:  i,
			Struct: structValue.Interface(),
		}

		result = append(result, wrapped)
	}

	return result, nil
}

func removeItems(slice interface{}, indicesToRemove []int) []interface{} {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice {
		fmt.Println("Provided value is not a slice")
		return nil
	}

	if len(indicesToRemove) == 0 {
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	}

	removeMap := make(map[int]bool)
	for _, idx := range indicesToRemove {
		removeMap[idx] = true
	}

	var result []interface{}
	for i := 0; i < v.Len(); i++ {
		if _, exists := removeMap[i]; !exists {
			result = append(result, v.Index(i).Interface())
		}
	}

	return result
}

func GetFilterIndices(input interface{}, filters []Filter) []int {
	wrapStructs, _ := WrapStructs(input)
	contents := common.ConvertStructToMaps(wrapStructs)
	contents = ApplyFilter(contents, filters)

	var indices []int
	for _, item := range contents {
		index, _ := common.ToInt(item["index"])
		indices = append(indices, index)
	}

	return indices
}
