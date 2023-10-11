package blueprints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/apply"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/circuitfactory"
)

var controllerConn utils.ControllerConn

func init() {
	controllerConn.InitFlags(generateCmd.PersistentFlags())
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory path where the generated Policy resources will be stored. If not provided, will use current directory")
	generateCmd.Flags().StringVar(&valuesFile, "values-file", "", "Path to the values file")
	generateCmd.Flags().StringVar(&valuesDir, "values-dir", "", "Directory path to the values file(s)")
	generateCmd.Flags().BoolVar(&applyPolicies, "apply", false, "Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed")
	generateCmd.Flags().BoolVar(&noYAMLModeline, "no-yaml-modeline", false, "Do not add YAML language server modeline to generated YAML files")
	generateCmd.Flags().BoolVar(&noValidate, "no-validation", false, "Do not validate values.yaml file")
	generateCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing output directory")
	generateCmd.Flags().IntVar(&graphDepth, "graph-depth", 1, "Max depth of the graph when generating DOT and Mermaid files")
	generateCmd.Flags().BoolVarP(&force, "force", "f", false, "Force apply policy even if it already exists")
	generateCmd.Flags().BoolVar(&selectAll, "select-all", false, "Select all blueprints")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Aperture Policy related resources from Aperture Blueprints",
	Long: `
Use this command to generate Aperture Policy related resources like Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.`,
	SilenceErrors: true,
	Example: `aperturectl blueprints generate --values-file=rate-limiting.yaml

aperturectl blueprints generate --values-file=rate-limiting.yaml --apply`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var valuesFiles []string
		var overrideBlueprintsURI, overrideBlueprintsVersion string

		// HACK: As we use global variables (unfortunately) we need to reset them
		// in each loop iteration.
		if blueprintsURI != "" {
			overrideBlueprintsURI = blueprintsURI
		}
		if blueprintsVersion != "" {
			overrideBlueprintsVersion = blueprintsVersion
		}

		// if outputDir is not provided and applyPolicy is false, return error
		if outputDir == "" && !applyPolicies {
			return fmt.Errorf("--output-dir must be provided")
		}

		updatedOutputDir, doRemove, e := setupOutputDir(outputDir)
		if e != nil {
			return e
		}
		if doRemove {
			// defer remove temp dir
			defer os.RemoveAll(updatedOutputDir)
		}

		// one of valuesFile or valuesDir must be provided
		if valuesFile == "" && valuesDir == "" {
			return fmt.Errorf("--values-file or --values-dir must be provided")
		} else if valuesFile != "" && valuesDir != "" {
			return fmt.Errorf("--values-file and --values-dir cannot be provided at the same time")
		} else if valuesFile != "" {
			valuesFiles = append(valuesFiles, valuesFile)
		} else if valuesDir != "" {
			// get all files in the directory
			allBlueprints, model, err := utils.GetBlueprintsTUIModel(valuesDir, selectAll)
			if err != nil {
				return err
			}
			for index := range model.Selected {
				valuesFiles = append(valuesFiles, allBlueprints[index])
			}
		}

		generate := func(vFile string) error {
			_, err := os.Stat(vFile)
			if err != nil {
				return err
			}
			vFile, err = filepath.Abs(vFile)
			if err != nil {
				return err
			}

			var valuesBytes []byte
			valuesBytes, err = os.ReadFile(vFile)
			if err != nil {
				return err
			}

			// decode values as yaml and retrieve blueprint and uri fields
			var values map[string]interface{}
			err = yaml.Unmarshal(valuesBytes, &values)
			if err != nil {
				return err
			}

			var ok bool
			blueprintName, ok = values["blueprint"].(string)
			if !ok {
				return fmt.Errorf("values file does not contain blueprint name field")
			}

			if overrideBlueprintsVersion == "" && overrideBlueprintsURI == "" {
				blueprintsURI, ok = values["uri"].(string)
				if !ok {
					return fmt.Errorf("values file does not contain uri field")
				}
			} else {
				// HACK: global variables suck
				blueprintsURI = overrideBlueprintsURI
				blueprintsVersion = overrideBlueprintsVersion
			}

			policy, ok := values["policy"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("values file does not contain policy field")
			}
			policyName, ok := policy["policy_name"].(string)
			if !ok {
				return fmt.Errorf("values file does not contain policy_name field")
			}

			// pull
			err = pullCmd.RunE(cmd, args)
			if err != nil {
				return err
			}
			defer func() {
				_ = pullCmd.PostRunE(cmd, args)
			}()

			err = blueprintExists(blueprintName)
			if err != nil {
				return err
			}
			// check whether the blueprint is deprecated
			blueprintDir := filepath.Join(blueprintsDir, blueprintName)
			ok, deprecated := utils.IsBlueprintDeprecated(blueprintDir)
			if ok {
				if utils.AllowDeprecated {
					log.Warn().Msgf("Blueprint %s is deprecated: %s", blueprintName, deprecated)
				} else {
					return fmt.Errorf("blueprint %s is deprecated: %s", blueprintName, deprecated)
				}
			}

			if !noValidate {
				if strings.Contains(string(valuesBytes), "__REQUIRED_FIELD__") {
					return fmt.Errorf("values file contains '__REQUIRED_FIELD__' placeholder value")
				}

				// validate values.yaml against the json schema
				schemaFile := filepath.Join(blueprintsDir, blueprintName, "gen/definitions.json")
				err = utils.ValidateWithJSONSchema(schemaFile, []string{}, vFile)
				if err != nil {
					log.Error().Msgf("Error validating values file: %s", err)
					return err
				}
			}

			buf := bytes.Buffer{}
			vm := jsonnet.MakeVM()
			vm.SetTraceOut(&buf)
			vm.Importer(&jsonnet.FileImporter{
				JPaths: []string{blueprintsURIRoot},
			})

			importPath := fmt.Sprintf("%s/%s", blueprintsDir, blueprintName)
			bundleStr, err := vm.EvaluateAnonymousSnippet("bundle.libsonnet", fmt.Sprintf(`
		local bundle = import '%s/bundle.libsonnet';
		local config = std.parseYaml(importstr '%s');
		bundle(config)
		`, importPath, vFile))
			if err != nil {
				return err
			}

			// change log.Error() to log.Info() for jsonnet trace output
			// Note use this with std.trace in jsonnet code to get the trace output
			log.Debug().Msgf("Jsonnet generation trace: %s", buf.String())

			var bundle map[string]interface{}
			err = json.Unmarshal([]byte(bundleStr), &bundle)
			if err != nil {
				return err
			}

			if err = processJsonnetOutput(bundle, updatedOutputDir, policyName); err != nil {
				return err
			}

			log.Info().Msgf("Generated manifests at %s", updatedOutputDir)

			return nil
		}

		for _, vFile := range valuesFiles {
			err := generate(vFile)
			if err != nil {
				return err
			}
		}

		if applyPolicies {
			apply.Controller = controllerConn
			err := apply.ApplyPolicyCmd.Flag("dir").Value.Set(updatedOutputDir)
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
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func processJsonnetOutput(bundle map[string]interface{}, outputDir, policyName string) error {
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
			if err := renderOutput(categoryName, updatedPath, fileName, policyName, contentMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func renderOutput(categoryName string, outputDir string, fileName, policyName string, content map[string]interface{}) error {
	outputFilePath := filepath.Join(outputDir, fileName)
	if strings.HasSuffix(fileName, ".yaml") {
		yamlBytes, err := saveYAMLFile(categoryName, outputFilePath, content)
		if err != nil {
			return err
		}

		if strings.HasSuffix(fileName, "-cr.yaml") {
			// prepare data
			circuit, componentsList, err := processPolicy(yamlBytes, outputFilePath)
			if err != nil {
				return err
			}

			err = generateDashboards(string(yamlBytes), componentsList, policyName, outputDir)
			if err != nil {
				return err
			}

			err = generateGraphs(circuit, outputDir, outputFilePath, graphDepth)
			if err != nil {
				return err
			}
		}
	} else if strings.HasSuffix(fileName, ".json") {
		if err := saveJSONFile(outputFilePath, content); err != nil {
			return err
		}
	}
	return nil
}

func saveYAMLFile(categoryName, outputFilePath string, content map[string]interface{}) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(content)
	if err != nil {
		return nil, err
	}

	if !noYAMLModeline {
		if categoryName == "policies" {
			var schemaURL string
			defPath := "gen/jsonschema/_definitions.json#/definitions"
			prefix := fmt.Sprintf("file:%s/%s", blueprintsDir, defPath)
			contentURL := utils.URIToRawContentURL(blueprintsURI)
			if contentURL != "" {
				prefix = fmt.Sprintf("%s/%s", contentURL, defPath)
			}

			// if content contains kind: Policy then use the policy schema
			if kind, ok := content["kind"]; ok && kind == "Policy" {
				schemaURL = fmt.Sprintf("%s/PolicyCustomResource", prefix)
			} else {
				schemaURL = fmt.Sprintf("%s/Policy", prefix)
			}
			yamlBytes = append([]byte(fmt.Sprintf("# yaml-language-server: $schema=%s\n", schemaURL)), yamlBytes...)
		}
	}

	err = os.WriteFile(outputFilePath, yamlBytes, 0o600)
	if err != nil {
		return nil, err
	}

	if !noValidate && categoryName == "policies" {
		// validate the generated yaml against the json schema
		schemaFile := filepath.Join(blueprintsDir, "gen/jsonschema/_definitions.json")
		err = utils.ValidateWithJSONSchema(schemaFile, []string{}, outputFilePath)
		if err != nil {
			log.Warn().Msgf("failed to validate generated policy yaml against the json schema: %s", err)
		}
	}

	return yamlBytes, nil
}

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

func blueprintExists(name string) error {
	blueprintsList, err := getBlueprints(blueprintsURIRoot, true)
	if err != nil {
		return err
	}

	if !slices.Contains(blueprintsList, name) {
		return fmt.Errorf("invalid blueprint name '%s'", name)
	}
	return nil
}

func setupOutputDir(outputDir string) (string, bool, error) {
	if outputDir == "" && !applyPolicies {
		// shouldn't happen
		return "", false, fmt.Errorf("output directory is empty")
	}
	if outputDir == "" {
		// Apply policy scenario
		var err error
		// use temp dir
		outputDir, err = os.MkdirTemp("", "aperturectl-generate-*")
		if err != nil {
			log.Error().Msgf("Error creating temp dir: %s", err)
			return "", false, err
		}
		log.Info().Msgf("Using temp dir: %s", outputDir)
		return outputDir, true, nil
	} else {
		// ask for user confirmation if the output directory already exists
		if !overwrite {
			if _, err := os.Stat(outputDir); err == nil {
				fmt.Printf("The output directory '%s' already exists. Do you want to merge the generated policy artifacts into the existing directory? [y/N]: ", outputDir)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) != "y" {
					return "", false, fmt.Errorf("output directory '%s' already exists", outputDir)
				}
			}
		}

		err := os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return "", false, err
		}
		absOutputDir, err := filepath.Abs(outputDir)
		if err != nil {
			return "", false, err
		}
		return absOutputDir, false, nil
	}
}

func generateGraphs(circuit *circuitfactory.Circuit, outputDir string, policyPath string, depth int) error {
	dir := filepath.Join(filepath.Dir(outputDir), "graphs")
	// create the directory if it does not exist
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))
	dotFilePath := filepath.Join(dir, fmt.Sprintf("%s.dot", fileName))
	mmdFilePath := filepath.Join(dir, fmt.Sprintf("%s.mmd", fileName))

	if err = utils.GenerateDotFile(circuit, dotFilePath, depth); err != nil {
		return err
	}

	if err = utils.GenerateMermaidFile(circuit, mmdFilePath, depth); err != nil {
		return err
	}

	return nil
}

