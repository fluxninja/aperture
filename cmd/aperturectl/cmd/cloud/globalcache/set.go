package globalcache

import (
	"time"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	value string
	ttl   int64
)

func init() {
	SetCommand.Flags().StringVarP(&agentGroup, "agent-group", "a", "", "Agent group")
	SetCommand.Flags().StringVarP(&key, "key", "k", "", "Key")
	SetCommand.Flags().StringVarP(&value, "value", "", "", "Value")
	SetCommand.Flags().Int64VarP(&ttl, "ttl", "t", 600000, "TTL in milliseconds")
	err := SetCommand.MarkFlagRequired("agent-group")
	if err != nil {
		panic(err)
	}
	err = SetCommand.MarkFlagRequired("key")
	if err != nil {
		panic(err)
	}
	err = SetCommand.MarkFlagRequired("value")
	if err != nil {
		panic(err)
	}
}

var SetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set a global cache entry",
	Long:  `Set a global cache entry`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.CacheUpsertInput{
			AgentGroup: agentGroup,
			Key:        key,
			Value:      value,
			TTL:        time.Duration(ttl) * time.Millisecond,
		}

		return utils.ParseGlobalCacheUpsert(client, input)
	},
}
