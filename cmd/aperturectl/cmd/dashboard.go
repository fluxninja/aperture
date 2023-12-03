package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
)

var (
	policyFile       string
	datasourceName   string
	uri              string
	dashboardVersion string
	outputDir        string
	overwrite        bool
	skipPull         bool
)

func init() {
	dashboardCmd.Flags().StringVar(&policyFile, "policy-file", "", "Path to the policy file to use")
	dashboardCmd.Flags().StringVar(&datasourceName, "datasource-name", "controller-prometheus", "Name of the datasource to use")
	dashboardCmd.Flags().StringVar(&uri, "uri", "", "URI of Custom Dashboards, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/dashboards/grafana@latest. This field should not be provided when the Version is provided.")
	dashboardCmd.Flags().StringVar(&dashboardVersion, "version", utils.LatestTag, "Version of official Aperture Dashboards, e.g. latest. This field should not be provided when the URI is provided")
	dashboardCmd.Flags().StringVar(&outputDir, "output-dir", "dashboards", "Output directory for the generated dashboards")
	dashboardCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite the output directory if it already exists")
	dashboardCmd.Flags().BoolVar(&skipPull, "skip-pull", false, "Skip pulling the dashboards from the remote repository")
}

var dashboardCmd = &cobra.Command{
	Use:           "dashboard",
	Short:         "Generate dashboards for Aperture",
	Long:          `Generate dashboards for Aperture`,
	SilenceErrors: true,
	RunE: func(*cobra.Command, []string) error {
		if policyFile == "" {
			return errors.New("policy-file is required")
		}

		if datasourceName == "" {
			return errors.New("datasource-name is required")
		}

		if uri == "" && dashboardVersion == "" {
			return errors.New("either uri or version is required")
		}

		// if uri is provided and version is default, set version to empty string
		if uri != "" && dashboardVersion == utils.LatestTag {
			dashboardVersion = ""
		}

		if outputDir == "" {
			return errors.New("output-dir is required")
		}

		absOutputDir, err := setupOutputDir(outputDir)
		if err != nil {
			return err
		}

		err = Generate(policyFile, absOutputDir)
		if err != nil {
			return err
		}

		return nil
	},
}

// Generate generates the dashboards for the given policy.
func Generate(policyFile, outputDir string) error {
	_, err := os.Stat(policyFile)
	if err != nil {
		log.Info().Msgf("Error reading values file: %s", err.Error())
		return err
	}

	_, dashboardsURIRoot, dashboardsDir, err := utils.Pull(uri, dashboardVersion, "dashboards", utils.DefaultDashboardsRepo, skipPull, true)
	if err != nil {
		return err
	}

	_, graph, policyName, policyBytes, err := utils.ProcessPolicy(policyFile)
	if err != nil {
		return err
	}

	err = generateDashboards(dashboardsURIRoot, dashboardsDir, string(policyBytes), graph, policyName, outputDir)
	if err != nil {
		return err
	}

	log.Info().Msgf("Generated dashboards for policy '%s' in '%s'", policyName, outputDir)

	return nil
}

// generateDashboards generates the dashboards for the given policy.
func generateDashboards(dashboardsURIRoot, dashboardsDir, policyStr, graph, policyName, outputDir string) error {
	// create the directory if it does not exist
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	vm := jsonnet.MakeVM()
	vm.SetTraceOut(&buf)
	vm.Importer(&jsonnet.FileImporter{
		JPaths: []string{dashboardsURIRoot},
	})

	vm.TLAReset()
	vm.TLAVar("policyFile", policyStr)
	vm.TLAVar("graph", graph)
	vm.TLAVar("policyName", policyName)
	vm.TLAVar("datasource", datasourceName)

	dashboardGroupFile := filepath.Join(dashboardsDir, "grafana", "group.libsonnet")
	dashboardsJSON, err := vm.EvaluateFile(dashboardGroupFile)
	log.Debug().Msgf("Jsonnet generation trace: %s", buf.String())
	if err != nil {
		return err
	}

	type dashboards struct {
		Dashboards map[string]interface{} `json:"dashboards"`
	}

	var dashboardsList dashboards
	err = json.Unmarshal([]byte(dashboardsJSON), &dashboardsList)
	if err != nil {
		return err
	}

	for key, val := range dashboardsList.Dashboards {
		outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s-%s.json", key, policyName))
		err = saveJSONFile(outputFilePath, val.(map[string]interface{}))
		if err != nil {
			return err
		}
	}

	return nil
}

// saveJSONFile saves the given content as a JSON file.
func saveJSONFile(outputFilePath string, content map[string]interface{}) error {
	jsonBytes, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFilePath, jsonBytes, 0o600)
	if err != nil {
		return err
	}

	return nil
}

// setupOutputDir creates the output directory if it does not exist and returns the absolute path to the output directory.
func setupOutputDir(outputDir string) (string, error) {
	// ask for user confirmation if the output directory already exists
	if !overwrite {
		if _, err := os.Stat(outputDir); err == nil {
			fmt.Printf("The output directory '%s' already exists. Do you want to merge the generated policy artifacts into the existing directory? [y/N]: ", outputDir)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return "", fmt.Errorf("output directory '%s' already exists", outputDir)
			}
		}
	}

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		return "", err
	}

	return absOutputDir, nil
}
