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
	components      []runtime.ConfiguredComponent
	outerComponents []runtime.ConfiguredComponent
}

// Components returns a list of CompiledComponents, ready to create runtime.Circuit.
func (circuit *Circuit) Components() []runtime.ConfiguredComponent { return circuit.components }

// CompileFromProto compiles a protobuf circuit definition into a Circuit.
//
// This is helper for CreateComponents + runtime.Compile.
func CompileFromProto(
	componentsProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	configuredComponents, outerComponents, option, err := CreateComponents(componentsProto, "root", policyReadAPI)
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
		components:      configuredComponents,
		outerComponents: outerComponents,
	}, option, nil
}

// CreateComponents creates circuit components along with their identifiers and fx options.
//
// Note that number of returned components might be greater than number of
// components in componentsProto, as some components may have their subcomponents.
func CreateComponents(
	componentsProto []*policylangv1.Component,
	circuitID string,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	var (
		configuredComponents []runtime.ConfiguredComponent
		outerComponents      []runtime.ConfiguredComponent
		options              []fx.Option
	)

	for compIndex, componentProto := range componentsProto {
		outerComponent, subComponents, compOption, err := NewComponentAndOptions(
			componentProto,
			ComponentID(circuitID, compIndex),
			policyReadAPI,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		options = append(options, compOption)

		// Add outerComponent to outerComponents
		outerComponents = append(outerComponents, outerComponent)

		// Add subComponents to configuredComponents
		configuredComponents = append(configuredComponents, subComponents...)
	}

	return configuredComponents, outerComponents, fx.Options(options...), nil
}

// ComponentID returns the ID for a component within a circuit at componentIndex.
func ComponentID(circuitID string, componentIndex int) string {
	return circuitID + "." + strconv.Itoa(componentIndex)
}
