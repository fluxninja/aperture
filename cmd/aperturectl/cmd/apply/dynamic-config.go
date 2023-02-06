package apply

import "github.com/spf13/cobra"

// ApplyDynamicConfigCmd is the command to apply DynamicConfig to a Policy.
var ApplyDynamicConfigCmd = &cobra.Command{
	Use:           "dyunamic-config",
	Short:         "Apply Aperture DynamicConfig to a Policy",
	Long:          `Use this command to apply the Aperture DynamicConfig to a Policy.`,
	SilenceErrors: true,
	Example:       `aperturectl apply dynamic-config --policy=static-rate-limiting`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return nil
	},
}
