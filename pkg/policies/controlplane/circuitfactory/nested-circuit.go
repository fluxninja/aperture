package circuitfactory

import (
	"encoding/json"
	"strings"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/mapstruct"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

// NestedCircuitDelimiter is the delimiter used to separate the parent circuit ID and the nested circuit ID.
const NestedCircuitDelimiter = "."

// ParseNestedCircuit parses a nested circuit and returns the parent, leaf components, and options.
func ParseNestedCircuit(
	nestedCircuitID string,
	nestedCircuit *policylangv1.NestedCircuit,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	// serialize to jsonBytes
	jsonBytes, err := json.Marshal(nestedCircuit)
	if err != nil {
		return nil, nil, nil, err
	}
	// unmarshal using our config layer to make sure defaults and validates happen
	nestedCircuitProto := &policylangv1.NestedCircuit{}
	err = config.UnmarshalJSON(jsonBytes, nestedCircuitProto)
	if err != nil {
		return nil, nil, nil, err
	}

	portMapping := runtime.NewPortMapping()
	parentCircuitID := ParentCircuitID(nestedCircuitID)

	inPortsMap := nestedCircuitProto.GetInPortsMap()
	ins, err := DecodePortMap(inPortsMap, parentCircuitID)
	if err != nil {
		return nil, nil, nil, err
	}

	outPortsMap := nestedCircuitProto.GetOutPortsMap()
	outs, err := DecodePortMap(outPortsMap, parentCircuitID)
	if err != nil {
		return nil, nil, nil, err
	}

	portMapping.Ins = ins
	portMapping.Outs = outs

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
	for i, configuredComponent := range leafComponents {
		component := configuredComponent.Component
		// dynamic cast to signal ingress or egress
		if nestedSignalIngress, ok := component.(*components.NestedSignalIngress); ok {
			portName := nestedSignalIngress.PortName()
			// tracking the port names in the nested circuit
			if _, ok := ingressPorts[portName]; ok {
				return nil, nil, nil, errors.Errorf("duplicate ingress port %s in nested circuit %s", portName, nestedCircuitProto.Name)
			}
			ingressPorts[portName] = nil
			signals, ok := portMapping.GetInPort(portName)
			if ok {
				// set the port mapping for the signal ingress component
				leafComponents[i].PortMapping.AddInPort(components.NestedSignalPortName, signals)
			}
		} else if nestedSignalEgress, ok := component.(*components.NestedSignalEgress); ok {
			portName := nestedSignalEgress.PortName()
			// tracking the port names in the nested circuit
			if _, ok := egressPorts[portName]; ok {
				return nil, nil, nil, errors.Errorf("duplicate egress port %s in nested circuit %s", portName, nestedCircuitProto.Name)
			}
			egressPorts[portName] = nil
			signals, ok := portMapping.GetOutPort(portName)
			if ok {
				// set the port mapping for the signal egress component
				leafComponents[i].PortMapping.AddOutPort(components.NestedSignalPortName, signals)
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

// DecodePortMap decodes a proto port map into a PortToSignals map.
func DecodePortMap(config any, circuitID string) (runtime.PortToSignals, error) {
	ports := make(runtime.PortToSignals)

	mapStruct, err := mapstruct.EncodeObject(config)
	if err != nil {
		return nil, err
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true, // So that singular ports will transparently be converted to lists.
		Result:           &ports,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(mapStruct)
	if err != nil {
		return nil, err
	}

	for _, signals := range ports {
		for i := range signals {
			signals[i].SubCircuitID = circuitID
		}
	}

	return ports, nil
}
