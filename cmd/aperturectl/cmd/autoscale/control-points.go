package autoscale

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
)

// ControlPointsCmd is the command to list control points.
var ControlPointsCmd = &cobra.Command{
	Use:           "control-points",
	Short:         "List AutoScale control points",
	Long:          `List AutoScale control points`,
	SilenceErrors: true,
	Example:       `aperturectl auto-scale control-points`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		resp, err := client.ListAutoScaleControlPoints(context.Background(), &cmdv1.ListAutoScaleControlPointsRequest{})
		if err != nil {
			return err
		}

		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}

		slices.SortFunc(resp.GlobalAutoScaleControlPoints, func(a, b *cmdv1.GlobalAutoScaleControlPoint) bool {
			if a.AgentGroup != b.AgentGroup {
				return a.AgentGroup < b.AgentGroup
			}
			if a.AutoScaleControlPoint.Name != b.AutoScaleControlPoint.Name {
				return a.AutoScaleControlPoint.Name < b.AutoScaleControlPoint.Name
			}
			if a.AutoScaleControlPoint.Namespace != b.AutoScaleControlPoint.Namespace {
				return a.AutoScaleControlPoint.Namespace < b.AutoScaleControlPoint.Namespace
			}
			return a.AutoScaleControlPoint.Kind < b.AutoScaleControlPoint.Kind
		})

		tabwriter := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
		fmt.Fprintln(tabwriter, "AGENT GROUP\tNAME\tNAMESPACE\tKIND")
		for _, cp := range resp.GlobalAutoScaleControlPoints {
			fmt.Fprintf(tabwriter, "%s\t%s\t%s\t%s\n",
				cp.AgentGroup,
				cp.AutoScaleControlPoint.Name,
				cp.AutoScaleControlPoint.Namespace,
				cp.AutoScaleControlPoint.Kind)
		}
		tabwriter.Flush()

		return nil
	},
}
