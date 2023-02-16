package circuitfactory

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/podscaler"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
)

// autoScaleModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func autoScaleModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		podscaler.Module(),
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

	if podScalerProto := autoScaleComponentProto.GetPodScaler(); podScalerProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			tree                 Tree
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		podScalerOptions, agentGroupName, podScalerErr := podscaler.NewPodScalerOptions(podScalerProto, componentID.String(), policyReadAPI)
		if podScalerErr != nil {
			return retErr(podScalerErr)
		}
		options = append(options, podScalerOptions)

		// Scale Reporter
		if scaleReporterProto := podScalerProto.GetScaleReporter(); scaleReporterProto != nil {
			scaleReporter, scaleReporterOptions, err := podscaler.NewScaleReporterAndOptions(scaleReporterProto, componentID.String(), policyReadAPI, agentGroupName)
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
		if scaleActuatorProto := podScalerProto.GetScaleActuator(); scaleActuatorProto != nil {
			scaleActuator, scaleActuatorOptions, err := podscaler.NewScaleActuatorAndOptions(scaleActuatorProto, componentID.String(), policyReadAPI, agentGroupName)
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

		kos := podScalerProto.KubernetesObjectSelector
		sd := fmt.Sprintf("%s/%s/%s/%s/%s",
			kos.GetAgentGroup(),
			kos.GetNamespace(),
			kos.GetApiVersion(),
			kos.GetKind(),
			kos.GetName(),
		)

		podScalerConfComp, err := prepareComponent(
			runtime.NewDummyComponent("PodScaler", sd, runtime.ComponentTypeSignalProcessor),
			podScalerProto,
			componentID,
		)
		if err != nil {
			return retErr(err)
		}

		podScalerConfComp.PortMapping = portMapping
		tree.Root = podScalerConfComp

		return tree, configuredComponents, fx.Options(options...), nil
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", autoScaleComponentProto))
}
