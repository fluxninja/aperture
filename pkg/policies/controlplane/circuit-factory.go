package controlplane

import (
	"errors"
	"fmt"

	"github.com/looplab/tarjan"
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// circuitFactoryModule for circuit factory run via the main app.
func circuitFactoryModule() fx.Option {
	return fx.Options(
		runtime.CircuitModule(),
	)
}

// compileCircuit takes a circuitProto and returns list of CompiledComponentAndPorts.
func compileCircuit(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) ([]runtime.CompiledComponentAndPorts, fx.Option, error) {
	// List of runtime.CompiledComponent. The index of runtime.CompiledComponents in compList is referred as graphNodeIndex.
	var compList []runtime.CompiledComponent
	// Map from signal name to a list of graphNodeIndex(es) which accept the signal as input.
	inSignals := make(map[string][]int)
	// Map from signal name to the graphNodeIndex which emits the signal as output.
	outSignals := make(map[string]int)

	// List of Fx options for components.
	componentOptions := []fx.Option{}

	for compIndex, componentProto := range circuitProto {
		// Create component
		compiledComp, compiledSubComps, compOption, compErr := NewComponentAndOptions(componentProto, compIndex, policyReadAPI)
		if compErr != nil {
			return nil, fx.Options(), compErr
		}
		componentOptions = append(componentOptions, compOption)

		// Add Component to compList
		if compiledComp.Component != nil {
			compList = append(compList, compiledComp)
		}

		// Add SubComponents to compList
		if len(compiledSubComps) > 0 {
			compList = append(compList, compiledSubComps...)
		}
	}
	log.Trace().Msgf("Comp list: %+v", compList)

	// Second pass to initialize port maps for each component
	compWithPortsList := make([]runtime.CompiledComponentAndPorts, len(compList))
	for graphNodeIndex, compiledComp := range compList {
		mapStruct := compiledComp.MapStruct
		log.Trace().Msgf("mapStruct: %+v", mapStruct)

		compWithPorts := runtime.CompiledComponentAndPorts{
			InPortToSignalsMap:  make(runtime.PortToSignal),
			OutPortToSignalsMap: make(runtime.PortToSignal),
			CompiledComponent:   compiledComp,
		}

		// Read in_ports in mapStruct
		inPorts, ok := mapStruct["in_ports"]
		log.Trace().Interface("inPorts", inPorts).Bool("ok", ok).Str("componentName", compiledComp.Name).Msg("mapStruct[in_ports]")
		if ok {
			// Convert in_ports to map[string]interface{}
			inPortsMap, castOk := inPorts.(map[string]interface{})
			if castOk {
				inPortToSignalsMap, err := getInPortSignals(inPortsMap, inSignals, graphNodeIndex)
				if err != nil {
					return nil, nil, err
				}
				log.Trace().Msgf("inPortToSignalsMap: %+v", inPortToSignalsMap)
				compWithPorts.InPortToSignalsMap = inPortToSignalsMap
			}
		}
		// Read out_ports in mapStruct
		outPorts, ok := mapStruct["out_ports"]
		log.Trace().Interface("outPorts", outPorts).Bool("ok", ok).Str("componentName", compiledComp.Name).Msg("mapStruct[out_ports]")
		if ok {
			// Convert out_ports to map[string]interface{}
			outPortsMap, castOk := outPorts.(map[string]interface{})
			if castOk {
				outPortToSignalsMap, err := getOutPortSignals(outPortsMap, outSignals, graphNodeIndex)
				if err != nil {
					return nil, nil, err
				}
				log.Trace().Msgf("inPortToSignalsMap: %+v", outPortToSignalsMap)
				compWithPorts.OutPortToSignalsMap = outPortToSignalsMap
			}
		}

		// Append compWithPorts to compWithPortsList
		compWithPortsList[graphNodeIndex] = compWithPorts
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

	for graphNodeIndex, compWithPorts := range compWithPortsList {
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
	log.Trace().Msgf("Tarjan Loops: %+v \nTarjan Graph: %+v", loops, graph)

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
			if removeToCompIndex >= len(compWithPortsList) {
				return nil, fx.Options(), errors.New("removeToCompId is out of range")
			}
			removeToComp := compWithPortsList[removeToCompIndex]
			if removeFromCompIndex >= len(compWithPortsList) {
				return nil, fx.Options(), errors.New("removeFromCompId is out of range")
			}
			removeFromComp := compWithPortsList[removeFromCompIndex]
			loopedSignals := make(map[string]bool)
			// Mark looped signals in InPortToSignalsMap
			for _, signals := range removeToComp.InPortToSignalsMap {
				for idx, signal := range signals {
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

	log.Trace().Msgf("comp with ports list: %+v", compWithPortsList)

	// Log compWithPortsList
	for _, compWithPorts := range compWithPortsList {
		log.Trace().Msgf("comp with ports: %+v", compWithPorts)
	}

	return compWithPortsList, fx.Options(componentOptions...), nil
}

type signalType int

const (
	inSignalType signalType = iota
	outSignalType
)

func getInPortSignals(portMapping map[string]interface{}, inSignals map[string][]int, graphNodeIndex int) (runtime.PortToSignal, error) {
	return getPortSignals(portMapping, inSignals, nil, inSignalType, graphNodeIndex)
}

func getOutPortSignals(portMapping map[string]interface{}, outSignals map[string]int, graphNodeIndex int) (runtime.PortToSignal, error) {
	return getPortSignals(portMapping, nil, outSignals, outSignalType, graphNodeIndex)
}

// getPortSignals takes a port mapping and returns a PortToSignal map and signals list.
func getPortSignals(portMapping map[string]interface{}, inSignals map[string][]int, outSignals map[string]int, sigType signalType, graphNodeIndex int) (runtime.PortToSignal, error) {
	// fillSignal takes a portSpec and fills portSignals at idx with Signal.
	// Returns the signal name from the port spec and a bool indicating if a valid signal name was found in the portSpec.
	fillSignal := func(portSpec map[string]interface{}, portSignals []runtime.Signal, idx int) (string, bool) {
		// Read signal_name
		signalName, ok := portSpec["signal_name"]
		if ok {
			signalNameStr, ok := signalName.(string)
			if ok {
				// Fill portSignals
				portSignals[idx] = runtime.Signal{
					Name:   signalNameStr,
					Looped: false,
				}
				// Return signalNameStr and true since signal_name is present
				return signalNameStr, true
			}
		}
		// Return empty string and false since signal_name is not present
		return "", false
	}

	fillInOutSignals := func(signalName string, inSignals map[string][]int, outSignals map[string]int, sigType signalType, graphNodeIndex int) error {
		if sigType == inSignalType {
			// Add signals from inSignalsList to inSignals
			inSignals[signalName] = append(inSignals[signalName], graphNodeIndex)
		} else if sigType == outSignalType {
			// Check if signal is already present in outSignals
			if _, ok := outSignals[signalName]; !ok {
				outSignals[signalName] = graphNodeIndex
			} else {
				return errors.New("duplicate signal definition for signal name: " + signalName)
			}
		}
		return nil
	}

	portToSignalMapping := make(runtime.PortToSignal)

	log.Trace().Msgf("portMapping: %+v", portMapping)

	// Iterate each port
	for port, portMap := range portMapping {
		// Convert portMap to map[string][]interface{}
		portList, isList := portMap.([]interface{})
		log.Trace().Bool("isList", isList).Msgf("portMap: %+v", portMap)
		// Convert portMap to map[string]interface{}
		portSpec, isSpec := portMap.(map[string]interface{})
		if isList {
			// Initialize portMapping for this port
			portSignals := make([]runtime.Signal, len(portList))
			portToSignalMapping[port] = portSignals
			// iterate each port spec
			for idx, innerPortSpec := range portList {
				innerPortSpec, isMapStruct := innerPortSpec.(map[string]interface{})
				if isMapStruct {
					signalName, filled := fillSignal(innerPortSpec, portSignals, idx)
					if filled {
						err := fillInOutSignals(signalName, inSignals, outSignals, sigType, graphNodeIndex)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		} else if isSpec {
			// Initialize portMapping for this port
			portSignals := make([]runtime.Signal, 1)
			portToSignalMapping[port] = portSignals
			signalName, filled := fillSignal(portSpec, portSignals, 0)
			if filled {
				err := fillInOutSignals(signalName, inSignals, outSignals, sigType, graphNodeIndex)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return portToSignalMapping, nil
}
