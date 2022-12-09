package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	policymonitoringv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/monitoring/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/status"
)

// CircuitModule returns fx options of Circuit for the main app.
func CircuitModule() fx.Option {
	return fx.Options(
		fx.Invoke(setupCircuitMetrics),
	)
}

// CircuitMetrics holds prometheus metrics related circuit.
type CircuitMetrics struct {
	SignalSummaryVec *prometheus.SummaryVec
}

var circuitMetrics = &CircuitMetrics{}

func setupCircuitMetrics(prometheusRegistry *prometheus.Registry, lifecycle fx.Lifecycle) {
	circuitMetricsLabels := []string{metrics.SignalNameLabel, metrics.PolicyNameLabel, metrics.ValidLabel}
	circuitMetrics.SignalSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       metrics.SignalReadingMetricName,
		Help:       "The reading from a signal",
		Objectives: map[float64]float64{0: 0, 0.01: 0.001, 0.05: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 1: 0},
	}, circuitMetricsLabels)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := prometheusRegistry.Register(circuitMetrics.SignalSummaryVec)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(context.Context) error {
			unregistered := prometheusRegistry.Unregister(circuitMetrics.SignalSummaryVec)
			if !unregistered {
				err := fmt.Errorf("failed to unregister metric %s", metrics.SignalReadingMetricName)
				return err
			}
			return nil
		},
	})
}

// SignalType enum.
type SignalType int

const (
	// SignalTypeNamed is a named signal.
	SignalTypeNamed = iota
	// SignalTypeConstant is a constant signal.
	SignalTypeConstant
)

// Signal represents a logical flow of Readings in a circuit and determines the linkage between Components.
type Signal struct {
	Name       string
	Value      float64
	Looped     bool
	SignalType SignalType
}

// MakeNamedSignal creates a new named Signal.
func MakeNamedSignal(name string, looped bool) Signal {
	return Signal{
		SignalType: SignalTypeNamed,
		Name:       name,
		Looped:     looped,
	}
}

// MakeConstantSignal creates a new constant Signal.
func MakeConstantSignal(value float64) Signal {
	return Signal{
		SignalType: SignalTypeConstant,
		Value:      value,
	}
}

// PortToSignal is a map from port name to a slice of Signals.
type PortToSignal map[string][]Signal

// ComponentType describes the type of a component based on its connectivity in the circuit.
type ComponentType string

const (
	// ComponentTypeStandAlone is a component that does not accept or emit any signals.
	ComponentTypeStandAlone ComponentType = "StandAlone"
	// ComponentTypeSource is a component that emits output signal(s) but does not accept an input signal.
	ComponentTypeSource ComponentType = "Source"
	// ComponentTypeSink is a component that accepts input signal(s) but does not emit an output signal.
	ComponentTypeSink ComponentType = "Sink"
	// ComponentTypeSignalProcessor is a component that accepts input signal(s) and emits output signal(s).
	ComponentTypeSignalProcessor ComponentType = "SignalProcessor"
)

// ConfiguredComponent consists of a Component, its PortMapping and its Config.
//
// This struct is not directly used by the runtime.
type ConfiguredComponent struct {
	Component
	// Which signals this component wants to have connected on its ports.
	PortMapping PortMapping
	// MapStruct representation of proto config that was used to create this
	// component.  This Config is used only for observability purposes.
	//
	// Note: PortMapping is also part of Config.
	Config map[string]any
}

// CompiledComponent consists of a Component and its In and Out ports.
//
// This representation of component is ready to be used by a circuit.
type CompiledComponent struct {
	Component
	// Note: These maps may seem to repeat information from
	// ConfiguredComponent.Ports, but the most important difference (apart from
	// different map value types), is that the following maps have the looped
	// signals resolved.
	InPortToSignals  PortToSignal
	OutPortToSignals PortToSignal
}

type signalToReading map[Signal]Reading

