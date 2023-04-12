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

// ControlPointsCmd is the command to list control points.
var ControlPointsCmd = &cobra.Command{
	Use:           "control-points",
	Short:         "List Flow Control control points",
	Long:          `List Flow Control control points`,
	SilenceErrors: true,
	Example:       `aperturectl flow-control control-points --kube`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		resp, err := client.ListFlowControlPoints(
			context.Background(),
			&cmdv1.ListFlowControlPointsRequest{},
		)
		if err != nil {
			return err
		}

		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}

		slices.SortFunc(resp.GlobalFlowControlPoints, func(a, b *cmdv1.GlobalFlowControlPoint) bool {
			if a.AgentGroup != b.AgentGroup {
				return a.AgentGroup < b.AgentGroup
			}
			if a.FlowControlPoint.Service != b.FlowControlPoint.Service {
				return a.FlowControlPoint.Service < b.FlowControlPoint.Service
			}
			return a.FlowControlPoint.ControlPoint < b.FlowControlPoint.ControlPoint
		})

		tabwriter := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
		fmt.Fprintln(tabwriter, "AGENT GROUP\tSERVICE\tNAME\tTYPE")
		for _, cp := range resp.GlobalFlowControlPoints {
			fmt.Fprintf(tabwriter, "%s\t%s\t%s\t%s\n",
				cp.AgentGroup,
				cp.FlowControlPoint.Service,
				cp.FlowControlPoint.ControlPoint,
				cp.FlowControlPoint.Type)
		}
		tabwriter.Flush()

		return nil
	},
}
