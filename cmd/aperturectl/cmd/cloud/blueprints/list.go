package blueprints

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// BlueprintsListCmd is the command to list blueprints from the Cloud Controller.
var BlueprintsListCmd = &cobra.Command{
	Use:           "list",
	Short:         "Cloud Blueprints List",
	Long:          `List cloud blueprints.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud blueprints list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		listResponse, err := client.List(context.Background(), &emptypb.Empty{})
		if err != nil {
			return fmt.Errorf("failed to list blueprints: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		for _, blueprint := range listResponse.GetBlueprints() {
			fmt.Fprintf(w, "%s\n", blueprint.GetPolicyName())
		}

		if err := w.Flush(); err != nil {
			return fmt.Errorf("failed to flush writer: %w", err)
		}

		return nil
	},
}
