package kubernetes

import (
	"context"
	"encoding/json"
	"sync"

	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/cenkalti/backoff"
	controlpointcachev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/controlpointcache/v1"
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/k8s"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/sourcegraph/conc/stream"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
)

// A ControlPoint is identified by Group, Version, Kind, Namespace and Name.
type ControlPoint struct {
	Group     string
	Version   string
	Kind      string
	Namespace string
	Name      string
}

// ToProto converts a ControlPoint to a KubernetesControlPoint.
func (cp *ControlPoint) ToProto() *controlpointcachev1.KubernetesControlPoint {
	groupVersion := schema.GroupVersion{
		Group:   cp.Group,
		Version: cp.Version,
	}

	return &controlpointcachev1.KubernetesControlPoint{
		ApiVersion: groupVersion.String(),
		Kind:       cp.Kind,
		Namespace:  cp.Namespace,
		Name:       cp.Name,
	}
}

// ControlPointFromSelector converts a policylangv1.KubernetesObjectSelector to a ControlPoint.
func ControlPointFromSelector(k8sObjectSelector *policylangv1.KubernetesObjectSelector) (ControlPoint, error) {
	// Convert Kubernetes APIVersion into Group and Version
	groupVersion, parseErr := schema.ParseGroupVersion(k8sObjectSelector.ApiVersion)
	if parseErr != nil {
		log.Error().Err(parseErr).Msgf("Unable to parse APIVersion: %s", k8sObjectSelector.ApiVersion)
		return ControlPoint{}, parseErr
	}

	return ControlPoint{
		Group:     groupVersion.Group,
		Version:   groupVersion.Version,
		Kind:      k8sObjectSelector.Kind,
		Namespace: k8sObjectSelector.Namespace,
		Name:      k8sObjectSelector.Name,
	}, nil
}

// AutoScaler provides an interface to invoke auto-scale.
type AutoScaler interface {
	Scale(ControlPoint, int32)
}

// ControlPointStore is the interface for Storing Kubernetes Control Points.
type ControlPointStore interface {
	Add(cp ControlPoint)
	Update(cp ControlPoint)
	Delete(cp ControlPoint)
}

// ControlPointCache is the interface for Reading or Watching Kubernetes Control Points.
type ControlPointCache interface {
	Keys() []ControlPoint
	AddKeyNotifier(notifiers.KeyNotifier) error
	RemoveKeyNotifier(notifiers.KeyNotifier) error
}

// controlPointCache is a cache of discovered Kubernetes control points and provides APIs to do CRUD on Scale type resources.
type controlPointCache struct {
	// RW controlPointsMutex
	controlPointsMutex sync.RWMutex
	k8sClient          k8s.K8sClient
	// Set of unique controlPoints
	controlPoints map[ControlPoint]*controlPointState
	trackers      notifiers.Trackers
	ctx           context.Context
	cancel        context.CancelFunc
	scaleStream   *stream.Stream
}

// controlPointCache implements the ControlPointStore interface.
var _ ControlPointStore = (*controlPointCache)(nil)

// controlPointCache implements the ControlPointCache interface.
var _ ControlPointCache = (*controlPointCache)(nil)

// newControlPointCache returns a new ControlPointCache.
func newControlPointCache(trackers notifiers.Trackers, k8sClient k8s.K8sClient) *controlPointCache {
	return &controlPointCache{
		controlPoints: make(map[ControlPoint]*controlPointState),
		trackers:      trackers,
		k8sClient:     k8sClient,
		scaleStream:   stream.New(),
	}
}

// start starts the autoScaler.
func (cpc *controlPointCache) start() {
	cpc.ctx, cpc.cancel = context.WithCancel(context.Background())
}

// stop stops the autoScaler.
func (cpc *controlPointCache) stop() {
	cpc.cancel()
}

// Add adds a ControlPoint to the cache.
func (cpc *controlPointCache) Add(cp ControlPoint) {
	log.Info().Msgf("Add called for %v", cp)
	// take write mutex before modifying map
	cpc.controlPointsMutex.Lock()
	defer cpc.controlPointsMutex.Unlock()
	// context for fetching scale subresource
	ctx, cancel := context.WithCancel(cpc.ctx)
	cps := &controlPointState{
		cancel: cancel,
		ctx:    ctx,
	}
	cpc.controlPoints[cp] = cps

	// Instead of launching a go routine, use sourcegraph/conc library to create a Stream and submit tasks to it.
	// This will allow us to call the WriteEvent from fetchScale in order of arrival.
	cpc.scaleStream.Go(func() stream.Callback {
		return cpc.fetchScale(cp, cps)
	})
}

