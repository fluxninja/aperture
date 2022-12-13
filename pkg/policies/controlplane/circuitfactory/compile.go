package circuitfactory

import (
	"strconv"
	"strings"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
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
// This is helper for CreateComponents + runtime.Compile.
func CompileFromProto(
	circuitProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	configuredComponents, componentIDs, option, err := CreateComponents(circuitProto, policyReadAPI)
	if err != nil {
		return nil, nil, err
	}

	compiledComponents, err := runtime.Compile(
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
	parentID, _, hasParent := strings.Cut(string(id), ".")
	if !hasParent {
		return ComponentID("")
	}
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
		component, subcomponents, compOption, err := NewComponentAndOptions(
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
