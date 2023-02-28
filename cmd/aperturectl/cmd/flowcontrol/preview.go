package flowcontrol

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	previewv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/preview/v1"
)

var (
	agentGroup    string
	isHTTPPreview bool
	numSamples    int
)

func init() {
	PreviewCmd.Flags().StringVar(&agentGroup, "agent-group", "default", "Agent group")
	PreviewCmd.Flags().IntVar(&numSamples, "samples", 10, "Number of samples to collect")
	PreviewCmd.Flags().BoolVar(
		&isHTTPPreview,
		"http",
		false,
		"Preview HTTP requests instead of flow labels",
	)
}

// PreviewCmd is the command to preview control points.
var PreviewCmd = &cobra.Command{
	Use:           "preview [--http] SERVICE CONTROL_POINT",
	Short:         "Preview control points",
	Long:          `Preview samples of flow labels or HTTP requests on control points`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := controller.Client()
		if err != nil {
			return err
		}

		previewReq := &previewv1.PreviewRequest{
			Samples:      int64(numSamples),
			Service:      args[0],
			ControlPoint: args[1],
			// FIXME LabelMatcher: Figure out how to represent label matcher in CLI.
		}

		if isHTTPPreview {
			resp, err := client.PreviewHTTPRequests(
				context.Background(),
				&cmdv1.PreviewHTTPRequestsRequest{
					AgentGroup: agentGroup,
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
					AgentGroup: agentGroup,
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
	},
}
