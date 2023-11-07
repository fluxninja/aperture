package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	loadscheduler "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/load-scheduler"
	ratelimiter "github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/rate-limiter"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/flowcontrol/sampler"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
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

	aimdLoadSchedulerProto := &policylangv1.AIMDLoadScheduler{}
	isAIMDLoadScheduler := false
	if proto := flowControlComponentProto.GetAimdLoadScheduler(); proto != nil {
		aimdLoadSchedulerProto = proto
		isAIMDLoadScheduler = true
	}

	rangeDrivenLoadSchedulerProto := &policylangv1.RangeDrivenLoadScheduler{}
	isRangeDrivenLoadScheduler := false
	if proto := flowControlComponentProto.GetRangeDrivenLoadScheduler(); proto != nil {
		rangeDrivenLoadSchedulerProto = proto
		isRangeDrivenLoadScheduler = true
	}

	aiadLoadSchedulerProto := &policylangv1.AIADLoadScheduler{}
	isAIADLoadScheduler := false
	if proto := flowControlComponentProto.GetAiadLoadScheduler(); proto != nil {
		aiadLoadSchedulerProto = proto
		isAIADLoadScheduler = true
	}

	loadRampProto := &policylangv1.LoadRamp{}
	isLoadRamp := false
	if proto := flowControlComponentProto.GetLoadRamp(); proto != nil {
		loadRampProto = proto
		isLoadRamp = true
	}

	rateLimiterProto := &policylangv1.RateLimiter{}
	isRateLimiter := false
	if proto := flowControlComponentProto.GetRateLimiter(); proto != nil {
		rateLimiterProto = proto
		isRateLimiter = true
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

		configuredComponent, nestedCircuit, err := loadscheduler.ParseLoadScheduler(loadSchedulerProto, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}

		tree, configuredComponents, nestedOptions, err := ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, nestedOptions)

		return tree, configuredComponents, fx.Options(options...), nil
	} else if isAdaptiveLoadScheduler {
		// convert to aimd load scheduler
		adaptiveLoadSchedulerProtoBytes, err := adaptiveLoadSchedulerProto.MarshalJSON()
		if err != nil {
			return retErr(err)
		}
		aimdProto := &policylangv1.AIMDLoadScheduler{}
		err = aimdProto.UnmarshalJSON(adaptiveLoadSchedulerProtoBytes)
		if err != nil {
			return retErr(err)
		}

		configuredComponent, nestedCircuit, err := loadscheduler.ParseAIMDLoadScheduler(aimdProto, componentID)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
	} else if isAIMDLoadScheduler {
		configuredComponent, nestedCircuit, err := loadscheduler.ParseAIMDLoadScheduler(aimdLoadSchedulerProto, componentID)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
	} else if isRangeDrivenLoadScheduler {
		configuredComponent, nestedCircuit, err := loadscheduler.ParseRangeDrivenLoadScheduler(rangeDrivenLoadSchedulerProto, componentID)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
	} else if isAIADLoadScheduler {
		configuredComponent, nestedCircuit, err := loadscheduler.ParseAIADLoadScheduler(aiadLoadSchedulerProto, componentID)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
	} else if isLoadRamp {
		configuredComponent, nestedCircuit, err := sampler.ParseLoadRamp(loadRampProto, componentID)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
	} else if isRateLimiter {
		var options []fx.Option
		// sync config
		configSyncOptions, err := ratelimiter.NewConfigSyncOptions(
			rateLimiterProto,
			componentID,
			policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, configSyncOptions)

		configuredComponent, nestedCircuit, err := ratelimiter.ParseRateLimiter(rateLimiterProto, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}

		tree, configuredComponents, nestedOptions, err := ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, nestedOptions)

		return tree, configuredComponents, fx.Options(options...), nil
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", flowControlComponentProto))
}
