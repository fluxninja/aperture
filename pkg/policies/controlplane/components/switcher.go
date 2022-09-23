package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Switcher switches between two inputs based on third input.
type Switcher struct{}

// Make sure Switcher complies with Component interface.
var _ runtime.Component = (*Switcher)(nil)

// NewSwitcherAndOptions creates a new Switcher Component.
func NewSwitcherAndOptions(switcherProto *policylangv1.Switcher, componentIndex int, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	switcher := Switcher{}
	return &switcher, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (dec *Switcher) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	onTrue := inPortReadings.ReadSingleValuePort("on_true")
	onFalse := inPortReadings.ReadSingleValuePort("on_false")
	switchValue := inPortReadings.ReadSingleValuePort("switch")

	var output runtime.Reading

	if switchValue.Valid() && switchValue.Value() != 0 {
		output = onTrue
	} else {
		output = onFalse
	}

	return runtime.PortToValue{
		"output": []runtime.Reading{output},
	}, nil
}
