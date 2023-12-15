package objectstorage

import (
	"cloud.google.com/go/storage"
	"context"
	olricstorage "github.com/buraksezer/olric/pkg/storage"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"go.uber.org/fx"
)

const (
	objectStorageOpPut    = iota
	objectStorageOpDelete = iota
)

type oper int
type Operation struct {
	op    oper
	entry PersistedEntry
}

type ObjectStorage struct {
	keyPrefix      string
	cancellableCtx context.Context
	cancel         context.CancelFunc
	client         *storage.Client
	bucket         *storage.BucketHandle

	operations chan *Operation
}

func (o *ObjectStorage) SetContextWithCancel(ctx context.Context, cancel context.CancelFunc) {
	o.cancellableCtx = ctx
	o.cancel = cancel
}

type Config struct {
	Backend   string `json:"backend" validate:"oneof=gcs" default:"gcs"`
	Bucket    string `json:"bucket" validate:"required"`
	KeyPrefix string `json:"key_prefix" validate:"required"`
}

func (o *ObjectStorage) Get(ctx context.Context, key string) (olricstorage.Entry, error) {
	obj := o.bucket.Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := reader.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close object storage reader")
		}
	}()

	var data []byte
	_, err = reader.Read(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read object storage object")
		return nil, err
	}

	return PersistedEntry{key: key, value: &data}, nil
}

func (o *ObjectStorage) Delete(_ context.Context, key string) error {
	o.operations <- &Operation{
		op: objectStorageOpDelete,
		entry: PersistedEntry{
			key:   key,
			value: nil,
		},
	}

	return nil
}

func (o *ObjectStorage) List(ctx context.Context, prefix string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ObjectStorage) Put(_ context.Context, key string, data []byte) error {
	o.operations <- &Operation{
		op: objectStorageOpPut,
		entry: PersistedEntry{
			key:   key,
			value: &data,
		},
	}

	return nil
}

func (o *ObjectStorage) KeyPrefix() string {
	return o.keyPrefix
}

var _ ObjectStorageIface = (*ObjectStorage)(nil)

func (o *ObjectStorage) Start(ctx context.Context) {
	for {
		select {
		case <-o.cancellableCtx.Done():
		case <-ctx.Done():
			return
		case op := <-o.operations:
			switch op.op {
			case objectStorageOpPut:
				obj := o.bucket.Object(op.entry.key)

				w := obj.NewWriter(ctx)
				_, err := w.Write(*op.entry.value)
				if err != nil {
					log.Error().Err(err).Msg("Failed to write cache object to the storage bucket")
				}
				err = w.Close()
				if err != nil {
					log.Error().Err(err).Msg("Failed to close writer for cache object")
				}
			case objectStorageOpDelete:
				obj := o.bucket.Object(op.entry.key)
				err := obj.Delete(o.cancellableCtx)
				if err != nil {
					log.Error().Err(err).Msg("Failed to delete cache object from the storage bucket")
				}
			}
		}
	}
}

func (o *ObjectStorage) Stop(_ context.Context) error {
	o.cancel()
	return o.client.Close()
}

type ProvideParams struct {
	fx.In

	Unmarshaller config.Unmarshaller
}

func Provide(in ProvideParams) (*ObjectStorage, error) {
	var cfg Config
	err := in.Unmarshaller.UnmarshalKey("object_storage", &cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal object_storage config")
		return nil, err
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GCS client")
		return nil, err
	}
	bucket := client.Bucket(cfg.Bucket)

	objStorage := &ObjectStorage{
		bucket:     bucket,
		keyPrefix:  cfg.KeyPrefix,
		operations: make(chan *Operation),
	}

	return objStorage, nil
}

type InvokeParams struct {
	fx.In

	Lifecycle     fx.Lifecycle
	ObjectStorage ObjectStorageIface
}

func Invoke(in InvokeParams) error {
	in.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			cancellableCtx, cancel := context.WithCancel(ctx)
			in.ObjectStorage.SetContextWithCancel(cancellableCtx, cancel)
			go in.ObjectStorage.Start(ctx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return in.ObjectStorage.Stop(ctx)
		},
	})

	return nil
}