func processPolicy(content []byte, policyPath string) (*circuitfactory.Circuit, string, error) {
	policy := &policyv1alpha1.Policy{}
	err := yaml.Unmarshal(content, policy)
	if err != nil || policy.Kind != "Policy" {
		return nil, "", err
	}

	policyFile, err := utils.FetchPolicyFromCR(policyPath)
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(policyFile)
	policyBytes, err := os.ReadFile(policyFile)
	if err != nil {
		return nil, "", err
	}
	circuit, _, err := utils.CompilePolicy(filepath.Base(policyFile), policyBytes)
	if err != nil {
		return nil, "", err
	}
	componentsList, err := utils.GetFlatComponentsList(circuit)
	if err != nil {
		return nil, "", err
	}

	return circuit, componentsList, nil
}

func generateDashboards(policyFile, componentsList, policyName, outputDir string) error {
	dir := filepath.Join(filepath.Dir(outputDir), "dashboards")
	// create the directory if it does not exist
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	vm := jsonnet.MakeVM()
	vm.SetTraceOut(&buf)
	vm.Importer(&jsonnet.FileImporter{
		JPaths: []string{blueprintsURIRoot},
	})

	vm.TLAReset()
	vm.TLAVar("policyFile", policyFile)
	vm.TLAVar("componentsList", componentsList)
	vm.TLAVar("policyName", policyName)
	vm.TLAVar("datasource", "controller-prometheus")

	dashboardGroupFile := filepath.Join(blueprintsDir, "grafana", "dashboard_group.libsonnet")
	dashboardsJSON, err := vm.EvaluateFile(dashboardGroupFile)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Jsonnet generation trace: %s", buf.String())

	type dashboards struct {
		Dashboards map[string]interface{} `json:"dashboards"`
	}

	var dashboardsList dashboards
	err = json.Unmarshal([]byte(dashboardsJSON), &dashboardsList)
	if err != nil {
		return err
	}

	for key, val := range dashboardsList.Dashboards {
		outputFilePath := filepath.Join(dir, fmt.Sprintf("%s-%s.json", key, policyName))
		err = saveJSONFile(outputFilePath, val.(map[string]interface{}))
		if err != nil {
			return err
		}
	}

	return nil
}
