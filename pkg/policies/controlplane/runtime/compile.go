package runtime

import (
	"errors"
	"fmt"

	"github.com/looplab/tarjan"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/mapstruct"
)

// ConfiguredComponent consists of a Component, its PortMapping and its Config.
type ConfiguredComponent struct {
	Component
	// Which signals this component wants to have connected on its ports.
	PortMapping PortMapping
	// Mapstruct representation of proto config that was used to create this
	// component.  This Config is used only for observability purposes.
	//
	// Note: PortMapping is also part of Config.
	Config mapstruct.Object
}

// Compile compiles list of configured components into a circuit and validates it.
func Compile(
	configuredComponents []ConfiguredComponent,
	logger *log.Logger,
) ([]CompiledComponent, error) {
	// A list of compiled components. The index of Component in the list is
	// referred as componentIndex. Order of components is the same as in
	// configuredComponents.
	compiledCircuit := make([]CompiledComponent, 0, len(configuredComponents))

	// Map from signal name to a list of componentIndex(es) which accept the signal as input.
	inSignals := make(map[string][]int)
	// Map from signal name to the componentIndex which emits the signal as output.
	outSignals := make(map[string]int)

	for componentIndex, comp := range configuredComponents {
		logger.Trace().Str("componentName", comp.Name()).Interface("ports", comp.PortMapping).Send()

		inPortToSignals := getInPortSignals(comp.PortMapping.Ins, inSignals, componentIndex)
		logger.Trace().Interface("inPortToSignals", inPortToSignals).Send()

		var err error
		outPortToSignals, err := getOutPortSignals(comp.PortMapping.Outs, outSignals, componentIndex)
		if err != nil {
			return nil, err
		}
		logger.Trace().Interface("outPortToSignals", outPortToSignals).Send()

		compiledCircuit = append(compiledCircuit, CompiledComponent{
			Component:        comp.Component,
			InPortToSignals:  inPortToSignals,
			OutPortToSignals: outPortToSignals,
		})
	}

	// Sanitization of inSignals i.e. all inSignals should be defined in outSignals
	for signal := range inSignals {
		if _, ok := outSignals[signal]; !ok {
			return nil, errors.New("undefined signal: " + signal)
		}
	}

	// Run loop detection and mark any looped signals
	// Create a graph for Tarjan's algorithm
	graph := make(map[interface{}][]interface{})

	for componentIndex, comp := range compiledCircuit {
		destCompIndexes := make([]interface{}, 0)
		for _, signals := range comp.OutPortToSignals {
			for _, signal := range signals {
				// Lookup signal in inSignals
				componentIndexes, ok := inSignals[signal.Name]
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
				return nil, errors.New("loop contains non-int component id")
			}
			smallestCompIndexLoopIdx := 0
			for loopIdx, compIndexIfc := range loop {
				componentIndex, ok := compIndexIfc.(int)
				if !ok {
					return nil, errors.New("loop contains non-int component id")
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
			if removeToCompIndex >= len(compiledCircuit) {
				return nil, errors.New("removeToCompId is out of range")
			}
			removeToComp := compiledCircuit[removeToCompIndex]
			if removeFromCompIndex >= len(compiledCircuit) {
				return nil, errors.New("removeFromCompId is out of range")
			}
			removeFromComp := compiledCircuit[removeFromCompIndex]
			loopedSignals := make(map[string]bool)
			// Mark looped signals in InPortToSignalsMap
			for _, signals := range removeToComp.InPortToSignals {
				for idx, signal := range signals {
					if signal.SignalType == SignalTypeNamed {
						outFromCompID, ok := outSignals[signal.Name]
						if !ok {
							return nil, fmt.Errorf("unexpected state: signal %s is not defined in outSignals", signal.Name)
						}
						if outFromCompID == removeFromCompIndex {
							// Mark signal as looped
							signals[idx].Looped = true
							loopedSignals[signal.Name] = true
						}
					}
				}
			}
			// Mark looped signals in OutPortToSignalsMap
			for _, signals := range removeFromComp.OutPortToSignals {
				for idx, signal := range signals {
					if _, ok := loopedSignals[signal.Name]; ok {
						// Mark signal as looped
						signals[idx].Looped = true
					}
				}
			}
		} else {
			// Loop is empty
			return nil, errors.New("got an empty loop from tarjan.Connections")
		}
	}

	// Log compiledCircuit
	for compIndex, compiledComp := range compiledCircuit {
		logger.Trace().Msgf("compIndex: %d, compiledComp: %+v", compIndex, compiledComp)
	}

	return compiledCircuit, nil
}

func getInPortSignals(
	portMapping map[string][]Port,
	signalConsumers map[string][]int,
	componentIndex int,
) PortToSignal {
	portToSignal := getPortSignals(portMapping, componentIndex)

	for _, signals := range portToSignal {
		for _, signal := range signals {
			if signal.SignalType == SignalTypeNamed {
				signalConsumers[signal.Name] = append(signalConsumers[signal.Name], componentIndex)
			}
		}
	}

	return portToSignal
}

func getOutPortSignals(
	portMapping map[string][]Port,
	signalProducers map[string]int,
	componentIndex int,
) (PortToSignal, error) {
	portToSignal := getPortSignals(portMapping, componentIndex)

	for _, signals := range portToSignal {
		for _, signal := range signals {
			if signal.SignalType == SignalTypeNamed {
				if _, ok := signalProducers[signal.Name]; !ok {
					signalProducers[signal.Name] = componentIndex
				} else {
					return nil, errors.New("duplicate signal definition for signal name: " + signal.Name)
				}
			}
		}
	}

	return portToSignal, nil
}

// getPortSignals takes a port mapping and returns a PortToSignal map and signals list.
func getPortSignals(portMapping map[string][]Port, componentIndex int) PortToSignal {
	portToSignalMapping := make(PortToSignal)

	for port, portList := range portMapping {
		portSignals := make([]Signal, 0, len(portList))
		for _, portSpec := range portList {
			if portSpec.SignalName != nil {
				portSignals = append(
					portSignals,
					MakeNamedSignal(*portSpec.SignalName, false),
				)
			} else if portSpec.ConstantSignal != nil {
				portSignals = append(
					portSignals,
					MakeConstantSignal(portSpec.ConstantSignal),
				)
			}
		}
		portToSignalMapping[port] = portSignals
	}

	return portToSignalMapping
}
