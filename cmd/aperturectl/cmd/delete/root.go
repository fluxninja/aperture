package delete

import (
	"errors"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
)

var (
	kubeConfig     string
	kubeRestConfig *rest.Config
	controller     utils.ControllerConn
	client         cmdv1.ControllerClient
	controllerNs   string
	policyName     string
)

func init() {
	controller.InitFlags(DeleteCmd.PersistentFlags())
	DeleteCmd.PersistentFlags().StringVar(&policyName, "policy", "", "Name of the Policy to delete")

	DeleteCmd.AddCommand(DeletePolicyCmd)
}

// DeleteCmd is the command to delete a policy from the cluster.
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Aperture Policies",
	Long: `
Use this command to delete the Aperture Policies.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if policyName == "" {
			return errors.New("policy name is required")
		}

		var err error
		kubeRestConfig, err = utils.GetKubeConfig(kubeConfig)
		if err != nil {
			return err
		}

		controllerNs, err = cmd.Flags().GetString("controller-ns")
		if err != nil {
			return err
		}

		controllerAddr, err := cmd.Flags().GetString("controller")
		if err != nil {
			return err
		}

		kube, err := cmd.Flags().GetBool("kube")
		if err != nil {
			return err
		}

		if controllerAddr == "" && !kube {
			err = cmd.Flags().Set("kube", "true")
			if err != nil {
				return err
			}
		}

		err = controller.PreRunE(cmd, args)
		if err != nil {
			return err
		}

		client, err = controller.Client()
		if err != nil {
			return err
		}
		return nil
	},
	PersistentPostRun: controller.PostRun,
}
