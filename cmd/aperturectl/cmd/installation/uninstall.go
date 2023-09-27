package installation

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/operator/api"
	_ "github.com/fluxninja/aperture/v2/operator/api/agent/v1alpha1"
	_ "github.com/fluxninja/aperture/v2/operator/api/controller/v1alpha1"
)

func init() {
	UnInstallCmd.PersistentFlags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernetes cluster config. Defaults to '~/.kube/config'")
	UnInstallCmd.PersistentFlags().StringVar(&version, "version", latestTag, "Version of the Aperture")
	UnInstallCmd.PersistentFlags().StringVar(&valuesFile, "values-file", "", "Values YAML file containing parameters to customize the installation")
	UnInstallCmd.PersistentFlags().StringVar(&namespace, "namespace", defaultNS, "Namespace from which the component will be uninstalled. Defaults to 'default' namespace")
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

		if version == "" {
			version = latestTag
		}
		return nil
	},
}
