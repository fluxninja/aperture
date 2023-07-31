package status

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	statusv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/status/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
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
		client, err := controller.Client()
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

		printLeafStatus("", statusResp)

		return nil
	},
	PersistentPostRun: controller.PostRun,
}

func printLeafStatus(parent string, resp *statusv1.GroupStatus) {
	if resp.Groups == nil && resp.GetStatus() != nil {
		if respErr := resp.Status.GetError(); respErr != nil {
			fmt.Printf("Error for %s: %s\n", parent, respErr.Message)
		}
		if respMsg := resp.Status.GetMessage(); respMsg != nil {
			value := respMsg.String()

			if respMsg.MessageIs(new(wrapperspb.StringValue)) {
				stringVal := &wrapperspb.StringValue{}
				err := respMsg.UnmarshalTo(stringVal)
				if err != nil {
					fmt.Printf("Error unmarshalling string value for key %s: %s\n", parent, err)
					return
				}
				value = stringVal.Value
			}
			if respMsg.MessageIs(new(wrapperspb.DoubleValue)) {
				doubleVal := &wrapperspb.DoubleValue{}
				err := respMsg.UnmarshalTo(doubleVal)
				if err != nil {
					fmt.Printf("Error unmarshalling double value for key %s: %s\n", parent, err)
					return
				}
				value = fmt.Sprintf("%v", doubleVal.Value)
			}
			if respMsg.MessageIs(new(emptypb.Empty)) {
				value = ""
			}

			fmt.Printf("Status for %s: %s\n", parent, value)
		}
	}

	for k, v := range resp.Groups {
		newParent := fmt.Sprintf("%s/%s", parent, k)
		printLeafStatus(newParent, v)
	}
}
