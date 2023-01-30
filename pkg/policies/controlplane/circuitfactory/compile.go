package circuitfactory

import (
	"strconv"

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
	leafComponents   []runtime.ConfiguredComponent
	parentComponents []runtime.ConfiguredComponent
}

// Components returns a list of CompiledComponents, ready to create runtime.Circuit.
func (circuit *Circuit) Components() []runtime.ConfiguredComponent { return circuit.leafComponents }

// CompileFromProto compiles a protobuf circuit definition into a Circuit.
//
// This is helper for CreateComponents + runtime.Compile.
func CompileFromProto(
	componentsProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	parentComponents, leafComponents, option, err := CreateComponents(componentsProto, runtime.NewComponentID("root"), policyReadAPI)
	if err != nil {
		return nil, nil, err
	}

	err = runtime.Compile(
		leafComponents,
		policyReadAPI.GetStatusRegistry().GetLogger(),
	)
	if err != nil {
		return nil, option, err
	}

	return &Circuit{
		parentComponents: parentComponents,
		leafComponents:   leafComponents,
	}, option, nil
}

// CreateComponents creates circuit components along with their identifiers and fx options.
//
// Note that number of returned components might be greater than number of
// components in componentsProto, as some components may be composite multi-component stacks or nested circuits.
func CreateComponents(
	componentsProto []*policylangv1.Component,
	circuitID runtime.ComponentID,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	var (
		leafComponents   []runtime.ConfiguredComponent
		parentComponents []runtime.ConfiguredComponent
		options          []fx.Option
	)

	for compIndex, componentProto := range componentsProto {
		parentComps, leafComps, compOption, err := NewComponentAndOptions(
			componentProto,
			circuitID.ChildID(strconv.Itoa(compIndex)),
			policyReadAPI,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		options = append(options, compOption)

		// Add outerComponent to outerComponents
		parentComponents = append(parentComponents, parentComps...)

		// Add subComponents to configuredComponents
		leafComponents = append(leafComponents, leafComps...)
	}

	return parentComponents, leafComponents, fx.Options(options...), nil
}
