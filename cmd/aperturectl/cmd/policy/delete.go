package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// DeleteCmd is the command to delete a policy from the Aperture Controller.
var DeleteCmd = &cobra.Command{
	Use:           "delete POLICY_NAME",
	Short:         "Delete Aperture Policy from the Aperture Controller",
	Long:          `Use this command to delete the Aperture Policy from the Aperture Controller.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl policy delete POLICY_NAME`,
	RunE: func(_ *cobra.Command, args []string) error {
		return deletePolicy(args[0])
	},
}

// deletePolicy deletes the policy from the cluster.
func deletePolicy(policyName string) error {
	if Controller.IsKube() {
		deployment, err := utils.GetControllerDeployment(Controller.GetKubeRestConfig(), controllerNs)
		if err != nil {
			return err
		}

		err = api.AddToScheme(scheme.Scheme)
		if err != nil {
			return fmt.Errorf("failed to connect to Kubernetes: %w", err)
		}

		kubeClient, err := k8sclient.New(Controller.GetKubeRestConfig(), k8sclient.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		policyCR := &policyv1alpha1.Policy{}
		policyCR.Name = policyName
		policyCR.Namespace = deployment.GetNamespace()
		err = kubeClient.Delete(context.Background(), policyCR)
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Info().Str("policy", policyName).Msg("Policy not found in Kubernetes")
				return nil
			}

			if utils.IsNoMatchError(err) {
				err = utils.DeletePolicyUsingAPI(client, policyName)
			}

			if err != nil {
				return fmt.Errorf("failed to delete policy from Kubernetes: %w", err)
			}
		}
	} else {
		err := utils.DeletePolicyUsingAPI(client, policyName)
		if err != nil {
			return fmt.Errorf("failed to delete policy: %w", err)
		}
	}

	log.Info().Str("policy", policyName).Msg("Deleted Policy successfully")
	return nil
}
