package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	controller.InitFlags(BlueprintsCmd.PersistentFlags())

	BlueprintsCmd.AddCommand(BlueprintsListCmd)
	BlueprintsCmd.AddCommand(BlueprintsGetCmd)
	BlueprintsCmd.AddCommand(BlueprintsDeleteCmd)
	BlueprintsCmd.AddCommand(BlueprintsApplyCmd)
	BlueprintsCmd.AddCommand(BlueprintsArchiveCmd)
}

// BlueprintsCmd is the command to apply a policy to the Cloud Controller.
var BlueprintsCmd = &cobra.Command{
	Use:           "blueprints",
	Short:         "Cloud Blueprints",
	Long:          `Interact with cloud blueprints.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to cloud controller setup: %w", err)
		}
		client, err = controller.CloudBlueprintsClient()
		if err != nil {
			return fmt.Errorf("failed to get cloud controller client: %w", err)
		}
		return nil
	},
	PersistentPostRun: controller.PostRun,
}
