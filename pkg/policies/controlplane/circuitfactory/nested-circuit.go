package circuitfactory

import (
	"strings"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// NestedCircuitDelimiter is the delimiter used to separate the parent circuit ID and the nested circuit ID.
const NestedCircuitDelimiter = "."

// ParseNestedCircuit parses a nested circuit and returns the parent, leaf components, and options.
func ParseNestedCircuit(
	nestedCircuitID string,
	nestedCircuitProto *policylangv1.NestedCircuit,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	portMapping := runtime.NewPortMapping()
	parentCircuitID := ParentCircuitID(nestedCircuitID)

	inPortsMap := nestedCircuitProto.GetInPortsMap()
	for portName, inPort := range inPortsMap {
		signals := []runtime.Signal{
			{
				SignalName:    inPort.GetSignalName(),
				ConstantValue: inPort.GetConstantValue(),
				CircuitID:     parentCircuitID,
			},
		}
		if portMapping.ExistsInPort(portName) {
			return nil, nil, nil, errors.Errorf("redefinition of port %s in nested circuit %s", portName, nestedCircuitProto.Name)
		}
		portMapping.AddInPort(portName, signals)
	}

	outPortsMap := nestedCircuitProto.GetOutPortsMap()
	for portName, outPort := range outPortsMap {
		signals := []runtime.Signal{
			{
				SignalName: outPort.GetSignalName(),
				CircuitID:  parentCircuitID,
			},
		}
		if portMapping.ExistsOutPort(portName) {
			return nil, nil, nil, errors.Errorf("redefinition of port %s in nested circuit %s", portName, nestedCircuitProto.Name)
		}
		portMapping.AddOutPort(portName, signals)
	}

	parentComponents, leafComponents, options, err := CreateComponents(
		nestedCircuitProto.GetComponents(),
		nestedCircuitID,
		policyReadAPI,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	// For tracking the ingress/egress port names in the nested circuit
	ingressPorts := make(map[string]interface{})
	egressPorts := make(map[string]interface{})

	// Set in and out ports at signal ingress and egress components based on the port mapping
	for _, configuredComponent := range leafComponents {
		component := configuredComponent.Component
		// dynamic cast to signal ingress or egress
		if nestedSignalIngress, ok := component.(*components.NestedSignalIngress); ok {
			portName := nestedSignalIngress.PortName()
			// for tracking the port names in the nested circuit
			ingressPorts[portName] = nil
			signals, ok := portMapping.Ins[portName]
			if ok {
				// set the port mapping for the signal ingress component
				configuredComponent.PortMapping.Ins[portName] = signals
			}
		} else if nestedSignalEgress, ok := component.(*components.NestedSignalEgress); ok {
			portName := nestedSignalEgress.PortName()
			// for tracking the port names in the nested circuit
			egressPorts[portName] = nil
			signals, ok := portMapping.Outs[portName]
			if ok {
				// set the port mapping for the signal egress component
				configuredComponent.PortMapping.Outs[portName] = signals
			}
		}

	}

	for portName := range portMapping.Ins {
		if _, ok := ingressPorts[portName]; !ok {
			return nil, nil, nil, errors.Errorf("port %s not found in nested circuit %s", portName, nestedCircuitProto.Name)
		}
	}

	for portName := range portMapping.Outs {
		if _, ok := egressPorts[portName]; !ok {
			return nil, nil, nil, errors.Errorf("port %s not found in nested circuit %s", portName, nestedCircuitProto.Name)
		}
	}

	nestedCircConfComp, err := prepareComponent(
		runtime.NewDummyComponent(nestedCircuitProto.Name, runtime.ComponentTypeSignalProcessor),
		nestedCircuitProto,
		nestedCircuitID,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	nestedCircConfComp.PortMapping = portMapping
	parentComponents = append(parentComponents, nestedCircConfComp)

	return parentComponents, leafComponents, options, err
}

// ParentCircuitID returns the parent circuit ID of the given child circuit ID assuming nested circuits are delimited by dots (".").
func ParentCircuitID(childCircuitID string) string {
	// Parent Child are delimited by dots. So, we split the child circuit ID by dots and return the first part.
	// For example, if the child circuit ID is "foo.bar.baz", then the parent circuit ID is "foo.bar".
	return childCircuitID[:strings.LastIndex(childCircuitID, NestedCircuitDelimiter)]
}
