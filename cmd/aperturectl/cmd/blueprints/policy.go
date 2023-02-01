package blueprints

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

var (
	policyType string
	outputDir  string
	apply      bool
	kubeConfig string
	valuesFile string
)

func init() {
	policyCmd.Flags().StringVar(&policyType, "policy_type", "", "Type of policy to generate e.g. static-rate-limiting, latency-aimd-concurrency-limiting")
	policyCmd.Flags().StringVar(&outputDir, "output_dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
	policyCmd.Flags().BoolVar(&apply, "apply", false, "Apply policy on the Kubernetes cluster")
	policyCmd.Flags().StringVar(&kubeConfig, "kube_config", "", "Path to the Kubernets cluster config. Defaults to '~/.kube/config'")
	policyCmd.Flags().StringVar(&valuesFile, "values_file", "", "Path to the values file for blueprints input")
	policyCmd.Flags().StringVar(&blueprintsVersion, "version", "main", "version of aperture to pull blueprints from")
}

var policyCmd = &cobra.Command{
	Use:           "generate-policy",
	Short:         "Generate policy",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if policy or cr is provided
		if policyType == "" {
			return errors.New("--policy_type must be provided")
		} else if valuesFile == "" {
			valuesFile = fmt.Sprintf("%s.yaml", policyType)
			fmt.Printf("Using '%s' for values file\n", valuesFile)
		}

		blueprintsList, err := getBlueprintsList()
		if err != nil {
			return err
		}

		// pullCmd.RunE(cmd, args)

		policies := blueprintsList[blueprintsVersion]
		if !slices.Contains(policies, policyType) {
			return fmt.Errorf("invalid value '%s' for --policy_type. Available policies are: %+v", policyType, policies)
		}

		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)

		vm := jsonnet.MakeVM()
		fmt.Println(apertureBlueprintsDir)

		vm.ExtVar("libPath", apertureBlueprintsDir)

		policyPath := filepath.Join(apertureBlueprintsDir, fmt.Sprintf("github.com/fluxninja/aperture/blueprints/lib/1.0/policies/%s/policy.libsonnet", policyType))

		jsonStr, err := vm.EvaluateAnonymousSnippet("policy.libsonnet", fmt.Sprintf(`
		local policy = import '%s';
		local config = std.parseYaml(importstr '%s');
		local policyResource = policy(config.common + config.policy).policyResource;

		policyResource
		`, policyPath, valuesFile))
		if err != nil {
			return err
		}

		var yamlData map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &yamlData)
		if err != nil {
			return err
		}

		yamlBytes, err := yaml.Marshal(yamlData)
		if err != nil {
			return err
		}

		if outputDir != "" {
			err := os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				return err
			}

			err = os.WriteFile(filepath.Join(outputDir, "policy.yaml"), yamlBytes, 0o600)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(string(yamlBytes))
		}

		return nil
	},
}
