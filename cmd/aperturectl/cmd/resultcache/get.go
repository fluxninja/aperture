package resultcache

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
	err := GetCommand.MarkFlagRequired("agent-group")
	if err != nil {
		panic(err)
	}
	err = GetCommand.MarkFlagRequired("control-point")
	if err != nil {
		panic(err)
	}
	err = GetCommand.MarkFlagRequired("key")
	if err != nil {
		panic(err)
	}
}

var GetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a result cache entry",
	Long:  `Get a result cache entry`,
	RunE: func(_ *cobra.Command, _ []string) error {
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
