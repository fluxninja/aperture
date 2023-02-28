package blueprints

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

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
	generateCmd.Flags().BoolVar(&noValidate, "no-validation", false, "Do not validate values.yaml file")
	generateCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing output directory")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Aperture Policy related resources from Aperture Blueprints",
	Long: `
Use this command to generate Aperture Policy related resources like Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.`,
	SilenceErrors: true,
	Example: `aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml

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

		err = blueprintExists(blueprintsDir, blueprintName)
		if err != nil {
			return err
		}
		if !noValidate {
			var valuesBytes []byte
			valuesBytes, err = os.ReadFile(valuesFile)
			if err != nil {
				return err
			}
			if strings.Contains(string(valuesBytes), "__REQUIRED_FIELD__") {
				return fmt.Errorf("values file contains '__REQUIRED_FIELD__' placeholder value")
			}

			// validate values.yaml against the json schema
			schemaFile := filepath.Join(blueprintsDir, blueprintName, "gen/definitions.json")
			definitionsFile := filepath.Join(blueprintsDir, "gen/jsonschema/_definitions.json")
			err = utils.ValidateWithJSONSchema(schemaFile, []string{definitionsFile}, valuesFile)
			if err != nil {
				return err
			}
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

func processJsonnetOutput(bundle map[string]interface{}, outputDir string) error {
	for categoryName, category := range bundle {
		categoriesMap, ok := category.(map[string]interface{})
		if !ok {
			log.Error().Msgf("failed to process output '%+v'", category)
			continue
		}
		var updatedPath string
		if outputDir != "" {
			updatedPath = filepath.Join(outputDir, categoryName)
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

func renderOutput(categoryName string, outputDir string, fileName string, content map[string]interface{}) error {
	if strings.HasSuffix(fileName, ".yaml") {
		if err := saveYAMLFile(categoryName, outputDir, fileName, content); err != nil {
			return err
		}
	} else if strings.HasSuffix(fileName, ".json") {
		if err := saveJSONFile(categoryName, outputDir, fileName, content); err != nil {
			return err
		}
	}
	return nil
}

func saveYAMLFile(categoryName, outputDir, filename string, content map[string]interface{}) error {
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

	outputFilePath := filepath.Join(outputDir, filename)

	err = os.WriteFile(outputFilePath, yamlBytes, 0o600)
	if err != nil {
		return err
	}

	if !noValidate && categoryName == "policies" {
		// validate the generated yaml against the json schema
		schemaFile := filepath.Join(blueprintsDir, "gen/jsonschema/_definitions.json")
		err = utils.ValidateWithJSONSchema(schemaFile, []string{}, outputFilePath)
		if err != nil {
			log.Warn().Msgf("failed to validate generated policy yaml against the json schema: %s", err)
		}
	}

	return generateGraphs(yamlBytes, outputDir, outputFilePath)
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

	return nil
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

func setupOutputDir(outputDir string) (string, error) {
	// ask for user confirmation if the output directory already exists
	if !overwrite {
		if _, err := os.Stat(outputDir); err == nil {
			fmt.Printf("The output directory '%s' already exists. Do you want to overwrite it? [y/N]: ", outputDir)
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				return "", fmt.Errorf("output directory '%s' already exists", outputDir)
			}
		}
	}

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return "", err
	}
	return outputDir, nil
}

func generateGraphs(content []byte, outputDir string, policyPath string) error {
	policy := &policyv1alpha1.Policy{}
	err := yaml.Unmarshal(content, policy)
	if err != nil || policy.Kind != "Policy" {
		return nil
	}

	fileName := strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))

	dir := filepath.Join(filepath.Dir(outputDir), "graphs")
	// create the directory if it doesn't exist
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	dotFilePath := filepath.Join(dir, fmt.Sprintf("%s.dot", fileName))
	mmdFilePath := filepath.Join(dir, fmt.Sprintf("%s.mmd", fileName))

	policyFile, err := utils.FetchPolicyFromCR(policyPath)
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
