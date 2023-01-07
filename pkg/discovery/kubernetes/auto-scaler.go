package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/scale"
)

// A ControlPoint is identified by Group, Version, Type, Namespace and Name.
type ControlPoint struct {
	Group     string
	Version   string
	Type      string
	Namespace string
	Name      string
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

// AutoScaler is a cache of discovered Kubernetes control points and provides APIs to do CRUD on Scale type resources.
type autoScaler struct {
	// RW controlPointsMutex
	controlPointsMutex sync.RWMutex
	scaleClient        scale.ScalesGetter
	// Set of unique controlPoints
	controlPoints map[ControlPoint]*controlPointState
	trackers      notifiers.Trackers
	ctx           context.Context
	cancel        context.CancelFunc
	scaleMutex    sync.Mutex
	scaleCancel   context.CancelFunc
}

// autoScaler implements the AutoScaler interface.
var _ AutoScaler = &autoScaler{}

// autoScaler implements the ControlPointStore interface.
var _ ControlPointStore = &autoScaler{}

// autoScaler implements the ControlPointCache interface.
var _ ControlPointCache = &autoScaler{}

// newAutoScaler returns a new ControlPointCache.
func newAutoScaler(scaleClient scale.ScalesGetter, trackers notifiers.Trackers) *autoScaler {
	return &autoScaler{
		controlPoints: make(map[ControlPoint]*controlPointState),
		scaleClient:   scaleClient,
		trackers:      trackers,
	}
}

// start starts the autoScaler.
func (as *autoScaler) start() {
	as.ctx, as.cancel = context.WithCancel(context.Background())
}

// stop stops the autoScaler.
func (as *autoScaler) stop() {
	as.cancel()
}

// Add adds a ControlPoint to the cache.
func (as *autoScaler) Add(cp ControlPoint) {
	log.Info().Msgf("Add called for %v", cp)
	// take write mutex before modifying map
	as.controlPointsMutex.Lock()
	defer as.controlPointsMutex.Unlock()
	// context for fetching scale subresource
	_, cancel := context.WithCancel(as.ctx)
	cps := &controlPointState{
		cancel: cancel,
	}
	as.controlPoints[cp] = cps

	// Fetch scale subresource in a goroutine
	panichandler.Go(func() {
		as.fetchScale(cp, cps)
	})
}

// Update updates a ControlPoint in the cache.
func (as *autoScaler) Update(cp ControlPoint) {
	log.Info().Msgf("Update called for %v", cp)
	// take write mutex before modifying map
	as.controlPointsMutex.Lock()
	defer as.controlPointsMutex.Unlock()

	// get current control point state
	cps, ok := as.controlPoints[cp]
	if !ok {
		log.Error().Msgf("Control point %v not found in cache", cp)
		return
	}

	// cancel current goroutine
	cps.cancel()

	// context for fetching scale subresource
	_, cancel := context.WithCancel(context.Background())
	cps.cancel = cancel

	// Fetch scale subresource in a goroutine
	panichandler.Go(func() {
		as.fetchScale(cp, cps)
	})
}

// Delete deletes a ControlPoint from the cache.
func (as *autoScaler) Delete(cp ControlPoint) {
	log.Info().Msgf("Delete called for %v", cp)
	// take write mutex before modifying map
	as.controlPointsMutex.Lock()
	defer as.controlPointsMutex.Unlock()
	cps, ok := as.controlPoints[cp]
	if !ok {
		log.Error().Msgf("Control point %v not found in cache", cp)
		return
	}
	cps.cancel()
	delete(as.controlPoints, cp)

	key, keyErr := json.Marshal(cp)
	if keyErr != nil {
		log.Error().Err(keyErr).Msgf("Unable to marshal key: %v", cp)
		return
	}

	as.trackers.RemoveEvent(notifiers.Key(key))
}

func (as *autoScaler) fetchScale(cp ControlPoint, cps *controlPointState) {
	scale, scaleErr := as.scaleClient.Scales(cp.Namespace).Get(context.Background(), schema.GroupResource{Group: cp.Group, Resource: cp.Type}, cp.Name, metav1.GetOptions{})
	if scaleErr != nil {
		log.Error().Err(scaleErr).Msg("Unable to get scale subresource")
		return
	}
	log.Info().Msgf("Scale subresource for %s/%s: %v", cp.Type, cp.Name, scale)

	// Write event to eventWriter
	reported := policysyncv1.KubernetesScaleReported{
		ConfiguredReplicas: scale.Spec.Replicas,
		ActualReplicas:     scale.Status.Replicas,
	}

	key, keyErr := json.Marshal(cp)
	if keyErr != nil {
		log.Error().Err(keyErr).Msgf("Unable to marshal key: %v", cp)
		return
	}

	value, valErr := proto.Marshal(&reported)
	if valErr != nil {
		log.Error().Err(valErr).Msg("Unable to marshal value")
		return
	}

	as.trackers.WriteEvent(notifiers.Key(key), value)
}

// Keys returns the list of ControlPoints in the cache.
func (as *autoScaler) Keys() []ControlPoint {
	// take read mutex before reading map
	as.controlPointsMutex.RLock()
	defer as.controlPointsMutex.RUnlock()
	var cps []ControlPoint
	for cp := range as.controlPoints {
		cps = append(cps, cp)
	}
	return cps
}

// AddKeyNotifier adds a KeyNotifier to the trackers.
func (as *autoScaler) AddKeyNotifier(notifier notifiers.KeyNotifier) error {
	return as.trackers.AddKeyNotifier(notifier)
}

// RemoveKeyNotifier removes a KeyNotifier from the trackers.
func (as *autoScaler) RemoveKeyNotifier(notifier notifiers.KeyNotifier) error {
	return as.trackers.RemoveKeyNotifier(notifier)
}

// Scale scales a Kubernetes resource.
func (as *autoScaler) Scale(cp ControlPoint, replicas int32) {
	// Take mutex to prevent concurrent scale operations
	as.scaleMutex.Lock()
	defer as.scaleMutex.Unlock()
	// Cancel any existing scale operation
	if as.scaleCancel != nil {
		as.scaleCancel()
	}
	ctx, cancel := context.WithCancel(as.ctx)
	as.scaleCancel = cancel
	panichandler.Go(func() {
		log.Info().Msgf("Scale called for %v with replicas %d", cp, replicas)

		data := []byte(fmt.Sprintf(`{"spec":{"replicas":%d}}`, replicas))
		_, err := as.scaleClient.Scales(cp.Namespace).Patch(ctx, schema.GroupVersionResource{Group: cp.Group, Version: cp.Version, Resource: cp.Type}, cp.Name, types.MergePatchType, data, metav1.PatchOptions{})
		if err != nil {
			log.Error().Err(err).Msg("Unable to patch scale subresource")
			return
		}
	})
}

type controlPointState struct {
	// cancel is the function to cancel the context for getting the scale
	cancel context.CancelFunc
}
