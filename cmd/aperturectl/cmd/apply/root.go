package apply

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
)

var (
	kubeConfig     string
	kubeRestConfig *rest.Config
)

func init() {
	ApplyCmd.Flags().StringVar(&kubeConfig, "kube-config", "", "Path to the Kubernetes cluster config. Defaults to '~/.kube/config'")

	ApplyCmd.AddCommand(ApplyPolicyCmd)
	ApplyCmd.AddCommand(ApplyDynamicConfigCmd)
}

// ApplyCmd is the command to apply a policy to the cluster.
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Aperture Policy to the cluster",
	Long: `
Use this command to apply the Aperture Policy to the cluster.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		kubeRestConfig, err = utils.GetKubeConfig(kubeConfig)
		if err != nil {
			return err
		}
		return nil
	},
	Example: `aperturectl apply --file=policy.yaml

aperturectl apply --dir=policy-dir`,
}
