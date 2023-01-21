package circuitfactory

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/mapstruct"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Module for circuit and component factory run via the main app.
func Module() fx.Option {
	return fx.Options(
		runtime.CircuitModule(),
		FactoryModule(),
	)
}

// GraphNode is a node in the circuit graph. It is used for representing the outermost components in the circuit.
type GraphNode struct {
	// Which signals this component wants to have connected on its ports.
	PortMapping runtime.PortMapping
	// Mapstruct representation of proto config that was used to create this
	// component.  This Config is used only for observability purposes.
	//
	// Note: PortMapping is also part of Config.
	Config mapstruct.Object
}

// Circuit is a compiled Circuit
//
// Circuit can also be converted to its graph view.
type Circuit struct {
	components []runtime.ConfiguredComponent
	graphNodes []GraphNode
}

// Components returns a list of CompiledComponents, ready to create runtime.Circuit.
func (circuit *Circuit) Components() []runtime.ConfiguredComponent { return circuit.components }

// CompileFromProto compiles a protobuf circuit definition into a Circuit.
//
// This is helper for CreateComponents + runtime.Compile.
func CompileFromProto(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	configuredComponents, graphNodes, option, err := CreateComponents(circuitProto, policyReadAPI)
	if err != nil {
		return nil, nil, err
	}

	err = runtime.Compile(
		configuredComponents,
		policyReadAPI.GetStatusRegistry().GetLogger(),
	)
	if err != nil {
		return nil, option, err
	}

	return &Circuit{
		components: configuredComponents,
		graphNodes: graphNodes,
	}, option, nil
}

// CreateComponents creates circuit components along with their identifiers and fx options.
//
// Note that number of returned components might be greater than number of
// components in circuitProto, as some components may have their subcomponents.
func CreateComponents(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []GraphNode, fx.Option, error) {
	var (
		configuredComponents []runtime.ConfiguredComponent
		graphNodes           []GraphNode
		options              []fx.Option
	)

	for compIndex, componentProto := range circuitProto {
		// Create graphNode
		graphNode, subComponents, compOption, err := NewComponentAndOptions(
			componentProto,
			compIndex,
			policyReadAPI,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		options = append(options, compOption)

		// Add graphNode to graphNodes
		graphNodes = append(graphNodes, graphNode)

		// Add subComponents to configuredComponents
		configuredComponents = append(configuredComponents, subComponents...)
	}

	return configuredComponents, graphNodes, fx.Options(options...), nil
}