// Update updates a ControlPoint in the cache.
func (cpc *controlPointCache) Update(cp ControlPoint) {
	log.Info().Msgf("Update called for %v", cp)
	// take write mutex before modifying map
	cpc.controlPointsMutex.Lock()
	defer cpc.controlPointsMutex.Unlock()

	// get current control point state
	cpsOld, ok := cpc.controlPoints[cp]
	if !ok {
		log.Error().Msgf("Control point %v not found in cache", cp)
		return
	}

	log.Info().Msgf("Canceling goroutine for %v", cp)
	// cancel current goroutine
	cpsOld.cancel()

	// context for fetching scale subresource
	ctx, cancel := context.WithCancel(cpc.ctx)
	// construct new control point state
	cpsNew := &controlPointState{
		cancel: cancel,
		ctx:    ctx,
	}
	// update control point state
	cpc.controlPoints[cp] = cpsNew

	// Fetch scale subresource in a goroutine
	cpc.scaleStream.Go(func() stream.Callback {
		return cpc.fetchScale(cp, cpsNew)
	})
}

// Delete deletes a ControlPoint from the cache.
func (cpc *controlPointCache) Delete(cp ControlPoint) {
	log.Info().Msgf("Delete called for %v", cp)
	// take write mutex before modifying map
	cpc.controlPointsMutex.Lock()
	defer cpc.controlPointsMutex.Unlock()
	cpsOld, ok := cpc.controlPoints[cp]
	if !ok {
		log.Error().Msgf("Control point %v not found in cache", cp)
		return
	}
	log.Info().Msgf("Canceling goroutine for %v", cp)
	cpsOld.cancel()
	delete(cpc.controlPoints, cp)

	key, keyErr := json.Marshal(cp)
	if keyErr != nil {
		log.Error().Err(keyErr).Msgf("Unable to marshal key: %v", cp)
		return
	}

	cpc.scaleStream.Go(func() stream.Callback {
		return func() { cpc.trackers.RemoveEvent(notifiers.Key(key)) }
	})
}

func (cpc *controlPointCache) fetchScale(cp ControlPoint, cps *controlPointState) stream.Callback {
	log.Info().Msgf("fetchScale called for %v", cp)
	noOp := func() {}

	targetGK := schema.GroupKind{
		Group: cp.Group,
		Kind:  cp.Kind,
	}

	// Fetch scale under backoff.Retry operation
	var (
		scale *autoscalingv1.Scale
		err   error
	)
	operation := func() error {
		scale, _, err = cpc.k8sClient.ScaleForGroupKind(cps.ctx, cp.Namespace, cp.Name, targetGK)
		// if cps.ctx is closed, return PermanentError
		if cps.ctx.Err() != nil {
			return backoff.Permanent(cps.ctx.Err())
		}
		if err != nil {
			// TODO: update status
			log.Error().Err(err).Msgf("Unable to get scale for %v", cp)
			return err
		}

		log.Info().Msgf("Scale subresource for %s/%s: %v", cp.Kind, cp.Name, scale)
		return nil
	}

	merr := backoff.Retry(operation, backoff.WithContext(backoff.NewExponentialBackOff(), cps.ctx))
	if merr != nil {
		log.Error().Err(merr).Msgf("Context canceled while fetching scale for %v", cp)
		return noOp
	}

	// Write event to eventWriter
	reported := policysyncv1.ScaleStatus{
		ConfiguredReplicas: scale.Spec.Replicas,
		ActualReplicas:     scale.Status.Replicas,
	}

	key, keyErr := json.Marshal(cp)
	if keyErr != nil {
		log.Error().Err(keyErr).Msgf("Unable to marshal key: %v", cp)
		return noOp
	}

	value, valErr := proto.Marshal(&reported)
	if valErr != nil {
		log.Error().Err(valErr).Msg("Unable to marshal value")
		return noOp
	}

	return func() {
		log.Info().Msgf("Writing event for %v, event: %v", cp, *scale)
		cpc.trackers.WriteEvent(notifiers.Key(key), value)
	}
}

// Keys returns the list of ControlPoints in the cache.
func (cpc *controlPointCache) Keys() []ControlPoint {
	// take read mutex before reading map
	cpc.controlPointsMutex.RLock()
	defer cpc.controlPointsMutex.RUnlock()
	var cps []ControlPoint
	for cp := range cpc.controlPoints {
		cps = append(cps, cp)
	}
	return cps
}

// AddKeyNotifier adds a KeyNotifier to the trackers.
func (cpc *controlPointCache) AddKeyNotifier(notifier notifiers.KeyNotifier) error {
	return cpc.trackers.AddKeyNotifier(notifier)
}

// RemoveKeyNotifier removes a KeyNotifier from the trackers.
func (cpc *controlPointCache) RemoveKeyNotifier(notifier notifiers.KeyNotifier) error {
	return cpc.trackers.RemoveKeyNotifier(notifier)
}

type controlPointState struct {
	// cancel is the function to cancel the context for getting the scale
	cancel context.CancelFunc
	// ctx is the context for getting the scale
	ctx context.Context
}
