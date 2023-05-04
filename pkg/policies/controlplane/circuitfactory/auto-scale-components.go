package circuitfactory

import (
	"fmt"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	policyprivatev1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/private/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/autoscale"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components/autoscale/podscaler"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// autoScaleModuleForPolicyApp for component factory run via the policy app. For singletons in the Policy scope.
func autoScaleModuleForPolicyApp(circuitAPI runtime.CircuitAPI) fx.Option {
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
		podScalerOptions, podScalerErr := podscaler.NewConfigSyncOptions(
			podScalerProto,
			componentID,
			policyReadAPI)
		if podScalerErr != nil {
			return retErr(podScalerErr)
		}
		options = append(options, podScalerOptions)

		nestedCircuit, err := podscaler.ParsePodScaler(podScalerProto, componentID, policyReadAPI)
		if err != nil {
			return retErr(err)
		}

		tree, configuredComponents, nestedOptions, err := ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
		if err != nil {
			return retErr(err)
		}
		options = append(options, nestedOptions)

		return tree, configuredComponents, fx.Options(options...), nil

	} else if autoScaler := autoScaleComponentProto.GetAutoScaler(); autoScaler != nil {
		nestedCircuit, err := autoscale.ParseAutoScaler(autoScaler)
		if err != nil {
			return retErr(err)
		}

		return ParseNestedCircuit(componentID, nestedCircuit, policyReadAPI)
	} else if private := autoScaleComponentProto.GetPrivate(); private != nil {
		switch private.TypeUrl {
		case "type.googleapis.com/aperture.policy.private.v1.PodScaleActuator":
			podScaleActuator := &policyprivatev1.PodScaleActuator{}
			if err := anypb.UnmarshalTo(private, podScaleActuator, proto.UnmarshalOptions{}); err != nil {
				return retErr(err)
			}
		}
	}
	return retErr(fmt.Errorf("unsupported/missing component type, proto: %+v", autoScaleComponentProto))
}
