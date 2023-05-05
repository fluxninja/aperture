// mapstruct is similar in idea to mitchellh/mapstructure, with a difference
// that values are assumed to always be created through json serialization.
//
// Eg. If one would try to directly convert some struct to map[string]any via
// mapstructure, resulting map wouldn't necessary match json representation of
// said object, especially if custom UnmarshalJSON are present. Object created
// by EncodeObject avoids such mismatch.
package utils

import "encoding/json"

// MapStruct is an in-memory representation of JSON object.
type MapStruct map[string]any

// ToMapStruct encodes any json-serializable struct as Object
//
// json.Unmarshal(obj) and json.Unmarshal(returnedObject) should be equivalent.
func ToMapStruct(obj any) (MapStruct, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var mapStruct MapStruct
	err = json.Unmarshal(b, &mapStruct)
	if err != nil {
		return nil, err
	}

	return mapStruct, nil
}
