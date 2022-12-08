package components

import (
	"encoding/json"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/rate"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/controller"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/promql"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// FactoryModule for component factory run via the main app.
func FactoryModule() fx.Option {
	return fx.Options(
		promql.Module(),
	)
}

// FactoryModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func FactoryModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
	return fx.Options(
		componentStackFactoryModuleForPolicyApp(circuitAPI),
		promql.ModuleForPolicyApp(circuitAPI),
	)
}

// NewComponentAndOptions creates component and its fx options.
func NewComponentAndOptions(
	componentProto *policylangv1.Component,
	componentIndex int,
	policyReadAPI iface.Policy,
) (runtime.CompiledComponent, []runtime.CompiledComponent, fx.Option, error) {
	var ctor componentConstructor
	switch config := componentProto.Component.(type) {
	case *policylangv1.Component_GradientController:
		ctor = mkCtor(config.GradientController, controller.NewGradientControllerAndOptions)
	case *policylangv1.Component_RateLimiter:
		ctor = mkCtor(config.RateLimiter, rate.NewRateLimiterAndOptions)
	case *policylangv1.Component_Ema:
		ctor = mkCtor(config.Ema, NewEMAAndOptions)
	case *policylangv1.Component_ArithmeticCombinator:
		ctor = mkCtor(config.ArithmeticCombinator, NewArithmeticCombinatorAndOptions)
	case *policylangv1.Component_Promql:
		ctor = mkCtor(config.Promql, promql.NewPromQLAndOptions)
	case *policylangv1.Component_Constant:
		ctor = mkCtor(config.Constant, NewConstantAndOptions)
	case *policylangv1.Component_Decider:
		ctor = mkCtor(config.Decider, NewDeciderAndOptions)
	case *policylangv1.Component_Switcher:
		ctor = mkCtor(config.Switcher, NewSwitcherAndOptions)
	case *policylangv1.Component_Sqrt:
		ctor = mkCtor(config.Sqrt, NewSqrtAndOptions)
	case *policylangv1.Component_Max:
		ctor = mkCtor(config.Max, NewMaxAndOptions)
	case *policylangv1.Component_Min:
		ctor = mkCtor(config.Min, NewMinAndOptions)
	case *policylangv1.Component_Extrapolator:
		ctor = mkCtor(config.Extrapolator, NewExtrapolatorAndOptions)
	case *policylangv1.Component_FirstValid:
		ctor = mkCtor(config.FirstValid, NewFirstValidAndOptions)
	case *policylangv1.Component_Sink:
		ctor = mkCtor(config.Sink, NewSinkAndOptions)
	case *policylangv1.Component_Alerter:
		ctor = mkCtor(config.Alerter, NewAlerterAndOptions)
	default:
		return newComponentStackAndOptions(componentProto, componentIndex, policyReadAPI)
	}

	component, config, option, err := ctor(componentIndex, policyReadAPI)
	if err != nil {
		return runtime.CompiledComponent{}, nil, nil, err
	}

	compiledComponent, err := prepareCompiledComponent(component, config)
	if err != nil {
		return runtime.CompiledComponent{}, nil, nil, err
	}

	return compiledComponent, nil, option, nil
}

type componentConstructor func(
	componentIdx int,
	policyReadAPI iface.Policy,
) (runtime.Component, any, fx.Option, error)

func mkCtor[Config any, Comp runtime.Component](
	config *Config,
	origCtor func(*Config, int, iface.Policy) (Comp, fx.Option, error),
) componentConstructor {
	return func(idx int, policy iface.Policy) (runtime.Component, any, fx.Option, error) {
		comp, opt, err := origCtor(config, idx, policy)
		return comp, config, opt, err
	}
}

func prepareCompiledComponent(
	component runtime.Component,
	config any,
) (runtime.CompiledComponent, error) {
	mapStruct, err := encodeMapStruct(config)
	if err != nil {
		return runtime.CompiledComponent{}, err
	}

	ports, err := runtime.PortsFromMapStruct(mapStruct)
	if err != nil {
		return runtime.CompiledComponent{}, err
	}

	return runtime.CompiledComponent{
		Component: component,
		MapStruct: mapStruct,
		Ports:     ports,
	}, nil
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
