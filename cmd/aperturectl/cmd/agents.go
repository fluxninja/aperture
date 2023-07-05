package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

func init() {
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
		client, err := controller.Client()
		if err != nil {
			return err
		}

		agents, err := client.ListAgents(context.Background(), &emptypb.Empty{})
		if err != nil {
			return err
		}

		for _, agent := range agents.Agents {
			fmt.Println(agent)
		}

		return nil
	},
}
