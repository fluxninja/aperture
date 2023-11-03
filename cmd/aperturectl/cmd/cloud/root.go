package cloud

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/autoscale"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/blueprints"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/decisions"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/discovery"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/dynamicconfig"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/flowcontrol"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/policy"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/status"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/utils"
)

// Version shows the version of ApertureCtl.
var (
	controller utils.ControllerConn
)

func init() {
	CloudCmd.AddCommand(dynamicconfig.DynamicConfigCmd)
	CloudCmd.AddCommand(autoscale.AutoScaleCmd)
	CloudCmd.AddCommand(decisions.DecisionsCmd)
	CloudCmd.AddCommand(agentsCmd)
	CloudCmd.AddCommand(status.StatusCmd)
	CloudCmd.AddCommand(flowcontrol.FlowControlCmd)
	CloudCmd.AddCommand(discovery.DiscoveryCmd)
	CloudCmd.AddCommand(blueprints.BlueprintsCmd)
	CloudCmd.AddCommand(policy.PolicyCmd)
}

// CloudCmd is the command to apply a policy to the Cloud Controller.
var CloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Commands to communicate with the Cloud Controller",
	Long: `
Use this command to talk to the Cloud Controller.`,
	SilenceErrors: true,
}
