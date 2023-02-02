package blueprints

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	generateAll   bool
	blueprintName string
	outputDir     string
	valuesFile    string
)

func init() {
	generateCmd.Flags().BoolVar(&generateAll, "all", false, "generate all the available blueprints")
	generateCmd.PersistentFlags().StringVar(&blueprintName, "name", "", "name of blueprint to generate")
	generateCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
	generateCmd.PersistentFlags().StringVar(&valuesFile, "values-file", "", "Path to the values file for blueprints input")
	// generateCmd.PersistentFlags().BoolVar(&apply, "apply", false, "Apply policy on the Kubernetes cluster")
	// generateCmd.PersistentFlags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernets cluster config. Defaults to '~/.kube/config'")

	generateCmd.AddCommand(generatePolicyCmd)
	generateCmd.AddCommand(generateDashboardCmd)
	generateCmd.AddCommand(generateGraphCmd)
}

var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate manifests for the given policy type",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}

		fileInfo, err := os.Stat(valuesFile)
		if err != nil {
			return err
		}
		valuesFile = filepath.Clean(fileInfo.Name())

		if generateAll {
			if outputDir == "" {
				outputDir = "."
			}

			if err := generatePolicyCmd.RunE(cmd, args); err != nil {
				return err
			}
			if err := generateGraphCmd.RunE(cmd, args); err != nil {
				return err
			}
			if err := generateDashboardCmd.RunE(cmd, args); err != nil {
				return err
			}
			return nil
		}
		return nil
	},
}

func blueprintExists(version string, name string) error {
	blueprintsList, err := getBlueprintsByVersion(version)
	if err != nil {
		return err
	}

	if !slices.Contains(blueprintsList, name) {
		return fmt.Errorf("invalid blueprint name '%s'", name)
	}
	return nil
}

func getOutputDir(outputPath string) (string, error) {
	newOutputPath := filepath.Join(outputPath, blueprintsVersion, blueprintName)
	err := os.MkdirAll(newOutputPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	newOutputPath, err = filepath.Abs(newOutputPath)
	if err != nil {
		return "", err
	}
	return newOutputPath, nil
}
