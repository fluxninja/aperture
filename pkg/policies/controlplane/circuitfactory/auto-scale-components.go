package circuitfactory

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/horizontalpodscaler"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
)

// autoScaleModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func autoScaleModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		horizontalpodscaler.Module(),
	)
}

// newIntegrationCompositeAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newAutoScaleCompositeAndOptions(
	autoScaleComponentProto *policylangv1.AutoScale,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
	retErr := func(err error) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
		return Tree{}, nil, nil, err
	}

	parentCircuitID, ok := componentID.ParentID()
	if !ok {
		return retErr(fmt.Errorf("parent circuit ID not found for component %s", componentID))
	}

	if horizontalPodScalerProto := autoScaleComponentProto.GetHorizontalPodScaler(); horizontalPodScalerProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			tree                 Tree
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		horizontalPodScalerOptions, agentGroupName, horizontalPodScalerErr := horizontalpodscaler.NewHorizontalPodScalerOptions(horizontalPodScalerProto, componentID.String(), policyReadAPI)
		if horizontalPodScalerErr != nil {
			return retErr(horizontalPodScalerErr)
		}
		options = append(options, horizontalPodScalerOptions)

		// Scale Reporter
		if scaleReporterProto := horizontalPodScalerProto.GetScaleReporter(); scaleReporterProto != nil {
			scaleReporter, scaleReporterOptions, err := horizontalpodscaler.NewScaleReporterAndOptions(scaleReporterProto, componentID.String(), policyReadAPI, agentGroupName)
			if err != nil {
				return retErr(err)
			}

			scaleReporterConfComp, err := prepareComponentInCircuit(scaleReporter, scaleReporterProto, componentID.ChildID("ScaleReporter"), parentCircuitID)
			if err != nil {
				return retErr(err)
			}

			configuredComponents = append(configuredComponents, scaleReporterConfComp)
			tree.Children = append(tree.Children, Tree{Root: scaleReporterConfComp})

			options = append(options, scaleReporterOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(scaleReporterConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		// Scale Actuator
		if scaleActuatorProto := horizontalPodScalerProto.GetScaleActuator(); scaleActuatorProto != nil {
			scaleActuator, scaleActuatorOptions, err := horizontalpodscaler.NewScaleActuatorAndOptions(scaleActuatorProto, componentID.String(), policyReadAPI, agentGroupName)
			if err != nil {
				return retErr(err)
			}

			scaleActuatorConfComp, err := prepareComponentInCircuit(scaleActuator, scaleActuatorProto, componentID.ChildID("ScaleActuator"), parentCircuitID)
			if err != nil {
				return retErr(err)
			}
			configuredComponents = append(configuredComponents, scaleActuatorConfComp)
			tree.Children = append(tree.Children, Tree{Root: scaleActuatorConfComp})

			options = append(options, scaleActuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(scaleActuatorConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		kos := horizontalPodScalerProto.KubernetesObjectSelector
		sd := fmt.Sprintf("%s/%s/%s/%s/%s",
			kos.GetAgentGroup(),
			kos.GetNamespace(),
			kos.GetApiVersion(),
			kos.GetKind(),
			kos.GetName(),
		)

		horizontalPodScalerConfComp, err := prepareComponent(
			runtime.NewDummyComponent("HorizontalPodScaler", sd, runtime.ComponentTypeSignalProcessor),
			horizontalPodScalerProto,
			componentID,
		)
		if err != nil {
			return retErr(err)
		}

		horizontalPodScalerConfComp.PortMapping = portMapping
		tree.Root = horizontalPodScalerConfComp

		return tree, configuredComponents, fx.Options(options...), nil
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", autoScaleComponentProto))
}
