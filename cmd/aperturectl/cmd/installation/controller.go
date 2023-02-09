package installation

import (
	"github.com/spf13/cobra"
)

// controllerInstallCmd is the command to install Aperture Controller on Kubernetes.
var controllerInstallCmd = &cobra.Command{
	Use:   "controller",
	Short: "Install Aperture Controller",
	Long: `
Use this command to install Aperture Controller and its dependencies on your Kubernetes cluster.
Refer https://artifacthub.io/packages/helm/aperture/aperture-controller#parameters for list of configurable parameters for preparing values file.`,
	SilenceErrors: true,
	Example: `aperturectl install controller --values-file=values.yaml

aperturectl install controller --values-file=values.yaml --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInstall(apertureController, controller)
	},
}

// controllerInstallCmd is the command to install Aperture Controller on Kubernetes.
var controllerUnInstallCmd = &cobra.Command{
	Use:   "controller",
	Short: "Uninstall Aperture Controller",
	Long: `
Use this command to uninstall Aperture Controller and its dependencies from your Kubernetes cluster`,
	SilenceErrors: true,
	Example: `aperturectl uninstall controller

aperturectl uninstall controller --namespace=aperture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleUnInstall(apertureController, controller)
	},
}
