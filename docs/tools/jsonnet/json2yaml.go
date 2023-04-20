package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

// simple program to convert json to yaml
// usage: json2yaml <json file> <yaml file>
func main() {
	jsonFile := os.Args[1]
	yamlFile := os.Args[2]
	// read json file
	jsonBytes, err := os.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	// convert json to map[string]interface{}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		panic(err)
	}
	// convert map to yaml
	yamlBytes, err := yaml.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}
	// write yaml file
	err = os.WriteFile(yamlFile, yamlBytes, 0o600)
	if err != nil {
		panic(err)
	}
}
