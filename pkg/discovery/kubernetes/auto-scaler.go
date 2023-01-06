package kubernetes

import (
	"context"
	"encoding/json"
	"sync"

	policysyncv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/sync/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/panichandler"
	"google.golang.org/protobuf/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

// AutoScaler is the interface for the auto-scaler.
type AutoScaler interface {
	Add(cp ControlPoint)
	Update(cp ControlPoint)
	Delete(cp ControlPoint)
	Keys() []ControlPoint
}

// AutoScaler is a cache of discovered Kubernetes control points and provides APIs to do CRUD on Scale type resources.
type autoScaler struct {
	// RW mutex
	mutex       sync.RWMutex
	scaleClient scale.ScalesGetter
	// Set of unique controlPoints
	controlPoints map[ControlPoint]*controlPointState
	eventWriter   notifiers.EventWriter
}

// autoScaler implements the AutoScaler interface.
var _ AutoScaler = &autoScaler{}

// newAutoScaler returns a new ControlPointCache.
func newAutoScaler(scaleClient scale.ScalesGetter, eventWriter notifiers.EventWriter) AutoScaler {
	return &autoScaler{
		controlPoints: make(map[ControlPoint]*controlPointState),
		scaleClient:   scaleClient,
		eventWriter:   eventWriter,
	}
}

// Add adds a ControlPoint to the cache.
func (as *autoScaler) Add(cp ControlPoint) {
	log.Info().Msgf("Add called for %v", cp)
	// take write mutex before modifying map
	as.mutex.Lock()
	defer as.mutex.Unlock()
	// context for fetching scale subresource
	_, cancel := context.WithCancel(context.Background())
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
	as.mutex.Lock()
	defer as.mutex.Unlock()

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
	as.mutex.Lock()
	defer as.mutex.Unlock()
	cps, ok := as.controlPoints[cp]
	if !ok {
		log.Error().Msgf("Control point %v not found in cache", cp)
		return
	}
	cps.cancel()
	delete(as.controlPoints, cp)
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
		ReplicasInSpec:   scale.Spec.Replicas,
		ReplicasInStatus: scale.Status.Replicas,
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

	as.eventWriter.WriteEvent(notifiers.Key(key), value)
}

// Keys returns the list of ControlPoints in the cache.
func (as *autoScaler) Keys() []ControlPoint {
	// take read mutex before reading map
	as.mutex.RLock()
	defer as.mutex.RUnlock()
	var cps []ControlPoint
	for cp := range as.controlPoints {
		cps = append(cps, cp)
	}
	return cps
}

type controlPointState struct {
	// cancel is the function to cancel the context for getting the scale
	cancel context.CancelFunc
}
