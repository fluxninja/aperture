package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/circuitfactory"
)

var (
	policy  string
	cr      string
	dot     string
	mermaid string
)

func init() {
	compileCmd.Flags().StringVar(&policy, "policy", "", "path to policy file")
	compileCmd.Flags().StringVar(&cr, "cr", "", "path to policy custom resource file")
	compileCmd.Flags().StringVar(&dot, "dot", "", "path to dot file")
	compileCmd.Flags().StringVar(&mermaid, "mermaid", "", "path to mermaid file")
	RootCmd.AddCommand(compileCmd)
}

var compileCmd = &cobra.Command{
	Use:           "compile",
	Short:         "Compile circuit from policy file",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if policy or cr is provided
		if policy == "" && cr == "" || policy != "" && cr != "" {
			errStr := "either --policy or --cr must be provided"
			return errors.New(errStr)
		}

		var policyFile string

		// check if cr is provided
		if cr != "" {
			crPath := cr
			// extract spec key from CR and save it a temp file
			// call compilePolicy with the temp file
			// delete the temp file
			crFile, err := os.ReadFile(crPath)
			if err != nil {
				log.Error().Err(err).Msg("failed to read CR file")
				return err
			}
			// unmarshal yaml to map struct and extract spec key
			var cr map[string]interface{}
			err = yaml.Unmarshal(crFile, &cr)
			if err != nil {
				log.Error().Err(err).Msg("failed to unmarshal CR file")
				return err
			}
			spec, ok := cr["spec"]
			if !ok {
				log.Error().Msg("failed to find spec key in CR file")
				return err
			}
			// marshal spec to yaml
			specYaml, err := yaml.Marshal(spec)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal spec key in CR file")
				return err
			}
			// get filename from path
			filename := filepath.Base(crPath)
			// create temp file
			tmpfile, err := os.CreateTemp("", filename)
			if err != nil {
				log.Error().Err(err).Msg("failed to create temp file")
				return err
			}
			defer os.Remove(tmpfile.Name())
			// write spec to temp file
			_, err = tmpfile.Write(specYaml)
			if err != nil {
				log.Error().Err(err).Msg("failed to write to temp file")
				return err
			}
			// close temp file
			err = tmpfile.Close()
			if err != nil {
				log.Error().Err(err).Msg("failed to close temp file")
				return err
			}
			// set policyFile to temp file
			policyFile = tmpfile.Name()
		} else {
			policyFile = policy
		}

		circuit, err := compilePolicy(policyFile)
		if err != nil {
			log.Error().Err(err).Msg("error reading policy spec")
			return err
		}

		log.Info().Msg("Compilation successful")

		// if --dot flag is set, write dotfile
		// check if the dot flag is set
		if dot != "" {
			dotFile := dot
			d := circuitfactory.DOT(circuit.ToGraphView())
			f, err := os.Create(dotFile)
			if err != nil {
				log.Error().Err(err).Msg("error creating file")
				return err
			}
			defer f.Close()

			_, err = f.WriteString(d)
			if err != nil {
				log.Error().Err(err).Msg("error writing to file")
				return err
			}
			log.Info().Msg("DOT file written")
		}
		// if --mermaid flag is set, write mermaid file
		if mermaid != "" {
			mermaidFile := mermaid
			m := circuitfactory.Mermaid(circuit.ToGraphView())
			f, err := os.Create(mermaidFile)
			if err != nil {
				log.Error().Err(err).Msg("error creating file")
				return err
			}
			defer f.Close()

			_, err = f.WriteString(m)
			if err != nil {
				log.Error().Err(err).Msg("error writing to file")
				return err
			}
			log.Info().Msg("Mermaid file written")
		}
		return nil
	},
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
