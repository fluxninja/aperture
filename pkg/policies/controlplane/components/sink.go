package components

import (
	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
	"go.uber.org/fx"
)

// Sink is a component that consumes array of input signals and does nothing with them.
type Sink struct{}

// Make sure Sink implements Component interface.
var _ runtime.Component = (*Sink)(nil)

// NewSinkAndOptions creates a new Sink component.
func NewSinkAndOptions(sinkProto *policylangv1.Sink, componentIndex int, policyReadAPI iface.Policy) (*Sink, fx.Option, error) {
	return &Sink{}, fx.Options(), nil
}

// Execute implements Component interface.
func (s *Sink) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	return nil, nil
}

// DynamicConfigUpdate is a no-op for Sink.
func (s *Sink) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {}
