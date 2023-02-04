package blueprints

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/log"
)

var (
	blueprintName string
	outputDir     string
	valuesFile    string
	graphDir      string
	apply         bool
	kubeConfig    string
	validPolicies []*policyv1alpha1.Policy
)

func init() {
	generateCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to generate Aperture Policy resources for. Can be skipped when '--custom-blueprint-path' is provided")
	generateCmd.Flags().StringVar(&outputDir, "output-dir", "", "Directory path where the generated Policy resources will be stored. If not provided, will use current directory")
	generateCmd.Flags().StringVar(&valuesFile, "values-file", "", "Path to the values file for Blueprint's input")
	generateCmd.Flags().BoolVar(&apply, "apply", false, "Apply generated policies on the Kubernetes cluster in the namespace where Aperture Controller is installed")
	generateCmd.Flags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernetes cluster config. Defaults to '~/.kube/config'")

	validPolicies = []*policyv1alpha1.Policy{}
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Aperture Policy related resources from Aperture Blueprints",
	Long: `
Use this command to generate Aperture Policy related resources like Kubernetes Custom Resource, Grafana Dashboards and graphs in DOT and Mermaid format.`,
	SilenceErrors: true,
	Example: `aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml

aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml --version v0.22.0

aperturectl blueprints generate --name=policies/static-rate-limiting --values-file=rate-limiting.yaml --apply

aperturectl blueprints generate --custom-blueprint-path=/path/to/blueprint/ --values-file=values.yaml`,
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

		apertureDir := filepath.Join(blueprintsDir, getRelPath(blueprintsDir))

		err = blueprintExists(apertureDir, blueprintName)
		if err != nil {
			return err
		}

		vm := jsonnet.MakeVM()
		vm.Importer(&jsonnet.FileImporter{
			JPaths: []string{blueprintsDir},
		})

		importPath := fmt.Sprintf("%s/%s", apertureDir, blueprintName)

		jsonStr, err := vm.EvaluateAnonymousSnippet("bundle.libsonnet", fmt.Sprintf(`
		local bundle = import '%s/bundle.libsonnet';
		local config = std.parseYaml(importstr '%s');
    bundle(config)
		`, importPath, valuesFile))
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
			updatedOutputDir, err = setupOutputDir(outputDir)
			if err != nil {
				return err
			}
		}

		if err = processContent(yamlData, updatedOutputDir); err != nil {
			return err
		}

		log.Info().Msgf("Generated manifests at %s.", updatedOutputDir)

		if apply {
			if err = applyPolicy(); err != nil {
				return err
			}
		}

		return nil
	},
}

func processContent(data map[string]interface{}, outputPath string) error {
	for name, content := range data {
		if strings.HasSuffix(name, ".yaml") {
			if err := saveYAMLFile(outputPath, name, content); err != nil {
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

func saveYAMLFile(path, filename string, content interface{}) error {
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

	validPolicies = append(validPolicies, policy)
	return nil
}

func applyPolicy() error {
	if kubeConfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		kubeConfig = filepath.Join(homeDir, ".kube", "config")
		fmt.Printf("Using '%s' as Kubernetes config\n", kubeConfig)
	}

	kubeRestConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes. Error: %s", err.Error())
	}

	return createPolicy(kubeRestConfig)
}

func getControllerDeployment(kubeRestConfig *rest.Config) (*appsv1.Deployment, error) {
	clientSet, err := kubernetes.NewForConfig(kubeRestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kubernetes. Error: %s", err.Error())
	}

	deployment, err := clientSet.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{
		FieldSelector: "metadata.name=aperture-controller",
	})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf(
				"no deployment with name 'aperture-controller' found on the Kubernetes cluster. The policy can be only applied in the namespace where the Aperture Controller is running")
		}
		return nil, fmt.Errorf("failed to fetch aperture-controller namespace in Kubernetes. Error: %s", err.Error())
	}

	if len(deployment.Items) != 1 {
		return nil, fmt.Errorf(
			"no deployment with name 'aperture-controller' found on the Kubernetes cluster. The policy can be only applied in the namespace where the Aperture Controller is running")
	}

	return &deployment.Items[0], nil
}

func createPolicy(kubeRestConfig *rest.Config) error {
	deployment, err := getControllerDeployment(kubeRestConfig)
	if err != nil {
		return err
	}

	err = api.AddToScheme(scheme.Scheme)
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes. Error: %s", err.Error())
	}
	client, err := client.New(kubeRestConfig, client.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes. Error: %s", err.Error())
	}

	for _, policy := range validPolicies {
		policy.Namespace = deployment.GetNamespace()
		spec := policy.Spec
		_, err = controllerutil.CreateOrUpdate(context.Background(), client, policy, func() error {
			policy.Spec = spec
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to apply policy in Kubernetes. Error: %s", err.Error())
		}

		log.Info().Msgf("Applied policy '%s' in '%s' namespace.\n", policy.GetName(), policy.GetNamespace())
	}
	return nil
}
