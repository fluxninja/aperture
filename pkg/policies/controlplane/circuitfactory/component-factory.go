package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/autoscale/podscaler"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/controller"
	loadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/load-scheduler"
	quotascheduler "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/quota-scheduler"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/sampler"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/query/promql"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
	"github.com/fluxninja/aperture/v2/pkg/utils"
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
) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
	var ctor componentConstructor
	switch config := componentProto.Component.(type) {
	case *policylangv1.Component_GradientController:
		ctor = mkCtor(config.GradientController, controller.NewGradientControllerAndOptions)
	case *policylangv1.Component_Ema:
		ctor = mkCtor(config.Ema, components.NewEMAAndOptions)
	case *policylangv1.Component_Sma:
		ctor = mkCtor(config.Sma, components.NewSMAAndOptions)
	case *policylangv1.Component_ArithmeticCombinator:
		ctor = mkCtor(config.ArithmeticCombinator, components.NewArithmeticCombinatorAndOptions)
	case *policylangv1.Component_Variable:
		ctor = mkCtor(config.Variable, components.NewVariableAndOptions)
	case *policylangv1.Component_BoolVariable:
		ctor = mkCtor(config.BoolVariable, components.NewBoolVariableAndOptions)
	case *policylangv1.Component_Decider:
		ctor = mkCtor(config.Decider, components.NewDeciderAndOptions)
	case *policylangv1.Component_Switcher:
		ctor = mkCtor(config.Switcher, components.NewSwitcherAndOptions)
	case *policylangv1.Component_UnaryOperator:
		ctor = mkCtor(config.UnaryOperator, components.NewUnaryOperatorAndOptions)
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
	case *policylangv1.Component_SignalGenerator:
		ctor = mkCtor(config.SignalGenerator, components.NewSignalGeneratorAndOptions)
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
		case *policylangv1.FlowControl_QuotaScheduler:
			ctor = mkCtor(flowControlConfig.QuotaScheduler, quotascheduler.NewQuotaSchedulerAndOptions)
		case *policylangv1.FlowControl_RateLimiter:
			ctor = mkCtor(flowControlConfig.RateLimiter, ratelimiter.NewRateLimiterAndOptions)
		case *policylangv1.FlowControl_Sampler:
			ctor = mkCtor(flowControlConfig.Sampler, sampler.NewSamplerAndOptions)
		case *policylangv1.FlowControl_Private:
			switch flowControlConfig.Private.TypeUrl {
			case "type.googleapis.com/aperture.policy.private.v1.LoadActuator":
				loadActuator := &policyprivatev1.LoadActuator{}
				if err := anypb.UnmarshalTo(flowControlConfig.Private, loadActuator, proto.UnmarshalOptions{}); err != nil {
					return Tree{}, nil, nil, err
				}
				ctor = mkCtor(loadActuator, loadscheduler.NewActuatorAndOptions)
			default:
				err := fmt.Errorf("unknown flow control type: %s", flowControlConfig.Private.TypeUrl)
				log.Error().Err(err).Msg("unknown flow control type")
				return Tree{}, nil, nil, err
			}

		default:
			return newFlowControlNestedAndOptions(flowControl, componentID, policyReadAPI)
		}
	case *policylangv1.Component_AutoScale:
		autoScale := componentProto.GetAutoScale()
		switch autoScaleConfig := autoScale.Component.(type) {
		case *policylangv1.AutoScale_Private:
			private := autoScaleConfig.Private
			switch private.TypeUrl {
			case "type.googleapis.com/aperture.policy.private.v1.PodScaleActuator":
				podScaleActuator := &policyprivatev1.PodScaleActuator{}
				if err := anypb.UnmarshalTo(private, podScaleActuator, proto.UnmarshalOptions{}); err != nil {
					return Tree{}, nil, nil, err
				}
				ctor = mkCtor(podScaleActuator, podscaler.NewScaleActuatorAndOptions)
			case "type.googleapis.com/aperture.policy.private.v1.PodScaleReporter":
				podScaleReporter := &policyprivatev1.PodScaleReporter{}
				if err := anypb.UnmarshalTo(private, podScaleReporter, proto.UnmarshalOptions{}); err != nil {
					return Tree{}, nil, nil, err
				}
				ctor = mkCtor(podScaleReporter, podscaler.NewScaleReporterAndOptions)
			default:
				err := fmt.Errorf("unknown auto scale type: %s", autoScaleConfig.Private.TypeUrl)
				log.Error().Err(err).Msg("unknown auto scale type")
				return Tree{}, nil, nil, err
			}
		default:
			return newAutoScaleNestedAndOptions(autoScale, componentID, policyReadAPI)
		}
	default:
		return Tree{}, nil, nil, fmt.Errorf("unknown component type: %T", config)
	}

	component, config, option, err := ctor(componentID, policyReadAPI)
	if err != nil {
		return Tree{}, nil, nil, err
	}

	configuredComponent, err := prepareComponent(component, config, componentID, true)
	if err != nil {
		return Tree{}, nil, nil, err
	}

	return Tree{Node: configuredComponent}, []*runtime.ConfiguredComponent{configuredComponent}, option, nil
}

type componentConstructor func(
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (runtime.Component, any, fx.Option, error)

func mkCtor[Config any, Comp runtime.Component](
	config *Config,
	origCtor func(*Config, runtime.ComponentID, iface.Policy) (Comp, fx.Option, error),
) componentConstructor {
	return func(componentID runtime.ComponentID, policy iface.Policy) (runtime.Component, any, fx.Option, error) {
		comp, opt, err := origCtor(config, componentID, policy)
		return comp, config, opt, err
	}
}

func prepareComponent(
	component runtime.Component,
	config any,
	componentID runtime.ComponentID,
	doParsePortMapping bool,
) (*runtime.ConfiguredComponent, error) {
	subCircuitID, ok := componentID.ParentID()
	if !ok {
		return nil, fmt.Errorf("component %s is not in a circuit", componentID.String())
	}

	return prepareComponentInCircuit(component, config, componentID, subCircuitID, doParsePortMapping)
}

func prepareComponentInCircuit(
	component runtime.Component,
	config any,
	componentID runtime.ComponentID,
	subCircuitID runtime.ComponentID,
	doParsePortMapping bool,
) (*runtime.ConfiguredComponent, error) {
	mapStruct, err := utils.ToMapStruct(config)
	if err != nil {
		return nil, err
	}

	ports := runtime.NewPortMapping()
	if doParsePortMapping {
		ports, err = runtime.PortsFromComponentConfig(mapStruct, subCircuitID.String())
		if err != nil {
			return nil, err
		}
	}

	return &runtime.ConfiguredComponent{
		Component:   component,
		PortMapping: ports,
		Config:      mapStruct,
		ComponentID: componentID,
	}, nil
}
