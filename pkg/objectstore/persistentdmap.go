package objectstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/buraksezer/olric"
)

// ObjectStoreBackedDMap is a wrapper around olric.DMap with a second layer persistent storage.
type ObjectStoreBackedDMap struct {
	dmap           olric.DMap
	backingStorage ObjectStorageIface

	getDistCacheLabels func() (prometheus.Labels, bool)
	getMissesTotal     *prometheus.CounterVec
	getHitsTotal       *prometheus.CounterVec
}

func (o *ObjectStoreBackedDMap) generateObjectKey(key string) string {
	return fmt.Sprintf("%s-%s-%s", o.backingStorage.KeyPrefix(), o.dmap.Name(), key)
}

// Delete deletes keys from in-memory and backed storage.
func (o *ObjectStoreBackedDMap) Delete(ctx context.Context, keys ...string) (int, error) {
	for _, key := range keys {
		err := o.backingStorage.Delete(ctx, o.generateObjectKey(key))
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete object from backing storage")
		}
	}

	return o.dmap.Delete(ctx, keys...)
}

// Destroy destroys persistent dmap.
func (o *ObjectStoreBackedDMap) Destroy(ctx context.Context) error {
	return o.dmap.Destroy(ctx)
}

// Expire deletes key from backing storage and expires it in in-memory storage.
func (o *ObjectStoreBackedDMap) Expire(ctx context.Context, key string, timeout time.Duration) error {
	err := o.backingStorage.Delete(ctx, o.generateObjectKey(key))
	if err != nil {
		return err
	}

	return o.dmap.Expire(ctx, key, timeout)
}

// Function executes function on dmap.
func (o *ObjectStoreBackedDMap) Function(ctx context.Context, label string, function string, arg []byte) ([]byte, error) {
	return o.dmap.Function(ctx, label, function, arg)
}

// Get tries to get key from in-memory storage and then from backing storage if that fails.
func (o *ObjectStoreBackedDMap) Get(ctx context.Context, key string) (*olric.GetResponse, error) {
	resp, err := o.dmap.Get(ctx, key)
	if err == nil {
		// Got hit in in-memory cache, no need to get from backing storage.
		metric, ready := o.getHitsTotalMetric(metrics.PersistentCacheTypeInMemory)
		if ready {
			metric.Inc()
		}
		return resp, nil
	}
	if !errors.Is(err, olric.ErrKeyNotFound) {
		// Some error from in-memory cache.
		return nil, err
	}
	// Key not found in in-memory cache. Need to check backing storage.
	metric, ready := o.getMissesTotalMetric(metrics.PersistentCacheTypeInMemory)
	if ready {
		metric.Inc()
	}

	objectKey := o.generateObjectKey(key)
	entry, innerErr := o.backingStorage.Get(ctx, objectKey)
	if innerErr != nil {
		if errors.Is(innerErr, ErrKeyNotFound) {
			metric, ready = o.getMissesTotalMetric(metrics.PersistentCacheTypeObjectStorage)
			if ready {
				metric.Inc()
			}
		}
		return nil, innerErr
	}

	// Got hit in backing cache, need to save this result in in-memory cache.
	metric, ready = o.getHitsTotalMetric(metrics.PersistentCacheTypeObjectStorage)
	if ready {
		metric.Inc()
	}

	expireAt := time.Duration(entry.TTL()) * time.Second
	innerErr = o.dmap.Put(ctx, key, entry.Value(), olric.EXAT(expireAt))
	if innerErr != nil {
		return nil, innerErr
	}

	return olric.NewResponse(entry), nil
}

// Lock locks dmap.
func (o *ObjectStoreBackedDMap) Lock(ctx context.Context, key string, deadline time.Duration) (olric.LockContext, error) {
	return o.dmap.Lock(ctx, key, deadline)
}

// LockWithTimeout locks dmap with timeout.
func (o *ObjectStoreBackedDMap) LockWithTimeout(ctx context.Context, key string, timeout time.Duration, deadline time.Duration) (olric.LockContext, error) {
	return o.dmap.LockWithTimeout(ctx, key, timeout, deadline)
}

// Name of the dmap.
func (o *ObjectStoreBackedDMap) Name() string {
	return o.dmap.Name()
}

// Put puts k/v pair to in-memory and backing storage.
func (o *ObjectStoreBackedDMap) Put(ctx context.Context, key string, value interface{}, options ...olric.PutOption) error {
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

func (o *ObjectStoreBackedDMap) getMissesTotalMetric(cacheType string) (prometheus.Counter, bool) {
	labels, ready := o.getDistCacheLabels()
	if !ready {
		return nil, false
	}
	labels[metrics.PersistentCacheTypeLabel] = cacheType
	return o.getMissesTotal.With(labels), true
}

func (o *ObjectStoreBackedDMap) getHitsTotalMetric(cacheType string) (prometheus.Counter, bool) {
	labels, ready := o.getDistCacheLabels()
	if !ready {
		return nil, false
	}
	labels[metrics.PersistentCacheTypeLabel] = cacheType
	return o.getHitsTotal.With(labels), true
}

// NewPersistentDMap returns new persistent dmap.
func NewPersistentDMap(
	dmap olric.DMap,
	backingStorage ObjectStorageIface,
	prometheusRegistry *prometheus.Registry,
	getDistCacheLabels func() (prometheus.Labels, bool),
) (*ObjectStoreBackedDMap, error) {
	labels := []string{
		metrics.DistCacheMemberIDLabel,
		metrics.DistCacheMemberNameLabel,
		metrics.PersistentCacheTypeLabel,
	}
	getMissesTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.PersistentCacheGetMissesMetricName,
		Help: "Cumulative number of persistent cache misses.",
	}, labels)
	getHitsTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: metrics.PersistentCacheGetHitsMetricName,
		Help: "Cumulative number of persistent cache hits.",
	}, labels)
	for _, m := range []prometheus.Collector{getMissesTotal, getHitsTotal} {
		err := prometheusRegistry.Register(m)
		if err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, fmt.Errorf("unable to register persistent cache metrics: %v", err)
			}
		}
	}
	return &ObjectStoreBackedDMap{
		dmap:               dmap,
		backingStorage:     backingStorage,
		getMissesTotal:     getMissesTotal,
		getHitsTotal:       getHitsTotal,
		getDistCacheLabels: getDistCacheLabels,
	}, nil
}

var _ olric.DMap = &ObjectStoreBackedDMap{}
