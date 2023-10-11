package utils

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/flowcontrol/preview/v1"
)

// ParseControlPoints parses the control points.
func ParseControlPoints(client IntrospectionClient) error {
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
}

type PreviewInput struct {
	AgentGroup    string
	IsHTTPPreview bool
	NumSamples    int
	Service       string
	ControlPoint  string
}

// ParsePreview parses the preview.
func ParsePreview(client IntrospectionClient, input PreviewInput) error {
	previewReq := &previewv1.PreviewRequest{
		Samples:      int64(input.NumSamples),
		Service:      input.Service,
		ControlPoint: input.ControlPoint,
		// FIXME LabelMatcher: Figure out how to represent label matcher in CLI.
	}

	if input.IsHTTPPreview {
		resp, err := client.PreviewHTTPRequests(
			context.Background(),
			&cmdv1.PreviewHTTPRequestsRequest{
				AgentGroup: input.AgentGroup,
				Request:    previewReq,
			},
		)
		if err != nil {
			return err
		}
		samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp.Response)
		if err != nil {
			return err
		}
		os.Stdout.Write(samplesJSON)
	} else {
		resp, err := client.PreviewFlowLabels(
			context.Background(),
			&cmdv1.PreviewFlowLabelsRequest{
				AgentGroup: input.AgentGroup,
				Request:    previewReq,
			},
		)
		if err != nil {
			return err
		}
		samplesJSON, err := protojson.MarshalOptions{Multiline: true}.Marshal(resp.Response)
		if err != nil {
			return err
		}
		os.Stdout.Write(samplesJSON)
	}

	return nil
}
