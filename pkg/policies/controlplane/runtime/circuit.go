package runtime

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"

	policymonitoringv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/monitoring/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/metrics"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/status"
)

// CircuitModule returns fx options of Circuit for the main app.
func CircuitModule() fx.Option {
	return fx.Options(
		fx.Invoke(setupCircuitMetrics),
	)
}

// CircuitMetrics holds prometheus metrics related circuit.
type CircuitMetrics struct {
	SignalSummaryVec            *prometheus.SummaryVec
	InvalidSignalReadingsTotal  *prometheus.CounterVec
	InfiniteSignalReadingsTotal *prometheus.CounterVec
}

var circuitMetrics = newCircuitMetrics()

func newCircuitMetrics() *CircuitMetrics {
	circuitMetricsLabels := []string{
		metrics.SignalNameLabel,
		metrics.SubCircuitIDLabel,
		metrics.PolicyNameLabel,
	}
	circuitMetrics := CircuitMetrics{
		SignalSummaryVec: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: metrics.SignalReadingMetricName,
				Help: "The reading from a signal",
				Objectives: map[float64]float64{
					0:    0,
					0.01: 0.001,
					0.05: 0.01,
					0.5:  0.05,
					0.9:  0.01,
					0.99: 0.001,
					1:    0,
				},
			},
			circuitMetricsLabels,
		),
		InvalidSignalReadingsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.InvalidSignalReadingsTotalMetricName,
				Help: "The number of invalid readings from a signal",
			},
			circuitMetricsLabels,
		),
		InfiniteSignalReadingsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metrics.InfiniteSignalReadingsTotalMetricName,
				Help: "The number of infinite readings from a signal",
			},
			circuitMetricsLabels,
		),
	}
	return &circuitMetrics
}

func setupCircuitMetrics(prometheusRegistry *prometheus.Registry, lifecycle fx.Lifecycle) {
	type metric struct {
		metric prometheus.Collector
		name   string
	}
	circuitMetrics := []metric{
		{circuitMetrics.SignalSummaryVec, metrics.SignalReadingMetricName},
		{circuitMetrics.InvalidSignalReadingsTotal, metrics.InvalidSignalReadingsTotalMetricName},
		{circuitMetrics.InfiniteSignalReadingsTotal, metrics.InfiniteSignalReadingsTotalMetricName},
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			for _, m := range circuitMetrics {
				err := prometheusRegistry.Register(m.metric)
				if err != nil {
					return fmt.Errorf("failed to register metric %s: %w", m.name, err)
				}
			}
			return nil
		},
		OnStop: func(context.Context) error {
			for _, m := range circuitMetrics {
				unregistered := prometheusRegistry.Unregister(m.metric)
				if !unregistered {
					return fmt.Errorf("failed to unregister metric %s", m.name)
				}
			}
			return nil
		},
	})
}

type signalToReading map[Signal]Reading

// Circuit manages the runtime state of a set of components and their inter linkages via signals.
type Circuit struct {
	// Policy Read API
	iface.Policy
	// Status registry to track signal values
	statusRegistry status.Registry
	// Current tick info
	tickInfo TickInfo
	// Looped signals persistence across ticks
	loopedSignals signalToReading
	// Background scheduler
	backgroundScheduler *backgroundScheduler
	// Components
	components []*ConfiguredComponent
	// Tick end callbacks
	tickEndCallbacks []TickEndCallback
	// Tick start callbacks
	tickStartCallbacks []TickStartCallback
	// Execution lock is taken when circuit needs to execute
	executionLock sync.Mutex
}

// Make sure Circuit complies with CircuitAPI interface.
var _ CircuitSuperAPI = &Circuit{}

// NewCircuitAndOptions create a new Circuit struct along with fx options.
func NewCircuitAndOptions(
	configuredComponents []*ConfiguredComponent,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option) {
	reg := policyReadAPI.GetStatusRegistry().Child("circuit", "circuit_signals")
	circuit := &Circuit{
		Policy:         policyReadAPI,
		loopedSignals:  make(signalToReading),
		components:     configuredComponents,
		statusRegistry: reg,
	}

	// Populate loopedSignals
	for _, component := range circuit.components {
		for _, outPort := range component.PortMapping.Outs {
			for _, signal := range outPort {
				if signal.Looped {
					circuit.loopedSignals[signal] = InvalidReading()
				}
			}
		}
	}

	return circuit, fx.Options(
		fx.Invoke(circuit.setup),
	)
}

