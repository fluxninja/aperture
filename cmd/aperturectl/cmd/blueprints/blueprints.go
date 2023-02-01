package blueprints

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

const apertureBlueprintsRepo = "github.com/fluxninja/aperture/blueprints"

var (
	blueprintsDir string

	// Args for `blueprints pull`.
	blueprintsVersion string
)

func init() {
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsVersion, "version", "main", "version of aperture blueprint")

	BlueprintsCmd.AddCommand(pullCmd)
	BlueprintsCmd.AddCommand(policyCmd)
	BlueprintsCmd.AddCommand(listCmd)
	BlueprintsCmd.AddCommand(removeCmd)
	BlueprintsCmd.AddCommand(dashboardCmd)
	BlueprintsCmd.AddCommand(graphCmd)
	BlueprintsCmd.AddCommand(allCmd)
}

var BlueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Manage blueprints",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		blueprintsDir = filepath.Join(userHomeDir, ".aperturectl", "blueprints")
		err = os.MkdirAll(blueprintsDir, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	},
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

func validatePolicyType(cmd *cobra.Command, args []string) error {
	blueprintsList, err := getBlueprintsList()
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
