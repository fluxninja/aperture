package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/actuators/concurrency"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// newFlowControlCompositeAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newFlowControlCompositeAndOptions(
	flowControlComponentProto *policylangv1.FlowControl,
	componentID string,
	policyReadAPI iface.Policy,
) ([]runtime.ConfiguredComponent, []runtime.ConfiguredComponent, fx.Option, error) {
	parentCircuitID := ParentCircuitID(componentID)
	// Factory parser to determine what kind of composite component to create
	if concurrencyLimiterProto := flowControlComponentProto.GetConcurrencyLimiter(); concurrencyLimiterProto != nil {
		var (
			configuredComponents []runtime.ConfiguredComponent
			options              []fx.Option
		)
		portMapping := runtime.NewPortMapping()
		concurrencyLimiterOptions, agentGroupName, concurrencyLimiterErr := concurrency.NewConcurrencyLimiterOptions(concurrencyLimiterProto, componentID, policyReadAPI)
		if concurrencyLimiterErr != nil {
			return nil, nil, nil, concurrencyLimiterErr
		}
		options = append(options, concurrencyLimiterOptions)

		// Scheduler
		if schedulerProto := concurrencyLimiterProto.GetScheduler(); schedulerProto != nil {
			// Use the same id as the component stack since agent sees only the component stack and generates metrics tagged with the component stack id
			scheduler, schedulerOptions, err := concurrency.NewSchedulerAndOptions(schedulerProto, componentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			// Need a unique ID for sub component since it's used for graph generation
			schedulerConfComp, err := prepareComponentInCircuit(scheduler, schedulerProto, componentID+".Scheduler", parentCircuitID)
			if err != nil {
				return nil, nil, nil, err
			}

			configuredComponents = append(configuredComponents, schedulerConfComp)

			options = append(options, schedulerOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(schedulerConfComp.PortMapping)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		// Actuation Strategy
		if loadActuatorProto := concurrencyLimiterProto.GetLoadActuator(); loadActuatorProto != nil {
			loadActuator, loadActuatorOptions, err := concurrency.NewLoadActuatorAndOptions(loadActuatorProto, componentID, policyReadAPI, agentGroupName)
			if err != nil {
				return nil, nil, nil, err
			}

			loadActuatorConfComp, err := prepareComponentInCircuit(loadActuator, loadActuatorProto, componentID+".LoadActuator", parentCircuitID)
			if err != nil {
				return nil, nil, nil, err
			}
			configuredComponents = append(configuredComponents, loadActuatorConfComp)

			options = append(options, loadActuatorOptions)

			// Merge port mapping for parent component
			err = portMapping.Merge(loadActuatorConfComp.PortMapping)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		concurrencyLimiterConfComp, err := prepareComponent(
			runtime.NewDummyComponent("ConcurrencyLimiter", runtime.ComponentTypeSignalProcessor),
			concurrencyLimiterProto,
			componentID,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		concurrencyLimiterConfComp.PortMapping = portMapping

		return []runtime.ConfiguredComponent{concurrencyLimiterConfComp}, configuredComponents, fx.Options(options...), nil
	} else if aimdConcurrencyController := flowControlComponentProto.GetAimdConcurrencyController(); aimdConcurrencyController != nil {
		return ParseAIMDConcurrencyController(componentID, aimdConcurrencyController, policyReadAPI)
	}
	return nil, nil, nil, fmt.Errorf("unsupported/missing component type, proto: %+v", flowControlComponentProto)
}
