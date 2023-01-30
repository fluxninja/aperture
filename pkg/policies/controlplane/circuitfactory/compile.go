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
	Tree           Tree
	LeafComponents []runtime.ConfiguredComponent
}

// Tree is a graph view of a Circuit.
type Tree struct {
	Root     runtime.ConfiguredComponent
	Children []Tree
}

// Components returns a list of CompiledComponents, ready to create runtime.Circuit.
func (circuit *Circuit) Components() []runtime.ConfiguredComponent { return circuit.LeafComponents }

// CompileFromProto compiles a protobuf circuit definition into a Circuit.
//
// This is helper for CreateComponents + runtime.Compile.
func CompileFromProto(
	componentsProto []*policylangv1.Component,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	tree, leafComponents, option, err := CreateComponents(componentsProto, runtime.NewComponentID("root"), policyReadAPI)
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
		Tree:           tree,
		LeafComponents: leafComponents,
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
) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
	var (
		leafComponents []runtime.ConfiguredComponent
		tree           Tree
		options        []fx.Option
	)

	for compIndex, componentProto := range componentsProto {
		subTree, leafComps, compOption, err := NewComponentAndOptions(
			componentProto,
			circuitID.ChildID(strconv.Itoa(compIndex)),
			policyReadAPI,
		)
		if err != nil {
			return Tree{}, nil, nil, err
		}
		options = append(options, compOption)

		// Append subTree to tree.Children
		tree.Children = append(tree.Children, subTree)

		// Add subComponents to configuredComponents
		leafComponents = append(leafComponents, leafComps...)
	}

	return tree, leafComponents, fx.Options(options...), nil
}
