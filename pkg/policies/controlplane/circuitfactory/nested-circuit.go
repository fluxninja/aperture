package circuitfactory

import (
	"encoding/json"

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

// ParseNestedCircuit parses a nested circuit and returns the parent, leaf components, and options.
func ParseNestedCircuit(
	nestedCircuitID runtime.ComponentID,
	nestedCircuit *policylangv1.NestedCircuit,
	policyReadAPI iface.Policy,
) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
	retErr := func(err error) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
		return Tree{}, nil, nil, err
	}

	// serialize to jsonBytes
	jsonBytes, err := json.Marshal(nestedCircuit)
	if err != nil {
		return retErr(err)
	}
	// unmarshal using our config layer to make sure defaults and validates happen
	nestedCircuitProto := &policylangv1.NestedCircuit{}
	err = config.UnmarshalJSON(jsonBytes, nestedCircuitProto)
	if err != nil {
		return retErr(err)
	}

	portMapping := runtime.NewPortMapping()
	parentCircuitID, ok := nestedCircuitID.ParentID()
	if !ok {
		return retErr(errors.Errorf("nested circuit %s does not have a parent circuit", nestedCircuitID))
	}

	inPortsMap := nestedCircuitProto.GetInPortsMap()
	ins, err := DecodePortMap(inPortsMap, parentCircuitID.String())
	if err != nil {
		return retErr(err)
	}

	outPortsMap := nestedCircuitProto.GetOutPortsMap()
	outs, err := DecodePortMap(outPortsMap, parentCircuitID.String())
	if err != nil {
		return retErr(err)
	}

	portMapping.Ins = ins
	portMapping.Outs = outs

	tree, leafComponents, options, err := CreateComponents(
		nestedCircuitProto.GetComponents(),
		nestedCircuitID,
		policyReadAPI,
	)
	if err != nil {
		return retErr(err)
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
				return retErr(errors.Errorf("duplicate ingress port %s in nested circuit %s", portName, nestedCircuitProto.Name))
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
				return retErr(errors.Errorf("duplicate egress port %s in nested circuit %s", portName, nestedCircuitProto.Name))
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
			return retErr(errors.Errorf("port %s not found in nested circuit %s", portName, nestedCircuitProto.Name))
		}
	}

	for portName := range portMapping.Outs {
		if _, ok := egressPorts[portName]; !ok {
			return retErr(errors.Errorf("port %s not found in nested circuit %s", portName, nestedCircuitProto.Name))
		}
	}

	nestedCircConfComp, err := prepareComponent(
		runtime.NewDummyComponent(nestedCircuitProto.Name,
			nestedCircuitProto.ShortDescription,
			runtime.ComponentTypeSignalProcessor),
		nestedCircuitProto,
		nestedCircuitID,
	)
	if err != nil {
		return retErr(err)
	}
	nestedCircConfComp.PortMapping = portMapping
	tree.Root = nestedCircConfComp

	return tree, leafComponents, options, err
}

// DecodePortMap decodes a proto port map into a PortToSignals map.
func DecodePortMap(config any, circuitID string) (runtime.PortToSignals, error) {
	ports := make(runtime.PortToSignals)
	if config == nil {
		return ports, nil
	}

	mapStruct, err := mapstruct.EncodeObject(config)
	if err != nil {
		return nil, err
	}
	if len(mapStruct) == 0 {
		return ports, nil
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
