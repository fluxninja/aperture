package blueprints

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/apply"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/log"
)

func init() {
	generateCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to generate Aperture Policy resources for")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory path where the generated Policy resources will be stored. If not provided, will use current directory")
	generateCmd.Flags().StringVar(&valuesFile, "values-file", "", "Path to the values file for Blueprint's input")
	generateCmd.Flags().BoolVar(&applyPolicy, "apply", false, "Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed")
	generateCmd.Flags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernetes cluster config. Defaults to '~/.kube/config'")
	generateCmd.Flags().BoolVar(&noYAMLModeline, "no-yaml-modeline", false, "Do not add YAML language server modeline to generated YAML files")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Aperture Policy related resources from Aperture Blueprints",
	Long: `
Use this command to generate Aperture Policy related resources like Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.`,
	SilenceErrors: true,
	Example: `aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml

aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml --version latest

aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml --apply`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := readerLock()
		if err != nil {
			return err
		}
		defer unlock()

		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}

		if valuesFile == "" {
			return fmt.Errorf("--values-file must be provided")
		}

		// if outputDir is not provided, default to current directory
		if outputDir == "" {
			if applyPolicy {
				// use temp dir
				outputDir = os.TempDir()
				// defer remove temp dir
				defer os.RemoveAll(outputDir)
			} else {
				return fmt.Errorf("--output-dir must be provided")
			}
		}

		_, err = os.Stat(valuesFile)
		if err != nil {
			return err
		}
		valuesFile, err = filepath.Abs(valuesFile)
		if err != nil {
			return err
		}

		blueprintsDir = filepath.Join(blueprintsURIRoot, getRelPath(blueprintsURIRoot))

		err = blueprintExists(blueprintsDir, blueprintName)
		if err != nil {
			return err
		}

		vm := jsonnet.MakeVM()
		vm.Importer(&jsonnet.FileImporter{
			JPaths: []string{blueprintsURIRoot},
		})

		importPath := fmt.Sprintf("%s/%s", blueprintsDir, blueprintName)

		bundleStr, err := vm.EvaluateAnonymousSnippet("bundle.libsonnet", fmt.Sprintf(`
		local bundle = import '%s/bundle.libsonnet';
		local config = std.parseYaml(importstr '%s');
    bundle(config)
		`, importPath, valuesFile))
		if err != nil {
			return err
		}

		var bundle map[string]interface{}
		err = json.Unmarshal([]byte(bundleStr), &bundle)
		if err != nil {
			return err
		}

		var updatedOutputDir string

		if outputDir != "" {
			updatedOutputDir, err = setupOutputDir(outputDir)
			if err != nil {
				return err
			}
		}

		if err = processJsonnetOutput(bundle, updatedOutputDir); err != nil {
			return err
		}

		log.Info().Msgf("Generated manifests at %s", updatedOutputDir)

		if applyPolicy {
			err = apply.ApplyPolicyCmd.Flag("dir").Value.Set(updatedOutputDir)
			if err != nil {
				return err
			}

			err = apply.ApplyCmd.PersistentPreRunE(cmd, args)
			if err != nil {
				return err
			}

			err = apply.ApplyPolicyCmd.RunE(cmd, args)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func processJsonnetOutput(bundle map[string]interface{}, outputPath string) error {
	for categoryName, category := range bundle {
		categoriesMap, ok := category.(map[string]interface{})
		if !ok {
			log.Error().Msgf("failed to process output '%+v'", category)
			continue
		}
		var updatedPath string
		if outputPath != "" {
			updatedPath = filepath.Join(outputPath, categoryName)
			err := os.MkdirAll(updatedPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		for fileName, content := range categoriesMap {
			contentMap, ok := content.(map[string]interface{})
			if !ok {
				log.Error().Msgf("failed to process output '%+v'", content)
				continue
			}
			if err := renderOutput(categoryName, updatedPath, fileName, contentMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderOutput(categoryName string, outputPath string, fileName string, content map[string]interface{}) error {
	if strings.HasSuffix(fileName, ".yaml") {
		if err := saveYAMLFile(categoryName, outputPath, fileName, content); err != nil {
			return err
		}
	} else if strings.HasSuffix(fileName, ".json") {
		if err := saveJSONFile(categoryName, outputPath, fileName, content); err != nil {
			return err
		}
	}
	return nil
}

func saveYAMLFile(categoryName, path, filename string, content map[string]interface{}) error {
	yamlBytes, err := yaml.Marshal(content)
	if err != nil {
		return err
	}

	if !noYAMLModeline {
		if categoryName == "policies" {
			var schemaURL string
			// if content contains kind: Policy then use the policy schema
			if kind, ok := content["kind"]; ok && kind == "Policy" {
				schemaURL = fmt.Sprintf("file:%s", filepath.Join(blueprintsDir, "gen/jsonschema/_definitions.json#/definitions/PolicyCustomResource"))
			} else {
				// prepend the file with yaml-language-server modeline that points to the schema
				schemaURL = fmt.Sprintf("file:%s", filepath.Join(blueprintsDir, "gen/jsonschema/_definitions.json#/definitions/Policy"))
			}
			yamlBytes = append([]byte(fmt.Sprintf("# yaml-language-server: $schema=%s\n", schemaURL)), yamlBytes...)
		}
	}

	filePath := filepath.Join(path, filename)

	err = os.WriteFile(filePath, yamlBytes, 0o600)
	if err != nil {
		return err
	}

	return generateGraphs(yamlBytes, filePath)
}

func saveJSONFile(_, path, filename string, content map[string]interface{}) error {
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

func blueprintExists(blueprintsDir, name string) error {
	blueprintsList, err := getBlueprints(blueprintsDir)
	if err != nil {
		return err
	}

	if !slices.Contains(blueprintsList, name) {
		return fmt.Errorf("invalid blueprints name '%s'", name)
	}
	return nil
}

func setupOutputDir(outputPath string) (string, error) {
	graphDir = filepath.Join(outputPath, "graphs")
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(graphDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	outputPath, err = filepath.Abs(outputPath)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

func generateGraphs(content []byte, contentPath string) error {
	policy := &policyv1alpha1.Policy{}
	err := yaml.Unmarshal(content, policy)
	if err != nil || policy.Kind == "" {
		return nil
	}

	fileName := strings.TrimSuffix(filepath.Base(contentPath), filepath.Ext(contentPath))
	dotFilePath := filepath.Join(graphDir, fmt.Sprintf("%s.dot", fileName))
	mmdFilePath := filepath.Join(graphDir, fmt.Sprintf("%s.mmd", fileName))

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

	if err = utils.GenerateMermaidFile(circuit, mmdFilePath); err != nil {
		return err
	}

	return nil
}
