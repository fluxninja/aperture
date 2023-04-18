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
	"k8s.io/client-go/kubernetes/scheme"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/tui"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/log"
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
	Example: `aperturectl apply policy --file=policies/static-rate-limiting.yaml

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
					log.Error().Err(err).Msgf("failed to apply policy '%s' on Kubernetes.", fileName)
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
	_, policy, err := utils.CompilePolicy(policyFile)
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

	return createAndApplyPolicy(policy, policyName)
}

func createAndApplyPolicy(policy *languagev1.Policy, name string) error {
	policyBytes, err := policy.MarshalJSON()
	if err != nil {
		return err
	}

	policyCR := &policyv1alpha1.Policy{}
	policyCR.Spec.Raw = policyBytes
	policyCR.Name = name

	deployment, err := utils.GetControllerDeployment(kubeRestConfig, controllerNs)
	if err != nil {
		return err
	}

	err = api.AddToScheme(scheme.Scheme)
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes: %w", err)
	}

	kubeClient, err := k8sclient.New(kubeRestConfig, k8sclient.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	policyCR.Namespace = deployment.GetNamespace()
	policyCR.Annotations = map[string]string{
		"fluxninja.com/validate": "true",
	}
	spec := policyCR.Spec
	_, err = controllerutil.CreateOrUpdate(context.Background(), kubeClient, policyCR, func() error {
		policyCR.Spec = spec
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "no matches for kind") {
			err = updatePolicyUsingAPI(name, policy)
		}

		if err != nil {
			return fmt.Errorf("failed to apply policy in Kubernetes: %w", err)
		}
	}

	log.Info().Str("policy", name).Str("namespace", deployment.GetNamespace()).Msg("Applied Policy successfully")
	return nil
}

// updatePolicyUsingAPI updates the policy using the API.
func updatePolicyUsingAPI(name string, policy *languagev1.Policy) error {
	request := languagev1.PostPoliciesRequest{
		Policies: []*languagev1.PostPoliciesRequest_PolicyRequest{
			{
				Name:   name,
				Policy: policy,
			},
		},
	}
	_, err := client.PostPolicies(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}
