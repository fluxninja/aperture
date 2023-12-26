package objectstorage

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	storageconfig "github.com/fluxninja/aperture/v2/pkg/objectstorage/config"
)

// ProvideParams for object storage.
type ProvideParams struct {
	fx.In

	AgentInfo    *agentinfo.AgentInfo
	Lifecycle    fx.Lifecycle
	Unmarshaller config.Unmarshaller
}

// Provide ObjectStorage.
func Provide(in ProvideParams) (ObjectStorageIface, error) {
	var cfg storageconfig.ObjectStorageConfig
	err := in.Unmarshaller.UnmarshalKey("object_storage", &cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal object_storage config")
		return nil, err
	}
	if !cfg.Enabled {
		log.Warn().Msg("Object storage not enabled. Creating persistent dmap will result in regular dmap")
		return nil, nil
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket cannot be empty")
	}
	if cfg.KeyPrefix == "" {
		return nil, fmt.Errorf("key prefix cannot be empty")
	}

	keyPrefix := fmt.Sprintf("%s-%s", cfg.KeyPrefix, in.AgentInfo.GetAgentGroup())

	objStorage := &ObjectStorage{
		bucketName:  cfg.Bucket,
		keyPrefix:   keyPrefix,
		operations:  make(chan *Operation, cfg.OperationsChannelSize),
		retryPolicy: cfg.RetryPolicy,
	}

	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cancellableCtx, cancel := context.WithCancel(context.Background())
			objStorage.SetContextWithCancel(cancellableCtx, cancel)
			return objStorage.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return objStorage.Stop(ctx)
		},
	})

	return objStorage, nil
}
