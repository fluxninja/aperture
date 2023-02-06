package apply

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/fluxninja/aperture/pkg/log"
)

var (
	kubeConfig     string
	kubeRestConfig *rest.Config
)

func init() {
	ApplyCmd.Flags().StringVar(&file, "file", "", "Path to Aperture Policy file")
	ApplyCmd.Flags().StringVar(&dir, "dir", "", "Path to directory containing Aperture Policy files")
	ApplyCmd.Flags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernetes cluster config. Defaults to '~/.kube/config'")

	ApplyCmd.AddCommand(ApplyPolicyCmd)
}

// ApplyCmd is the command to apply a policy to the cluster.
var ApplyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Apply Aperture Policy to the cluster",
	Long:          `Use this command to apply the Aperture Policy to the cluster.`,
	SilenceErrors: true,
	Example: `aperturectl apply --file=policy.yaml

aperturectl apply --dir=policy-dir`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if kubeConfig == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			kubeConfig = filepath.Join(homeDir, ".kube", "config")
			log.Info().Msgf("Using Kubernetes config '%s'", kubeConfig)
		}
		restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to Kubernetes. Error: %s", err.Error())
		}
		kubeRestConfig = restConfig
		return nil
	},
}
