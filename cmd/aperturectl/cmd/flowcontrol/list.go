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

		resp, err := client.ListFlowControlControlPoints(
			context.Background(),
			&cmdv1.ListFlowControlControlPointsRequest{},
		)
		if err != nil {
			return err
		}

		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}

		slices.SortFunc(resp.GlobalFlowControlControlPoints, func(a, b *cmdv1.GlobalFlowControlControlPoint) bool {
			if a.AgentGroup != b.AgentGroup {
				return a.AgentGroup < b.AgentGroup
			}
			if a.FlowControlControlPoint.Service != b.FlowControlControlPoint.Service {
				return a.FlowControlControlPoint.Service < b.FlowControlControlPoint.Service
			}
			return a.FlowControlControlPoint.ControlPoint < b.FlowControlControlPoint.ControlPoint
		})

		tabwriter := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
		fmt.Fprintln(tabwriter, "AGENT GROUP\tSERVICE\tNAME")
		for _, cp := range resp.GlobalFlowControlControlPoints {
			fmt.Fprintf(tabwriter, "%s\t%s\t%s\n",
				cp.AgentGroup,
				cp.FlowControlControlPoint.Service,
				cp.FlowControlControlPoint.ControlPoint)
		}
		tabwriter.Flush()

		return nil
	},
}
