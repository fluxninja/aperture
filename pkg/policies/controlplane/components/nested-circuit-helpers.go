package components

import (
	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
)

// AddNestedIngress adds a nested ingress to the nested circuit.
func AddNestedIngress(nestedCircuit *policylangv1.NestedCircuit, portName, signalName string) {
	nestedCircuit.Components = append(nestedCircuit.Components, &policylangv1.Component{
		Component: &policylangv1.Component_NestedSignalIngress{
			NestedSignalIngress: &policylangv1.NestedSignalIngress{
				PortName: portName,
				OutPorts: &policylangv1.NestedSignalIngress_Outs{
					Signal: &policylangv1.OutPort{
						SignalName: signalName,
					},
				},
			},
		},
	})
}

// AddNestedEgress adds a nested egress to the nested circuit.
func AddNestedEgress(nestedCircuit *policylangv1.NestedCircuit, portName, signalName string) {
	nestedCircuit.Components = append(nestedCircuit.Components, &policylangv1.Component{
		Component: &policylangv1.Component_NestedSignalEgress{
			NestedSignalEgress: &policylangv1.NestedSignalEgress{
				PortName: portName,
				InPorts: &policylangv1.NestedSignalEgress_Ins{
					Signal: &policylangv1.InPort{
						Value: &policylangv1.InPort_SignalName{
							SignalName: signalName,
						},
					},
				},
			},
		},
	})
}
