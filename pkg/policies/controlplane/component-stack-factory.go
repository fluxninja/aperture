package controlplane

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component/actuator/concurrency"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// ComponentStackFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func ComponentStackFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options()
}

// NewComponentStackAndOptions creates components for component stack, sub components and their fx options.
func NewComponentStackAndOptions(
	componentStackProto *policylangv1.Component,
	componentStackIndex int,
	policyReadAPI iface.PolicyRead,
) (compiledComponent, []compiledComponent, fx.Option, error) {
	// Factory parser to determine what kind of component stack to create
	if concurrencyLimiter := componentStackProto.GetConcurrencyLimiter(); concurrencyLimiter != nil {
		var (
			compiledComponents []compiledComponent
			options            []fx.Option
		)
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiter, componentStackIndex, policyReadAPI)
		if concurrencyLimiterErr != nil {
			return compiledComponent{}, nil, nil, concurrencyLimiterErr
		}
		// Append concurrencyLimiter options
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiter.GetScheduler(); schedulerProto != nil {
			scheduler, schedulerOptions, schedulerErr := concurrency.NewSchedulerAndOptions(schedulerProto, componentStackIndex, policyReadAPI, agentGroupName)
			if schedulerErr != nil {
				return compiledComponent{}, nil, nil, schedulerErr
			}
			schedulerMapStruct, err := encodeMapStruct(schedulerProto)
			if err != nil {
				return compiledComponent{}, nil, nil, err
			}
			// Append scheduler as a compiledComponent
			compiledComponents = append(compiledComponents, compiledComponent{
				component: scheduler,
				mapStruct: schedulerMapStruct,
				name:      "Scheduler",
			})

			// Append scheduler options
			options = append(options, schedulerOptions)
		}

		// Actuation Strategy
		if loadShedActuatorProto := concurrencyLimiter.GetLoadShedActuator(); loadShedActuatorProto != nil {
			loadShedActuator, loadShedActuatorOptions, loadShedActuatorErr := concurrency.NewLoadShedActuatorAndOptions(loadShedActuatorProto, componentStackIndex, policyReadAPI, agentGroupName)
			if loadShedActuatorErr != nil {
				return compiledComponent{}, nil, nil, loadShedActuatorErr
			}
			loadShedActuatorMapStruct, err := encodeMapStruct(loadShedActuatorProto)
			if err != nil {
				return compiledComponent{}, nil, nil, err
			}
			// Append loadShedActuator as a compiledComponent
			compiledComponents = append(compiledComponents, compiledComponent{
				component: loadShedActuator,
				mapStruct: loadShedActuatorMapStruct,
				name:      "LoadShedActuator",
			})

			// Append loadShedActuator options
			options = append(options, loadShedActuatorOptions)
		}

		return compiledComponent{}, compiledComponents, fx.Options(options...), nil
	}
	return compiledComponent{}, nil, nil, fmt.Errorf("unsupported/missing component type")
}
