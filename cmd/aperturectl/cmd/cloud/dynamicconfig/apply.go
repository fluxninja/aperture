package dynamicconfig

import (
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/yaml"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var (
	policyName         string
	dynamicConfigFile  string
	dynamicConfigBytes []byte
)

func init() {
	ApplyCmd.Flags().StringVar(&policyName, "policy", "", "Name of the Policy to apply the DynamicConfig to")
	ApplyCmd.Flags().StringVar(&dynamicConfigFile, "file", "", "Path to the dynamic config file")
}

// ApplyCmd is the command to apply DynamicConfig to a Policy.
var ApplyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Apply Aperture DynamicConfig to a Policy",
	Long:          `Use this command to apply the Aperture DynamicConfig to a Policy.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud dynamic-config apply --policy=rate-limiting --file=dynamic-config.yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// read the dynamic config file
		var err error
		dynamicConfigBytes, err = utils.GetDynamicConfigBytes(policyName, dynamicConfigFile)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		dynamicConfigYAML := make(map[string]interface{})
		err := yaml.Unmarshal(dynamicConfigBytes, &dynamicConfigYAML)
		if err != nil {
			return fmt.Errorf("failed to parse DynamicConfig YAML: %w", err)
		}

		err = utils.ApplyDynamicConfigUsingAPI(client, dynamicConfigYAML, policyName)
		if err != nil {
			return err
		}

		log.Info().Str("policy", policyName).Msg("Updated DynamicConfig successfully")

		return nil
	},
}
