package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fluxninja/aperture/pkg/policies/controlplane"
)

func main() {
	// use flag to parse flags and args
	// 2 args - input file and output file
	fs := flag.NewFlagSet("circuit-generator", flag.ExitOnError)
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if fs.NArg() != 2 {
		log.Fatal("Usage: circuit-generator <input-file> <output-file>")
	}
	inFile := fs.Arg(0)
	outFile := fs.Arg(1)

	circuit, err := compile(inFile)
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

func compile(path string) (controlplane.CompiledCircuit, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	circuit, valid, msg, err := controlplane.ValidateAndCompile(ctx, filepath.Base(path), yamlFile)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("invalid circuit: %s", msg)
	}
	return circuit, nil
}