// Circuit manages the runtime state of a set of components and their inter linkages via signals.
type Circuit struct {
	// Execution lock is taken when circuit needs to execute
	executionLock sync.Mutex
	// Policy Read API
	iface.Policy
	// Status registry to track signal values
	statusRegistry status.Registry
	// Looped signals persistence across ticks
	loopedSignals signalToReading
	// Components
	components []CompiledComponent
	// Tick end callbacks
	tickEndCallbacks []TickEndCallback
	// Tick start callbacks
	tickStartCallbacks []TickStartCallback
}

// Make sure Circuit complies with CircuitAPI interface.
var _ CircuitAPI = &Circuit{}

// NewCircuitAndOptions create a new Circuit struct along with fx options.
func NewCircuitAndOptions(
	compiledComponents []CompiledComponent,
	policyReadAPI iface.Policy,
	registry status.Registry,
) (*Circuit, fx.Option) {
	reg := registry.Child("circuit_signals")
	circuit := &Circuit{
		Policy:         policyReadAPI,
		loopedSignals:  make(signalToReading),
		components:     make([]CompiledComponent, 0),
		statusRegistry: reg,
	}

	// setComponents sets the Components for a Circuit
	setComponents := func(components []CompiledComponent) error {
		if len(circuit.components) > 0 {
			return errors.New("circuit already has components")
		}
		circuit.components = components
		// Populate loopedSignals
		for _, component := range components {
			for _, outPort := range component.OutPortToSignals {
				for _, signal := range outPort {
					if signal.Looped {
						circuit.loopedSignals[signal] = InvalidReading()
					}
				}
			}
		}
		return nil
	}

	// Set components in circuit
	err := setComponents(compiledComponents)
	if err != nil {
		// FIXME don't hide this error contents
		return nil, fx.Options()
	}

	return circuit, fx.Options(
		fx.Invoke(circuit.setup),
	)
}

// Setup handle lifecycle of the inner metrics of Circuit.
func (circuit *Circuit) setup(lifecycle fx.Lifecycle) {
	var circuitMetricsLabels []prometheus.Labels

	for _, component := range circuit.components {
		for _, outPort := range component.OutPortToSignals {
			for _, signal := range outPort {
				circuitMetricsLabels = append(circuitMetricsLabels,
					prometheus.Labels{
						metrics.SignalNameLabel: signal.Name,
						metrics.PolicyNameLabel: circuit.GetPolicyName(),
						metrics.ValidLabel:      metrics.ValidFalse,
					},
					prometheus.Labels{
						metrics.SignalNameLabel: signal.Name,
						metrics.PolicyNameLabel: circuit.GetPolicyName(),
						metrics.ValidLabel:      metrics.ValidTrue,
					},
				)
			}
		}
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var merr error
			for _, labels := range circuitMetricsLabels {
				_, err := circuitMetrics.SignalSummaryVec.GetMetricWith(labels)
				if err != nil {
					err = errors.Wrapf(err, "failed to create metrics for %+v", labels)
					merr = multierr.Append(merr, err)
				}
			}
			return merr
		},
		OnStop: func(context.Context) error {
			var merr error
			for _, labels := range circuitMetricsLabels {
				deleted := circuitMetrics.SignalSummaryVec.Delete(labels)
				if !deleted {
					err := fmt.Errorf("failed to delete metrics for %+v", labels)
					merr = multierr.Append(merr, err)
				}
			}
			return merr
		},
	})
}

