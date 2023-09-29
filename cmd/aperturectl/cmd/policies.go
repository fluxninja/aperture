package cmd

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	controller.InitFlags(policiesCmd.PersistentFlags())
}

var policiesCmd = &cobra.Command{
	Use:               "policies",
	Short:             "List applied policies",
	Long:              `List applied policies`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
	RunE: func(*cobra.Command, []string) error {
		client, err := controller.PolicyClient()
		if err != nil {
			return err
		}

		return utils.ListPolicies(client)
	},
}
