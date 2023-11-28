package statecache

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCommand.Flags().StringVarP(&agentGroup, "agent-group", "a", "", "Agent group")
	DeleteCommand.Flags().StringVarP(&controlPoint, "control-point", "c", "", "Control point")
	DeleteCommand.Flags().StringVarP(&key, "key", "k", "", "Key")
	err := DeleteCommand.MarkFlagRequired("agent-group")
	if err != nil {
		panic(err)
	}
	err = DeleteCommand.MarkFlagRequired("control-point")
	if err != nil {
		panic(err)
	}
	err = DeleteCommand.MarkFlagRequired("key")
	if err != nil {
		panic(err)
	}
}

var DeleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a state cache entry",
	Long:  `Delete a state cache entry`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.CacheDeleteInput{
			AgentGroup:   agentGroup,
			ControlPoint: controlPoint,
			Key:          key,
		}

		return utils.ParseStateCacheDelete(client, input)
	},
}
