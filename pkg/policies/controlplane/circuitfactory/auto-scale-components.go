package circuitfactory

import (
	"fmt"

	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/autoscale"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components/autoscale/podscaler"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// autoScaleModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func autoScaleModuleForPolicyApp(circuitAPI runtime.CircuitSuperAPI) fx.Option {
	return fx.Options(
		podscaler.Module(),
	)
}

// newIntegrationCompositeAndOptions creates parent and leaf components and their fx options for a component stack spec.
func newAutoScaleNestedAndOptions(
	autoScaleComponentProto *policylangv1.AutoScale,
	componentID runtime.ComponentID,
	policyReadAPI iface.Policy,
) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
	retErr := func(err error) (Tree, []*runtime.ConfiguredComponent, fx.Option, error) {
		return Tree{}, nil, nil, err
	}

	if podScalerProto := autoScaleComponentProto.GetPodScaler(); podScalerProto != nil {
		var options []fx.Option
		// sync config
		podScalerOptions, podScalerErr := podscaler.NewConfigSyncOptions(podScalerProto, componentID, policyReadAPI)
		if podScalerErr != nil {
			return retErr(podScalerErr)
		}
		options = append(options, podScalerOptions)

		configuredComponent, nestedCircuit, err := podscaler.ParsePodScaler(podScalerProto, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}

		tree, configuredComponents, nestedOptions, err := ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, nestedOptions)

		return tree, configuredComponents, fx.Options(options...), nil
	} else if autoScaler := autoScaleComponentProto.GetAutoScaler(); autoScaler != nil {
		configuredComponent, nestedCircuit, err := autoscale.ParseAutoScaler(autoScaler, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(configuredComponent, nestedCircuit, componentID, policyReadAPI)
	}

	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", autoScaleComponentProto))
}
