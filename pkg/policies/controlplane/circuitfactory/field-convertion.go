package circuitfactory

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ConvertOldComponentToNew converts an old component protobuf message to a new component protobuf message using the provided field mappings.
// It takes in the oldComponentProto, newComponentProto, and fieldMappings as parameters.
// Param: oldComponentProto is the protobuf message that needs to be converted to the new format.
// Param: newComponentProto is the protobuf message that will be populated with the converted data.
// Param: fieldMappings is a map of old field names to new field names. Example format:
//
//	{
//	    "A_field": "Y_field",
//	    "A_field.B_field": "Z_field",
//	    "C_field.D_field": "G_field",
//	    "C_field.D_field.E_field": "H_field",
//	}
//
// In this example, the following renames will be performed:
// 1. A_field -> Y_field
// 2. A_field.B_field -> Y_field.Z_field
// 3. C_field.D_field -> C_field.G_field
// 4. C_field.D_field.E_field -> C_field.G_field.H_field
// The function returns an error if there is an issue with marshaling or unmarshaling the protobuf messages or JSON data.
func ConvertOldComponentToNew(oldComponentProto proto.Message, newComponentProto proto.Message, fieldMappings map[string]string) error {
	// Marshal the oldComponentProto to JSON
	jsonStr, err := json.Marshal(oldComponentProto)
	if err != nil {
		return fmt.Errorf("error marshaling old component to JSON: %v", err)
	}

	// Unmarshal the JSON data into a map
	var data map[string]interface{}
	if err2 := json.Unmarshal(jsonStr, &data); err2 != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Define a recursive function to replace old field names with new field names in the map
	var replaceField func(data map[string]interface{}, oldField, newField string) bool
	replaceField = func(data map[string]interface{}, oldField, newField string) bool {
		parts := strings.SplitN(oldField, ".", 2)
		if len(parts) == 1 {
			if _, ok := data[oldField]; ok {
				data[newField] = data[oldField]
				delete(data, oldField)
				return true
			}
		} else {
			if data[parts[0]] != nil {
				if nestedData, ok := data[parts[0]].(map[string]interface{}); ok {
					return replaceField(nestedData, parts[1], newField)
				}
			}
		}
		return false
	}

	// Sort the fieldMappings by depth (deepest to shallowest)
	sortedKeys := make([]string, 0, len(fieldMappings))
	for k := range fieldMappings {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return strings.Count(sortedKeys[i], ".") > strings.Count(sortedKeys[j], ".")
	})

	// Replace old field names with new field names in the map using the sorted fieldMappings
	for _, oldField := range sortedKeys {
		newField := fieldMappings[oldField]
		replaceField(data, oldField, newField)
	}

	// Marshal the modified map back to JSON
	newJSONStr, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling modified JSON: %v", err)
	}

	// Unmarshal the new JSON data into the newComponentProto
	err = protojson.Unmarshal(newJSONStr, newComponentProto)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON to new component: %v", err)
	}

	return nil
}
