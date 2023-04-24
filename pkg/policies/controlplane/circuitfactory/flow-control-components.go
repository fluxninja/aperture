package circuitfactory

import (
	"encoding/json"
	"fmt"

	"go.uber.org/fx"
	"google.golang.org/protobuf/encoding/protojson"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/flowcontrol/loadregulator"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/flowcontrol/loadscheduler"
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
	} else if proto := flowControlComponentProto.GetConcurrencyLimiter(); proto != nil {
		// Convert from *policylangv1.FlowControl_ConcurrencyLimiter to *policylangv1.FlowControl_LoadScheduler since they have the same fields
		jsonStr, err := json.Marshal(proto)
		if err != nil {
			return Tree{}, nil, nil, fmt.Errorf("error marshaling ConcurrencyLimiter to JSON: %v", err)
		}

		err = protojson.Unmarshal(jsonStr, loadSchedulerProto)
		if err != nil {
			return Tree{}, nil, nil, fmt.Errorf("error unmarshaling JSON to LoadScheduler: %v", err)
		}
		isLoadScheduler = true
	}

	// Factory parser to determine what kind of composite component to create
	if isLoadScheduler {
		var (
			configuredComponents []*runtime.ConfiguredComponent
			tree                 Tree
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		loadSchedulerOptions, agentGroupName, loadSchedulerErr := loadscheduler.NewLoadSchedulerOptions(loadSchedulerProto, componentID.String(), policyReadAPI)
		if loadSchedulerErr != nil {
			return retErr(loadSchedulerErr)
		}
		options = append(options, loadSchedulerOptions)

		// Scheduler
		if schedulerProto := loadSchedulerProto.GetScheduler(); schedulerProto != nil {
			// Use the same id as the component stack since agent sees only the component stack and generates metrics tagged with the component stack id
			scheduler, schedulerOptions, err := loadscheduler.NewSchedulerAndOptions(schedulerProto, componentID.String(), policyReadAPI, agentGroupName)
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
		if loadActuatorProto := loadSchedulerProto.GetLoadActuator(); loadActuatorProto != nil {
			loadActuator, loadActuatorOptions, err := loadscheduler.NewLoadActuatorAndOptions(loadActuatorProto, componentID.String(), policyReadAPI, agentGroupName)
			if err != nil {
				return retErr(err)
			}

			loadActuatorConfComp, err := prepareComponentInCircuit(loadActuator, loadActuatorProto, componentID.ChildID("LoadActuator"), parentCircuitID, true)
			if err != nil {
				return retErr(err)
			}
			configuredComponents = append(configuredComponents, loadActuatorConfComp)
			tree.Children = append(tree.Children, Tree{Node: loadActuatorConfComp})

			options = append(options, loadActuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(loadActuatorConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		loadSchedulerConfComp, err := prepareComponent(
			runtime.NewDummyComponent("LoadScheduler",
				iface.GetServiceShortDescription(loadSchedulerProto.FlowSelector.ServiceSelector),
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
	} else if aimdConcurrencyController := flowControlComponentProto.GetAimdConcurrencyController(); aimdConcurrencyController != nil {
		nestedCircuit, err := loadscheduler.ParseAIMDConcurrencyController(aimdConcurrencyController)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
	} else if loadShaper := flowControlComponentProto.GetLoadShaper(); loadShaper != nil {
		nestedCircuit, err := loadregulator.ParseLoadShaper(loadShaper)
		if err != nil {
			return retErr(err)
		}

		tree, configuredComponents, options, err := ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
		return tree, configuredComponents, options, err
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", flowControlComponentProto))
}
