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

// ObjectStorageBackedDMap is a wrapper around olric.DMap with a second layer persistent storage.
type ObjectStorageBackedDMap struct {
	defaultTTL     time.Duration
	dmap           olric.DMap
	backingStorage ObjectStorageIface

	getDistCacheLabels func() (prometheus.Labels, bool)
	getMissesTotal     *prometheus.CounterVec
	getHitsTotal       *prometheus.CounterVec
	operationDuration  *prometheus.SummaryVec
}

func (o *ObjectStorageBackedDMap) generateObjectKey(key string) string {
	return fmt.Sprintf("%s-%s-%s", o.backingStorage.KeyPrefix(), o.dmap.Name(), key)
}

// Delete deletes keys from in-memory and backed storage.
func (o *ObjectStorageBackedDMap) Delete(ctx context.Context, keys ...string) (int, error) {
	startTime := time.Now()
	defer func() {
		durationMetric, ready := o.getOperationDurationMetric(metrics.PersistentCacheOperationDelete)
		if ready {
			durationMetric.Observe(float64(time.Since(startTime).Milliseconds()))
		}
	}()

	for _, key := range keys {
		err := o.backingStorage.Delete(ctx, o.generateObjectKey(key))
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete object from backing storage")
		}
	}

	return o.dmap.Delete(ctx, keys...)
}

// Destroy destroys persistent dmap.
func (o *ObjectStorageBackedDMap) Destroy(ctx context.Context) error {
	return o.dmap.Destroy(ctx)
}

// Expire deletes key from backing storage and expires it in in-memory storage.
func (o *ObjectStorageBackedDMap) Expire(ctx context.Context, key string, timeout time.Duration) error {
	err := o.backingStorage.Delete(ctx, o.generateObjectKey(key))
	if err != nil {
		return err
	}

	return o.dmap.Expire(ctx, key, timeout)
}

// Function executes function on dmap.
func (o *ObjectStorageBackedDMap) Function(ctx context.Context, label string, function string, arg []byte) ([]byte, error) {
	return o.dmap.Function(ctx, label, function, arg)
}

// Get tries to get key from in-memory storage and then from backing storage if that fails.
func (o *ObjectStorageBackedDMap) Get(ctx context.Context, key string) (*olric.GetResponse, error) {
	startTime := time.Now()
	defer func() {
		durationMetric, ready := o.getOperationDurationMetric(metrics.PersistentCacheOperationGet)
		if ready {
			durationMetric.Observe(float64(time.Since(startTime).Milliseconds()))
		}
	}()
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
	log.Trace().Str("expireAt", expireAt.String()).Msg("Entry from storage expiration time/TTL")
	_, innerErr = o.dmap.Put(ctx, key, entry.Value(), olric.EXAT(expireAt), olric.TS(entry.Timestamp()))
	if innerErr != nil {
		return nil, innerErr
	}

	return olric.NewResponse(entry), nil
}

// Lock locks dmap.
func (o *ObjectStorageBackedDMap) Lock(ctx context.Context, key string, deadline time.Duration) (olric.LockContext, error) {
	return o.dmap.Lock(ctx, key, deadline)
}

// LockWithTimeout locks dmap with timeout.
func (o *ObjectStorageBackedDMap) LockWithTimeout(ctx context.Context, key string, timeout time.Duration, deadline time.Duration) (olric.LockContext, error) {
	return o.dmap.LockWithTimeout(ctx, key, timeout, deadline)
}

// Name of the dmap.
func (o *ObjectStorageBackedDMap) Name() string {
	return o.dmap.Name()
}

