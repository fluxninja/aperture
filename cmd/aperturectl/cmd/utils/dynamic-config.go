package utils

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/protobuf/types/known/structpb"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// GetDynamicConfigBytes returns the bytes of the dynamic config file.
func GetDynamicConfigBytes(policyName, dynamicConfigFile string) ([]byte, error) {
	if policyName == "" {
		return nil, errors.New("policy name is required")
	}
	if dynamicConfigFile == "" {
		return nil, errors.New("dynamic config file is required")
	}
	// read the dynamic config file
	return os.ReadFile(dynamicConfigFile)
}

// ApplyDynamicConfig applies the dynamic config.
func ApplyDynamicConfigUsingAPI(client SelfHostedPolicyClient, dynamicConfigYAML map[string]interface{}, policyName string) error {
	var dynamicConfigStruct *structpb.Struct
	var err error
	dynamicConfigStruct, err = structpb.NewStruct(dynamicConfigYAML)
	if err != nil {
		return fmt.Errorf("failed to parse DynamicConfig Struct: %w", err)
	}
	request := languagev1.PostDynamicConfigRequest{
		PolicyName:    policyName,
		DynamicConfig: dynamicConfigStruct,
	}
	_, err = client.PostDynamicConfig(context.Background(), &request)
	if err != nil {
		return fmt.Errorf("failed to update DynamicConfig: %w", err)
	}

	return nil
}

// GetDynamicConfigUsingAPI gets the dynamic config.
func GetDynamicConfigUsingAPI(client SelfHostedPolicyClient, policyName string) error {
	request := languagev1.GetDynamicConfigRequest{
		PolicyName: policyName,
	}
	resp, err := client.GetDynamicConfig(context.Background(), &request)
	if err != nil {
		return fmt.Errorf("failed to update DynamicConfig: %w", err)
	}

	if resp.DynamicConfig == nil {
		log.Info().Str("policy-name", policyName).Msg("DynamicConfig is not set for the given Policy")
		return nil
	}

	j, err := resp.DynamicConfig.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	yamlString, err := GetYAMLString(j)
	if err != nil {
		return err
	}
	fmt.Print(yamlString)
	return nil
}
