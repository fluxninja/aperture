package cloud

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/apply"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/delete"
)

func init() {
	CloudCmd.AddCommand(apply.ApplyCmd)
	CloudCmd.AddCommand(delete.DeleteCmd)
}

// CloudCmd is the command to apply a policy to the Cloud Controller.
var CloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Commands to communicate with the Cloud Controller",
	Long: `
Use this command to talk to the Cloud Controller.`,
	SilenceErrors: true,
}
