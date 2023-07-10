package apply

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/protobuf/field_mask"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/tui"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var (
	file string
	dir  string
)

func init() {
	ApplyPolicyCmd.Flags().StringVar(&file, "file", "", "Path to Aperture Policy file")
	ApplyPolicyCmd.Flags().StringVar(&dir, "dir", "", "Path to directory containing Aperture Policy files")
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
			policies, err := getPolicies(dir)
			if err != nil {
				return err
			}

			model := tui.InitialCheckboxModel(policies, "Which policies to apply?")
			p := tea.NewProgram(model)
			if _, err := p.Run(); err != nil {
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

// getPolicies applies all policies in a directory to the cluster.
func getPolicies(policyDir string) ([]string, error) {
	policies := []string{}
	// walk the directory and apply all policies
	return policies, filepath.Walk(policyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileBase := info.Name()[:len(info.Name())-len(filepath.Ext(info.Name()))]
		if filepath.Ext(info.Name()) == ".yaml" && !strings.HasSuffix(fileBase, "-cr") {
			_, err := getPolicy(path)
			if err != nil {
				return err
			}
			policies = append(policies, path)
		}
		return nil
	})
}

func getPolicy(policyFile string) (*languagev1.Policy, error) {
	policyBytes, err := os.ReadFile(policyFile)
	if err != nil {
		return nil, err
	}
	_, policy, err := utils.CompilePolicy(filepath.Base(policyFile), policyBytes)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// applyPolicy applies a policy to the cluster.
func applyPolicy(policyFile string) error {
	policyFileBase := filepath.Base(policyFile)
	policyName := policyFileBase[:len(policyFileBase)-len(filepath.Ext(policyFileBase))]

	policy, err := getPolicy(policyFile)
	if err != nil {
		return err
	}

	return createAndApplyPolicy(policyName, policy)
}

func createAndApplyPolicy(name string, policy *languagev1.Policy) error {
	policyBytes, err := policy.MarshalJSON()
	if err != nil {
		return err
	}

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
			if apimeta.IsNoMatchError(err) {
				var isUpdated bool
				isUpdated, updatePolicyUsingAPIErr := updatePolicyUsingAPI(name, policy)
				if !isUpdated {
					return updatePolicyUsingAPIErr
				}
			} else if apierrors.IsAlreadyExists(err) {
				var update bool
				update, checkForUpdateErr := checkForUpdate(name)
				if checkForUpdateErr != nil {
					return fmt.Errorf("failed to check for update: %w", checkForUpdateErr)
				}
				if !update {
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
		isUpdated, updatePolicyUsingAPIErr := updatePolicyUsingAPI(name, policy)
		if !isUpdated {
			return updatePolicyUsingAPIErr
		}
	}

	log.Info().Str("policy", name).Msg("Applied Policy successfully")
	return nil
}

// updatePolicyUsingAPI updates the policy using the API.
func updatePolicyUsingAPI(name string, policy *languagev1.Policy) (bool, error) {
	request := languagev1.UpsertPolicyRequest{
		PolicyName: name,
		Policy:     policy,
	}
	_, err := client.UpsertPolicy(context.Background(), &request)
	if err != nil {
		if strings.Contains(err.Error(), "Use UpsertPolicy with PATCH call to update it.") {
			var update bool
			update, err = checkForUpdate(name)
			if err != nil {
				return false, fmt.Errorf("failed to check for update: %w", err)
			}

			if !update {
				log.Info().Str("policy", name).Str("namespace", controllerNs).Msg("Skipping update of Policy")
				return false, nil
			}

			request.UpdateMask = &field_mask.FieldMask{
				Paths: []string{"all"},
			}
			_, err = client.UpsertPolicy(context.Background(), &request)
			if err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}
	return true, nil
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

// checkForUpdate checks if the user wants to update the policy.
func checkForUpdate(name string) (bool, error) {
	model := tui.InitialRadioButtonModel([]string{"Yes", "No"}, fmt.Sprintf("Policy '%s' already exists. Do you want to update it?", name))
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return false, err
	}

	return *model.Selected == 0, nil
}
