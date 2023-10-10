package blueprints

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// BlueprintsListCmd is the command to list blueprints from the Cloud Controller.
var BlueprintsListCmd = &cobra.Command{
	Use:           "list",
	Short:         "Cloud Blueprints List",
	Long:          `List cloud blueprints.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		listResponse, err := client.List(context.Background(), &emptypb.Empty{}, nil)
		if err != nil {
			return err
		}

		for _, blueprint := range listResponse.GetBlueprints() {
			fmt.Println(blueprint.GetValues())
		}

		return nil
	},
}
