package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// ListCmd is the command to list policies from the Aperture Cloud Controller.
var ListCmd = &cobra.Command{
	Use:           "list",
	Short:         "List all Aperture Policies from the Aperture Cloud Controller",
	Long:          `Use this command to list all the Aperture Policies from the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud policy list`,
	RunE: func(_ *cobra.Command, args []string) error {
		policies, err := client.ListPolicies(context.Background(), new(emptypb.Empty))
		if err != nil {
			return fmt.Errorf("failed to list policies: %w", err)
		}

		for policyName := range policies.GetPolicies().GetPolicies() {
			fmt.Println(policyName)
		}

		return nil
	},
}
