package objectstorage

import (
	"context"
	"errors"
	"strconv"

	"cloud.google.com/go/storage"

	olricstorage "github.com/buraksezer/olric/pkg/storage"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	storageconfig "github.com/fluxninja/aperture/v2/pkg/objectstorage/config"
	"go.uber.org/fx"
)

const (
	objectStorageOpPut    = iota
	objectStorageOpDelete = iota
)

type (
	oper int
	// Operation executed on object storage.
	Operation struct {
		op    oper
		entry *PersistentEntry
	}
)

// ErrKeyNotFound means that given key is not present in the object storage.
var ErrKeyNotFound = errors.New("key not found")

// ObjectStorageIface is an abstract over persistent storage for Olric DMap.
type ObjectStorageIface interface {
	KeyPrefix() string
	SetContextWithCancel(ctx context.Context, cancel context.CancelFunc)
	Start(ctx context.Context)
	Stop(ctx context.Context) error
	Put(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) (olricstorage.Entry, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) (string, error)
}

// ObjectStorage is an ObjectStorageIface implementation using GCP storage.
type ObjectStorage struct {
	keyPrefix      string
	cancellableCtx context.Context
	cancel         context.CancelFunc
	client         *storage.Client
	bucket         *storage.BucketHandle

	operations chan *Operation
}

// SetContextWithCancel sets long running context and cancel function.
func (o *ObjectStorage) SetContextWithCancel(ctx context.Context, cancel context.CancelFunc) {
	o.cancellableCtx = ctx
	o.cancel = cancel
}

// Get gets object from object storage.
func (o *ObjectStorage) Get(ctx context.Context, key string) (olricstorage.Entry, error) {
	obj := o.bucket.Object(key)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	defer func() {
		closeErr := reader.Close()
		if closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close object storage reader")
		}
	}()

	var data []byte
	_, err = reader.Read(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read object storage object")
		return nil, err
	}

	entry := &PersistentEntry{key: key, value: &data}
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get object storage object attributes")
	} else {
		timestamp, err := strconv.ParseInt(attrs.Metadata["timestamp"], 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse object storage object timestamp")
			// XXX: Should we use current time instead?
			entry.SetTimestamp(0)
		} else {
			entry.SetTimestamp(timestamp)
		}

		ttl, err := strconv.ParseInt(attrs.Metadata["ttl"], 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse object storage object ttl")
			// XXX: What timestamp can we use here as default?
			entry.SetTTL(60)
		} else {
			entry.SetTTL(ttl)
		}
	}

	return entry, nil
}

// Delete queues delete operation from object storage.
func (o *ObjectStorage) Delete(_ context.Context, key string) error {
	o.operations <- &Operation{
		op: objectStorageOpDelete,
		entry: &PersistentEntry{
			key:   key,
			value: nil,
		},
	}

	return nil
}

// List lists object storage.
func (o *ObjectStorage) List(ctx context.Context, prefix string) (string, error) {
	panic("implement me")
}

// Put queues put operation to object storage.
func (o *ObjectStorage) Put(_ context.Context, key string, data []byte) error {
	o.operations <- &Operation{
		op: objectStorageOpPut,
		entry: &PersistentEntry{
			key:   key,
			value: &data,
		},
	}

	return nil
}

// KeyPrefix getter.
func (o *ObjectStorage) KeyPrefix() string {
	return o.keyPrefix
}

var _ ObjectStorageIface = (*ObjectStorage)(nil)

// Start starts a goroutine which performs operations on object storage.
func (o *ObjectStorage) Start(ctx context.Context) {
	go func() {
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
						continue
					}

					err = w.Close()
					if err != nil {
						log.Error().Err(err).Msg("Failed to close writer for cache object")
						continue
					}

					_, err = obj.Update(ctx, storage.ObjectAttrsToUpdate{
						Metadata: map[string]string{
							"timestamp": strconv.FormatInt(op.entry.Timestamp(), 10),
							"ttl":       strconv.FormatInt(op.entry.TTL(), 10),
						},
					})

					if err != nil {
						log.Error().Err(err).Msg("Failed to set object metadata")
						continue
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
	}()
}

// Stop kills the goroutine started in Start().
// TODO - add graceful shutdown to wait for the operations' queue to be empty.
func (o *ObjectStorage) Stop(_ context.Context) error {
	o.cancel()
	return o.client.Close()
}

// ProvideParams for object storage.
type ProvideParams struct {
	fx.In

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
			cancellableCtx, cancel := context.WithCancel(ctx)
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
