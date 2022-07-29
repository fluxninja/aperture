package utils

import (
	"encoding/base64"

	goObjectHash "github.com/benlaurie/objecthash/go/objecthash"
	"github.com/ghodss/yaml"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	configv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/config/v1"
	"github.com/FluxNinja/aperture/pkg/log"
)

// HashAndWrapWithConfProps wraps a proto message with a config properties wrapper and hashes it.
func HashAndWrapWithConfProps(message proto.Message, policyName string) (*configv1.ConfigPropertiesWrapper, error) {
	dat, marshalErr := yaml.Marshal(message)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msgf("Failed to marshal proto message %+v", message)
		return nil, marshalErr
	}
	hashBytes, hashErr := goObjectHash.ObjectHash(dat)
	if hashErr != nil {
		log.Warn().Err(hashErr).Msgf("Failed to hash json serialized proto message %s", string(dat))
		return nil, hashErr
	}
	hash := base64.StdEncoding.EncodeToString(hashBytes[:])

	return WrapWithConfProps(message, "", policyName, hash, 0)
}

// WrapWithConfProps wraps a proto message with a config properties wrapper.
func WrapWithConfProps(message proto.Message, agentGroupName, policyName, policyHash string, componentIndex int) (*configv1.ConfigPropertiesWrapper, error) {
	anyPb, anyErr := anypb.New(message)
	if anyErr != nil {
		log.Warn().Err(anyErr).Msgf("Failed to create anypb from proto message %+v", message)
		return nil, anyErr
	}

	wrapper := configv1.ConfigPropertiesWrapper{
		Config:         anyPb,
		AgentGroupName: agentGroupName,
		PolicyName:     policyName,
		PolicyHash:     policyHash,
		ComponentIndex: int64(componentIndex),
	}

	return &wrapper, nil
}