// Put puts k/v pair to in-memory and backing storage.
func (o *ObjectStorageBackedDMap) Put(
	ctx context.Context,
	key string,
	value interface{},
	options ...olric.PutOption,
) (*olric.PutConfig, error) {
	startTime := time.Now()

	defer func() {
		durationMetric, ready := o.getOperationDurationMetric(metrics.PersistentCacheOperationPut)
		if ready {
			durationMetric.Observe(float64(time.Since(startTime).Milliseconds()))
		}
	}()

	bytes, ok := value.([]byte)
	if !ok {
		log.Error().Msg("Object storage backed cache only supports []byte values")
		return nil, fmt.Errorf("invalid type for object storage backed cache: %T", value)
	}

	entryCfg, err := o.dmap.Put(ctx, key, value, options...)
	if err != nil {
		return nil, err
	}

	timestamp := entryCfg.Timestamp
	ttl := o.prepareTTL(entryCfg)

	objectKey := o.generateObjectKey(key)
	err = o.backingStorage.Put(ctx, objectKey, bytes, timestamp, ttl)
	if err != nil {
		return nil, err
	}

	return entryCfg, nil
}

func (o *ObjectStorageBackedDMap) prepareTTL(putConfig *olric.PutConfig) int64 {
	var ttl int64
	switch {
	case putConfig.HasEX:
		ttl = (putConfig.EX.Nanoseconds() + time.Now().UnixNano()) / 1000000
	case putConfig.HasPX:
		ttl = (putConfig.PX.Nanoseconds() + time.Now().UnixNano()) / 1000000
	case putConfig.HasEXAT:
		ttl = putConfig.EXAT.Nanoseconds() / 1000000
	case putConfig.HasPXAT:
		ttl = putConfig.PXAT.Nanoseconds() / 1000000
	default:
		ns := o.defaultTTL.Nanoseconds()
		if ns != 0 {
			ttl = (ns + time.Now().UnixNano()) / 1000000
		}
	}
	return ttl
}

func (o *ObjectStorageBackedDMap) getMissesTotalMetric(cacheType string) (prometheus.Counter, bool) {
	labels, ready := o.getDistCacheLabels()
	if !ready {
		return nil, false
	}
	labels[metrics.PersistentCacheTypeLabel] = cacheType
	return o.getMissesTotal.With(labels), true
}

func (o *ObjectStorageBackedDMap) getHitsTotalMetric(cacheType string) (prometheus.Counter, bool) {
	labels, ready := o.getDistCacheLabels()
	if !ready {
		return nil, false
	}
	labels[metrics.PersistentCacheTypeLabel] = cacheType
	return o.getHitsTotal.With(labels), true
}

func (o *ObjectStorageBackedDMap) getOperationDurationMetric(operation string) (prometheus.Observer, bool) {
	labels, ready := o.getDistCacheLabels()
	if !ready {
		return nil, false
	}
	labels[metrics.PersistentCacheOperationLabel] = operation
	return o.operationDuration.With(labels), true
}

// NewPersistentDMap returns new persistent dmap.
func NewPersistentDMap(
	dmap olric.DMap,
	defaultTTL time.Duration,
	backingStorage ObjectStorageIface,
	prometheusRegistry *prometheus.Registry,
	getDistCacheLabels func() (prometheus.Labels, bool),
) (*ObjectStorageBackedDMap, error) {
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
	durationLabels := []string{
		metrics.DistCacheMemberIDLabel,
		metrics.DistCacheMemberNameLabel,
		metrics.PersistentCacheOperationLabel,
	}
	operationDuration := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: metrics.PersistentCacheOperationDurationMetricName,
		Help: "Duration of persistent cache operations.",
	}, durationLabels)
	for _, m := range []prometheus.Collector{getMissesTotal, getHitsTotal, operationDuration} {
		err := prometheusRegistry.Register(m)
		if err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return nil, fmt.Errorf("unable to register persistent cache metrics: %v", err)
			}
		}
	}
	return &ObjectStorageBackedDMap{
		dmap:               dmap,
		backingStorage:     backingStorage,
		getMissesTotal:     getMissesTotal,
		getHitsTotal:       getHitsTotal,
		operationDuration:  operationDuration,
		getDistCacheLabels: getDistCacheLabels,
	}, nil
}

var _ olric.DMap = &ObjectStorageBackedDMap{}
