package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"golang.org/x/exp/slices"

	cmdv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/cmd/v1"
)

// ParseAutoScaleControlPoints parses the control points.
func ParseAutoScaleControlPoints(client IntrospectionClient) error {
	resp, err := client.ListAutoScaleControlPoints(context.Background(), &cmdv1.ListAutoScaleControlPointsRequest{})
	if err != nil {
		return err
	}

	if resp.ErrorsCount != 0 {
		fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
	}

	slices.SortFunc(resp.GlobalAutoScaleControlPoints, func(a, b *cmdv1.GlobalAutoScaleControlPoint) int {
		if a.AgentGroup != b.AgentGroup {
			return strings.Compare(a.AgentGroup, b.AgentGroup)
		}
		if a.AutoScaleControlPoint.Name != b.AutoScaleControlPoint.Name {
			return strings.Compare(a.AutoScaleControlPoint.Name, b.AutoScaleControlPoint.Name)
		}
		if a.AutoScaleControlPoint.Namespace != b.AutoScaleControlPoint.Namespace {
			return strings.Compare(a.AutoScaleControlPoint.Namespace, b.AutoScaleControlPoint.Namespace)
		}
		return strings.Compare(a.AutoScaleControlPoint.Kind, b.AutoScaleControlPoint.Kind)
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
}
