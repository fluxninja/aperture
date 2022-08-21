package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/metrics"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
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
	circuitMetricsLabels := []string{metrics.SignalNameLabel, metrics.PolicyNameLabel}
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

// Signal represents a logical flow of Readings in a circuit and determines the linkage between Components.
type Signal struct {
	Name   string
	Looped bool
}

type signalToReading map[Signal]reading.Reading

// PortToSignal is a map from port name to a slice of Signals.
type PortToSignal map[string][]Signal

// ComponentWithPorts is a Component along with In and Out ports.
type ComponentWithPorts struct {
	Component           Component
	InPortToSignalsMap  PortToSignal
	OutPortToSignalsMap PortToSignal
	ComponentName       string
}

// Circuit manages the runtime state of a set of components and their inter linkages via signals.
type Circuit struct {
	// Policy Read API
	iface.PolicyRead
	// Execution lock is taken when circuit needs to execute
	executionLock sync.Mutex
	// Looped signals persistence across ticks
	loopedSignals signalToReading
	// Components
	components []ComponentWithPorts
	// Tick end callbacks
	tickEndCallbacks []TickEndCallback
}

// Make sure Circuit complies with CircuitAPI interface.
var _ CircuitAPI = &Circuit{}

// NewCircuitAndOptions create a new Circuit struct along with fx options.
func NewCircuitAndOptions(policyReadAPI iface.PolicyRead, compWithPortsList []ComponentWithPorts) (*Circuit, fx.Option) {
	circuit := &Circuit{
		PolicyRead:    policyReadAPI,
		loopedSignals: make(signalToReading),
		components:    make([]ComponentWithPorts, 0),
	}

	// setComponents sets the Components for a Circuit
	setComponents := func(components []ComponentWithPorts) error {
		if len(circuit.components) > 0 {
			return errors.New("circuit already has components")
		}
		circuit.components = components
		// Populate loopedSignals
		for _, component := range components {
			for _, outPort := range component.OutPortToSignalsMap {
				for _, signal := range outPort {
					if signal.Looped {
						circuit.loopedSignals[signal] = reading.NewInvalid()
					}
				}
			}
		}
		return nil
	}

	// Set components in circuit
	err := setComponents(compWithPortsList)
	if err != nil {
		return nil, fx.Options()
	}

	return circuit, fx.Options(
		fx.Invoke(circuit.setup),
	)
}

// Setup handle lifecycle of the inner metrics of Circuit.
func (circuit *Circuit) setup(lifecycle fx.Lifecycle) {
	circuitMetricsLabels := make(map[string]prometheus.Labels, 0)

	for _, component := range circuit.components {
		for _, outPort := range component.OutPortToSignalsMap {
			for _, signal := range outPort {
				l := prometheus.Labels{
					metrics.SignalNameLabel: signal.Name,
					metrics.PolicyNameLabel: circuit.GetPolicyName(),
				}
				circuitMetricsLabels[signal.Name] = l
			}
		}
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			var merr error
			for signalName, labels := range circuitMetricsLabels {
				_, err := circuitMetrics.SignalSummaryVec.GetMetricWith(labels)
				if err != nil {
					err = errors.Wrapf(err, "failed to create metrics for %s", signalName)
					merr = multierr.Append(merr, err)
				}
			}
			return merr
		},
		OnStop: func(context.Context) error {
			var merr error
			for signalName, labels := range circuitMetricsLabels {
				deleted := circuitMetrics.SignalSummaryVec.Delete(labels)
				if !deleted {
					err := fmt.Errorf("failed to delete metrics for %s", signalName)
					merr = multierr.Append(merr, err)
				}
			}
			return merr
		},
	})
}

// Execute runs one tick of computations of all the Components in the Circuit.
func (circuit *Circuit) Execute(tickInfo TickInfo) error {
	// Lock execution
	circuit.LockExecution()
	// Defer unlock
	defer circuit.UnlockExecution()
	// errMulti appends errors from Executing all the components
	var errMulti error
	// Signals for this tick
	circuitSignalReadings := make(signalToReading)
	defer func() {
		// log all circuitSignalReadings
		for signal, reading := range circuitSignalReadings {
			circuitMetricsLabels := prometheus.Labels{
				metrics.SignalNameLabel: signal.Name,
				metrics.PolicyNameLabel: circuit.PolicyRead.GetPolicyName(),
			}
			signalSummaryVecMetric, err := circuitMetrics.SignalSummaryVec.GetMetricWith(circuitMetricsLabels)
			if err != nil {
				log.Error().Err(err).Msg("Failed to initialize SummaryVec")
			}
			if reading.Valid {
				signalSummaryVecMetric.Observe(reading.Value)
			}
		}
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
			for port, sigs := range cmp.InPortToSignalsMap {
				// Reading list for this port
				readingList := make([]reading.Reading, len(sigs))
				// Check if all the sig(s) in sigs are ready
				for index, sig := range sigs {
					if sigReading, ok := circuitSignalReadings[sig]; ok {
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
			log.Trace().Str("component", cmp.ComponentName).Int("tick", tickInfo.Tick).Interface("in_ports", componentInPortReadings).Interface("InPortToSignalsMap", cmp.InPortToSignalsMap).Msg("Executing component")
			// If control reaches this point, the component is ready to execute
			componentOutPortReadings, err := cmp.Component.Execute(
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
			for port, signals := range cmp.OutPortToSignalsMap {
				if _, ok := componentOutPortReadings[port]; !ok {
					// Fill with invalid readings
					componentOutPortReadings[port] = make([]reading.Reading, len(signals))
					for index := range signals {
						componentOutPortReadings[port][index] = reading.NewInvalid()
					}
				} else if len(componentOutPortReadings[port]) < len(signals) {
					// The reading list has fewer readings compared to portOutsSpec
					// Fill with invalid readings
					for index := len(componentOutPortReadings[port]); index < len(signals); index++ {
						componentOutPortReadings[port][index] = reading.NewInvalid()
					}
				}
			}
			// Update circuitSignalReadings with componentOutPortReadings while iterating through outPortsMapping
			for port, signals := range cmp.OutPortToSignalsMap {
				for index, sig := range signals {
					readings, ok := componentOutPortReadings[port]
					if !ok {
						// Create error message
						errMsg := fmt.Sprintf("unexpected state: port %s is not defined in componentOutPortReadings. abort circuit execution", port)
						// Log error
						log.Error().Msg(errMsg)
						return errors.New(errMsg)
					}
					// check precense of index in readings
					if index >= len(readings) {
						// Create error message
						errMsg := fmt.Sprintf("unexpected state: index %d is out of range in port %s. abort circuit execution", index, port)
						// Log error
						log.Error().Msg(errMsg)
						return errors.New(errMsg)
					}
					if sig.Looped {
						// Looped signals are stored in circuit.loopedSignals for the next round
						circuit.loopedSignals[sig] = readings[index]
						// Store the reading in circuitSignalReadings under the same signal name without the looped flag
						circuitSignalReadings[Signal{Name: sig.Name, Looped: false}] = readings[index]
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
			log.Error().Msg(errMsg)
			return errors.New(errMsg)
		}
	}
	// Invoke TickEndCallback(s)
	for _, cb := range circuit.tickEndCallbacks {
		err := cb(tickInfo)
		errMulti = multierr.Append(errMulti, err)
	}
	return errMulti
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
func (circuit *Circuit) RegisterTickEndCallback(cb TickEndCallback) {
	circuit.tickEndCallbacks = append(circuit.tickEndCallbacks, cb)
}
