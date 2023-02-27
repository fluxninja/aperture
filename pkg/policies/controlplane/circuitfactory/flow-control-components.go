package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/flowcontrol/concurrency"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// newFlowControlCompositeAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newFlowControlCompositeAndOptions(
	flowControlComponentProto *policylangv1.FlowControl,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
	retErr := func(err error) (Tree, []runtime.ConfiguredComponent, fx.Option, error) {
		return Tree{}, nil, nil, err
	}

	parentCircuitID, ok := componentID.ParentID()
	if !ok {
		return retErr(fmt.Errorf("parent circuit ID not found for component %s", componentID))
	}
	// Factory parser to determine what kind of composite component to create
	if concurrencyLimiterProto := flowControlComponentProto.GetConcurrencyLimiter(); concurrencyLimiterProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			tree                 Tree
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiterProto, componentID.String(), policyReadAPI)
		if concurrencyLimiterErr != nil {
			return retErr(concurrencyLimiterErr)
		}
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiterProto.GetScheduler(); schedulerProto != nil {
			// Use the same id as the component stack since agent sees only the component stack and generates metrics tagged with the component stack id
			scheduler, schedulerOptions, err := concurrency.NewSchedulerAndOptions(schedulerProto, componentID.String(), policyReadAPI, agentGroupName)
			if err != nil {
				return retErr(err)
			}

			// Need a unique ID for sub component since it's used for graph generation
			schedulerConfComp, err := prepareComponentInCircuit(scheduler, schedulerProto, componentID.ChildID("Scheduler"), parentCircuitID)
			if err != nil {
				return retErr(err)
			}

			configuredComponents = append(configuredComponents, schedulerConfComp)
			tree.Children = append(tree.Children, Tree{Root: schedulerConfComp})

			options = append(options, schedulerOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(schedulerConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		// Actuation Strategy
		if loadActuatorProto := concurrencyLimiterProto.GetLoadActuator(); loadActuatorProto != nil {
			loadActuator, loadActuatorOptions, err := concurrency.NewLoadActuatorAndOptions(loadActuatorProto, componentID.String(), policyReadAPI, agentGroupName)
			if err != nil {
				return retErr(err)
			}

			loadActuatorConfComp, err := prepareComponentInCircuit(loadActuator, loadActuatorProto, componentID.ChildID(".LoadActuator"), parentCircuitID)
			if err != nil {
				return retErr(err)
			}
			configuredComponents = append(configuredComponents, loadActuatorConfComp)
			tree.Children = append(tree.Children, Tree{Root: loadActuatorConfComp})

			options = append(options, loadActuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(loadActuatorConfComp.PortMapping)
			if err != nil {
				return retErr(err)
			}
		}

		concurrencyLimiterConfComp, err := prepareComponent(
			runtime.NewDummyComponent("ConcurrencyLimiter",
				iface.GetServiceShortDescription(concurrencyLimiterProto.FlowSelector.ServiceSelector),
				runtime.ComponentTypeSignalProcessor),
			concurrencyLimiterProto,
			componentID,
		)
		if err != nil {
			return retErr(err)
		}

		concurrencyLimiterConfComp.PortMapping = portMapping
		tree.Root = concurrencyLimiterConfComp

		return tree, configuredComponents, fx.Options(options...), nil
	} else if aimdConcurrencyController := flowControlComponentProto.GetAimdConcurrencyController(); aimdConcurrencyController != nil {
		nestedCircuit, err := concurrency.ParseAIMDConcurrencyController(aimdConcurrencyController)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", flowControlComponentProto))
}
