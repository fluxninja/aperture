package apply

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/fluxninja/aperture/operator/api"
	policyv1alpha1 "github.com/fluxninja/aperture/operator/api/policy/v1alpha1"
	"github.com/fluxninja/aperture/pkg/log"
)

var (
	policyName         string
	dynamicConfigFile  string
	dynamicConfigBytes []byte
)

func init() {
	ApplyDynamicConfigCmd.Flags().StringVar(&policyName, "policy", "", "Name of the Policy to apply the DynamicConfig to")
	ApplyDynamicConfigCmd.Flags().StringVar(&dynamicConfigFile, "file", "", "Path to the dynamic config file")
}

// ApplyDynamicConfigCmd is the command to apply DynamicConfig to a Policy.
var ApplyDynamicConfigCmd = &cobra.Command{
	Use:           "dynamic-config",
	Short:         "Apply Aperture DynamicConfig to a Policy",
	Long:          `Use this command to apply the Aperture DynamicConfig to a Policy.`,
	SilenceErrors: true,
	Example:       `aperturectl apply dynamic-config --policy=static-rate-limiting --file=dynamic-config.yaml`,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		if policyName == "" {
			return errors.New("policy name is required")
		}
		if dynamicConfigFile == "" {
			return errors.New("dynamic config file is required")
		}
		// read the dynamic config file
		var err error
		dynamicConfigBytes, err = os.ReadFile(dynamicConfigFile)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(_ *cobra.Command, _ []string) error {
		err := api.AddToScheme(scheme.Scheme)
		if err != nil {
			return fmt.Errorf("failed to connect to Kubernetes: %w", err)
		}

		c, err := client.New(kubeRestConfig, client.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		deployment, err := getControllerDeployment()
		if err != nil {
			return err
		}

		policy := &policyv1alpha1.Policy{}
		err = c.Get(context.Background(), client.ObjectKey{
			Namespace: deployment.Namespace,
			Name:      policyName,
		}, policy)
		if err != nil {
			return fmt.Errorf("failed to get Policy '%s': %w", policyName, err)
		}

		dynamicConfigYAML := make(map[string]interface{})
		err = yaml.Unmarshal(dynamicConfigBytes, &dynamicConfigYAML)
		if err != nil {
			return fmt.Errorf("failed to parse DynamicConfig YAML: %w", err)
		}
		dynamicConfigBytes, err := json.Marshal(dynamicConfigYAML)
		if err != nil {
			return fmt.Errorf("failed to parse DynamicConfig JSON: %w", err)
		}

		policy.DynamicConfig.Raw = dynamicConfigBytes
		err = c.Update(context.Background(), policy)
		if err != nil {
			return fmt.Errorf("failed to update Policy '%s': %w", policyName, err)
		}

		log.Info().Str("policy", policyName).Str("namespace", deployment.Namespace).Msg("Updated DynamicConfig successfully")

		return nil
	},
}
