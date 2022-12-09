package circuitcompiler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/looplab/tarjan"
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
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

// Component is runtime.CompiledComponent annotated with IDs.
type Component struct {
	runtime.CompiledComponent
	ComponentID       string
	ParentComponentID string
}

// Circuit is a compiled Circuit
//
// Circuit can also be converted to its graph view.
type Circuit struct {
	components   []runtime.CompiledComponent
	configs      []runtime.ConfiguredComponent
	componentIDs []ComponentID
}

// Components returns a list of CompiledComponents, ready to create runtime.Circuit.
func (circuit *Circuit) Components() []runtime.CompiledComponent { return circuit.components }

// CompileFromProto compiles a protobuf circuit definition into a Circuit.
//
// This is helper for CreateComponents + Compile.
func CompileFromProto(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	configuredComponents, componentIDs, option, err := CreateComponents(circuitProto, policyReadAPI)
	if err != nil {
		return nil, nil, err
	}

	compiledComponents, err := Compile(
		configuredComponents,
		policyReadAPI.GetStatusRegistry().GetLogger(),
	)
	if err != nil {
		return nil, option, err
	}

	return &Circuit{
		configs:      configuredComponents,
		components:   compiledComponents,
		componentIDs: componentIDs,
	}, option, nil
}

// ComponentID is a component identifier based on position in original proto list of components.
type ComponentID string

func subcomponentID(parentID ComponentID, subcomponent runtime.Component) ComponentID {
	return ComponentID(string(parentID) + "." + subcomponent.Name())
}

// ParentID returns ID of parent component.
func (id ComponentID) ParentID() ComponentID {
	parentID, _, _ := strings.Cut(string(id), ".")
	return ComponentID(parentID)
}

// CreateComponents creates circuit components along with their identifiers and fx options.
//
// Note that number of returned components might be greater than number of
// components in circuitProto, as some components may have their subcomponents.
func CreateComponents(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []ComponentID, fx.Option, error) {
	var configuredComponents []runtime.ConfiguredComponent
	var ids []ComponentID
	var options []fx.Option

	for compIndex, componentProto := range circuitProto {
		// Create component
		component, subcomponents, compOption, err := components.NewComponentAndOptions(
			componentProto,
			compIndex,
			policyReadAPI,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		options = append(options, compOption)

		compID := ComponentID(strconv.Itoa(compIndex))
		// Add Component to compiledCircuit
		configuredComponents = append(configuredComponents, component)
		ids = append(ids, compID)

		// Add SubComponents to compiledCircuit
		for _, subComp := range subcomponents {
			configuredComponents = append(configuredComponents, subComp)
			ids = append(ids, subcomponentID(compID, subComp))
		}
	}

	return configuredComponents, ids, fx.Options(options...), nil
}

// Compile compiles list of prepared components into a circuit and validates it.
func Compile(configuredComponents []runtime.ConfiguredComponent, logger *log.Logger) ([]runtime.CompiledComponent, error) {
	// A list of compiled components. The index of Component in the list is
	// referred as componentIndex. Order of components is the same as in
	// configuredComponents.
	compiledCircuit := make([]runtime.CompiledComponent, 0, len(configuredComponents))

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

		compiledCircuit = append(compiledCircuit, runtime.CompiledComponent{
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
					if signal.SignalType == runtime.SignalTypeNamed {
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
	portMapping map[string][]runtime.Port,
	signalConsumers map[string][]int,
	componentIndex int,
) runtime.PortToSignal {
	portToSignal := getPortSignals(portMapping, componentIndex)

	for _, signals := range portToSignal {
		for _, signal := range signals {
			if signal.SignalType == runtime.SignalTypeNamed {
				signalConsumers[signal.Name] = append(signalConsumers[signal.Name], componentIndex)
			}
		}
	}

	return portToSignal
}

func getOutPortSignals(
	portMapping map[string][]runtime.Port,
	signalProducers map[string]int,
	componentIndex int,
) (runtime.PortToSignal, error) {
	portToSignal := getPortSignals(portMapping, componentIndex)

	for _, signals := range portToSignal {
		for _, signal := range signals {
			if signal.SignalType == runtime.SignalTypeNamed {
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
func getPortSignals(portMapping map[string][]runtime.Port, componentIndex int) runtime.PortToSignal {
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
