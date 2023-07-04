package decisions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
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

		if !all {
			if decisionType == "" {
				return errors.New("decision type is required or use --all to get all decisions")
			} else {
				var found bool
				for _, v := range decisionTypes {
					if v == decisionType {
						found = true
						break
					}
				}
				if !found {
					return errors.New("invalid decision type, use one of the valid types (" + strings.Join(decisionTypes, ", ") + ") or use --all to get all decisions")
				}
			}
		} else {
			decisionType = ""
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		getDecisionsReq := &languagev1.GetDecisionsRequest{
			DecisionType: decisionType,
		}
		if all {
			getDecisionsReq = nil
		}
		decisionsResp, err := client.GetDecisions(
			context.Background(),
			getDecisionsReq,
		)
		if err != nil {
			return err
		}

		for k, v := range decisionsResp.Decisions {
			cmd.Printf("%s:\n%s\n\n", k, v)
		}

		return nil
	},
	PersistentPostRun: controller.PostRun,
}
