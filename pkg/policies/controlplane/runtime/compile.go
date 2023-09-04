package runtime

import (
	"errors"
	"fmt"

	"github.com/looplab/tarjan"

	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Compile compiles list of configured components into a circuit and validates it.
func Compile(
	configuredComponents []*ConfiguredComponent,
	logger *log.Logger,
) error {
	// Map from signal name to a list of componentIndex(es) which accept the signal as input.
	inSignals := make(map[SignalID][]int)
	// Map from signal name to the componentIndex which emits the signal as output.
	outSignals := make(map[SignalID]int)

	for componentIndex, comp := range configuredComponents {
		logger.Trace().Str("componentName", comp.Name()).Interface("ports", comp.PortMapping).Send()

		updateSignalConsumers(comp.PortMapping.Ins, inSignals, componentIndex)

		err := updateSignalProducers(comp.PortMapping.Outs, outSignals, componentIndex)
		if err != nil {
			return err
		}
	}

	// Sanitization of inSignals i.e. all inSignals should be defined in outSignals
	for signalID := range inSignals {
		if _, ok := outSignals[signalID]; !ok {
			return fmt.Errorf("undefined signal: %+v", signalID)
		}
	}

	// Run loop detection and mark any looped signals
	// Create a graph for Tarjan's algorithm
	graph := make(map[interface{}][]interface{})

	for componentIndex, comp := range configuredComponents {
		destCompIndexes := make([]interface{}, 0)
		for _, signals := range comp.PortMapping.Outs {
			for _, signal := range signals {
				// Lookup signal in inSignals
				componentIndexes, ok := inSignals[signal.SignalID()]
				// Convert componentIndexes to []interface{}
				componentIndexesIfc := make([]interface{}, len(componentIndexes))
				for i, componentIndex := range componentIndexes {
					componentIndexesIfc[i] = componentIndex
				}
				if ok {
					// Add componentIndexes to destCompIndexes
					destCompIndexes = append(destCompIndexes, componentIndexesIfc...)
				}
			}
		}
		// Add componentIndex:destCompIndexes in graph
		if len(destCompIndexes) > 0 {
			graph[componentIndex] = destCompIndexes
		}
	}

	// Run Tarjan's algorithm for detecting loops
	loops := tarjan.Connections(graph)
	// Log loops and graph
	logger.Trace().Msgf("Tarjan Loops: %+v \nTarjan Graph: %+v", loops, graph)

	// Iterate over loops
	for _, loop := range loops {
		// Need to break loop at smallest componentIndex. Find smallest componentIndex in loop
		if len(loop) > 0 {
			smallestCompIndex, ok := loop[0].(int)
			if !ok {
				return errors.New("loop contains non-int component id")
			}
			smallestCompIndexLoopIdx := 0
			for loopIdx, compIndexIfc := range loop {
				componentIndex, ok := compIndexIfc.(int)
				if !ok {
					return errors.New("loop contains non-int component id")
				}
				if componentIndex < smallestCompIndex {
					smallestCompIndex = componentIndex
					smallestCompIndexLoopIdx = loopIdx
				}
			}
			// Break loop at smallest compId.
			removeToCompIndex := smallestCompIndex
			// Remove connections from the next component in the loop
			removeFromCompIndexLoopIdx := (smallestCompIndexLoopIdx + 1) % len(loop)
			removeFromCompIndex := loop[removeFromCompIndexLoopIdx].(int)
			// Remove connections from components at removeFromCompIndex to removeToCompIndex
			if removeToCompIndex >= len(configuredComponents) {
				return errors.New("removeToCompId is out of range")
			}
			removeToComp := configuredComponents[removeToCompIndex]
			if removeFromCompIndex >= len(configuredComponents) {
				return errors.New("removeFromCompId is out of range")
			}
			removeFromComp := configuredComponents[removeFromCompIndex]
			loopedSignals := make(map[string]bool)
			// Mark looped signals in Ins
			for _, signals := range removeToComp.PortMapping.Ins {
				for idx, signal := range signals {
					if signal.SignalType() == SignalTypeNamed {
						outFromCompID, ok := outSignals[signal.SignalID()]
						if !ok {
							return fmt.Errorf("unexpected state: signal %s is not defined in outSignals", signal.SignalName)
						}
						if outFromCompID == removeFromCompIndex {
							// Mark signal as looped
							signals[idx].Looped = true
							loopedSignals[signal.SignalName] = true
						}
					}
				}
			}
			// Mark looped signals in Outs
			for _, signals := range removeFromComp.PortMapping.Outs {
				for idx, signal := range signals {
					if _, ok := loopedSignals[signal.SignalName]; ok {
						// Mark signal as looped
						signals[idx].Looped = true
					}
				}
			}
		} else {
			// Loop is empty
			return errors.New("got an empty loop from tarjan.Connections")
		}
	}

	// Log components
	for compIndex, comp := range configuredComponents {
		logger.Trace().Msgf("compIndex: %d, comp: %+v", compIndex, comp)
	}

	return nil
}

func updateSignalConsumers(
	portMapping map[string][]Signal,
	signalConsumers map[SignalID][]int,
	componentIndex int,
) {
	for _, signals := range portMapping {
		for _, signal := range signals {
			if signal.SignalType() == SignalTypeNamed {
				signalConsumers[signal.SignalID()] = append(signalConsumers[signal.SignalID()], componentIndex)
			}
		}
	}
}

func updateSignalProducers(
	portMapping map[string][]Signal,
	signalProducers map[SignalID]int,
	componentIndex int,
) error {
	for _, signals := range portMapping {
		for _, signal := range signals {
			if signal.SignalType() == SignalTypeNamed {
				if _, ok := signalProducers[signal.SignalID()]; !ok {
					signalProducers[signal.SignalID()] = componentIndex
				} else {
					return errors.New("duplicate signal definition for signal name: " + signal.SignalName)
				}
			}
		}
	}

	return nil
}
