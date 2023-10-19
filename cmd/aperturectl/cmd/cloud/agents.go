package cloud

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

var agentGroup string

func init() {
	agentsCmd.Flags().StringVar(&agentGroup, "agent-group", "", "Name of the agent group to list agents for")

	controller.InitFlags(agentsCmd.PersistentFlags())
}

var agentsCmd = &cobra.Command{
	Use:               "agents",
	Short:             "List connected agents",
	Long:              `List connected agents`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
	RunE: func(*cobra.Command, []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		return utils.ListAgents(client, agentGroup)
	},
}
