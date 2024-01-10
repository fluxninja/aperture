package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/objectstorage/config"
	"github.com/googleapis/gax-go/v2"

	"cloud.google.com/go/storage"
	olricstorage "github.com/buraksezer/olric/pkg/storage"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/sourcegraph/conc/pool"
)

const (
	objectStorageOpPut = iota
	objectStorageOpDelete
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
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Put(ctx context.Context, key string, data []byte, timestamp, ttl int64) error
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
	bucketName     string
	bucket         *storage.BucketHandle
	retryPolicy    config.ObjectStorageRetryPolicy

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
	if !o.isStarted() {
		return nil, fmt.Errorf("storage not yet started")
	}
	// If the object is missing timestamp, we will use timestamp of when the Get() was called.
	timestampDefault := time.Now().UnixNano()

	backoff := o.retryPolicy.Backoff
	timeout := o.retryPolicy.Timeout
	obj := o.bucket.Object(key).Retryer(
		storage.WithBackoff(gax.Backoff{
			Initial:    backoff.Initial.AsDuration(),
			Multiplier: backoff.Multiplier,
			Max:        backoff.Maximum.AsDuration(),
		}),
		storage.WithPolicy(storage.RetryIdempotent),
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout.AsDuration())
	defer cancel()

	reader, err := obj.NewReader(timeoutCtx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, ErrKeyNotFound
		}
		log.Error().Err(err).Msg("Failed to create object storage reader")
		return nil, err
	}

	defer func() {
		closeErr := reader.Close()
		if closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close object storage reader")
		}
	}()

	data := make([]byte, reader.Attrs.Size)
	_, err = reader.Read(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read object storage object")
		return nil, err
	}

	entry := &PersistentEntry{key: key, value: data}
	attrs, err := obj.Attrs(timeoutCtx)
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
			// We handle deletion of stale cache entries outside the main context with timeout.
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

				now := time.Now()
				if time.UnixMilli(ttl).Before(now) {
					log.Warn().
						Time("ttl", time.UnixMilli(ttl)).
						Time("now", now).
						Msg("Object in storage has expired, deleting it and returning ErrKeyNotFound")
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
	if !o.isStarted() {
		return fmt.Errorf("storage not yet started")
	}
	return o.internalDelete(ctx, &PersistentEntry{
		key:       key,
		value:     nil,
		timestamp: time.Now().UnixNano(),
	})
}

// List lists object storage.
func (o *ObjectStorage) List(ctx context.Context, prefix string) (string, error) {
	if !o.isStarted() {
		return "", fmt.Errorf("storage not yet started")
	}
	panic("implement me")
}

// Put queues put operation to object storage.
func (o *ObjectStorage) Put(
	_ context.Context,
	key string,
	data []byte,
	timestamp int64,
	ttl int64,
) error {
	if !o.isStarted() {
		return fmt.Errorf("storage not yet started")
	}
	entry := &PersistentEntry{
		key:       key,
		value:     data,
		timestamp: timestamp,
		ttl:       ttl,
	}
	o.operations <- &Operation{
		op:    objectStorageOpPut,
		entry: entry,
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
	if err != nil {
		if !errors.Is(err, storage.ErrObjectNotExist) {
			log.Error().Err(err).Msg("Failed to query storage bucket for object attributes")
		}
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
	backoff := o.retryPolicy.Backoff
	timeout := o.retryPolicy.Timeout
	obj := o.bucket.Object(entry.key).Retryer(
		storage.WithBackoff(gax.Backoff{
			Initial:    backoff.Initial.AsDuration(),
			Multiplier: backoff.Multiplier,
			Max:        backoff.Maximum.AsDuration(),
		}),
		storage.WithPolicy(storage.RetryIdempotent),
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout.AsDuration())
	defer cancel()

	timestamp, err := o.getObjectTimestamp(timeoutCtx, obj)
	if err != nil && !errors.Is(err, errObjectTimestampMissing) && !errors.Is(err, storage.ErrObjectNotExist) {
		return err
	}

	if timestamp > entry.Timestamp() {
		log.Debug().Msg("Object in storage is more recent than the entry being created")
		return nil
	}

	w := obj.NewWriter(timeoutCtx)
	_, err = w.Write(entry.value)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write cache object to the storage bucket")
		return err
	}

	err = w.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close writer for cache object")
		return nil
	}

	_, err = obj.Update(timeoutCtx, storage.ObjectAttrsToUpdate{
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
	backoff := o.retryPolicy.Backoff
	timeout := o.retryPolicy.Timeout
	obj := o.bucket.Object(entry.key).Retryer(
		storage.WithBackoff(gax.Backoff{
			Initial:    backoff.Initial.AsDuration(),
			Multiplier: backoff.Multiplier,
			Max:        backoff.Maximum.AsDuration(),
		}),
		storage.WithPolicy(storage.RetryIdempotent),
	)

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout.AsDuration())
	defer cancel()

	timestamp, err := o.getObjectTimestamp(timeoutCtx, obj)
	if err != nil {
		// If the object is not found in the storage, just return
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil
		}

		// Ignore errObjectTimestampMissing and return all the other errors. If timestamp is missing, we'll
		// just use 0 instead.
		if !errors.Is(err, errObjectTimestampMissing) {
			return err
		}
	}

	if entry.Timestamp() > 0 && timestamp > entry.Timestamp() {
		log.Debug().Msg("Object in storage is more recent than the entry being deleted")
		return nil
	}

	err = obj.Delete(timeoutCtx)
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
func (o *ObjectStorage) Start(ctx context.Context) error {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GCS client")
		return err
	}
	bucket := client.Bucket(o.bucketName)

	o.client = client
	o.bucket = bucket

	go func() {
		p := pool.New().WithMaxGoroutines(10)

		o.inFlightOpsMutex.Lock()

		defer o.inFlightOpsMutex.Unlock()
		defer p.Wait()

		for oper := range o.operations {
			p.Go(func() {
				err := o.handleOp(o.cancellableCtx, oper)
				if err != nil {
					log.Warn().Err(err).Msg("Failed handling operation")
				}
			})
		}
	}()

	go func() {
		<-o.cancellableCtx.Done()
		close(o.operations)
	}()
	return nil
}

func (o *ObjectStorage) isStarted() bool {
	return o.bucket != nil
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
