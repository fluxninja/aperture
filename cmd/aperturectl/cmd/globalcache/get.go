package globalcache

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	agentGroup string
	key        string
)

func init() {
	GetCommand.Flags().StringVarP(&agentGroup, "agent-group", "a", "", "Agent group")
	GetCommand.Flags().StringVarP(&key, "key", "k", "", "Key")
	err := GetCommand.MarkFlagRequired("agent-group")
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
	Short: "Get a global cache entry",
	Long:  `Get a global cache entry`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.CacheLookupInput{
			AgentGroup: agentGroup,
			Key:        key,
		}

		return utils.ParseGlobalCacheLookup(client, input)
	},
}
