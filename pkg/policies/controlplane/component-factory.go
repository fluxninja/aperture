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

// componentFactoryModule for component factory run via the main app.
func componentFactoryModule() fx.Option {
	return fx.Options(
		component.PromQLModule(),
	)
}

// componentFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func componentFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		componentStackFactoryModuleForPolicyApp(circuitAPI),
		component.PromQLModuleForPolicyApp(circuitAPI),
	)
}

// NewComponentAndOptions creates component and its fx options.
func NewComponentAndOptions(
	componentProto *policylangv1.Component,
	componentIndex int,
	policyReadAPI iface.PolicyRead,
) (runtime.CompiledComponent, []runtime.CompiledComponent, fx.Option, error) {
	// Factory parser to determine what kind of component to create
	if gradientController := componentProto.GetGradientController(); gradientController != nil {
		component, option, err := controller.NewGradientControllerAndOptions(gradientController, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(gradientController, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Gradient",
		}, nil, option, err
	} else if limiter := componentProto.GetRateLimiter(); limiter != nil {
		component, option, err := rate.NewRateLimiterAndOptions(limiter, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(limiter, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "RateLimiter",
		}, nil, option, err
	} else if ema := componentProto.GetEma(); ema != nil {
		component, option, err := component.NewEMAAndOptions(ema, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(ema, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "EMA",
		}, nil, option, err
	} else if arithmeticCombinator := componentProto.GetArithmeticCombinator(); arithmeticCombinator != nil {
		component, option, err := component.NewArithmeticCombinatorAndOptions(arithmeticCombinator, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(arithmeticCombinator, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "ArithmeticCombinator",
		}, nil, option, err
	} else if promQL := componentProto.GetPromql(); promQL != nil {
		component, option, err := component.NewPromQLAndOptions(promQL, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(promQL, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "PromQL",
		}, nil, option, err
	} else if constant := componentProto.GetConstant(); constant != nil {
		component, option, err := component.NewConstantAndOptions(constant, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(constant, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Constant",
		}, nil, option, err
	} else if decider := componentProto.GetDecider(); decider != nil {
		component, option, err := component.NewDeciderAndOptions(decider, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(decider, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Decider",
		}, nil, option, err
	} else if sqrt := componentProto.GetSqrt(); sqrt != nil {
		component, option, err := component.NewSqrtAndOptions(sqrt, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(sqrt, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Sqrt",
		}, nil, option, err
	} else if max := componentProto.GetMax(); max != nil {
		component, option, err := component.NewMaxAndOptions(max, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(max, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Max",
		}, nil, option, err
	} else if min := componentProto.GetMin(); min != nil {
		component, option, err := component.NewMinAndOptions(min, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(min, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Min",
		}, nil, option, err
	} else if extrapolator := componentProto.GetExtrapolator(); extrapolator != nil {
		component, option, err := component.NewExtrapolatorAndOptions(extrapolator, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(extrapolator, err)
		return runtime.CompiledComponent{
			Component: component,
			MapStruct: mapStruct,
			Name:      "Extrapolator",
		}, nil, option, err
	} else {
		// Try Component Stack Factory
		return newComponentStackAndOptions(componentProto, componentIndex, policyReadAPI)
	}
}

func encodeMapStructOnNilErr(comp any, err error) (map[string]any, error) {
	if err != nil {
		return nil, err
	}
	return encodeMapStruct(comp)
}

func encodeMapStruct(obj any) (map[string]any, error) {
	// TODO: mapstruct functionality needs to be moved to a common package
	var mapStruct map[string]interface{}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &mapStruct)
	if err != nil {
		return nil, err
	}
	return mapStruct, nil
}
