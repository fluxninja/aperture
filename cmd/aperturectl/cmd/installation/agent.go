package installation

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agentInstallCmd is the command to install Aperture Agent on Kubernetes.
var agentInstallCmd = &cobra.Command{
	Use:   "agent",
	Short: "Install Aperture Agent",
	Long: `
Use this command to install Aperture Agent on your Kubernetes cluster.
Refer https://artifacthub.io/packages/helm/aperture/aperture-agent#parameters for list of configurable parameters for preparing values file.`,
	SilenceErrors: true,
	Example: `aperturectl install agent --values-file=values.yaml

aperturectl install agent --values-file=values.yaml --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if valuesFile == "" {
			return fmt.Errorf("--values-file must be provided")
		}

		return handleInstall(apertureAgent, agent)
	},
}

// agentUnInstallCmd is the command to uninstall Aperture Agent from Kubernetes.
var agentUnInstallCmd = &cobra.Command{
	Use:   "agent",
	Short: "Uninstall Aperture Agent",
	Long: `
Use this command to uninstall Aperture Agent from your Kubernetes cluster`,
	SilenceErrors: true,
	Example: `aperturectl uninstall agent

aperturectl uninstall agent --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleUnInstall(apertureAgent, agent)
	},
}
