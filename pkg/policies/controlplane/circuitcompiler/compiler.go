package circuitcompiler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/looplab/tarjan"
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Module for circuit compiler run via the main app.
func Module() fx.Option {
	return fx.Options(
		runtime.CircuitModule(),
	)
}

// Component is composed of runtime.Component, ComponentID and ParentComponentID.
type Component struct {
	runtime.CompiledComponentAndPorts
	ComponentID       string
	ParentComponentID string
}

// Circuit is a list of CompiledComponent(s).
type Circuit []*Component

// ToComponentsWithPorts converts circuit to list of CompiledComponentAndPorts
// (dropping information about ComponentIDs).
func (c Circuit) ToComponentsWithPorts() []runtime.CompiledComponentAndPorts {
	componentsWithPorts := make([]runtime.CompiledComponentAndPorts, 0, len(c))
	for _, component := range c {
		componentsWithPorts = append(componentsWithPorts, component.CompiledComponentAndPorts)
	}
	return componentsWithPorts
}

// Compile compiles a protobuf circuit definition into a Circuit.
func Compile(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (Circuit, fx.Option, error) {
	logger := policyReadAPI.GetStatusRegistry().GetLogger()
	// List of runtime.CompiledComponent. The index of CompiledComponent in compiledCircuit is referred as graphNodeIndex.
	var compiledCircuit Circuit

	// List of Fx options for components.
	componentOptions := []fx.Option{}

	for compIndex, componentProto := range circuitProto {
		// Create component
		compiledComp, compiledSubComps, compOption, compErr := components.NewComponentAndOptions(componentProto, compIndex, policyReadAPI)
		if compErr != nil {
			return nil, fx.Options(), compErr
		}
		componentOptions = append(componentOptions, compOption)

		compID := strconv.Itoa(compIndex)
		// Add Component to compiledCircuit
		compiledCircuit = append(compiledCircuit, &Component{
			CompiledComponentAndPorts: runtime.CompiledComponentAndPorts{
				CompiledComponent:   compiledComp,
				InPortToSignalsMap:  make(runtime.PortToSignal),
				OutPortToSignalsMap: make(runtime.PortToSignal),
			},
			ComponentID: compID,
		})

		// Add SubComponents to compiledCircuit
		for _, subComp := range compiledSubComps {
			compiledCircuit = append(compiledCircuit, &Component{
				CompiledComponentAndPorts: runtime.CompiledComponentAndPorts{
					CompiledComponent:   subComp,
					InPortToSignalsMap:  make(runtime.PortToSignal),
					OutPortToSignalsMap: make(runtime.PortToSignal),
				},
				ComponentID:       compID + "." + subComp.Component.Name(),
				ParentComponentID: compID,
			})
		}
	}

	// Map from signal name to a list of graphNodeIndex(es) which accept the signal as input.
	inSignals := make(map[string][]int)
	// Map from signal name to the graphNodeIndex which emits the signal as output.
	outSignals := make(map[string]int)

	for graphNodeIndex, compiledComp := range compiledCircuit {
		logger.Trace().
			Str("componentName", compiledComp.Component.Name()).
			Interface("ports", compiledComp.Ports).
			Send()

		compiledComp.InPortToSignalsMap = getInPortSignals(
			compiledComp.Ports.InPorts,
			inSignals,
			graphNodeIndex,
		)
		logger.Trace().Interface("inPortToSignalsMap", compiledComp.InPortToSignalsMap).Send()

		var err error
		compiledComp.OutPortToSignalsMap, err = getOutPortSignals(
			compiledComp.Ports.OutPorts,
			outSignals,
			graphNodeIndex,
		)
		if err != nil {
			return nil, nil, err
		}
		logger.Trace().Interface("outPortToSignalsMap", compiledComp.OutPortToSignalsMap).Send()
	}

	// Sanitization of inSignals i.e. all inSignals should be defined in outSignals
	for signal := range inSignals {
		if _, ok := outSignals[signal]; !ok {
			return nil, fx.Options(), errors.New("undefined signal: " + signal)
		}
	}

	// Run loop detection and mark any looped signals
	// Create a graph for Tarjan's algorithm
	graph := make(map[interface{}][]interface{})

	for graphNodeIndex, compWithPorts := range compiledCircuit {
		destCompIndexes := make([]interface{}, 0)
		for _, signals := range compWithPorts.OutPortToSignalsMap {
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
		// Add graphNodeIndex:destCompIndexes in graph
		if len(destCompIndexes) > 0 {
			graph[graphNodeIndex] = destCompIndexes
		}
	}

	// Run Tarjan's algorithm for detecting loops
	loops := tarjan.Connections(graph)
	// Log loops and graph
	logger.Trace().Msgf("Tarjan Loops: %+v \nTarjan Graph: %+v", loops, graph)

	// Iterate over loops
	for _, loop := range loops {
		// Need to break loop at smallest graphNodeIndex. Find smallest graphNodeIndex in loop
		if len(loop) > 0 {
			smallestCompIndex, ok := loop[0].(int)
			if !ok {
				return nil, fx.Options(), errors.New("loop contains non-int component id")
			}
			smallestCompIndexLoopIdx := 0
			for loopIdx, compIndexIfc := range loop {
				graphNodeIndex, ok := compIndexIfc.(int)
				if !ok {
					return nil, fx.Options(), errors.New("loop contains non-int component id")
				}
				if graphNodeIndex < smallestCompIndex {
					smallestCompIndex = graphNodeIndex
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
				return nil, fx.Options(), errors.New("removeToCompId is out of range")
			}
			removeToComp := compiledCircuit[removeToCompIndex]
			if removeFromCompIndex >= len(compiledCircuit) {
				return nil, fx.Options(), errors.New("removeFromCompId is out of range")
			}
			removeFromComp := compiledCircuit[removeFromCompIndex]
			loopedSignals := make(map[string]bool)
			// Mark looped signals in InPortToSignalsMap
			for _, signals := range removeToComp.InPortToSignalsMap {
				for idx, signal := range signals {
					if signal.SignalType == runtime.SignalTypeNamed {
						outFromCompID, ok := outSignals[signal.Name]
						if !ok {
							return nil, fx.Options(), fmt.Errorf("unexpected state: signal %s is not defined in outSignals", signal.Name)
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
			for _, signals := range removeFromComp.OutPortToSignalsMap {
				for idx, signal := range signals {
					if _, ok := loopedSignals[signal.Name]; ok {
						// Mark signal as looped
						signals[idx].Looped = true
					}
				}
			}
		} else {
			// Loop is empty
			return nil, fx.Options(), errors.New("got an empty loop from tarjan.Connections")
		}
	}

	// Log compiledCircuit
	for compIndex, compiledComp := range compiledCircuit {
		logger.Trace().Msgf("compIndex: %d, compiledComp: %+v", compIndex, compiledComp)
	}

	return compiledCircuit, fx.Options(componentOptions...), nil
}

func getInPortSignals(
	portMapping map[string][]runtime.Port,
	signalConsumers map[string][]int,
	graphNodeIndex int,
) runtime.PortToSignal {
	portToSignal := getPortSignals(portMapping, graphNodeIndex)

	for _, signals := range portToSignal {
		for _, signal := range signals {
			if signal.SignalType == runtime.SignalTypeNamed {
				signalConsumers[signal.Name] = append(signalConsumers[signal.Name], graphNodeIndex)
			}
		}
	}

	return portToSignal
}

func getOutPortSignals(
	portMapping map[string][]runtime.Port,
	signalProducers map[string]int,
	graphNodeIndex int,
) (runtime.PortToSignal, error) {
	portToSignal := getPortSignals(portMapping, graphNodeIndex)

	for _, signals := range portToSignal {
		for _, signal := range signals {
			if signal.SignalType == runtime.SignalTypeNamed {
				if _, ok := signalProducers[signal.Name]; !ok {
					signalProducers[signal.Name] = graphNodeIndex
				} else {
					return nil, errors.New("duplicate signal definition for signal name: " + signal.Name)
				}
			}
		}
	}

	return portToSignal, nil
}

// getPortSignals takes a port mapping and returns a PortToSignal map and signals list.
func getPortSignals(portMapping map[string][]runtime.Port, graphNodeIndex int) runtime.PortToSignal {
	portToSignalMapping := make(runtime.PortToSignal)

	for port, portList := range portMapping {
		portSignals := make([]runtime.Signal, 0, len(portList))
		for _, portSpec := range portList {
			if portSpec.SignalName != nil {
				portSignals = append(
					portSignals,
					runtime.MakeNamedSignal(*portSpec.SignalName, false),
				)
			} else if portSpec.ConstantValue != nil {
				portSignals = append(
					portSignals,
					runtime.MakeConstantSignal(*portSpec.ConstantValue),
				)
			}
		}
		portToSignalMapping[port] = portSignals
	}

	return portToSignalMapping
}
