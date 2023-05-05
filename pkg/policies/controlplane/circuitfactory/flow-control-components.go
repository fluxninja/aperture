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

// newFlowControlNestedAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newFlowControlNestedAndOptions(
	flowControlComponentProto *policylangv1.FlowControl,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
	retErr := func(err error) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
		return Tree{}, nil, nil, err
	}

	loadSchedulerProto := &policylangv1.LoadScheduler{}
	isLoadScheduler := false
	if proto := flowControlComponentProto.GetLoadScheduler(); proto != nil {
		loadSchedulerProto = proto
		isLoadScheduler = true
	}

	adaptiveLoadSchedulerProto := &policylangv1.AdaptiveLoadScheduler{}
	isAdaptiveLoadScheduler := false
	if proto := flowControlComponentProto.GetAdaptiveLoadScheduler(); proto != nil {
		adaptiveLoadSchedulerProto = proto
		isAdaptiveLoadScheduler = true
	}

	loadRampProto := &policylangv1.LoadRamp{}
	isLoadRamp := false
	if proto := flowControlComponentProto.GetLoadRamp(); proto != nil {
		loadRampProto = proto
		isLoadRamp = true
	}

	// Factory parser to determine what kind of composite component to create
	if isLoadScheduler {
		var options []fx.Option
		// sync config
		configSyncOptions, err := loadscheduler.NewConfigSyncOptions(loadSchedulerProto,
			componentID,
			policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, configSyncOptions)

		nestedCircuit, err := loadscheduler.ParseLoadScheduler(loadSchedulerProto, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}

		tree, configuredComponents, nestedOptions, err := ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, nestedOptions)

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
