package main

import (
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// process swaagger extensions such as x-go-tag-default and x-go-tag-validate
// and translate them to swagger's defaults and validations.
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

	definitions, ok := swagger["definitions"].(map[string]interface{})
	if !ok {
		panic("no definitions")
	}

	processDefinitions(definitions)

	// write back to file
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

func processDefinitions(definitions map[string]interface{}) {
	for _, definition := range definitions {
		m, ok := definition.(map[string]interface{})
		if !ok {
			continue
		}
		for k, v := range m {
			if k == "properties" {
				properties, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				required := processProperties(properties)
				if len(required) > 0 {
					m["required"] = required
				}
			}
		}
	}
}

func processProperties(properties map[string]interface{}) []string {
	var required []string
	// walk swagger to look for x-go-tag-validate key
	// if found, add swagger's validation based on that
	for p, pv := range properties {
		pmap, ok := pv.(map[string]interface{})
		if !ok {
			continue
		}
		for k, v := range pmap {
			if k == "x-go-tag-validate" {
				// extract validation rules which are comma separated
				// example: "oneof=info warn crit,required"
				rules, ok := v.(string)
				if !ok {
					continue
				}
				// extract each rule
				req := processValidateRules(pmap, rules)
				if req {
					required = append(required, p)
				}
			}
		}
	}
	// sort required
	sort.Strings(required)
	return required
}

func processValidateRules(_ map[string]interface{}, rules string) (required bool) {
	// extract each rule, for now stop when we encounter "dive"
	for _, rule := range strings.Split(rules, ",") {
		switch rule {
		case "required":
			required = true
		case "dive":
			return
		}
	}
	return
}
