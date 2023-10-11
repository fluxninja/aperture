package apply

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var (
	file      string
	dir       string
	force     bool
	selectAll bool

	client utils.CloudPolicyClient
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
	Example: `aperturectl cloud apply policy --file=policies/rate-limiting.yaml --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key PERSONAL_API_KEY

aperturectl cloud apply policy --dir=policies --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key PERSONAL_API_KEY`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		client, err = Controller.CloudPolicyClient()
		if err != nil {
			return fmt.Errorf("failed to get cloud controller client: %w", err)
		}

		return nil
	},
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
	updatePolicyUsingAPIErr := utils.UpdatePolicyUsingAPI(client, name, policyBytes, force)
	if updatePolicyUsingAPIErr != nil {
		return updatePolicyUsingAPIErr
	}

	log.Info().Str("policy", name).Msg("Applied Policy successfully")
	return nil
}
