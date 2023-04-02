package flowregulator

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/components"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
)

const (
	loadShaperForwardPortName          = "forward"
	loadShaperBackwardPortName         = "backward"
	loadShaperResetPortName            = "reset"
	loadShaperAcceptPercentagePortName = "accept_percentage"
	loadShaperStartSignalPortName      = "start_signal"
	loadShaperEndSignalPortName        = "end_signal"
)

// ParseLoadShaper parses a LoadShaper from the given proto and returns its nested circuit representation.
func ParseLoadShaper(loadShaper *policylangv1.LoadShaper) (*policylangv1.NestedCircuit, error) {
	nestedInPortsMap := make(map[string]*policylangv1.InPort)
	inPorts := loadShaper.InPorts
	if inPorts != nil {
		forwardPort := inPorts.Forward
		if forwardPort != nil {
			nestedInPortsMap[loadShaperForwardPortName] = forwardPort
		}
		backwardPort := inPorts.Backward
		if backwardPort != nil {
			nestedInPortsMap[loadShaperBackwardPortName] = backwardPort
		}
		resetPort := inPorts.Forward
		if resetPort != nil {
			nestedInPortsMap[loadShaperForwardPortName] = resetPort
		}
	}

	nestedOutPortsMap := make(map[string]*policylangv1.OutPort)
	outPorts := loadShaper.OutPorts
	if outPorts != nil {
		acceptPercentagePort := outPorts.AcceptPercentage
		if acceptPercentagePort != nil {
			nestedOutPortsMap[loadShaperAcceptPercentagePortName] = acceptPercentagePort
		}
		startSignalPort := outPorts.StartSignal
		if startSignalPort != nil {
			nestedOutPortsMap[loadShaperStartSignalPortName] = startSignalPort
		}
		endSignalPort := outPorts.EndSignal
		if endSignalPort != nil {
			nestedOutPortsMap[loadShaperEndSignalPortName] = endSignalPort
		}
	}

	// Populate the steps
	steps := make([]*policylangv1.SignalGenerator_Parameters_Step, 0)
	for _, step := range loadShaper.Parameters.Steps {
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
		Name:             "LoadShaper",
		ShortDescription: iface.GetServiceShortDescription(loadShaper.Parameters.FlowRegulatorParameters.FlowSelector.ServiceSelector),
		InPortsMap:       nestedInPortsMap,
		OutPortsMap:      nestedOutPortsMap,
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
							StartSignal: &policylangv1.OutPort{
								SignalName: "START_SIGNAL",
							},
							EndSignal: &policylangv1.OutPort{
								SignalName: "END_SIGNAL",
							},
						},
					},
				},
			},
			{
				Component: &policylangv1.Component_FlowControl{
					FlowControl: &policylangv1.FlowControl{
						Component: &policylangv1.FlowControl_FlowRegulator{
							FlowRegulator: &policylangv1.FlowRegulator{
								Parameters: &policylangv1.FlowRegulator_Parameters{
									FlowSelector: loadShaper.Parameters.FlowRegulatorParameters.FlowSelector,
								},
								InPorts: &policylangv1.FlowRegulator_Ins{
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

	components.AddNestedIngress(nestedCircuit, loadShaperForwardPortName, "FORWARD")
	components.AddNestedIngress(nestedCircuit, loadShaperBackwardPortName, "BACKWARD")
	components.AddNestedIngress(nestedCircuit, loadShaperResetPortName, "RESET")
	components.AddNestedEgress(nestedCircuit, loadShaperAcceptPercentagePortName, "ACCEPT_PERCENTAGE")
	components.AddNestedEgress(nestedCircuit, loadShaperStartSignalPortName, "START_SIGNAL")
	components.AddNestedEgress(nestedCircuit, loadShaperEndSignalPortName, "END_SIGNAL")

	return nestedCircuit, nil
}
