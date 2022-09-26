package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
)

func main() {
	info.Service = "circuit-compiler"

	logger := log.NewLogger(log.GetPrettyConsoleWriter(), "debug")
	log.SetGlobalLogger(logger)

	// flags:
	// 1. Required: --policy - path to policy file
	// 2. Optional: --dot - path to dot file
	fs := flag.NewFlagSet("circuit-compiler", flag.ExitOnError)
	policy := fs.String("policy", "", "path to policy file")
	dot := fs.String("dot", "", "path to dot file")
	// parse flags
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
		os.Exit(1)
	}

	// check if policy flag is set
	if *policy == "" {
		log.Error().Msg("policy flag is required")
		os.Exit(1)
	}
	policyFile := *policy

	circuit, err := compile(policyFile)
	if err != nil {
		log.Error().Err(err).Msg("error reading policy spec")
		os.Exit(1)
	}

	log.Info().Msg("Compilation successful")

	// if --dot flag is set, write dotfile
	// check if the dot flag is set
	if *dot != "" {
		dotFile := *dot
		dot := controlplane.DOT(controlplane.ComponentDTO(circuit))
		f, err := os.Create(dotFile)
		if err != nil {
			log.Error().Err(err).Msg("error creating file")
			os.Exit(1)
		}
		defer f.Close()

		_, err = f.WriteString(dot)
		if err != nil {
			log.Error().Err(err).Msg("error writing to file")
			os.Exit(1)
		}
		log.Info().Msg("DOT file written")
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
