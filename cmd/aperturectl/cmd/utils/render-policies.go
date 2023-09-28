package utils

import (
	"errors"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/tui"
	policyv1alpha1 "github.com/fluxninja/aperture/v2/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/ghodss/yaml"
)

// GetPolicyTUIModel prepares the TUI model for selecting policies to apply from the given directory path.
func GetPolicyTUIModel(policyDir string, selectAll bool) ([]string, *tui.CheckBoxModel, error) {
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
		if filepath.Ext(info.Name()) == ".yaml" {
			_, policyName, err := GetPolicy(path)
			if err != nil {
				return err
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

func GetPolicy(policyFile string) (*languagev1.Policy, string, error) {
	policyFileBase := filepath.Base(policyFile)
	policyName := policyFileBase[:len(policyFileBase)-len(filepath.Ext(policyFileBase))]

	policyBytes, err := os.ReadFile(policyFile)
	if err != nil {
		return nil, policyName, err
	}
	_, policy, err := CompilePolicy(filepath.Base(policyFile), policyBytes)
	if err != nil {
		policyCR, err := GetPolicyCR(policyFile)
		if err != nil {
			return nil, policyName, err
		}

		policy = &languagev1.Policy{}
		err = config.UnmarshalYAML(policyCR.Spec.Raw, policy)
		if err != nil {
			return nil, policyName, err
		}

		policyName = policyCR.Name
		return policy, policyName, nil
	}

	return policy, policyName, nil
}

func GetPolicyCR(policyFile string) (*policyv1alpha1.Policy, error) {
	policyBytes, err := os.ReadFile(policyFile)
	if err != nil {
		return nil, err
	}

	policyCR := &policyv1alpha1.Policy{}
	err = yaml.Unmarshal(policyBytes, policyCR)
	if err != nil {
		return nil, err
	}

	return policyCR, nil
}
