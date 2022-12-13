// mapstruct is similar in idea to mitchellh/mapstructure, with a difference
// that values are assumed to always be created through json serialization.
//
// Eg. If one would try to directly convert some struct to map[string]any via
// mapstructure, resulting map wouldn't necessary match json representation of
// said object, especially if custom UnmarshalJSON are present. Object created
// by EncodeObject avoids such mismatch.
package mapstruct

import "encoding/json"

// Object is an in-memory representation of JSON object.
type Object map[string]any

// EncodeObject encodes any json-serializable struct as Object
//
// json.Unmarshal(obj) and json.Unmarshal(returnedObject) should be equivalent.
func EncodeObject(obj any) (Object, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var mapStruct Object
	err = json.Unmarshal(b, &mapStruct)
	if err != nil {
		return nil, err
	}

	return mapStruct, nil
}
