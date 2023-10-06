package apply

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var (
	file      string
	dir       string
	force     bool
	selectAll bool
)

func init() {
	ApplyPolicyCmd.Flags().StringVar(&file, "file", "", "Path to Aperture Policy file")
	ApplyPolicyCmd.Flags().StringVar(&dir, "dir", "", "Path to directory containing Aperture Policy files")
	ApplyPolicyCmd.Flags().BoolVarP(&force, "force", "f", false, "Force apply policy even if it already exists")
	ApplyPolicyCmd.Flags().BoolVarP(&selectAll, "select-all", "s", false, "Apply all policies in the directory")
}

// ApplyPolicyCmd is the command to apply a policy to the cluster.
var ApplyPolicyCmd = &cobra.Command{
	Use:           "policy",
	Short:         "Apply Aperture Policy to the cluster",
	Long:          `Use this command to apply the Aperture Policy to the cluster.`,
	SilenceErrors: true,
	Example: `aperturectl apply policy --file=policies/rate-limiting.yaml

aperturectl apply policy --dir=policies`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if file != "" {
			return applyPolicy(file)
		} else if dir != "" {
			policies, model, err := utils.GetPoliciesTUIModel(dir, selectAll)
			if err != nil {
				return err
			}

			for policyIndex := range model.Selected {
				fileName := policies[policyIndex]
				if err := applyPolicy(fileName); err != nil {
					log.Error().Err(err).Msgf("failed to apply policy '%s'.", fileName)
				}
			}
			return nil
		} else {
			return errors.New("either --file or --dir must be provided")
		}
	},
}

// applyPolicy applies a policy to the cluster.
func applyPolicy(policyFile string) error {
	policyBytes, policyName, err := utils.GetPolicy(policyFile)
	if err != nil {
		return err
	}

	return createAndApplyPolicy(policyName, policyBytes)
}

func createAndApplyPolicy(name string, policyBytes []byte) error {
	if Controller.IsKube() {
		policyCR := &policyv1alpha1.Policy{}
		policyCR.Spec.Raw = policyBytes
		policyCR.Name = name

		deployment, err := utils.GetControllerDeployment(Controller.GetKubeRestConfig(), controllerNs)
		if err != nil {
			return err
		}
		controllerNs = deployment.GetNamespace()

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

		policyCR.Namespace = deployment.GetNamespace()
		policyCR.Annotations = map[string]string{
			"fluxninja.com/validate": "true",
		}
		err = kubeClient.Create(context.Background(), policyCR)
		if err != nil {
			if utils.IsNoMatchError(err) {
				var isUpdated bool
				isUpdated, updatePolicyUsingAPIErr := utils.UpdatePolicyUsingAPI(client, name, policyBytes, force)
				if !isUpdated {
					return updatePolicyUsingAPIErr
				}
			} else if apierrors.IsAlreadyExists(err) {
				var shouldUpdate bool
				shouldUpdate, checkForUpdateErr := utils.CheckForUpdate(name, force)
				if checkForUpdateErr != nil {
					return fmt.Errorf("failed to check for update: %w", checkForUpdateErr)
				}
				if !shouldUpdate {
					log.Info().Str("policy", name).Str("namespace", deployment.GetNamespace()).Msg("Skipping update of Policy")
					return nil
				}
				updatePolicyCRErr := updatePolicyCR(name, policyCR, kubeClient)
				if updatePolicyCRErr != nil {
					return updatePolicyCRErr
				}
			} else {
				return fmt.Errorf("failed to apply policy in Kubernetes: %w", err)
			}
		}
	} else {
		isUpdated, updatePolicyUsingAPIErr := utils.UpdatePolicyUsingAPI(client, name, policyBytes, force)
		if !isUpdated {
			return updatePolicyUsingAPIErr
		}
	}

	log.Info().Str("policy", name).Msg("Applied Policy successfully")
	return nil
}

// updatePolicyCR updates the policy CR.
func updatePolicyCR(name string, policy *policyv1alpha1.Policy, kubeClient k8sclient.Client) error {
	existingPolicy := &policyv1alpha1.Policy{}
	err := kubeClient.Get(context.Background(), types.NamespacedName{Name: name, Namespace: controllerNs}, existingPolicy)
	if err != nil {
		return err
	}

	existingPolicy.Spec = policy.Spec
	if existingPolicy.GetAnnotations() == nil {
		existingPolicy.Annotations = map[string]string{}
	}
	existingPolicy.Annotations["fluxninja.com/validate"] = "true"

	err = kubeClient.Update(context.Background(), existingPolicy)
	if err != nil && apierrors.IsConflict(err) {
		return updatePolicyCR(name, policy, kubeClient)
	}
	return err
}
