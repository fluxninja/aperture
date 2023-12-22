package objectstorage

import (
	"context"
	"fmt"

	agentinfo "github.com/fluxninja/aperture/v2/pkg/agent-info"

	"cloud.google.com/go/storage"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	storageconfig "github.com/fluxninja/aperture/v2/pkg/objectstorage/config"
	"go.uber.org/fx"
)

// ProvideParams for object storage.
type ProvideParams struct {
	fx.In

	AgentInfo    *agentinfo.AgentInfo
	Unmarshaller config.Unmarshaller
}

// Provide ObjectStorage.
func Provide(in ProvideParams) (*ObjectStorage, error) {
	var cfg storageconfig.Config
	err := in.Unmarshaller.UnmarshalKey("object_storage", &cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal object_storage config")
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket cannot be empty")
	}
	if cfg.KeyPrefix == "" {
		return nil, fmt.Errorf("key prefix cannot be empty")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GCS client")
		return nil, err
	}
	bucket := client.Bucket(cfg.Bucket)

	keyPrefix := fmt.Sprintf("%s-%s", cfg.KeyPrefix, in.AgentInfo.GetAgentGroup())

	objStorage := &ObjectStorage{
		bucket:      bucket,
		keyPrefix:   keyPrefix,
		operations:  make(chan *Operation),
		retryPolicy: cfg.RetryPolicy,
	}

	return objStorage, nil
}

// InvokeParams for object storage.
type InvokeParams struct {
	fx.In

	Lifecycle     fx.Lifecycle
	ObjectStorage ObjectStorageIface
}

// Invoke ObjectStorage.
func Invoke(in InvokeParams) error {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cancellableCtx, cancel := context.WithCancel(context.Background())
			in.ObjectStorage.SetContextWithCancel(cancellableCtx, cancel)
			in.ObjectStorage.Start(ctx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return in.ObjectStorage.Stop(ctx)
		},
	})

	return nil
}
