package apply

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var (
	// Controller is the controller connection object.
	Controller utils.ControllerConn

	kubeConfig     string
	kubeRestConfig *rest.Config
	client         cmdv1.ControllerClient
	controllerNs   string
	isKube         bool
	controllerAddr string
)

func init() {
	Controller.InitFlags(ApplyCmd.PersistentFlags())

	ApplyCmd.AddCommand(ApplyPolicyCmd)
	ApplyCmd.AddCommand(ApplyDynamicConfigCmd)
}

// ApplyCmd is the command to apply a policy to the cluster.
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Aperture Policies",
	Long: `
Use this command to apply the Aperture Policies.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		controllerAddr, err = cmd.Flags().GetString("controller")
		if err != nil {
			return fmt.Errorf("failed to get controller address flag: %w", err)
		}

		isKube, err = cmd.Flags().GetBool("kube")
		if err != nil {
			return fmt.Errorf("failed to get kube flag: %w", err)
		}

		if controllerAddr == "" && !isKube {
			log.Info().Msg("Both --controller and --kube flags are not set. Assuming --kube=true.")
			err = cmd.Flags().Set("kube", "true")
			if err != nil {
				return fmt.Errorf("failed to set kube flag: %w", err)
			}

			Controller.IsKube = true
			isKube = true
		}

		if isKube {
			kubeRestConfig, err = utils.GetKubeConfig(kubeConfig)
			if err != nil {
				return fmt.Errorf("failed to get kube config: %w", err)
			}

			controllerNs, err = cmd.Flags().GetString("controller-ns")
			if err != nil {
				return fmt.Errorf("failed to get controller namespace flag: %w", err)
			}
		}

		err = Controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		client, err = Controller.Client()
		if err != nil {
			return fmt.Errorf("failed to get controller client: %w", err)
		}

		return nil
	},
	PersistentPostRun: Controller.PostRun,
}
