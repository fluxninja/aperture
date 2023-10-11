package apply

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var (
	policyName         string
	dynamicConfigFile  string
	dynamicConfigBytes []byte
)

func init() {
	ApplyDynamicConfigCmd.Flags().StringVar(&policyName, "policy", "", "Name of the Policy to apply the DynamicConfig to")
	ApplyDynamicConfigCmd.Flags().StringVar(&dynamicConfigFile, "file", "", "Path to the dynamic config file")
}

// ApplyDynamicConfigCmd is the command to apply DynamicConfig to a Policy.
var ApplyDynamicConfigCmd = &cobra.Command{
	Use:           "dynamic-config",
	Short:         "Apply Aperture DynamicConfig to a Policy",
	Long:          `Use this command to apply the Aperture DynamicConfig to a Policy.`,
	SilenceErrors: true,
	Example:       `aperturectl apply dynamic-config --policy=rate-limiting --file=dynamic-config.yaml`,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		var err error
		dynamicConfigBytes, err = utils.GetDynamicConfigBytes(policyName, dynamicConfigFile)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		err := api.AddToScheme(scheme.Scheme)
		if err != nil {
			return fmt.Errorf("failed to connect to Kubernetes: %w", err)
		}

		dynamicConfigYAML := make(map[string]interface{})
		err = yaml.Unmarshal(dynamicConfigBytes, &dynamicConfigYAML)
		if err != nil {
			return fmt.Errorf("failed to parse DynamicConfig YAML: %w", err)
		}

		if Controller.IsKube() {
			var kubeClient k8sclient.Client
			kubeClient, err = k8sclient.New(Controller.GetKubeRestConfig(), k8sclient.Options{
				Scheme: scheme.Scheme,
			})
			if err != nil {
				return fmt.Errorf("failed to create Kubernetes client: %w", err)
			}

			var deployment *appsv1.Deployment
			deployment, err = utils.GetControllerDeployment(Controller.GetKubeRestConfig(), controllerNs)
			if err != nil {
				return err
			}

			dynamicConfigBytes, err = json.Marshal(dynamicConfigYAML)
			if err != nil {
				return fmt.Errorf("failed to parse DynamicConfig JSON: %w", err)
			}

			policy := &policyv1alpha1.Policy{}
			err = kubeClient.Get(context.Background(), k8sclient.ObjectKey{
				Namespace: deployment.Namespace,
				Name:      policyName,
			}, policy)
			if err != nil {
				if utils.IsNoMatchError(err) {
					err = utils.ApplyDynamicConfigUsingAPI(client, dynamicConfigYAML, policyName)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("failed to get Policy '%s': %w", policyName, err)
				}
			} else {
				policy.DynamicConfig.Raw = dynamicConfigBytes
				err = kubeClient.Update(context.Background(), policy)
				if err != nil {
					return fmt.Errorf("failed to update DynamicConfig for policy '%s': %w", policyName, err)
				}
			}
		} else {
			err = utils.ApplyDynamicConfigUsingAPI(client, dynamicConfigYAML, policyName)
			if err != nil {
				return err
			}
		}

		log.Info().Str("policy", policyName).Msg("Updated DynamicConfig successfully")

		return nil
	},
}