// Setup handle lifecycle of the inner metrics of Circuit.
func (circuit *Circuit) setup(backgroundScheduler *backgroundScheduler, lifecycle fx.Lifecycle) {
	circuit.backgroundScheduler = backgroundScheduler

	var circuitMetricsLabels []prometheus.Labels

	for _, component := range circuit.components {
		for _, outPort := range component.PortMapping.Outs {
			for _, signal := range outPort {
				circuitMetricsLabels = append(circuitMetricsLabels,
					prometheus.Labels{
						metrics.SignalNameLabel:   signal.SignalName,
						metrics.SubCircuitIDLabel: signal.SubCircuitID,
						metrics.PolicyNameLabel:   circuit.GetPolicyName(),
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
					err = fmt.Errorf("%w: failed to create metrics for %+v", err, labels)
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
			if merr != nil {
				log.Info().Msgf("stopping circuit: %v", merr)
			}
			return nil
		},
	})
}

// Execute runs one tick of computations of all the Components in the Circuit.
func (circuit *Circuit) Execute(tickInfo TickInfo) error {
	logger := circuit.GetStatusRegistry().GetLogger()
	// Save tickInfo
	circuit.tickInfo = tickInfo
	// Lock execution
	circuit.LockExecution()
	// Defer unlock
	defer circuit.UnlockExecution()
	// errMulti appends errors from Executing all the components
	var errMulti error
	// Invoke TickStartCallback(s)
	for _, sc := range circuit.tickStartCallbacks {
		err := sc(circuit)
		errMulti = multierr.Append(errMulti, err)
	}
	policyID := fmt.Sprintf("%s-%s", circuit.GetPolicyName(), circuit.GetPolicyHash())
	reg := circuit.statusRegistry.Child("policy", policyID)

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
			circuitMetricsLabels := prometheus.Labels{
				metrics.SignalNameLabel:   signal.SignalName,
				metrics.SubCircuitIDLabel: signal.SubCircuitID,
				metrics.PolicyNameLabel:   circuit.Policy.GetPolicyName(),
			}
			signalReadingProto := &policymonitoringv1.SignalReading{
				SignalName: signal.SignalName,
				Valid:      reading.Valid(),
				Value:      reading.Value(),
			}
			signalInfo.SignalReading = append(signalInfo.SignalReading, signalReadingProto)

			if reading.Valid() {
				if math.IsInf(reading.Value(), 0) {
					infReadings, err := circuitMetrics.InfiniteSignalReadingsTotal.GetMetricWith(circuitMetricsLabels)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to get InfiniteSignalReadingsTotal metric")
						panic(err)
					}
					infReadings.Inc()
				} else {
					signalSummary, err := circuitMetrics.SignalSummaryVec.GetMetricWith(circuitMetricsLabels)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to get signal metric")
						panic(err)
					}
					signalSummary.Observe(reading.Value())
				}
			} else {
				invReadings, err := circuitMetrics.InvalidSignalReadingsTotal.GetMetricWith(circuitMetricsLabels)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to get InvalidSignalReadingsTotal metric")
					panic(err)
				}
				invReadings.Inc()
			}

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
			componentInPortReadings := make(PortToReading)
			// Skip if already executed
			if executedComponents[cmpIdx] {
				continue
			}
			// Check readiness of cmp by checking in_ports
			for port, sigs := range cmp.PortMapping.Ins {
				// Reading list for this port
				readingList := make([]Reading, len(sigs))
				// Check if all the sig(s) in sigs are ready
				for index, sig := range sigs {
					if sig.SignalType() == SignalTypeConstant {
						readingList[index] = NewReading(sig.ConstantSignalValue())
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
				Interface("InPortToSignalsMap", cmp.PortMapping.Ins).
				Msg("Executing component")
			// If control reaches this point, the component is ready to execute
			componentOutPortReadings, err := cmp.Execute(
				/* pass signal */
				componentInPortReadings,
				/* pass circuit api */
				circuit,
			)
			if componentOutPortReadings == nil {
				componentOutPortReadings = make(PortToReading)
			}
			// Update executedComponents
			executedComponents[cmpIdx] = true

			if err != nil {
				// Append err to errMulti
				errMulti = multierr.Append(errMulti, err)
			}
			// Fill any missing values from cmp.outPortsMapping in componentOutPortReadings with invalid readings
			for port, signals := range cmp.PortMapping.Outs {
				if _, ok := componentOutPortReadings[port]; !ok {
					// Fill with invalid readings
					componentOutPortReadings[port] = make([]Reading, len(signals))
					for index := range signals {
						componentOutPortReadings[port][index] = InvalidReading()
					}
				} else if len(componentOutPortReadings[port]) < len(signals) {
					// The reading list has fewer reading, currentTick, ticksPerExecution ints compared to portOutsSpec
					// Fill with invalid readings
					for index := len(componentOutPortReadings[port]); index < len(signals); index++ {
						componentOutPortReadings[port][index] = InvalidReading()
					}
				}
			}
			// Update circuitSignalReadings with componentOutPortReadings while iterating through outPortsMapping
			for port, signals := range cmp.PortMapping.Outs {
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
						sigNoLoop := sig
						sigNoLoop.Looped = false
						circuitSignalReadings[sigNoLoop] = readings[index]
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
		err := ec(circuit)
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

// GetTickInfo returns the current tick info.
func (circuit *Circuit) GetTickInfo() TickInfo {
	return circuit.tickInfo
}

// ScheduleConditionalBackgroundJob schedules a background job for one time execution. The job gets scheduled only if currentTick is a multiple of ticksPerExecution. There can be at most a single job with a certain name pending to be run at a time. Subsequent invocations with the same job name overwrite the previous one.
// Warning: This method should only be called by the Components during a tick execution.
func (circuit *Circuit) ScheduleConditionalBackgroundJob(job BackgroundJob, ticksPerExecution int) {
	// Do not schedule if currentTick is not a multiple of ticksPerExecution
	if circuit.tickInfo.Tick()%ticksPerExecution != 0 {
		return
	}
	circuit.backgroundScheduler.scheduleJob(job)
}
