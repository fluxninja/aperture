package main

import (
	"log"
	"os"

	"github.com/ghodss/yaml"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("expected 2 arguments (input and output file), received %d", len(os.Args))
	}
	inFile := os.Args[1]
	outFile := os.Args[2]

	circuit, err := readFromFile(inFile)
	if err != nil {
		log.Fatalf("error reading policy spec: %v", err)
	}

	dot := controlplane.DOT(controlplane.ComponentDTO(circuit))

	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(dot)
	if err != nil {
		log.Fatalf("error writing to file: %v", err)
	}
}

func readFromFile(path string) (controlplane.CompiledCircuit, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var policy policylangv1.Policy
	err = yaml.Unmarshal(yamlFile, &policy)
	if err != nil {
		return nil, err
	}
	circuit, err := controlplane.CompilePolicy(&policy)
	if err != nil {
		return nil, err
	}

	return circuit, nil
}
