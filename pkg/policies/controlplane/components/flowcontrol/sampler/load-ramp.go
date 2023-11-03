package sampler

import (
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

const (
	loadRampForwardPortName          = "forward"
	loadRampBackwardPortName         = "backward"
	loadRampResetPortName            = "reset"
	loadRampAcceptPercentagePortName = "accept_percentage"
	loadRampAtStartPortName          = "at_start"
	loadRampAtEndPortName            = "at_end"
)

// ParseLoadRamp parses a LoadRamp from the given proto and returns its nested circuit representation.
func ParseLoadRamp(
	loadRamp *policylangv1.LoadRamp,
	componentID runtime.ComponentID,
) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
	retErr := func(err error) (*runtime.ConfiguredComponent, *policylangv1.NestedCircuit, error) {
		return nil, nil, err
	}

	loadRamp.Parameters.Sampler.RampMode = true

	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := loadRamp.InPorts
	if inPorts != nil {
		forwardPort := inPorts.Forward
		if forwardPort != nil {
			nestedInPortsMap[loadRampForwardPortName] = forwardPort
		}
		backwardPort := inPorts.Backward
		if backwardPort != nil {
			nestedInPortsMap[loadRampBackwardPortName] = backwardPort
		}
		resetPort := inPorts.Reset_
		if resetPort != nil {
			nestedInPortsMap[loadRampResetPortName] = resetPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := loadRamp.OutPorts
	if outPorts != nil {
		acceptPercentagePort := outPorts.AcceptPercentage
		if acceptPercentagePort != nil {
			nestedOutPortsMap[loadRampAcceptPercentagePortName] = acceptPercentagePort
		}
		startSignalPort := outPorts.AtStart
		if startSignalPort != nil {
			nestedOutPortsMap[loadRampAtStartPortName] = startSignalPort
		}
		endSignalPort := outPorts.AtEnd
		if endSignalPort != nil {
			nestedOutPortsMap[loadRampAtEndPortName] = endSignalPort
		}
	}

	// Populate the steps
	steps := make([]*policylangv1.SignalGenerator_Parameters_Step, 0)
	for _, step := range loadRamp.Parameters.Steps {
		steps = append(steps, &policylangv1.SignalGenerator_Parameters_Step{
			Duration: step.Duration,
			TargetOutput: &policylangv1.ConstantSignal{
				Const: &policylangv1.ConstantSignal_Value{
					Value: step.TargetAcceptPercentage,
				},
			},
		})
	}

	nestedCircuit := &policylangv1.NestedCircuit{
		InPortsMap:  nestedInPortsMap,
		OutPortsMap: nestedOutPortsMap,
		Components: []*policylangv1.Component{
			{
				Component: &policylangv1.Component_SignalGenerator{
					SignalGenerator: &policylangv1.SignalGenerator{
						Parameters: &policylangv1.SignalGenerator_Parameters{
							Steps: steps,
						},
						InPorts: &policylangv1.SignalGenerator_Ins{
							Forward: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "FORWARD",
								},
							},
							Backward: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "BACKWARD",
								},
							},
							Reset_: &policylangv1.InPort{
								Value: &policylangv1.InPort_SignalName{
									SignalName: "RESET",
								},
							},
						},
						OutPorts: &policylangv1.SignalGenerator_Outs{
							Output: &policylangv1.OutPort{
								SignalName: "ACCEPT_PERCENTAGE",
							},
							AtStart: &policylangv1.OutPort{
								SignalName: "START_SIGNAL",
							},
							AtEnd: &policylangv1.OutPort{
								SignalName: "END_SIGNAL",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_Sampler{
							Sampler: &policylangv1.Sampler{
								Parameters: loadRamp.Parameters.Sampler,
								InPorts: &policylangv1.Sampler_Ins{
									AcceptPercentage: &policylangv1.InPort{
										Value: &policylangv1.InPort_SignalName{
											SignalName: "ACCEPT_PERCENTAGE",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	components.AddNestedIngress(nestedCircuit, loadRampForwardPortName, "FORWARD")
	components.AddNestedIngress(nestedCircuit, loadRampBackwardPortName, "BACKWARD")
	components.AddNestedIngress(nestedCircuit, loadRampResetPortName, "RESET")
	components.AddNestedEgress(nestedCircuit, loadRampAcceptPercentagePortName, "ACCEPT_PERCENTAGE")
	components.AddNestedEgress(nestedCircuit, loadRampAtStartPortName, "START_SIGNAL")
	components.AddNestedEgress(nestedCircuit, loadRampAtEndPortName, "END_SIGNAL")

	configuredComponent, err := runtime.NewConfiguredComponent(
		runtime.NewDummyComponent("LoadScheduler",
			iface.GetSelectorsShortDescription(loadRamp.Parameters.Sampler.GetSelectors()),
			runtime.ComponentTypeSignalProcessor),
		loadRamp,
		componentID,
		false,
	)
	if err != nil {
		return retErr(err)
	}

	return configuredComponent, nestedCircuit, nil
}
