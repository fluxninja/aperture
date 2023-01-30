package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/mapstruct"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
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
		autoScaleModuleForPolicyApp(circuitAPI),
		promql.ModuleForPolicyApp(circuitAPI),
	)
}

// NewComponentAndOptions creates parent and leaf components and their fx options for a component spec.
func NewComponentAndOptions(
	componentProto *policylangv1.Component,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
	var ctor componentConstructor
	switch config := componentProto.Component.(type) {
	case *policylangv1.Component_GradientController:
		ctor = mkCtor(config.GradientController, controller.NewGradientControllerAndOptions)
	case *policylangv1.Component_Ema:
		ctor = mkCtor(config.Ema, components.NewEMAAndOptions)
	case *policylangv1.Component_ArithmeticCombinator:
		ctor = mkCtor(config.ArithmeticCombinator, components.NewArithmeticCombinatorAndOptions)
	case *policylangv1.Component_Variable:
		ctor = mkCtor(config.Variable, components.NewVariableAndOptions)
	case *policylangv1.Component_Decider:
		ctor = mkCtor(config.Decider, components.NewDeciderAndOptions)
	case *policylangv1.Component_Switcher:
		ctor = mkCtor(config.Switcher, components.NewSwitcherAndOptions)
	case *policylangv1.Component_Sqrt:
		ctor = mkCtor(config.Sqrt, components.NewSqrtAndOptions)
	case *policylangv1.Component_Max:
		ctor = mkCtor(config.Max, components.NewMaxAndOptions)
	case *policylangv1.Component_Min:
		ctor = mkCtor(config.Min, components.NewMinAndOptions)
	case *policylangv1.Component_Extrapolator:
		ctor = mkCtor(config.Extrapolator, components.NewExtrapolatorAndOptions)
	case *policylangv1.Component_FirstValid:
		ctor = mkCtor(config.FirstValid, components.NewFirstValidAndOptions)
	case *policylangv1.Component_Alerter:
		ctor = mkCtor(config.Alerter, components.NewAlerterAndOptions)
	case *policylangv1.Component_Integrator:
		ctor = mkCtor(config.Integrator, components.NewIntegratorAndOptions)
	case *policylangv1.Component_Differentiator:
		ctor = mkCtor(config.Differentiator, components.NewDifferentiatorAndOptions)
	case *policylangv1.Component_And:
		ctor = mkCtor(config.And, components.NewAndAndOptions)
	case *policylangv1.Component_Or:
		ctor = mkCtor(config.Or, components.NewOrAndOptions)
	case *policylangv1.Component_Inverter:
		ctor = mkCtor(config.Inverter, components.NewInverterAndOptions)
	case *policylangv1.Component_PulseGenerator:
		ctor = mkCtor(config.PulseGenerator, components.NewPulseGeneratorAndOptions)
	case *policylangv1.Component_Holder:
		ctor = mkCtor(config.Holder, components.NewHolderAndOptions)
	case *policylangv1.Component_NestedSignalIngress:
		ctor = mkCtor(config.NestedSignalIngress, components.NewNestedSignalIngressAndOptions)
	case *policylangv1.Component_NestedSignalEgress:
		ctor = mkCtor(config.NestedSignalEgress, components.NewNestedSignalEgressAndOptions)
	case *policylangv1.Component_NestedCircuit:
		return ParseNestedCircuit(componentID, config.NestedCircuit, policyReadAPI)
	case *policylangv1.Component_Query:
		query := componentProto.GetQuery()
		switch queryConfig := query.Component.(type) {
		case *policylangv1.Query_Promql:
			ctor = mkCtor(queryConfig.Promql, promql.NewPromQLAndOptions)
		}
	case *policylangv1.Component_FlowControl:
		flowControl := componentProto.GetFlowControl()
		switch flowControlConfig := flowControl.Component.(type) {
		case *policylangv1.FlowControl_RateLimiter:
			ctor = mkCtor(flowControlConfig.RateLimiter, rate.NewRateLimiterAndOptions)
		default:
			return newFlowControlCompositeAndOptions(flowControl, componentID, policyReadAPI)
		}
	case *policylangv1.Component_AutoScale:
		autoScale := componentProto.GetAutoScale()
		return newAutoScaleCompositeAndOptions(autoScale, componentID, policyReadAPI)
	default:
		return Tree{}, nil, nil, fmt.Errorf("unknown component type: %T", config)
	}

	component, config, option, err := ctor(componentID.String(), policyReadAPI)
	if err != nil {
		return Tree{}, nil, nil, err
	}

	configuredComponent, err := prepareComponent(component, config, componentID)
	if err != nil {
		return Tree{}, nil, nil, err
	}

	return Tree{Root: configuredComponent}, []runtime.ConfiguredComponent{configuredComponent}, option, nil
}

type componentConstructor func(
	componentID string,
	policyReadAPI iface.Policy,
) (runtime.Component, any, fx.Option, error)

func mkCtor[Config any, Comp runtime.Component](
	config *Config,
	origCtor func(*Config, string, iface.Policy) (Comp, fx.Option, error),
) componentConstructor {
	return func(componentID string, policy iface.Policy) (runtime.Component, any, fx.Option, error) {
		comp, opt, err := origCtor(config, componentID, policy)
		return comp, config, opt, err
	}
}

func prepareComponent(
	component runtime.Component,
	config any,
	componentID runtime.ComponentID,
) (runtime.ConfiguredComponent, error) {
	subCircuitID, ok := componentID.ParentID()
	if !ok {
		return runtime.ConfiguredComponent{}, fmt.Errorf("component %s is not in a circuit", componentID.String())
	}

	return prepareComponentInCircuit(component, config, componentID, subCircuitID)
}

func prepareComponentInCircuit(
	component runtime.Component,
	config any,
	componentID runtime.ComponentID,
	subCircuitID runtime.ComponentID,
) (runtime.ConfiguredComponent, error) {
	mapStruct, err := mapstruct.EncodeObject(config)
	if err != nil {
		return runtime.ConfiguredComponent{}, err
	}

	ports, err := runtime.PortsFromComponentConfig(mapStruct, subCircuitID.String())
	if err != nil {
		return runtime.ConfiguredComponent{}, err
	}

	return runtime.ConfiguredComponent{
		Component:   component,
		PortMapping: ports,
		Config:      mapStruct,
		ComponentID: componentID,
	}, nil
}
