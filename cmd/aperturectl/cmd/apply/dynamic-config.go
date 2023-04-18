package apply

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/client-go/kubernetes/scheme"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	languagev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
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

		c, err := k8sclient.New(kubeRestConfig, k8sclient.Options{
			Scheme: scheme.Scheme,
		})
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %w", err)
		}

		deployment, err := utils.GetControllerDeployment(kubeRestConfig, controllerNs)
		if err != nil {
			return err
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

		policy := &policyv1alpha1.Policy{}
		err = c.Get(context.Background(), k8sclient.ObjectKey{
			Namespace: deployment.Namespace,
			Name:      policyName,
		}, policy)
		if err != nil {
			if strings.Contains(err.Error(), "no matches for kind") {
				var dynamicConfigStruct *structpb.Struct
				dynamicConfigStruct, err = structpb.NewStruct(dynamicConfigYAML)
				if err != nil {
					return fmt.Errorf("failed to parse DynamicConfig Struct: %w", err)
				}
				request := languagev1.PostDynamicConfigsRequest{
					DynamicConfigs: []*languagev1.PostDynamicConfigsRequest_DynamicConfigRequest{
						{
							PolicyName:    policyName,
							DynamicConfig: dynamicConfigStruct,
						},
					},
				}
				_, err = client.PostDynamicConfigs(context.Background(), &request)
				if err != nil {
					return fmt.Errorf("failed to update DynamicConfig: %w", err)
				}
			} else {
				return fmt.Errorf("failed to get Policy '%s': %w", policyName, err)
			}
		} else {
			policy.DynamicConfig.Raw = dynamicConfigBytes
			err = c.Update(context.Background(), policy)
			if err != nil {
				return fmt.Errorf("failed to update DynamicConfig for policy '%s': %w", policyName, err)
			}
		}

		log.Info().Str("policy", policyName).Str("namespace", deployment.Namespace).Msg("Updated DynamicConfig successfully")

		return nil
	},
}
