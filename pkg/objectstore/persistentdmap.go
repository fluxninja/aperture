package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"time"

	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/pkg/storage"
)

type ObjectStorageIface interface {
	KeyPrefix() string
	SetContextWithCancel(ctx context.Context, cancel context.CancelFunc)
	Start(ctx context.Context)
	Stop(ctx context.Context) error
	Put(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) (storage.Entry, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) (string, error)
}

type ObjectStoreBackedDMap struct {
	dmap           olric.DMap
	backingStorage ObjectStorageIface
}

func (o ObjectStoreBackedDMap) generateObjectKey(key string) string {
	return fmt.Sprintf("%s-%s-%s", o.backingStorage.KeyPrefix(), o.dmap.Name(), key)
}

func (o ObjectStoreBackedDMap) Delete(ctx context.Context, keys ...string) (int, error) {
	for _, key := range keys {
		err := o.backingStorage.Delete(ctx, o.generateObjectKey(key))
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete object from backing storage")
		}
	}

	return o.dmap.Delete(ctx, keys...)
}

func (o ObjectStoreBackedDMap) Destroy(ctx context.Context) error {
	return o.dmap.Destroy(ctx)
}

func (o ObjectStoreBackedDMap) Expire(ctx context.Context, key string, timeout time.Duration) error {
	err := o.backingStorage.Delete(ctx, o.generateObjectKey(key))
	if err != nil {
		return err
	}

	return o.dmap.Expire(ctx, key, timeout)
}

func (o ObjectStoreBackedDMap) Function(ctx context.Context, label string, function string, arg []byte) ([]byte, error) {
	return o.dmap.Function(ctx, label, function, arg)
}

func (o ObjectStoreBackedDMap) Get(ctx context.Context, key string) (*olric.GetResponse, error) {
	resp, err := o.dmap.Get(ctx, key)
	if err != nil {
		if errors.Is(err, olric.ErrKeyNotFound) {
			objectKey := o.generateObjectKey(key)
			entry, innerErr := o.backingStorage.Get(ctx, objectKey)
			if innerErr != nil {
				return nil, innerErr
			}

			expireAt := time.Duration(entry.TTL()) * time.Second
			innerErr = o.dmap.Put(ctx, key, entry.Value(), olric.EXAT(expireAt))
			if innerErr != nil {
				return nil, innerErr
			}

			resp = olric.NewResponse(entry)
		} else {
			return nil, err
		}
	}

	return resp, nil
}

func (o ObjectStoreBackedDMap) Lock(ctx context.Context, key string, deadline time.Duration) (olric.LockContext, error) {
	return o.Lock(ctx, key, deadline)
}

func (o ObjectStoreBackedDMap) LockWithTimeout(ctx context.Context, key string, timeout time.Duration, deadline time.Duration) (olric.LockContext, error) {
	return o.LockWithTimeout(ctx, key, timeout, deadline)
}

func (o ObjectStoreBackedDMap) Name() string {
	return o.dmap.Name()
}

func (o ObjectStoreBackedDMap) Put(ctx context.Context, key string, value interface{}, options ...olric.PutOption) error {
	bytes, ok := value.([]byte)
	if !ok {
		log.Error().Msg("Object storage backed cache only supports []byte values")
		return fmt.Errorf("invalid type for object storage backed cache: %T", value)
	}

	objectKey := o.generateObjectKey(key)
	err := o.backingStorage.Put(ctx, objectKey, bytes)
	if err != nil {
		return err
	}

	return o.dmap.Put(ctx, key, value, options...)
}

var _ olric.DMap = ObjectStoreBackedDMap{}