// Execute runs one tick of computations of all the Components in the Circuit.
func (circuit *Circuit) Execute(tickInfo TickInfo) error {
	logger := circuit.GetStatusRegistry().GetLogger()
	// Lock execution
	circuit.LockExecution()
	// Defer unlock
	defer circuit.UnlockExecution()
	// errMulti appends errors from Executing all the components
	var errMulti error
	// Invoke TickStartCallback(s)
	for _, sc := range circuit.tickStartCallbacks {
		err := sc(tickInfo)
		errMulti = multierr.Append(errMulti, err)
	}
	policyID := fmt.Sprintf("%s-%s", circuit.GetPolicyName(), circuit.GetPolicyHash())
	reg := circuit.statusRegistry.Child(policyID)

	// Signals for this tick
	circuitSignalReadings := make(signalToReading)
	defer func() {
		signalInfo := &policymonitoringv1.SignalMetricsInfo{
			PolicyName:    circuit.GetPolicyName(),
			PolicyHash:    circuit.GetPolicyHash(),
			SignalReading: make([]*policymonitoringv1.SignalReading, 0),
		}
		// log all circuitSignalReadings
		for signal, reading := range circuitSignalReadings {
			signalReadingProto := &policymonitoringv1.SignalReading{
				SignalName: signal.Name,
				Valid:      reading.Valid(),
				Value:      reading.Value(),
			}
			signalInfo.SignalReading = append(signalInfo.SignalReading, signalReadingProto)

			circuitMetricsLabels := prometheus.Labels{
				metrics.SignalNameLabel: signal.Name,
				metrics.PolicyNameLabel: circuit.Policy.GetPolicyName(),
			}
			if reading.Valid() {
				circuitMetricsLabels[metrics.ValidLabel] = metrics.ValidTrue
			} else {
				circuitMetricsLabels[metrics.ValidLabel] = metrics.ValidFalse
			}
			signalSummaryMetric, err := circuitMetrics.SignalSummaryVec.GetMetricWith(circuitMetricsLabels)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get signal metric")
				panic(err)
			}
			signalSummaryMetric.Observe(reading.Value())
		}
		signalStatus := status.NewStatus(signalInfo, nil)
		reg.SetStatus(signalStatus)
	}()

	// Populate with last run's looped signal
	for sig, reading := range circuit.loopedSignals {
		circuitSignalReadings[sig] = reading
	}
	// Clear looped signals for next tick
	circuit.loopedSignals = signalToReading{}
	// Map of executed components for this tick
	executedComponents := make(map[int]bool)
	// Number of components executed
	var numExecutedBefore, numExecutedAfter int
	// Loop rounds until no components are executed
	for len(executedComponents) < len(circuit.components) {
		numExecutedBefore = len(executedComponents)

	OUTER:
		// Check readiness by component and execute if ready
		for cmpIdx, cmp := range circuit.components {
			componentInPortReadings := make(PortToValue)
			// Skip if already executed
			if executedComponents[cmpIdx] {
				continue
			}
			// Check readiness of cmp by checking in_ports
			for port, sigs := range cmp.InPortToSignals {
				// Reading list for this port
				readingList := make([]Reading, len(sigs))
				// Check if all the sig(s) in sigs are ready
				for index, sig := range sigs {
					if sig.SignalType == SignalTypeConstant {
						readingList[index] = NewReading(sig.Value)
					} else if sigReading, ok := circuitSignalReadings[sig]; ok {
						// Set sigReading in readingList at index
						readingList[index] = sigReading
					} else {
						// not ready yet
						continue OUTER
					}
				}
				// All the sig(s) in sigs are ready, set readingList in componentInPortReadings
				componentInPortReadings[port] = readingList
			}
			// log the component being executed
			logger.Trace().Str("component", cmp.Name()).
				Int("tick", tickInfo.Tick()).
				Interface("in_ports", componentInPortReadings).
				Interface("InPortToSignalsMap", cmp.InPortToSignals).
				Msg("Executing component")
			// If control reaches this point, the component is ready to execute
			componentOutPortReadings, err := cmp.Execute(
				/* pass signal */
				componentInPortReadings,
				/* pass tick info */
				tickInfo,
			)
			if componentOutPortReadings == nil {
				componentOutPortReadings = make(PortToValue)
			}
			// Update executedComponents
			executedComponents[cmpIdx] = true

			if err != nil {
				// Append err to errMulti
				errMulti = multierr.Append(errMulti, err)
			}
			// Fill any missing values from cmp.outPortsMapping in componentOutPortReadings with invalid readings
			for port, signals := range cmp.OutPortToSignals {
				if _, ok := componentOutPortReadings[port]; !ok {
					// Fill with invalid readings
					componentOutPortReadings[port] = make([]Reading, len(signals))
					for index := range signals {
						componentOutPortReadings[port][index] = InvalidReading()
					}
				} else if len(componentOutPortReadings[port]) < len(signals) {
					// The reading list has fewer readings compared to portOutsSpec
					// Fill with invalid readings
					for index := len(componentOutPortReadings[port]); index < len(signals); index++ {
						componentOutPortReadings[port][index] = InvalidReading()
					}
				}
			}
			// Update circuitSignalReadings with componentOutPortReadings while iterating through outPortsMapping
			for port, signals := range cmp.OutPortToSignals {
				for index, sig := range signals {
					readings, ok := componentOutPortReadings[port]
					if !ok {
						// Create error message
						errMsg := fmt.Sprintf("unexpected state: port %s is not defined in componentOutPortReadings. abort circuit execution", port)
						// Log error
						logger.Error().Msg(errMsg)
						return errors.New(errMsg)
					}
					// check presence of index in readings
					if index >= len(readings) {
						// Create error message
						errMsg := fmt.Sprintf("unexpected state: index %d is out of range in port %s. abort circuit execution", index, port)
						// Log error
						logger.Error().Msg(errMsg)
						return errors.New(errMsg)
					}
					if sig.Looped {
						// Looped signals are stored in circuit.loopedSignals for the next round
						circuit.loopedSignals[sig] = readings[index]
						// Store the reading in circuitSignalReadings under the same signal name without the looped flag
						circuitSignalReadings[MakeNamedSignal(sig.Name, false)] = readings[index]
					} else {
						// Store the reading in circuitSignalReadings
						circuitSignalReadings[sig] = readings[index]
					}
				}
			}
		} // this is the component for loop
		numExecutedAfter = len(executedComponents)

		// Return with error if there is no change in number of executed components.
		if numExecutedBefore == numExecutedAfter {
			errMsg := fmt.Sprintf("circuit execution failed. number of executed components (%d) remained same across consecutive rounds", numExecutedBefore)
			logger.Error().Msg(errMsg)
			return errors.New(errMsg)
		}
	}
	// Invoke TickEndCallback(s)
	for _, ec := range circuit.tickEndCallbacks {
		err := ec(tickInfo)
		errMulti = multierr.Append(errMulti, err)
	}
	return errMulti
}

// DynamicConfigUpdate updates the circuit with the new dynamic config.
func (circuit *Circuit) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	// Take Circuit execution lock
	circuit.LockExecution()
	defer circuit.UnlockExecution()
	// loop through all the components
	for _, cmp := range circuit.components {
		// update the dynamic config
		cmp.DynamicConfigUpdate(event, unmarshaller)
	}
}

// LockExecution locks the execution of the circuit.
func (circuit *Circuit) LockExecution() {
	circuit.executionLock.Lock()
}

// UnlockExecution unlocks the execution of the circuit.
func (circuit *Circuit) UnlockExecution() {
	circuit.executionLock.Unlock()
}

// RegisterTickEndCallback adds a callback function to be called when a tick ends.
func (circuit *Circuit) RegisterTickEndCallback(ec TickEndCallback) {
	circuit.tickEndCallbacks = append(circuit.tickEndCallbacks, ec)
}

// RegisterTickStartCallback adds a callback function to be called when a tick starts.
func (circuit *Circuit) RegisterTickStartCallback(sc TickStartCallback) {
	circuit.tickStartCallbacks = append(circuit.tickStartCallbacks, sc)
}
