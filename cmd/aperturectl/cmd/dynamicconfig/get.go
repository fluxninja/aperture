package dynamicconfig

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// GetCmd is a command to get a policy's dynamic config.
var GetCmd = &cobra.Command{
	Use:           "get POLICY_NAME",
	Short:         "Get Aperture DynamicConfig for a Policy.",
	Long:          "Use this command to get the Aperture DynamicConfig of a Policy.",
	Example:       "aperture dynamic-config get rate-limiting",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		policyName := args[0]
		if controller.IsKube() {
			err := api.AddToScheme(scheme.Scheme)
			if err != nil {
				return fmt.Errorf("failed to connect to Kubernetes: %w", err)
			}

			var kubeClient k8sclient.Client
			kubeClient, err = k8sclient.New(controller.GetKubeRestConfig(), k8sclient.Options{
				Scheme: scheme.Scheme,
			})
			if err != nil {
				return fmt.Errorf("failed to create Kubernetes client: %w", err)
			}

			var deployment *appsv1.Deployment
			deployment, err = utils.GetControllerDeployment(controller.GetKubeRestConfig(), controllerNs)
			if err != nil {
				return fmt.Errorf("failed to get controller deployment: %w", err)
			}

			policy := &policyv1alpha1.Policy{}
			err = kubeClient.Get(context.Background(), k8sclient.ObjectKey{
				Namespace: deployment.Namespace,
				Name:      policyName,
			}, policy)
			if err != nil {
				if utils.IsNoMatchError(err) {
					err = utils.GetDynamicConfigUsingAPI(client, policyName)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("failed to get Dynamic Config '%s': %w", policyName, err)
				}
			}
			if len(policy.DynamicConfig.Raw) == 0 {
				log.Info().Str("policy-name", policyName).Msg("DynamicConfig is not set for the given Policy")
				return nil
			}
			j, err := policy.DynamicConfig.MarshalJSON()
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}

			yamlString, err := utils.GetYAMLString(j)
			if err != nil {
				return fmt.Errorf("failed to convert JSON to YAML: %w", err)
			}
			fmt.Print(yamlString)
		} else {
			err := utils.GetDynamicConfigUsingAPI(client, policyName)
			if err != nil {
				return fmt.Errorf("failed to get dynamic config using API: %w", err)
			}
		}
		return nil
	},
}
