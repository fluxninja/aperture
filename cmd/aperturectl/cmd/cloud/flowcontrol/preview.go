package flowcontrol

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
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
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		input := utils.PreviewInput{
			AgentGroup:    agentGroup,
			IsHTTPPreview: isHTTPPreview,
			NumSamples:    numSamples,
			Service:       args[0],
			ControlPoint:  args[1],
		}

		return utils.ParsePreview(client, input)
	},
}
