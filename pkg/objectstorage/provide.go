package objectstorage

import (
	"context"
	"fmt"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	storageconfig "github.com/fluxninja/aperture/v2/pkg/objectstorage/config"
	"go.uber.org/fx"
)

// ProvideParams for object storage.
type ProvideParams struct {
	fx.In

	Unmarshaller config.Unmarshaller
}

// Provide ObjectStorage.
func Provide(in ProvideParams) (*ObjectStorage, error) {
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

	objStorage := &ObjectStorage{
		bucketName: cfg.Bucket,
		keyPrefix:  cfg.KeyPrefix,
		operations: make(chan *Operation, cfg.OperationsChannelSize),
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
			return in.ObjectStorage.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return in.ObjectStorage.Stop(ctx)
		},
	})

	return nil
}
