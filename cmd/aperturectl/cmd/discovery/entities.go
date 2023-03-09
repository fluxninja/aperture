package discovery

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
)

// EntitiesCmd is the command to list control points.
var EntitiesCmd = &cobra.Command{
	Use:           "entities",
	Short:         "List AutoScale control points",
	Long:          `List AutoScale control points`,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		resp, err := client.ListDiscoveryEntities(context.Background(), &cmdv1.ListDiscoveryEntitiesRequest{})
		if err != nil {
			return err
		}

		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}

		// fmt.Fprintln(tabwriter, "AGENT GROUP\tNAME\tNAMESPACE\tKIND")
		// for _, cp := range resp.Entities.GetEntities() {
		// fmt.Fprintf(tabwriter, "%s\t%s\t%s\t%s\n",
		// 	cp.AgentGroup,
		// 	cp.AutoScaleControlPoint.Name,
		// 	cp.AutoScaleControlPoint.Namespace,
		// 	cp.AutoScaleControlPoint.Kind)
		// }

		return nil
	},
}
