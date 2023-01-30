package components

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
)

// NestedSignalPortName is the port name for nested signal ingress and egress components.
const NestedSignalPortName = "signal"

// nestedSignal is the base component for signal ingress and egress.
type nestedSignal struct {
	NoOp
	portName string
}

// Make sure NoOp complies with Component interface.
var _ runtime.Component = (*nestedSignal)(nil)

// PortName returns the port name.
func (nsi *nestedSignal) PortName() string {
	return nsi.portName
}

// NestedSignalIngress is a component that ingresses a signal into a nested circuit.
type NestedSignalIngress struct {
	nestedSignal
}

// NewNestedSignalIngressAndOptions creates a new NestedSignalIngress and its options.
func NewNestedSignalIngressAndOptions(nestedSignalIngressProto *policylangv1.NestedSignalIngress, _ string, _ iface.Policy) (runtime.Component, fx.Option, error) {
	comp := NestedSignalIngress{}
	comp.portName = nestedSignalIngressProto.PortName
	return &comp, fx.Options(), nil
}

// Name returns the name of the component.
func (nsi *NestedSignalIngress) Name() string {
	return "NestedSignalIngress"
}

// NestedSignalEgress is a component that ingresses a signal into a nested circuit.
type NestedSignalEgress struct {
	nestedSignal
}

// NewNestedSignalEgressAndOptions creates a new NestedSignalEgress and its options.
func NewNestedSignalEgressAndOptions(nestedSignalEgressProto *policylangv1.NestedSignalEgress, _ string, _ iface.Policy) (runtime.Component, fx.Option, error) {
	comp := NestedSignalEgress{}
	comp.portName = nestedSignalEgressProto.PortName
	return &comp, fx.Options(), nil
}

// Name returns the name of the component.
func (nse *NestedSignalEgress) Name() string {
	return "NestedSignalEgress"
}
