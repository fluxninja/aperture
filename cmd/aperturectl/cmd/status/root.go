package status

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(StatusCmd.PersistentFlags())
}

// StatusCmd is the command to get a status of jobs.
var StatusCmd = &cobra.Command{
	Use:           "status",
	Short:         "Get Jobs status",
	Long:          `Use this command to get the status of internal jobs.`,
	SilenceErrors: true,
	Example: `
	aperturectl status
	`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := controller.StatusClient()
		if err != nil {
			return err
		}

		getStatusReq := &statusv1.GroupStatusRequest{
			Path: "",
		}
		statusResp, err := client.GetStatus(
			context.Background(),
			getStatusReq,
		)
		if err != nil {
			return err
		}

		result, err := status.ParseGroupStatus(make(map[string]string), "", statusResp)
		if err != nil {
			return err
		}
		for k, v := range result {
			fmt.Printf("%s: %s\n", k, v)
		}
		return nil
	},
	PersistentPostRun: controller.PostRun,
}
