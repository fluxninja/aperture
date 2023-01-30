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
	componentID string,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	parentCircuitID := ParentCircuitID(componentID)

	if horizontalPodScalerProto := autoScaleComponentProto.GetHorizontalPodScaler(); horizontalPodScalerProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		horizontalPodScalerOptions, agentGroupName, horizontalPodScalerErr := horizontalpodscaler.NewHorizontalPodScalerOptions(horizontalPodScalerProto, componentID, policyReadAPI)
		if horizontalPodScalerErr != nil {
			return nil, nil, nil, horizontalPodScalerErr
		}
		options = append(options, horizontalPodScalerOptions)

		// Scale Reporter
		if scaleReporterProto := horizontalPodScalerProto.GetScaleReporter(); scaleReporterProto != nil {
			scaleReporter, scaleReporterOptions, err := horizontalpodscaler.NewScaleReporterAndOptions(scaleReporterProto, componentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			scaleReporterConfComp, err := prepareComponentInCircuit(scaleReporter, scaleReporterProto, componentID+".ScaleReporter", parentCircuitID)
			if err != nil {
				return nil, nil, nil, err
			}

			configuredComponents = append(configuredComponents, scaleReporterConfComp)

			options = append(options, scaleReporterOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(scaleReporterConfComp.PortMapping)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		// Scale Actuator
		if scaleActuatorProto := horizontalPodScalerProto.GetScaleActuator(); scaleActuatorProto != nil {
			scaleActuator, scaleActuatorOptions, err := horizontalpodscaler.NewScaleActuatorAndOptions(scaleActuatorProto, componentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			scaleActuatorConfComp, err := prepareComponentInCircuit(scaleActuator, scaleActuatorProto, componentID+".ScaleActuator", parentCircuitID)
			if err != nil {
				return nil, nil, nil, err
			}
			configuredComponents = append(configuredComponents, scaleActuatorConfComp)

			options = append(options, scaleActuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(scaleActuatorConfComp.PortMapping)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		horizontalPodScalerConfComp, err := prepareComponent(
			runtime.NewDummyComponent("HorizontalPodScaler", runtime.ComponentTypeSignalProcessor),
			horizontalPodScalerProto,
			componentID,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		horizontalPodScalerConfComp.PortMapping = portMapping

		return []runtime.ConfiguredComponent{horizontalPodScalerConfComp}, configuredComponents, fx.Options(options...), nil
	}
	return nil, nil, nil, fmt.Errorf("unsupported/missing component type, proto: %+v", autoScaleComponentProto)
}
