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
) (runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	// Factory parser to determine what kind of component stack to create
	if concurrencyLimiterProto := componentStackProto.GetConcurrencyLimiter(); concurrencyLimiterProto != nil {
		var (
			compiledComponents []runtime.ConfiguredComponent
			options            []fx.Option
		)
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiterProto, componentStackIndex, policyReadAPI)
		if concurrencyLimiterErr != nil {
			return runtime.ConfiguredComponent{}, nil, nil, concurrencyLimiterErr
		}
		// Append concurrencyLimiter options
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiterProto.GetScheduler(); schedulerProto != nil {
			scheduler, schedulerOptions, err := concurrency.NewSchedulerAndOptions(schedulerProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}

			compiledScheduler, err := prepareConfiguredComponent(scheduler, schedulerProto)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}

			// Append scheduler as a runtime.ConfiguredComponent
			compiledComponents = append(compiledComponents, compiledScheduler)

			// Append scheduler options
			options = append(options, schedulerOptions)
		}

		// Actuation Strategy
		if loadActuatorProto := concurrencyLimiterProto.GetLoadActuator(); loadActuatorProto != nil {
			loadActuator, loadActuatorOptions, err := concurrency.NewLoadActuatorAndOptions(loadActuatorProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}

			compiledLoadActuator, err := prepareConfiguredComponent(loadActuator, loadActuatorProto)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}
			// Append loadActuator as a runtime.ConfiguredComponent
			compiledComponents = append(compiledComponents, compiledLoadActuator)

			// Append loadActuator options
			options = append(options, loadActuatorOptions)
		}

		compiledConcurrencyLimiter, err := prepareConfiguredComponent(
			runtime.NewDummyComponent("ConcurrencyLimiter"),
			concurrencyLimiterProto,
		)
		if err != nil {
			return runtime.ConfiguredComponent{}, nil, nil, err
		}

		return compiledConcurrencyLimiter, compiledComponents, fx.Options(options...), nil
	} else if horizontalPodScalerProto := componentStackProto.GetHorizontalPodScaler(); horizontalPodScalerProto != nil {
		var (
			compiledComponents []runtime.ConfiguredComponent
			options            []fx.Option
		)
		horizontalPodScalerOptions, agentGroupName, horizontalPodScalerErr := horizontalpodscaler.NewHorizontalPodScalerOptions(horizontalPodScalerProto, componentStackIndex, policyReadAPI)
		if horizontalPodScalerErr != nil {
			return runtime.ConfiguredComponent{}, nil, nil, horizontalPodScalerErr
		}
		// Append horizontalPodScaler options
		options = append(options, horizontalPodScalerOptions)

		// Scale Reporter
		if scaleReporterProto := horizontalPodScalerProto.GetScaleReporter(); scaleReporterProto != nil {
			scaleReporter, scaleReporterOptions, err := horizontalpodscaler.NewScaleReporterAndOptions(scaleReporterProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}

			compiledScaleReporter, err := prepareConfiguredComponent(scaleReporter, scaleReporterProto)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}

			// Append scaleReporter as a runtime.ConfiguredComponent
			compiledComponents = append(compiledComponents, compiledScaleReporter)

			// Append scaleReporter options
			options = append(options, scaleReporterOptions)
		}

		// Scale Actuator
		if scaleActuatorProto := horizontalPodScalerProto.GetScaleActuator(); scaleActuatorProto != nil {
			scaleActuator, scaleActuatorOptions, err := horizontalpodscaler.NewScaleActuatorAndOptions(scaleActuatorProto, componentStackIndex, policyReadAPI, agentGroupName)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}

			compiledScaleActuator, err := prepareConfiguredComponent(scaleActuator, scaleActuatorProto)
			if err != nil {
				return runtime.ConfiguredComponent{}, nil, nil, err
			}
			// Append scaleActuator as a runtime.ConfiguredComponent
			compiledComponents = append(compiledComponents, compiledScaleActuator)

			// Append scaleActuator options
			options = append(options, scaleActuatorOptions)
		}

		compiledHorizontalPodScaler, err := prepareConfiguredComponent(
			runtime.NewDummyComponent("HorizontalPodScaler"),
			horizontalPodScalerProto,
		)
		if err != nil {
			return runtime.ConfiguredComponent{}, nil, nil, err
		}

		return compiledHorizontalPodScaler, compiledComponents, fx.Options(options...), nil
	}
	return runtime.ConfiguredComponent{}, nil, nil, fmt.Errorf("unsupported/missing component type")
}
