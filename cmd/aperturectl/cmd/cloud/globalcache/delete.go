package globalcache

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	DeleteCommand.Flags().StringVarP(&agentGroup, "agent-group", "a", "", "Agent group")
	DeleteCommand.Flags().StringVarP(&key, "key", "k", "", "Key")
	err := DeleteCommand.MarkFlagRequired("agent-group")
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
	Short: "Delete a global cache entry",
	Long:  `Delete a global cache entry`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.CacheDeleteInput{
			AgentGroup: agentGroup,
			Key:        key,
		}

		return utils.ParseGlobalCacheDelete(client, input)
	},
}
