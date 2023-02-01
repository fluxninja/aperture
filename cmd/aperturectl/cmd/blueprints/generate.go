package blueprints

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func init() {
	generateCmd.Flags().StringVar(&policyType, "policy_type", "", "Type of policy to generate e.g. static-rate-limiting, latency-aimd-concurrency-limiting")
	generateCmd.Flags().StringVar(&outputDir, "output_dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
	generateCmd.Flags().BoolVar(&apply, "apply", false, "Apply policy on the Kubernetes cluster")
	generateCmd.Flags().StringVar(&kubeConfig, "kube_config", "", "Path to the Kubernets cluster config. Defaults to '~/.kube/config'")
	generateCmd.Flags().StringVar(&valuesFile, "values_file", "", "Path to the values file for blueprints input")

	generateCmd.AddCommand(generatePolicyCmd)
	generateCmd.AddCommand(generateDashboardCmd)
	generateCmd.AddCommand(generateGraphCmd)
}

var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate manifests for the given policy type",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateGraphCmd.RunE(cmd, args); err != nil {
			return err
		}
		if err := generateDashboardCmd.RunE(cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func validatePolicyType(cmd *cobra.Command, args []string) error {
	blueprintsList, err := getBlueprints()
	if err != nil {
		return err
	}

	err = pullCmd.RunE(cmd, args)
	if err != nil {
		return err
	}

	policies := blueprintsList[blueprintsVersion]
	if !slices.Contains(policies, policyType) {
		return fmt.Errorf("invalid value '%s' for --policy_type. Available policies are: %+v", policyType, policies)
	}

	return nil
}

func getOutputDir(oldPath string) (string, error) {
	newPath := filepath.Join(oldPath, blueprintsVersion, policyType)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	newPath, err = filepath.Abs(newPath)
	if err != nil {
		return "", err
	}

	return newPath, nil
}

func getValuesFilePath() error {
	if valuesFile == "" {
		valuesFile = fmt.Sprintf("%s.yaml", policyType)
	}

	absoluteValuesFile, err := filepath.Abs(valuesFile)
	if err != nil {
		return err
	}

	valuesFile = absoluteValuesFile
	if _, err = os.Stat(valuesFile); err != nil {
		return fmt.Errorf("values_file '%s' doesn't exist", valuesFile)
	}

	return nil
}
