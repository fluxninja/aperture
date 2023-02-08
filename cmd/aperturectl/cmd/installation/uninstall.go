package installation

import (
	"fmt"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/operator/api"
	_ "github.com/fluxninja/aperture/operator/api/agent/v1alpha1"
	_ "github.com/fluxninja/aperture/operator/api/controller/v1alpha1"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	UnInstallCmd.PersistentFlags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernetes cluster config. Defaults to '~/.kube/config'")
	UnInstallCmd.PersistentFlags().StringVar(&version, "version", apertureLatestVersion, "Version of the Aperture")
	UnInstallCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "Namespace from which the component will be uninstalled. Defaults to component name")
	UnInstallCmd.PersistentFlags().IntVar(&timeout, "timeout", 300, "Timeout of waiting for uninstallation hooks completion")

	UnInstallCmd.AddCommand(controllerUnInstallCmd)
	UnInstallCmd.AddCommand(agentUnInstallCmd)
	UnInstallCmd.AddCommand(istioConfigUnInstallCmd)
}

// UnInstallCmd is the command to uninstall Aperture Controller and Aperture Agent from Kubernetes.
var UnInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Aperture components",
	Long: `
Use this command to uninstall Aperture Controller and Agent from your Kubernetes cluster.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if namespace == "" {
			namespace = defaultNS
		}

		kubeRestConfig, err = utils.GetKubeConfig(kubeConfig)
		if err != nil {
			return err
		}

		err = api.SchemeBuilder.AddToScheme(scheme.Scheme)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}
		kubeClient, err = client.New(kubeRestConfig, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		latestVersion, err = utils.ResolveLatestVersion()
		if err != nil {
			return err
		}

		if version == "" || version == apertureLatestVersion {
			version = latestVersion
		}
		return nil
	},
}
