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
	"github.com/fluxninja/aperture/pkg/policies/controlplane/circuitfactory"
	"gopkg.in/yaml.v2"
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
	cr := fs.String("cr", "", "path to policy custom resource file")
	dot := fs.String("dot", "", "path to dot file")
	mermaid := fs.String("mermaid", "", "path to mermaid file")

	// parse flags
	e := fs.Parse(os.Args[1:])
	if e != nil {
		log.Error().Err(e).Msg("failed to parse flags")
		os.Exit(1)
	}

	// check if policy or cr is provided
	if *policy == "" && *cr == "" || *policy != "" && *cr != "" {
		log.Error().Msg("either --policy or --cr must be provided")
		os.Exit(1)
	}

	var policyFile string

	// check if cr is provided
	if *cr != "" {
		crPath := *cr
		// extract spec key from CR and save it a temp file
		// call compilePolicy with the temp file
		// delete the temp file
		crFile, err := os.ReadFile(crPath)
		if err != nil {
			log.Error().Err(err).Msg("failed to read CR file")
			os.Exit(1)
		}
		// unmarshal yaml to map struct and extract spec key
		var cr map[string]interface{}
		err = yaml.Unmarshal(crFile, &cr)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal CR file")
			os.Exit(1)
		}
		spec, ok := cr["spec"]
		if !ok {
			log.Error().Msg("failed to find spec key in CR file")
			os.Exit(1)
		}
		// marshal spec to yaml
		specYaml, err := yaml.Marshal(spec)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal spec key in CR file")
			os.Exit(1)
		}
		// get filename from path
		filename := filepath.Base(crPath)
		// create temp file
		tmpfile, err := os.CreateTemp("", filename)
		if err != nil {
			log.Error().Err(err).Msg("failed to create temp file")
			os.Exit(1)
		}
		defer os.Remove(tmpfile.Name())
		// write spec to temp file
		_, err = tmpfile.Write(specYaml)
		if err != nil {
			log.Error().Err(err).Msg("failed to write to temp file")
			os.Exit(1)
		}
		// close temp file
		err = tmpfile.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close temp file")
			os.Exit(1)
		}
		// set policyFile to temp file
		policyFile = tmpfile.Name()
	} else {
		policyFile = *policy
	}

	circuit, err := compilePolicy(policyFile)
	if err != nil {
		log.Error().Err(err).Msg("error reading policy spec")
		os.Exit(1)
	}

	log.Info().Msg("Compilation successful")

	// if --dot flag is set, write dotfile
	// check if the dot flag is set
	if *dot != "" {
		dotFile := *dot
		dot := circuitfactory.DOT(circuit.ToGraphView())
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
	// if --mermaid flag is set, write mermaid file
	if *mermaid != "" {
		mermaidFile := *mermaid
		mermaid := circuitfactory.Mermaid(circuit.ToGraphView())
		f, err := os.Create(mermaidFile)
		if err != nil {
			log.Error().Err(err).Msg("error creating file")
			os.Exit(1)
		}
		defer f.Close()

		_, err = f.WriteString(mermaid)
		if err != nil {
			log.Error().Err(err).Msg("error writing to file")
			os.Exit(1)
		}
		log.Info().Msg("Mermaid file written")
	}
}

func compilePolicy(path string) (*circuitfactory.Circuit, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	// FIXME This ValidateAndCompile function validates the policy as a whole –
	// circuit, but also the other resource classifiers, fluxmeters.  This
	// command is called "circuit-compiler" though, so it's bit... surprising.
	// If we compiled just a circuit, we could drop dependency on
	// `controlplane` package.
	circuit, valid, msg, err := controlplane.ValidateAndCompile(ctx, filepath.Base(path), yamlFile)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("invalid circuit: %s", msg)
	}
	return circuit, nil
}
