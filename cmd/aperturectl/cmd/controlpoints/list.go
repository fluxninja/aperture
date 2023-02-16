package controlpoints

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cmdv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/cmd/v1"
)

// ListCmd is the command to list control points.
var ListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List control points",
	Long:          `List control points`,
	SilenceErrors: true,
	Example:       `aperturectl control-points list`,
	RunE: func(_ *cobra.Command, _ []string) error {
		var opts []grpc.DialOption

		if insecure {
			opts = append(opts, grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				}),
			),
			)
		}

		conn, err := grpc.Dial(controllerAddr, opts...)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := cmdv1.NewControllerClient(conn)

		resp, err := client.ListControlPoints(
			context.Background(),
			&cmdv1.ListControlPointsRequest{},
		)
		if err != nil {
			return err
		}

		if resp.ErrorsCount != 0 {
			fmt.Fprintf(os.Stderr, "Could not get answer from %d agents", resp.ErrorsCount)
		}

		slices.SortFunc(resp.ControlPoints, func(a, b *cmdv1.ServiceControlPoint) bool {
			if a.ServiceName < b.ServiceName {
				return true
			}
			return a.Name < b.Name
		})

		tabwriter := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
		fmt.Fprintln(tabwriter, "SERVICE\tNAME")
		for _, cp := range resp.ControlPoints {
			fmt.Fprintf(tabwriter, "%s\t%s\n", cp.ServiceName, cp.Name)
		}
		tabwriter.Flush()

		return nil
	},
}
