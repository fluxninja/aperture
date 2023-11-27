package statecache

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	agentGroup   string
	controlPoint string
	key          string
)

func init() {
	GetCommand.Flags().StringVarP(&agentGroup, "agent-group", "a", "", "Agent group")
	GetCommand.Flags().StringVarP(&controlPoint, "control-point", "c", "", "Control point")
	GetCommand.Flags().StringVarP(&key, "key", "k", "", "Key")
}

var GetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a state cache entry",
	Long:  `Get a state cache entry`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.CacheLookupInput{
			AgentGroup:   agentGroup,
			ControlPoint: controlPoint,
			Key:          key,
		}

		return utils.ParseResultCacheLookup(client, input)
	},
}
