package statecache

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
	GetCommand.Flags().StringVarP(&agentGroup, "agent-group", "a", "", "Agent group")
	GetCommand.Flags().StringVarP(&controlPoint, "control-point", "c", "", "Control point")
	GetCommand.Flags().StringVarP(&key, "key", "k", "", "Key")
	GetCommand.Flags().StringVarP(&value, "value", "v", "", "Value")
	GetCommand.Flags().Int64VarP(&ttl, "ttl", "t", 600000, "TTL in milliseconds")
}

var SetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set a state cache entry",
	Long:  `Set a state cache entry`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.CacheUpsertInput{
			AgentGroup:   agentGroup,
			ControlPoint: controlPoint,
			Key:          key,
			Value:        value,
			TTL:          time.Duration(time.Duration(ttl) * time.Millisecond),
		}

		return utils.ParseResultCacheUpsert(client, input)
	},
}
