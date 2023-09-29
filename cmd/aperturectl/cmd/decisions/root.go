package decisions

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	all           bool
	decisionType  string
	controller    utils.ControllerConn
	decisionTypes = []string{
		"load_scheduler",
		"rate_limiter",
		"quota_scheduler",
		"pod_scaler",
		"sampler",
	}
)

func init() {
	controller.InitFlags(DecisionsCmd.PersistentFlags())

	DecisionsCmd.PersistentFlags().BoolVar(&all, "all", false, "Get all decisions")
	DecisionsCmd.PersistentFlags().StringVar(&decisionType, "decision-type", "", fmt.Sprintf("Type of the decision to get (%s)", strings.Join(decisionTypes, ", ")))
}

// DecisionsCmd is the command to apply a policy to the cluster.
var DecisionsCmd = &cobra.Command{
	Use:           "decisions",
	Short:         "Get Aperture Decisions",
	Long:          `Use this command to get the Aperture Decisions.`,
	SilenceErrors: true,
	Example: `
	aperturectl decisions --all
	aperturectl decisions --decision-type="load_scheduler"`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		decisionType, err = utils.DecisionsPreRun(all, decisionType)
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := controller.PolicyClient()
		if err != nil {
			return err
		}

		return utils.ParseDecisions(cmd, client, all, decisionType)
	},
	PersistentPostRun: controller.PostRun,
}
