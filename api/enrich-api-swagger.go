package main

import (
	"os"
	"strconv"
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
			if i == -1 {
				panic("cannot find v1 in key: all aperture definitions must have v1 in their name")
			}
			// remove these characters from key including v1
			key = key[i+len(v1):]
		}
		// remove dots in key
		key = strings.ReplaceAll(key, ".", "")
		if key == k {
			continue
		}
		// replace this key
		// first check if the new key already exists
		if _, ok := definitions[key]; ok {
			panic("please provide unique definition name as the key already exists: " + key)
		}
		definitions[key] = v
		delete(definitions, k)
		replaceRefs["#/definitions/"+k] = "#/definitions/" + key
	}

	// search for $ref key in entire swagger and replace with new name
	replaceRef(swagger, replaceRefs)

	// process gotags
	processGoTags(swagger)

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
		// dive into the map
		if m, ok := v.(map[string]interface{}); ok {
			replaceRef(m, replacements)
		}
		// dive into the array
		if a, ok := v.([]interface{}); ok {
			for _, v1 := range a {
				if m, ok := v1.(map[string]interface{}); ok {
					replaceRef(m, replacements)
				}
			}
		}
	}
}

func processGoTags(content map[string]interface{}) {
	// look for "@gotags: " line in the description
	// example of gotags annotation: "@gotags: default:"info" validate:"oneof=info warn crit"
	for k, v := range content {
		if k == "description" {
			desc, ok := v.(string)
			if !ok {
				continue
			}
			// split desc into separate lines
			lines := strings.Split(desc, "\n")
			// look for "@gotags: " line
			for _, line := range lines {
				prefix := "@gotags: "
				// look for "@gotags:" line
				if !strings.HasPrefix(line, prefix) {
					continue
				}

				// remove line from description
				desc = strings.ReplaceAll(desc, line, "")
				content["description"] = desc

				// remove "@gotags: " prefix
				tags := line[len(prefix):]
				tagMap := parseStructTag(tags)
				// add each tag as "x-go-tag-<tagname>" key
				for k1, v1 := range tagMap {
					content["x-go-tag-"+k1] = v1
				}
			}
		}
		// dive into the map
		if m, ok := v.(map[string]interface{}); ok {
			processGoTags(m)
		}
		// dive into the array
		if a, ok := v.([]interface{}); ok {
			for _, v1 := range a {
				if m, ok := v1.(map[string]interface{}); ok {
					processGoTags(m)
				}
			}
		}
	}
}

func parseStructTag(tags string) map[string]string {
	tagMap := make(map[string]string)

	for tags != "" {
		// Skip leading space.
		i := 0
		for i < len(tags) && tags[i] == ' ' {
			i++
		}
		tags = tags[i:]
		if tags == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tags) && tags[i] > ' ' && tags[i] != ':' && tags[i] != '"' && tags[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tags) || tags[i] != ':' || tags[i+1] != '"' {
			break
		}
		name := string(tags[:i])
		tags = tags[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tags) && tags[i] != '"' {
			if tags[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tags) {
			break
		}
		qvalue := string(tags[:i+1])
		tags = tags[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			continue
		}

		tagMap[name] = value
	}

	return tagMap
}
