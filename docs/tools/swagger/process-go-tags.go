package main

import (
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// process swagger extensions such as x-go-tag-default and x-go-tag-validate
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
			} else if k == "x-go-tag-default" {
				// extract default value
				defaultValue, ok := v.(string)
				if !ok {
					continue
				}
				// add default value to swagger
				processDefault(pmap, defaultValue)
			}
		}
	}
	// sort required
	sort.Strings(required)
	return required
}

func processDefault(m map[string]interface{}, d string) {
	dv := processValue(m, d)
	if dv != nil {
		m["default"] = dv
	}
}

func processValue(m map[string]interface{}, d string) interface{} {
	vtype, ok := m["type"].(string)
	if !ok {
		return nil
	}
	switch vtype {
	case "array":
		if strings.HasPrefix(d, "[") && strings.HasSuffix(d, "]") {
			// extract items
			items, ok := m["items"].(map[string]interface{})
			if !ok {
				panic("no items")
			}
			d = d[1 : len(d)-1]
			elements := strings.Split(d, ",")
			var values []interface{}
			for _, e := range elements {
				values = append(values, processValue(items, e))
			}
			return values
		}
	case "string":
		return d
	case "boolean":
		return d == "true"
	case "integer":
		format, ok := m["format"].(string)
		if !ok {
			return nil
		}
		switch format {
		case "int64":
			f, err := strconv.ParseInt(d, 10, 64)
			if err != nil {
				panic(err)
			}
			return f
		case "int32":
			f, err := strconv.ParseInt(d, 10, 32)
			if err != nil {
				panic(err)
			}
			return f
		}
	case "number":
		format, ok := m["format"].(string)
		if !ok {
			return nil
		}
		switch format {
		case "float":
			f, err := strconv.ParseFloat(d, 32)
			if err != nil {
				panic(err)
			}
			return f
		case "double":
			f, err := strconv.ParseFloat(d, 64)
			if err != nil {
				panic(err)
			}
			return f
		}
	}
	return nil
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
