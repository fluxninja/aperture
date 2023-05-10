package circuitfactory

import (
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policymonitoringv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/monitoring/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
)

// Circuit is a compiled Circuit
//
// Circuit can also be converted to its graph view.
type Circuit struct {
	Tree           Tree
	LeafComponents []*runtime.ConfiguredComponent
}

// Components returns a list of CompiledComponents, ready to create runtime.Circuit.
func (circuit *Circuit) Components() []*runtime.ConfiguredComponent { return circuit.LeafComponents }

// CompileFromProto compiles a protobuf circuit definition into a Circuit.
//
// This is helper for CreateComponents + runtime.Compile.
func CompileFromProto(
	policyProto *policylangv1.Policy,
	policyReadAPI iface.Policy,
) (*Circuit, fx.Option, error) {
	componentsProto := policyProto.GetCircuit().Components
	tree, err := RootTree(policyProto.GetCircuit())
	if err != nil {
		return nil, nil, err
	}
	leafComponents, option, err := tree.CreateComponents(componentsProto, runtime.NewComponentID(runtime.RootComponentID), policyReadAPI)
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

// CircuitView returns a CircuitView of the circuit.
func (circuit *Circuit) CircuitView() (*policymonitoringv1.CircuitView, error) {
	tree, err := circuit.Tree.TreeGraph()
	if err != nil {
		log.Errorf("Error getting circuit tree: %v", err)
		return nil, err
	}
	return &policymonitoringv1.CircuitView{
		Tree: tree,
	}, nil
}
