package controlplane

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/FluxNinja/aperture/pkg/policies/apis/policyapi"
	"github.com/FluxNinja/aperture/pkg/policies/controlplane/component/actuator/concurrency"
	"github.com/FluxNinja/aperture/pkg/policies/controlplane/runtime"
)

// ComponentStackFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func ComponentStackFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options()
}

// NewComponentStackAndOptions creates components for component stack, sub components and their fx options.
func NewComponentStackAndOptions(
	componentStackProto *policylangv1.Component,
	componentStackIndex int,
	policyReadAPI policyapi.PolicyReadAPI,
) (string, map[string]any, map[string]runtime.Component, fx.Option, error) {
	// Factory parser to determine what kind of component stack to create
	if concurrencyLimiter := componentStackProto.GetConcurrencyLimiter(); concurrencyLimiter != nil {
		subComponents := make(map[string]runtime.Component)
		var options []fx.Option
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiter, componentStackIndex, policyReadAPI)
		if concurrencyLimiterErr != nil {
			return "", nil, nil, nil, concurrencyLimiterErr
		}
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiter.GetScheduler(); schedulerProto != nil {
			subComponentKey := "scheduler"
			scheduler, schedulerOptions, schedulerErr := concurrency.NewSchedulerAndOptions(schedulerProto, componentStackIndex, policyReadAPI, agentGroupName)
			if schedulerErr != nil {
				return "", nil, nil, nil, schedulerErr
			}
			subComponents[subComponentKey] = scheduler
			options = append(options, schedulerOptions)
		}

		// Actuation Strategy
		if loadShedActuatorProto := concurrencyLimiter.GetLoadShedActuator(); loadShedActuatorProto != nil {
			subComponentKey := "load_shed_actuator"
			loadShedActuator, loadShedActuatorOptions, loadShedActuatorErr := concurrency.NewLoadShedActuatorAndOptions(loadShedActuatorProto, componentStackIndex, policyReadAPI, agentGroupName)
			if loadShedActuatorErr != nil {
				return "", nil, nil, nil, loadShedActuatorErr
			}
			subComponents[subComponentKey] = loadShedActuator
			options = append(options, loadShedActuatorOptions)
		}

		mapStruct, err := getComponentMapStruct(concurrencyLimiter, nil)
		return "ConcurrencyLimiter", mapStruct, subComponents, fx.Options(options...), err
	}
	return "", nil, nil, nil, fmt.Errorf("unsupported/missing component type")
}
