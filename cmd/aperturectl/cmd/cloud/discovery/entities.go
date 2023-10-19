package discovery

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	findBy     string
	agentGroup string
)

func init() {
	EntitiesCmd.Flags().StringVar(&findBy, "find-by", "", "Find entity by [name|ip]")
	EntitiesCmd.Flags().StringVar(&agentGroup, "agent-group", "", "Name of the agent group to list agents for")
}

// EntitiesCmd is the command to list control points.
var EntitiesCmd = &cobra.Command{
	Use:           "entities",
	Short:         "List AutoScale control points",
	Long:          `List AutoScale control points`,
	SilenceErrors: true,
	Example: `aperturectl cloud discovery entities

aperturectl cloud discovery entities --find-by="name=service1-demo-app-7dfdf9c698-4wmlt"

aperturectl cloud discovery entities --find-by=“ip=10.244.1.24”`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		return utils.ParseEntities(client, findBy, agentGroup)
	},
}
