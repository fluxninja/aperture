package installation

import (
	"github.com/spf13/cobra"
)

// istioConfigInstallCmd is the command to install example Istio EnvoyFilter on Kubernetes.
var istioConfigInstallCmd = &cobra.Command{
	Use:   "istioconfig",
	Short: "Install example Istio EnvoyFilter",
	Long: `
Use this command to install example Istio EnvoyFilter on your Kubernetes cluster.
Refer https://artifacthub.io/packages/helm/aperture/istioconfig#parameters for list of configurable parameters for preparing values file.`,
	SilenceErrors: true,
	Example: `aperturectl install istioconfig --values-file=values.yaml

aperturectl install istioconfig --values-file=values.yaml --namespace=istio-system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleInstall(istioConfig, istioConfigReleaseName)
	},
}

// istioConfigUnInstallCmd is the command to uninstall example Istio EnvoyFilter from Kubernetes.
var istioConfigUnInstallCmd = &cobra.Command{
	Use:   "istioconfig",
	Short: "Install example Istio EnvoyFilter",
	Long: `
Use this command to uninstall example Istio EnvoyFilter from your Kubernetes cluster.`,
	SilenceErrors: true,
	Example: `aperturectl uninstall istioconfig

aperturectl uninstall istioconfig --namespace=istio-system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleUnInstall(istioConfig, istioConfigReleaseName)
	},
}
