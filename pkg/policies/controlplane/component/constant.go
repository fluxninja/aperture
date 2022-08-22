package component

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/reading"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Constant is a constant signal.
type Constant struct {
	// The value of the constant setpoint.
	value float64
}

// NewConstantAndOptions creates constant setpoint and its fx options.
func NewConstantAndOptions(constant *policylangv1.Constant, componentIndex int, policyReadAPI iface.PolicyRead) (runtime.Component, fx.Option, error) {
	con := Constant{
		value: constant.Value,
	}
	return &con, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (con *Constant) Execute(inPortReadings runtime.PortToValue, tickInfo runtime.TickInfo) (runtime.PortToValue, error) {
	// Always emit the value.
	return runtime.PortToValue{
		"output": []reading.Reading{reading.New(con.value)},
	}, nil
}
