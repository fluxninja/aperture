package blueprints

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/fluxninja/aperture/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
			return fmt.Errorf("--policy_type must be provided")
		} else if valuesFile == "" {
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
		fmt.Printf("Using '%s' as values file\n", valuesFile)

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

		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)

		vm := jsonnet.MakeVM()

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
			outputDir = filepath.Join(outputDir, blueprintsVersion, policyType)
			err = os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				return err
			}

			outputDir, err = filepath.Abs(outputDir)
			if err != nil {
				return err
			}
			policyFilePath := filepath.Join(outputDir, "policy.yaml")
			err = os.WriteFile(policyFilePath, yamlBytes, 0o600)
			if err != nil {
				return err
			}
			fmt.Printf("Stored policy YAML at '%s'\n", policyFilePath)
		} else {
			fmt.Println(string(yamlBytes))
		}

		if apply {
			err = applyPolicy(kubeConfig, yamlBytes)
			if err != nil {
				return fmt.Errorf("failed to apply policy to Kubernetes. Error: '%s'", err.Error())
			}
		}

		return nil
	},
}

func applyPolicy(kubeConfig string, policyBytes []byte) error {
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

	return createPolicy(kubeRestConfig, policyBytes)
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

func createPolicy(kubeRestConfig *rest.Config, policyBytes []byte) error {
	policy := &policyv1alpha1.Policy{}
	err := config.UnmarshalYAML(policyBytes, policy)
	if err != nil {
		return fmt.Errorf("failed to prepare policy object for Kubernetes. Error: %s", err.Error())
	}

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

	policy.Namespace = deployment.GetNamespace()
	spec := policy.Spec
	_, err = controllerutil.CreateOrUpdate(context.Background(), client, policy, func() error {
		policy.Spec = spec
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to apply policy in Kubernetes. Error: %s", err.Error())
	}

	fmt.Printf("Applied policy '%s' in '%s' namespace.\n", policy.GetName(), policy.GetNamespace())
	return nil
}
