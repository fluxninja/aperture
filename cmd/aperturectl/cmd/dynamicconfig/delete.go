package dynamicconfig

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	languagev1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
)

// DelCmd is a command to delete a policy's dynamic config.
var DelCmd = &cobra.Command{
	Use:           "delete POLICY_NAME",
	Short:         "Delete Aperture DynamicConfig of a Policy.",
	Long:          "Use this command to delete the Aperture DynamicConfig to a Policy.",
	Example:       "aperture dynamic-config delete rate-limiting",
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
				return err
			}

			policy := &policyv1alpha1.Policy{}
			err = kubeClient.Get(context.Background(), k8sclient.ObjectKey{
				Namespace: deployment.Namespace,
				Name:      policyName,
			}, policy)
			if err != nil {
				if utils.IsNoMatchError(err) {
					_, err = client.DeleteDynamicConfig(context.Background(), &languagev1.DeleteDynamicConfigRequest{
						PolicyName: policyName,
					})
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("failed to get Dynamic Config '%s': %w", policyName, err)
				}
			}
			if policy.DynamicConfig.Raw != nil {
				policy.DynamicConfig = runtime.RawExtension{}
				err = kubeClient.Update(context.Background(), policy)
				if err != nil {
					return fmt.Errorf("failed to delete DynamicConfig for policy '%s': %w", policyName, err)
				}
			} else {
				fmt.Println("No DynamicConfig found for policy", policyName)
				return nil
			}
		} else {
			_, err := client.DeleteDynamicConfig(context.Background(), &languagev1.DeleteDynamicConfigRequest{
				PolicyName: policyName,
			})
			if err != nil {
				return err
			}
		}
		fmt.Println("Deleted DynamicConfig successfully")
		return nil
	},
}
