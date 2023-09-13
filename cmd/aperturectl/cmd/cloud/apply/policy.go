package apply

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/protobuf/field_mask"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/tui"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
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

// ApplyPolicyCmd is the command to apply a policy to the Aperture Cloud Controller.
var ApplyPolicyCmd = &cobra.Command{
	Use:           "policy",
	Short:         "Apply Aperture Policy to the Aperture Cloud Controller",
	Long:          `Use this command to apply the Aperture Policy to the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Example: `aperturectl cloud apply policy --file=policies/rate-limiting.yaml --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key API_KEY

aperturectl cloud apply policy --dir=policies --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key API_KEY`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if file != "" {
			return applyPolicy(file)
		} else if dir != "" {
			policies, model, err := utils.GetPolicyTUIModel(dir, selectAll)
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
	policy, policyName, err := utils.GetPolicy(policyFile)
	if err != nil {
		return err
	}

	return createAndApplyPolicy(policyName, policy)
}

func createAndApplyPolicy(name string, policy *languagev1.Policy) error {
	isUpdated, updatePolicyUsingAPIErr := updatePolicyUsingAPI(name, policy)
	if !isUpdated {
		return updatePolicyUsingAPIErr
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

// checkForUpdate checks if the user wants to update the policy.
func checkForUpdate(name string) (bool, error) {
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
