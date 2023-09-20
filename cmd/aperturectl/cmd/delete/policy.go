package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// DeletePolicyCmd is the command to delete a policy from the cluster.
var DeletePolicyCmd = &cobra.Command{
	Use:           "policy",
	Short:         "Delete Aperture Policy from the cluster",
	Long:          `Use this command to delete the Aperture Policy from the cluster.`,
	SilenceErrors: true,
	Example:       `aperturectl delete policy --policy=rate-limiting`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return deletePolicy()
	},
}

// deletePolicy deletes the policy from the cluster.
func deletePolicy() error {
	if controller.IsKube() {
		deployment, err := utils.GetControllerDeployment(controller.GetKubeRestConfig(), controllerNs)
		if err != nil {
			return err
		}

		err = api.AddToScheme(scheme.Scheme)
		if err != nil {
			return fmt.Errorf("failed to connect to Kubernetes: %w", err)
		}

		kubeClient, err := k8sclient.New(controller.GetKubeRestConfig(), k8sclient.Options{
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

			if apimeta.IsNoMatchError(err) {
				err = deletePolicyUsingAPI()
			}

			if err != nil {
				return fmt.Errorf("failed to delete policy from Kubernetes: %w", err)
			}
		}
	} else {
		err := deletePolicyUsingAPI()
		if err != nil {
			return fmt.Errorf("failed to delete policy: %w", err)
		}
	}

	log.Info().Str("policy", policyName).Msg("Deleted Policy successfully")
	return nil
}

// deletePolicyUsingAPI deletes the policy using the API.
func deletePolicyUsingAPI() error {
	policyRequest := languagev1.DeletePolicyRequest{
		Name: policyName,
	}
	_, err := client.DeletePolicy(context.Background(), &policyRequest)
	if err != nil {
		log.Warn().Err(err).Str("policy", policyName).Msg("failed to delete Policy")
	}

	return nil
}
