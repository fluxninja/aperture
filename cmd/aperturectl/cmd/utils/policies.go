package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ghodss/yaml"
	"google.golang.org/protobuf/types/known/emptypb"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/tui"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// GetPoliciesTUIModel prepares the TUI model for selecting policies to apply from the given directory path.
func GetPoliciesTUIModel(policyDir string, selectAll bool) ([]string, *tui.CheckBoxModel, error) {
	policies, err := GetPolicies(policyDir)
	if err != nil {
		return nil, nil, err
	}

	if len(policies) == 0 {
		return nil, nil, errors.New("no policies found in the directory")
	}

	model := tui.InitialCheckboxModel(policies, "Which policies to apply?")
	if !selectAll {
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			return nil, nil, err
		}
	} else {
		for i := range policies {
			model.Selected[i] = struct{}{}
		}
	}

	return policies, model, nil
}

// GetPolicies returns path of valid files having a valid Aperture Policy .
func GetPolicies(policyDir string) ([]string, error) {
	policies := []string{}
	policyMap := map[string]string{}
	// walk the directory and apply all policies
	return policies, filepath.Walk(policyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(info.Name()) == ".yaml" || filepath.Ext(info.Name()) == ".yml" {
			_, policyName, err := GetPolicy(path)
			if err != nil {
				log.Info().Str("file", path).Msg("Invalid policy found. Skipping...")
				return nil
			}
			if _, ok := policyMap[policyName]; !ok {
				policyMap[policyName] = path
				policies = append(policies, path)
			} else {
				log.Info().Str("policy", policyName).Msg("Duplicate policy found. Skipping...")
			}
		}
		return nil
	})
}

// GetPolicy returns the policy from the policy file.
func GetPolicy(policyFile string) ([]byte, string, error) {
	policyFileBase := filepath.Base(policyFile)
	policyName := policyFileBase[:len(policyFileBase)-len(filepath.Ext(policyFileBase))]

	var err error
	policyBytes, err := os.ReadFile(policyFile)
	if err != nil {
		return nil, policyName, err
	}

	var policyCR *policyv1alpha1.Policy
	policyCR, err = GetPolicyCR(policyBytes)
	if err == nil {
		policyBytes = policyCR.Spec.Raw
		policyName = policyCR.Name
	}
	return policyBytes, policyName, nil
}

// GetPolicyCR returns the policy CR from the policy bytes.
func GetPolicyCR(policyBytes []byte) (*policyv1alpha1.Policy, error) {
	policyCR := &policyv1alpha1.Policy{}
	err := yaml.Unmarshal(policyBytes, policyCR)
	if err != nil {
		return nil, err
	}

	if policyCR.Name == "" {
		return nil, fmt.Errorf("policy name is missing in the policy file")
	}

	return policyCR, nil
}

// UpdatePolicyUsingAPI updates the policy using the API.
func UpdatePolicyUsingAPI(client PolicyClient, listClient SelfHostedPolicyClient, name string, policyBytes []byte, force bool) error {
	request := policylangv1.UpsertPolicyRequest{
		PolicyName:   name,
		PolicyString: string(policyBytes),
	}

	if !force {
		// If directly using controller API, we can call GetPolicies to
		// verify that we're not accidentally overwriting the policy.
		// Cloud API doesn't have this method and we always allow
		// overwriting, even with force=false.
		existingPolicies, err := listClient.ListPolicies(context.Background(), new(emptypb.Empty))

		needsConfirmation := false
		if err != nil {
			return err
		}
		for policyName, policy := range existingPolicies.GetPolicies().GetPolicies() {
			if policyName == name && policy.Status != policylangv1.GetPolicyResponse_STALE {
				needsConfirmation = true
				break
			}
		}

		if needsConfirmation {
			update, err := CheckForUpdate(name, force)
			if err != nil {
				return fmt.Errorf("failed to check for update: %w", err)
			}

			if !update {
				log.Info().Str("policy", name).Str("namespace", controllerNs).Msg("Skipping update of Policy")
				return errors.New("policy already exists")
			}
		}
	}

	_, err := client.UpsertPolicy(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}

// CheckForUpdate checks if the user wants to update the policy.
func CheckForUpdate(name string, force bool) (bool, error) {
	if force {
		return true, nil
	}

	model := tui.InitialRadioButtonModel([]string{"Yes", "No"}, fmt.Sprintf("Policy '%s' already exists. Do you want to update it?", name))
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return false, err
	}

	return *model.Selected == 0, nil
}

// DeletePolicyUsingAPI deletes the policy using the API.
func DeletePolicyUsingAPI(client PolicyClient, policyName string) error {
	policyRequest := policylangv1.DeletePolicyRequest{
		Name: policyName,
	}
	_, err := client.DeletePolicy(context.Background(), &policyRequest)
	if err != nil {
		return fmt.Errorf("failed to delete policy '%s' using API: %w", policyName, err)
	}

	return nil
}

// ListPolicies lists the policies using the API.
func ListPolicies(client SelfHostedPolicyClient) error {
	policies, err := client.ListPolicies(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	for name, body := range policies.GetPolicies().Policies {
		fmt.Printf("%v:\n", name)
		if body.GetStatus() != policylangv1.GetPolicyResponse_VALID {
			fmt.Println("\tStatus:", body.GetStatus())
			reason := strings.ReplaceAll(body.GetReason(), "\n", "\n\n\t\t")
			reason = strings.ReplaceAll(reason, " Error", "\n\t\tError")
			fmt.Printf("\tReason: %s\n", reason)
			fmt.Println("\t\t---")
		}

		// Note: We try to print policy details also if status is not
		// VALID, because statuses like like OUTDATED or STALE can contain
		// policy details.

		resources := body.GetPolicy().GetResources()
		if ims := resources.GetInfraMeters(); len(ims) > 0 {
			fmt.Println("\tInfra Meters:")
			for im := range ims {
				fmt.Printf("\t\t%v\n", im)
			}
		}
		if fms := resources.GetFlowControl().GetFluxMeters(); len(fms) > 0 {
			fmt.Println("\tFlux Meters:")
			for fm := range fms {
				fmt.Printf("\t\t%v\n", fm)
			}
		}
		if cs := resources.GetFlowControl().GetClassifiers(); len(cs) > 0 {
			fmt.Println("\tClassifiers:")
			for _, c := range body.Policy.Resources.FlowControl.Classifiers {
				if len(c.Selectors) > 0 {
					fmt.Println("\t\tSelectors:")
					for _, s := range c.Selectors {
						fmt.Printf("\t\t\t%v\n", s)
					}
				}
				if len(c.Rules) > 0 {
					fmt.Println("\t\tRules:")
					for r := range c.Rules {
						fmt.Printf("\t\t\t%v\n", r)
					}
				}
			}
			fmt.Println("\t\t---")
		}
	}

	return nil
}
