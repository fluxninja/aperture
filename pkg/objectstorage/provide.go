package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sourcegraph/conc/pool"

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

var (
	// ErrKeyNotFound means that given key is not present in the object storage.
	ErrKeyNotFound = errors.New("key not found")
	// ErrObjectTimestampMissing means that either object is not present, or its metadata is incomplete.
	errObjectTimestampMissing = errors.New("object or its timestamp does not exist")
)

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

	inFlightOpsMutex sync.Mutex
	operations       chan *Operation
}

// SetContextWithCancel sets long-running context and cancel function.
func (o *ObjectStorage) SetContextWithCancel(ctx context.Context, cancel context.CancelFunc) {
	o.cancellableCtx = ctx
	o.cancel = cancel
}

// Get gets object from object storage.
func (o *ObjectStorage) Get(ctx context.Context, key string) (olricstorage.Entry, error) {
	// If the object is missing timestamp, we will use timestamp of when the Get() was called.
	timestampDefault := time.Now().UnixNano()

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
		timestampMetadata, ok := attrs.Metadata["timestamp"]
		if !ok {
			log.Error().Msg("Missing object timestamp in metadata, using time of the Get() call instead")
			entry.SetTimestamp(timestampDefault)
		} else {
			timestamp, err := strconv.ParseInt(timestampMetadata, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse object timestamp, using time of the Get() call instead")
				entry.SetTimestamp(timestampDefault)
			} else {
				entry.SetTimestamp(timestamp)
			}
		}

		deleteStaleCacheEntry := func(entry *PersistentEntry) {
			err := o.internalDelete(ctx, entry)
			if err != nil {
				log.Error().Err(err).Msg("Failed to queue expired cache entry for deletion from object storage")
			}
		}

		ttlMetadata, ok := attrs.Metadata["ttl"]
		// If the TTL is either missing from metadata, or cannot be parsed we must assume that the entry is expired.
		// Otherwise, we would be always returning potentially stale entry until it's either overwritten, or bucket
		// lifecycle policy deletes it for us.
		// XXX: Another approach would be to use a default TTL (but what would be a good default?) and update the
		//      metadata.
		if !ok {
			log.Error().Msg("Missing object TTL in metadata, assume that the persisted cache entry is stale")
			deleteStaleCacheEntry(entry)
			return nil, ErrKeyNotFound
		} else {
			ttl, err := strconv.ParseInt(ttlMetadata, 10, 64)
			if err != nil {
				log.Error().Err(err).Msg("Failed to parse object TTL, assuming that the persisted cache entry is stale")
				deleteStaleCacheEntry(entry)
			} else {
				entry.SetTTL(ttl)

				if time.UnixMilli(ttl).Before(time.Now()) {
					log.Debug().Msg("Object in storage has expired, deleting it and returning ErrKeyNotFound")
					deleteStaleCacheEntry(entry)
					return nil, ErrKeyNotFound
				}
			}
		}
	}

	return entry, nil
}

func (o *ObjectStorage) internalDelete(_ context.Context, entry *PersistentEntry) error {
	o.operations <- &Operation{
		op:    objectStorageOpDelete,
		entry: entry,
	}

	return nil
}

// Delete queues delete operation from object storage.
func (o *ObjectStorage) Delete(ctx context.Context, key string) error {
	return o.internalDelete(ctx, &PersistentEntry{
		key:   key,
		value: nil,
	})
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

func (o *ObjectStorage) getObjectTimestamp(ctx context.Context, obj *storage.ObjectHandle) (int64, error) {
	attrs, err := obj.Attrs(ctx)
	if err != nil && !errors.Is(err, storage.ErrObjectNotExist) {
		log.Error().Err(err).Msg("Failed to query storage bucket for object attributes")
		return 0, err
	}
	if attrs == nil {
		return 0, nil
	}

	timeStamp, ok := attrs.Metadata["timestamp"]
	if ok {
		i, err := strconv.ParseInt(timeStamp, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse object timestamp")
			return 0, err
		}
		return i, nil
	} else {
		return 0, errObjectTimestampMissing
	}
}

func (o *ObjectStorage) handleOpPut(ctx context.Context, entry *PersistentEntry) error {
	obj := o.bucket.Object(entry.key)

	timestamp, err := o.getObjectTimestamp(ctx, obj)
	if err != nil && !errors.Is(err, errObjectTimestampMissing) {
		return err
	}

	if timestamp > entry.Timestamp() {
		log.Debug().Msg("Object in storage is more recent than the entry being created")
		return nil
	}

	w := obj.NewWriter(ctx)
	_, err = w.Write(*entry.value)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write cache object to the storage bucket")
		return err
	}

	err = w.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close writer for cache object")
		return nil
	}

	_, err = obj.Update(ctx, storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"timestamp": strconv.FormatInt(entry.Timestamp(), 10),
			"ttl":       strconv.FormatInt(entry.TTL(), 10),
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to set object metadata")
		return nil
	}

	return nil
}

func (o *ObjectStorage) handleOpDelete(ctx context.Context, entry *PersistentEntry) error {
	obj := o.bucket.Object(entry.key)

	timestamp, err := o.getObjectTimestamp(ctx, obj)
	if err != nil && !errors.Is(err, errObjectTimestampMissing) {
		return err
	}

	if timestamp > entry.Timestamp() {
		log.Debug().Msg("Object in storage is more recent than the entry being deleted")
		return nil
	}

	err = obj.Delete(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete cache object from the storage bucket")
	}

	return nil
}

func (o *ObjectStorage) handleOp(ctx context.Context, op *Operation) error {
	switch op.op {
	case objectStorageOpPut:
		return o.handleOpPut(ctx, op.entry)

	case objectStorageOpDelete:
		return o.handleOpDelete(ctx, op.entry)
	}

	return nil
}

// Start starts a goroutine which performs operations on object storage.
func (o *ObjectStorage) Start(_ context.Context) {
	go func() {
		p := pool.New().WithMaxGoroutines(10)

		o.inFlightOpsMutex.Lock()

		defer o.inFlightOpsMutex.Unlock()
		defer p.Wait()

		for oper := range o.operations {
			_ = o.handleOp(o.cancellableCtx, oper)
		}
	}()

	go func() {
		<-o.cancellableCtx.Done()
		close(o.operations)
	}()
}

// Stop kills the goroutine started in Start().
func (o *ObjectStorage) Stop(_ context.Context) error {
	o.cancel()
	// The lock is held by the go-routine responsible for handling operations. Once the channel is closed (by calling o.cancel())
	// the go-routine will finish processing operations, and then exit, releasing the lock in the process.
	o.inFlightOpsMutex.Lock()
	defer o.inFlightOpsMutex.Unlock()

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
