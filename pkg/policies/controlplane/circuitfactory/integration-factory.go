package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/concurrency"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/horizontalpodscaler"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// compositeComponentFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func compositeComponentFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		horizontalpodscaler.Module(),
	)
}

// newIntegrationCompositeAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newIntegrationCompositeAndOptions(
	compositeComponentProto *policylangv1.Integration,
	compositeComponentID string,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	parentCircuitID := ParentCircuitID(compositeComponentID)
	// Factory parser to determine what kind of composite component to create
	if concurrencyLimiterProto := compositeComponentProto.GetConcurrencyLimiter(); concurrencyLimiterProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiterProto, compositeComponentID, policyReadAPI)
		if concurrencyLimiterErr != nil {
			return nil, nil, nil, concurrencyLimiterErr
		}
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiterProto.GetScheduler(); schedulerProto != nil {
			// Use the same id as the component stack since agent sees only the component stack and generates metrics tagged with the component stack id
			scheduler, schedulerOptions, err := concurrency.NewSchedulerAndOptions(schedulerProto, compositeComponentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			// Need a unique ID for sub component since it's used for graph generation
			schedulerConfComp, err := prepareComponentInCircuit(scheduler, schedulerProto, compositeComponentID+".Scheduler", parentCircuitID)
			if err != nil {
				return nil, nil, nil, err
			}

			configuredComponents = append(configuredComponents, schedulerConfComp)

			options = append(options, schedulerOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(schedulerConfComp.PortMapping)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		// Actuation Strategy
		if loadActuatorProto := concurrencyLimiterProto.GetLoadActuator(); loadActuatorProto != nil {
			loadActuator, loadActuatorOptions, err := concurrency.NewLoadActuatorAndOptions(loadActuatorProto, compositeComponentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			loadActuatorConfComp, err := prepareComponentInCircuit(loadActuator, loadActuatorProto, compositeComponentID+".LoadActuator", parentCircuitID)
			if err != nil {
				return nil, nil, nil, err
			}
			configuredComponents = append(configuredComponents, loadActuatorConfComp)

			options = append(options, loadActuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(loadActuatorConfComp.PortMapping)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		concurrencyLimiterConfComp, err := prepareComponent(
			runtime.NewDummyComponent("ConcurrencyLimiter", runtime.ComponentTypeSignalProcessor),
			concurrencyLimiterProto,
			compositeComponentID,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		concurrencyLimiterConfComp.PortMapping = portMapping

		return []runtime.ConfiguredComponent{concurrencyLimiterConfComp}, configuredComponents, fx.Options(options...), nil
	} else if horizontalPodScalerProto := compositeComponentProto.GetHorizontalPodScaler(); horizontalPodScalerProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		horizontalPodScalerOptions, agentGroupName, horizontalPodScalerErr := horizontalpodscaler.NewHorizontalPodScalerOptions(horizontalPodScalerProto, compositeComponentID, policyReadAPI)
		if horizontalPodScalerErr != nil {
			return nil, nil, nil, horizontalPodScalerErr
		}
		options = append(options, horizontalPodScalerOptions)

		// Scale Reporter
		if scaleReporterProto := horizontalPodScalerProto.GetScaleReporter(); scaleReporterProto != nil {
			scaleReporter, scaleReporterOptions, err := horizontalpodscaler.NewScaleReporterAndOptions(scaleReporterProto, compositeComponentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			scaleReporterConfComp, err := prepareComponentInCircuit(scaleReporter, scaleReporterProto, compositeComponentID+".ScaleReporter", parentCircuitID)
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
			scaleActuator, scaleActuatorOptions, err := horizontalpodscaler.NewScaleActuatorAndOptions(scaleActuatorProto, compositeComponentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			scaleActuatorConfComp, err := prepareComponentInCircuit(scaleActuator, scaleActuatorProto, compositeComponentID+".ScaleActuator", parentCircuitID)
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
			compositeComponentID,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		horizontalPodScalerConfComp.PortMapping = portMapping

		return []runtime.ConfiguredComponent{horizontalPodScalerConfComp}, configuredComponents, fx.Options(options...), nil
	} else if aimdConcurrencyController := compositeComponentProto.GetAimdConcurrencyController(); aimdConcurrencyController != nil {
		return ParseAIMDConcurrencyController(compositeComponentID, aimdConcurrencyController, policyReadAPI)
	}
	return nil, nil, nil, fmt.Errorf("unsupported/missing component type, proto: %+v", compositeComponentProto)
}
