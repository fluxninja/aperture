package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// ArchiveCmd is the command to archive a policy from the Aperture Cloud Controller.
var ArchiveCmd = &cobra.Command{
	Use:           "archive POLICY_NAME",
	Short:         "Archive Aperture Policy from the Aperture Cloud Controller",
	Long:          `Use this command to archive the Aperture Policy from the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl cloud policy archive POLICY_NAME`,
	RunE: func(_ *cobra.Command, args []string) error {
		_, err := cloudClient.ArchivePolicy(context.Background(), &policylangv1.DeletePolicyRequest{
			Name: args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to archive policy: %w", err)
		}

		log.Info().Str("policy", args[0]).Msg("Archived Policy successfully")
		return nil
	},
}
