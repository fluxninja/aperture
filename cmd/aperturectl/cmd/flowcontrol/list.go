package flowcontrol

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
)

// ListCmd is the command to list control points.
var ListCmd = &cobra.Command{
	Use:           "control-points",
	Short:         "List Flow Control control points",
	Long:          `List Flow Control control points`,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		resp, err := client.ListControlPoints(
			context.Background(),
			&cmdv1.ListControlPointsRequest{},
		)
		if err != nil {
			return err
		}

		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}

		slices.SortFunc(resp.ControlPoints, func(a, b *cmdv1.ServiceControlPoint) bool {
			if a.ServiceName != b.ServiceName {
				return a.ServiceName < b.ServiceName
			}
			return a.Name < b.Name
		})

		tabwriter := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
		fmt.Fprintln(tabwriter, "SERVICE\tNAME")
		for _, cp := range resp.ControlPoints {
			fmt.Fprintf(tabwriter, "%s\t%s\n", cp.ServiceName, cp.Name)
		}
		tabwriter.Flush()

		return nil
	},
}
