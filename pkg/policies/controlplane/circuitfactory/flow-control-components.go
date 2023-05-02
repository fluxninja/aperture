package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/flowcontrol/loadscheduler"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/flowcontrol/regulator"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// newFlowControlCompositeAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newFlowControlCompositeAndOptions(
	flowControlComponentProto *policylangv1.FlowControl,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
	retErr := func(err error) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
		return Tree{}, nil, nil, err
	}

	parentCircuitID, ok := componentID.ParentID()
	if !ok {
		return retErr(fmt.Errorf("parent circuit ID not found for component %s", componentID))
	}
	loadSchedulerProto := &policylangv1.LoadScheduler{}
	isLoadScheduler := false
	if proto := flowControlComponentProto.GetLoadScheduler(); proto != nil {
		loadSchedulerProto = proto
		isLoadScheduler = true
	} else if concurrencyLimiterProto := flowControlComponentProto.GetConcurrencyLimiter(); concurrencyLimiterProto != nil {
		// Convert from *policylangv1.FlowControl_ConcurrencyLimiter to *policylangv1.FlowControl_LoadScheduler
		fieldMappings := map[string]string{
			"load_actuator": "actuator",
			"scheduler.out_ports.accepted_concurrency": "accepted_token_rate",
			"scheduler.out_ports.incoming_concurrency": "incoming_token_rate",
		}
		err := convertOldComponentToNew(concurrencyLimiterProto, loadSchedulerProto, fieldMappings)
		if err != nil {
			return Tree{}, nil, nil, err
		}
		isLoadScheduler = true
	}

	adaptiveLoadSchedulerProto := &policylangv1.AdaptiveLoadScheduler{}
	isAdaptiveLoadScheduler := false
	if proto := flowControlComponentProto.GetAdaptiveLoadScheduler(); proto != nil {
		adaptiveLoadSchedulerProto = proto
		isAdaptiveLoadScheduler = true
	} else if aimdConcurrencyControllerProto := flowControlComponentProto.GetAimdConcurrencyController(); aimdConcurrencyControllerProto != nil {
		// Convert from *policylangv1.FlowControl_AimdConcurrencyController to *policylangv1.FlowControl_AdaptiveLoadScheduler
		fieldMappings := map[string]string{
			"out_ports.accepted_concurrency": "accepted_token_rate",
			"out_ports.incoming_concurrency": "incoming_token_rate",
		}
		err := convertOldComponentToNew(aimdConcurrencyControllerProto, adaptiveLoadSchedulerProto, fieldMappings)
		if err != nil {
			return Tree{}, nil, nil, err
		}
		isAdaptiveLoadScheduler = true
	}

	loadRampProto := &policylangv1.LoadRamp{}
	isLoadRamp := false
	if proto := flowControlComponentProto.GetLoadRamp(); proto != nil {
		loadRampProto = proto
		isLoadRamp = true
	} else if loadShaperProto := flowControlComponentProto.GetLoadShaper(); loadShaperProto != nil {
		// Convert from *policylangv1.FlowControl_LoadShaper to *policylangv1.FlowControl_LoadRamp
		fieldMappings := map[string]string{
			"flow_regulator_parameters": "regulator_parameters",
		}
		err := convertOldComponentToNew(loadShaperProto, loadRampProto, fieldMappings)
		if err != nil {
			return Tree{}, nil, nil, err
		}
		isLoadRamp = true
	}

	// Factory parser to determine what kind of composite component to create
	if isLoadScheduler {
		var (
			configuredComponents []*runtime.ConfiguredComponent
			tree                 Tree
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		loadSchedulerOptions, agentGroups, loadSchedulerErr := loadscheduler.NewLoadSchedulerOptions(loadSchedulerProto, componentID.String(), policyReadAPI)
		if loadSchedulerErr != nil {
			return retErr(loadSchedulerErr)
		}
		options = append(options, loadSchedulerOptions)

		// Scheduler
		if schedulerProto := loadSchedulerProto.GetScheduler(); schedulerProto != nil {
			// Use the same id as the component stack since agent sees only the component stack and generates metrics tagged with the component stack id
			scheduler, schedulerOptions, err := loadscheduler.NewSchedulerAndOptions(schedulerProto, componentID.String(), policyReadAPI, agentGroups)
			if err != nil {
				return retErr(err)
			}

			// Need a unique ID for sub component since it is used for graph generation
			schedulerConfComp, err := prepareComponentInCircuit(scheduler, schedulerProto, componentID.ChildID("Scheduler"), parentCircuitID, true)
			if err != nil {
				return retErr(err)
			}

			configuredComponents = append(configuredComponents, schedulerConfComp)
			tree.Children = append(tree.Children, Tree{Node: schedulerConfComp})

			options = append(options, schedulerOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(schedulerConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		// Actuation Strategy
		if actuatorProto := loadSchedulerProto.GetActuator(); actuatorProto != nil {
			actuator, actuatorOptions, err := loadscheduler.NewActuatorAndOptions(actuatorProto, componentID.String(), policyReadAPI, agentGroups)
			if err != nil {
				return retErr(err)
			}

			actuatorConfComp, err := prepareComponentInCircuit(actuator, actuatorProto, componentID.ChildID("Actuator"), parentCircuitID, true)
			if err != nil {
				return retErr(err)
			}
			configuredComponents = append(configuredComponents, actuatorConfComp)
			tree.Children = append(tree.Children, Tree{Node: actuatorConfComp})

			options = append(options, actuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(actuatorConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		loadSchedulerConfComp, err := prepareComponent(
			runtime.NewDummyComponent("LoadScheduler",
				iface.GetSelectorsShortDescription(loadSchedulerProto.GetSelectors()),
				runtime.ComponentTypeSignalProcessor),
			loadSchedulerProto,
			componentID,
			false,
		)
		if err != nil {
			return retErr(err)
		}

		loadSchedulerConfComp.PortMapping = portMapping
		tree.Node = loadSchedulerConfComp

		return tree, configuredComponents, fx.Options(options...), nil
	} else if isAdaptiveLoadScheduler {
		nestedCircuit, err := loadscheduler.ParseAdaptiveLoadScheduler(adaptiveLoadSchedulerProto)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
	} else if isLoadRamp {
		nestedCircuit, err := regulator.ParseLoadRamp(loadRampProto)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", flowControlComponentProto))
}
