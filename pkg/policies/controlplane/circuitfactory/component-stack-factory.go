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

// componentStackFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func componentStackFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		horizontalpodscaler.Module(),
	)
}

// newComponentStackAndOptions creates components for component stack, sub components and their fx options.
func newComponentStackAndOptions(
	componentStackProto *policylangv1.Component,
	componentStackIndex int,
	policyReadAPI iface.Policy,
) (GraphNode, []runtime.ConfiguredComponent, fx.Option, error) {
	// Factory parser to determine what kind of component stack to create
	if concurrencyLimiterProto := componentStackProto.GetConcurrencyLimiter(); concurrencyLimiterProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiterProto, componentStackIndex, policyReadAPI)
		if concurrencyLimiterErr != nil {
			return GraphNode{}, nil, nil, concurrencyLimiterErr
		}
		// Append concurrencyLimiter options
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiterProto.GetScheduler(); schedulerProto != nil {
			scheduler, schedulerOptions, err := concurrency.NewSchedulerAndOptions(schedulerProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}

			configuredScheduler, graphNodeScheduler, err := prepareComponent(scheduler, schedulerProto)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}

			// Append scheduler as a runtime.ConfiguredComponent
			configuredComponents = append(configuredComponents, configuredScheduler)

			// Append scheduler options
			options = append(options, schedulerOptions)

			// Merge port mapping for graph node
			err = portMapping.Merge(graphNodeScheduler.PortMapping)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}
		}

		// Actuation Strategy
		if loadActuatorProto := concurrencyLimiterProto.GetLoadActuator(); loadActuatorProto != nil {
			loadActuator, loadActuatorOptions, err := concurrency.NewLoadActuatorAndOptions(loadActuatorProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}

			configuredLoadActuator, graphNodeLoadActuator, err := prepareComponent(loadActuator, loadActuatorProto)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}
			// Append loadActuator as a runtime.ConfiguredComponent
			configuredComponents = append(configuredComponents, configuredLoadActuator)

			// Append loadActuator options
			options = append(options, loadActuatorOptions)

			// Merge port mapping for graph node
			err = portMapping.Merge(graphNodeLoadActuator.PortMapping)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}
		}

		_, graphNodeConcurrencyLimiter, err := prepareComponent(
			runtime.NewDummyComponent("ConcurrencyLimiter"),
			concurrencyLimiterProto,
		)
		if err != nil {
			return GraphNode{}, nil, nil, err
		}

		return GraphNode{
			PortMapping: portMapping,
			Config:      graphNodeConcurrencyLimiter.Config,
		}, configuredComponents, fx.Options(options...), nil
	} else if horizontalPodScalerProto := componentStackProto.GetHorizontalPodScaler(); horizontalPodScalerProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		horizontalPodScalerOptions, agentGroupName, horizontalPodScalerErr := horizontalpodscaler.NewHorizontalPodScalerOptions(horizontalPodScalerProto, componentStackIndex, policyReadAPI)
		if horizontalPodScalerErr != nil {
			return GraphNode{}, nil, nil, horizontalPodScalerErr
		}
		// Append horizontalPodScaler options
		options = append(options, horizontalPodScalerOptions)

		// Scale Reporter
		if scaleReporterProto := horizontalPodScalerProto.GetScaleReporter(); scaleReporterProto != nil {
			scaleReporter, scaleReporterOptions, err := horizontalpodscaler.NewScaleReporterAndOptions(scaleReporterProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}

			configuredScaleReporter, graphNodeScaleReporter, err := prepareComponent(scaleReporter, scaleReporterProto)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}

			// Append scaleReporter as a runtime.ConfiguredComponent
			configuredComponents = append(configuredComponents, configuredScaleReporter)

			// Append scaleReporter options
			options = append(options, scaleReporterOptions)

			// Merge port mapping for graph node
			err = portMapping.Merge(graphNodeScaleReporter.PortMapping)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}
		}

		// Scale Actuator
		if scaleActuatorProto := horizontalPodScalerProto.GetScaleActuator(); scaleActuatorProto != nil {
			scaleActuator, scaleActuatorOptions, err := horizontalpodscaler.NewScaleActuatorAndOptions(scaleActuatorProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}

			configuredScaleActuator, graphNodeScaleActuator, err := prepareComponent(scaleActuator, scaleActuatorProto)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}
			// Append scaleActuator as a runtime.ConfiguredComponent
			configuredComponents = append(configuredComponents, configuredScaleActuator)

			// Append scaleActuator options
			options = append(options, scaleActuatorOptions)

			// Merge port mapping for graph node
			err = portMapping.Merge(graphNodeScaleActuator.PortMapping)
			if err != nil {
				return GraphNode{}, nil, nil, err
			}
		}

		_, graphNodeHorizontalPodScaler, err := prepareComponent(
			runtime.NewDummyComponent("HorizontalPodScaler"),
			horizontalPodScalerProto,
		)
		if err != nil {
			return GraphNode{}, nil, nil, err
		}

		return GraphNode{
			PortMapping: portMapping,
			Config:      graphNodeHorizontalPodScaler.Config,
		}, configuredComponents, fx.Options(options...), nil
	}
	return GraphNode{}, nil, nil, fmt.Errorf("unsupported/missing component type")
}
