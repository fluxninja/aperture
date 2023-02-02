package blueprints

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/goccy/go-graphviz"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"sigs.k8s.io/yaml"
)

var (
	blueprintName string
	outputDir     string
	valuesFile    string
	graphDir      string
)

func init() {
	generateCmd.PersistentFlags().StringVar(&blueprintName, "name", "", "Name of the blueprints to generate manifests for")
	generateCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Directory path where the generated manifests will be stored. If not provided, will use current directory")
	generateCmd.PersistentFlags().StringVar(&valuesFile, "values-file", "", "Path to the values file for blueprints input")
}

var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate manifests for the given blueprint",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}

		if valuesFile == "" {
			return fmt.Errorf("--values-file must be provided")
		}

		// if outputDir is not provided, default to current directory
		if outputDir == "" {
			outputDir = "."
		}

		_, err := os.Stat(valuesFile)
		if err != nil {
			return err
		}
		valuesFile, err = filepath.Abs(valuesFile)
		if err != nil {
			return err
		}

		if err = pullCmd.RunE(cmd, args); err != nil {
			return err
		}

		err = blueprintExists(blueprintsVersion, blueprintName)
		if err != nil {
			return err
		}

		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)

		vm := jsonnet.MakeVM()
		vm.Importer(&jsonnet.FileImporter{
			JPaths: []string{apertureBlueprintsDir},
		})

		jsonStr, err := vm.EvaluateAnonymousSnippet("bundle.libsonnet", fmt.Sprintf(`
		local bundle = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/%s/bundle.libsonnet';
		local config = std.parseYaml(importstr '%s');
		bundle { _config+:: config }
		`, blueprintName, valuesFile))
		if err != nil {
			return err
		}

		var yamlData map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &yamlData)
		if err != nil {
			return err
		}

		var updatedOutputDir string

		if outputDir != "" {
			updatedOutputDir, err = getOutputDir(outputDir)
			if err != nil {
				return err
			}
		}

		if err = processContent(yamlData, updatedOutputDir); err != nil {
			return err
		}

		log.Info().Msgf("Stored all the manifests at '%s'.", updatedOutputDir)
		return nil
	},
}

func processContent(data map[string]interface{}, outputPath string) error {
	for name, content := range data {
		if strings.HasSuffix(name, ".yaml") {
			if err := saveYamlFile(outputPath, name, content); err != nil {
				return err
			}
		} else if strings.HasSuffix(name, ".json") {
			if err := saveJSONFile(outputPath, name, content); err != nil {
				return err
			}
		} else {
			contentMap, ok := content.(map[string]interface{})
			if !ok {
				log.Error().Msgf("failed to process output '%+v'", content)
				continue
			}
			var updatedPath string
			if outputPath != "" {
				updatedPath = filepath.Join(outputPath, name)
				err := os.MkdirAll(updatedPath, os.ModePerm)
				if err != nil {
					return err
				}
			}
			if err := processContent(contentMap, updatedPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func saveYamlFile(path, filename string, content interface{}) error {
	yamlBytes, err := yaml.Marshal(content)
	if err != nil {
		return err
	}

	filePath := filepath.Join(path, filename)

	err = os.WriteFile(filePath, yamlBytes, 0o600)
	if err != nil {
		return err
	}

	return generateGraphs(yamlBytes, filePath)
}

func saveJSONFile(path, filename string, content interface{}) error {
	jsonBytes, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(path, filename)

	err = os.WriteFile(filePath, jsonBytes, 0o600)
	if err != nil {
		return err
	}

	return generateGraphs(jsonBytes, filePath)
}

func blueprintExists(version string, name string) error {
	blueprintsList, err := getBlueprintsByVersion(version)
	if err != nil {
		return err
	}

	if !slices.Contains(blueprintsList, name) {
		return fmt.Errorf("invalid blueprints name '%s'", name)
	}
	return nil
}

func getOutputDir(outputPath string) (string, error) {
	newOutputPath := filepath.Join(outputPath, blueprintsVersion, blueprintName)
	graphDir = filepath.Join(newOutputPath, "graphs")
	err := os.MkdirAll(newOutputPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(graphDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	newOutputPath, err = filepath.Abs(newOutputPath)
	if err != nil {
		return "", err
	}
	return newOutputPath, nil
}

func generateGraphs(content []byte, contentPath string) error {
	policy := &policyv1alpha1.Policy{}
	err := yaml.Unmarshal(content, policy)
	if err != nil || policy.Kind == "" {
		return nil
	}

	fileName := strings.TrimSuffix(filepath.Base(contentPath), filepath.Ext(contentPath))
	dotFilePath := filepath.Join(graphDir, fmt.Sprintf("%s.dot", fileName))
	svgFilePath := filepath.Join(graphDir, fmt.Sprintf("%s.svg", fileName))

	policyFile, err := utils.FetchPolicyFromCR(contentPath)
	if err != nil {
		return nil
	}
	defer os.Remove(policyFile)

	circuit, err := utils.CompilePolicy(policyFile)
	if err != nil {
		return nil
	}

	if err = utils.GenerateDotFile(circuit, dotFilePath); err != nil {
		return err
	}

	graphBytes, err := os.ReadFile(dotFilePath)
	if err != nil {
		return err
	}

	g := graphviz.New()
	graph, err := graphviz.ParseBytes(graphBytes)
	if err != nil {
		return err
	}

	if err := g.RenderFilename(graph, graphviz.SVG, svgFilePath); err != nil {
		return err
	}
	return nil
}
