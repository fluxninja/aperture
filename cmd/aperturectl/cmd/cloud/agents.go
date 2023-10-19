package cloud

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
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

		agents, err := client.ListAgents(context.Background(), &cmdv1.ListAgentsRequest{
			AgentGroup: agentGroup,
		})
		if err != nil {
			return err
		}

		for _, agent := range agents.Agents {
			fmt.Println(agent)
		}

		return nil
	},
}
