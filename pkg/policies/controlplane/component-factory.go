package controlplane

import (
	"encoding/json"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component/actuator/rate"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/component/controller"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// ComponentFactoryModule for component factory run via the main app.
func ComponentFactoryModule() fx.Option {
	return fx.Options(
		component.PromQLModule(),
	)
}

// ComponentFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func ComponentFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		ComponentStackFactoryModuleForPolicyApp(circuitAPI),
		component.PromQLModuleForPolicyApp(circuitAPI),
	)
}

// NewComponentAndOptions creates component and its fx options.
func NewComponentAndOptions(
	componentProto *policylangv1.Component,
	componentIndex int,
	policyReadAPI iface.PolicyRead,
) (string, map[string]any, map[string]runtime.Component, runtime.Component, fx.Option, error) {
	// Factory parser to determine what kind of component to create
	if gradientController := componentProto.GetGradientController(); gradientController != nil {
		component, option, err := controller.NewGradientControllerAndOptions(gradientController, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(gradientController, err)
		return "Gradient", mapStruct, nil, component, option, err
	} else if limiter := componentProto.GetRateLimiter(); limiter != nil {
		component, option, err := rate.NewRateLimiterAndOptions(limiter, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(limiter, err)
		return "RateLimiter", mapStruct, nil, component, option, err
	} else if ema := componentProto.GetEma(); ema != nil {
		component, option, err := component.NewEMAAndOptions(ema, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(ema, err)
		return "EMA", mapStruct, nil, component, option, err
	} else if arithmeticCombinator := componentProto.GetArithmeticCombinator(); arithmeticCombinator != nil {
		component, option, err := component.NewArithmeticCombinatorAndOptions(arithmeticCombinator, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(arithmeticCombinator, err)
		return "ArithmeticCombinator", mapStruct, nil, component, option, err
	} else if promQL := componentProto.GetPromql(); promQL != nil {
		component, option, err := component.NewPromQLAndOptions(promQL, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(promQL, err)
		return "PromQL", mapStruct, nil, component, option, err
	} else if constant := componentProto.GetConstant(); constant != nil {
		component, option, err := component.NewConstantAndOptions(constant, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(constant, err)
		return "Constant", mapStruct, nil, component, option, err
	} else if decider := componentProto.GetDecider(); decider != nil {
		component, option, err := component.NewDeciderAndOptions(decider, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(decider, err)
		return "Decider", mapStruct, nil, component, option, err
	} else if sqrt := componentProto.GetSqrt(); sqrt != nil {
		component, option, err := component.NewSqrtAndOptions(sqrt, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(sqrt, err)
		return "Sqrt", mapStruct, nil, component, option, err
	} else if max := componentProto.GetMax(); max != nil {
		component, option, err := component.NewMaxAndOptions(max, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(max, err)
		return "Max", mapStruct, nil, component, option, err
	} else if min := componentProto.GetMin(); min != nil {
		component, option, err := component.NewMinAndOptions(min, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(min, err)
		return "Min", mapStruct, nil, component, option, err
	} else if extrapolator := componentProto.GetExtrapolator(); extrapolator != nil {
		component, option, err := component.NewExtrapolatorAndOptions(extrapolator, componentIndex, policyReadAPI)
		mapStruct, err := getComponentMapStruct(extrapolator, err)
		return "Extrapolator", mapStruct, nil, component, option, err
	} else {
		// Try Component Stack Factory
		componentName, mapStruct, subComponents, option, err := NewComponentStackAndOptions(componentProto, componentIndex, policyReadAPI)
		return componentName, mapStruct, subComponents, nil, option, err
	}
}

func getComponentMapStruct(comp any, err error) (map[string]any, error) {
	if err != nil {
		return nil, err
	}
	// TODO: mapstruct functionality needs to be moved to a common package
	var mapStruct map[string]interface{}
	b, err := json.Marshal(comp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &mapStruct)
	if err != nil {
		return nil, err
	}
	return mapStruct, nil
}
