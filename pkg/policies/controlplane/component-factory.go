package controlplane

import (
	"encoding/json"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/rate"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/controller"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// componentFactoryModule for component factory run via the main app.
func componentFactoryModule() fx.Option {
	return fx.Options(
		components.PromQLModule(),
	)
}

// componentFactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func componentFactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		componentStackFactoryModuleForPolicyApp(circuitAPI),
		components.PromQLModuleForPolicyApp(circuitAPI),
	)
}

// NewComponentAndOptions creates component and its fx options.
func NewComponentAndOptions(
	componentProto *policylangv1.Component,
	componentIndex int,
	policyReadAPI iface.Policy,
) (runtime.CompiledComponent, []runtime.CompiledComponent, fx.Option, error) {
	// Factory parser to determine what kind of component to create
	if gradientController := componentProto.GetGradientController(); gradientController != nil {
		component, option, err := controller.NewGradientControllerAndOptions(gradientController, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(gradientController, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Gradient",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if limiter := componentProto.GetRateLimiter(); limiter != nil {
		component, option, err := rate.NewRateLimiterAndOptions(limiter, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(limiter, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "RateLimiter",
			ComponentType: runtime.ComponentTypeSink,
		}, nil, option, err
	} else if ema := componentProto.GetEma(); ema != nil {
		component, option, err := components.NewEMAAndOptions(ema, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(ema, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "EMA",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if arithmeticCombinator := componentProto.GetArithmeticCombinator(); arithmeticCombinator != nil {
		component, option, err := components.NewArithmeticCombinatorAndOptions(arithmeticCombinator, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(arithmeticCombinator, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "ArithmeticCombinator",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if promQL := componentProto.GetPromql(); promQL != nil {
		component, option, err := components.NewPromQLAndOptions(promQL, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(promQL, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "PromQL",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if constant := componentProto.GetConstant(); constant != nil {
		component, option, err := components.NewConstantAndOptions(constant, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(constant, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Constant",
			ComponentType: runtime.ComponentTypeSource,
		}, nil, option, err
	} else if decider := componentProto.GetDecider(); decider != nil {
		component, option, err := components.NewDeciderAndOptions(decider, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(decider, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Decider",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if switcher := componentProto.GetSwitcher(); switcher != nil {
		component, option, err := components.NewSwitcherAndOptions(switcher, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(switcher, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Switcher",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if sqrt := componentProto.GetSqrt(); sqrt != nil {
		component, option, err := components.NewSqrtAndOptions(sqrt, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(sqrt, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Sqrt",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if max := componentProto.GetMax(); max != nil {
		component, option, err := components.NewMaxAndOptions(max, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(max, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Max",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if min := componentProto.GetMin(); min != nil {
		component, option, err := components.NewMinAndOptions(min, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(min, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Min",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if extrapolator := componentProto.GetExtrapolator(); extrapolator != nil {
		component, option, err := components.NewExtrapolatorAndOptions(extrapolator, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(extrapolator, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Extrapolator",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if firstValid := componentProto.GetFirstValid(); firstValid != nil {
		component, option, err := components.NewFirstValidAndOptions(firstValid, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(firstValid, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "FirstValid",
			ComponentType: runtime.ComponentTypeSignalProcessor,
		}, nil, option, err
	} else if sink := componentProto.GetSink(); sink != nil {
		component, option, err := components.NewSinkAndOptions(sink, componentIndex, policyReadAPI)
		mapStruct, err := encodeMapStructOnNilErr(sink, err)
		return runtime.CompiledComponent{
			Component:     component,
			MapStruct:     mapStruct,
			Name:          "Sink",
			ComponentType: runtime.ComponentTypeSink,
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
