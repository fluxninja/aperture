package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	swaggerFile := os.Args[1]

	swaggerBytes, err := os.ReadFile(swaggerFile)
	if err != nil {
		panic(err)
	}

	// decode as YAML
	var swagger map[string]interface{}
	err = yaml.Unmarshal(swaggerBytes, &swagger)
	if err != nil {
		panic(err)
	}
	// look for definitions
	definitions, ok := swagger["definitions"].(map[string]interface{})
	if !ok {
		panic("no definitions")
	}
	replaceRefs := make(map[string]string)
	// read all definitions to look for keys with prefix
	// these deinitions have to replaced with new names that do not contain the prefix and all the dot separators are removed
	for k, v := range definitions {
		key := k
		if strings.HasPrefix(key, "aperture.") {
			v1 := ".v1."
			// find the first occurrence of v1 in key
			i := strings.Index(key, v1)
			// remove these characters from key including v1
			key = key[i+len(v1):]
		}
		// remove dots in key
		key = strings.ReplaceAll(key, ".", "")
		if key == k {
			continue
		}
		// replace this key
		definitions[key] = v
		delete(definitions, k)
		replaceRefs["#/definitions/"+k] = "#/definitions/" + key
	}

	// search for $ref key in entire swagger and replace with new name
	replaceRef(swagger, replaceRefs)

	// encode as YAML
	swaggerBytes, err = yaml.Marshal(swagger)
	if err != nil {
		panic(err)
	}
	// write to file
	err = os.WriteFile(swaggerFile, swaggerBytes, 0o600)
	if err != nil {
		panic(err)
	}
}

func replaceRef(content map[string]interface{}, replacements map[string]string) {
	for k, v := range content {
		if k == "$ref" {
			ref, ok := v.(string)
			if !ok {
				continue
			}
			if newRef, ok := replacements[ref]; ok {
				content[k] = newRef
			}
		}
		if m, ok := v.(map[string]interface{}); ok {
			replaceRef(m, replacements)
		}
	}
}
