package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// BoolVariable is a dynamically configurable boolean variable signal.
type BoolVariable struct {
	policyReadAPI  iface.Policy
	variableProto  *policylangv1.BoolVariable
	constantOutput bool
}

// Name returns the name of the BoolVariable component.
func (*BoolVariable) Name() string { return "BoolVariable" }

// Type returns the type of the BoolVariable component.
func (*BoolVariable) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription returns a short description of the BoolVariable component.
func (v *BoolVariable) ShortDescription() string {
	return "Boolean Variable"
}

// IsActuator returns whether the BoolVariable component is an actuator.
func (*BoolVariable) IsActuator() bool { return false }

// NewBoolVariableAndOptions creates a new BoolVariable component and its fx options.
func NewBoolVariableAndOptions(variableProto *policylangv1.BoolVariable, _ runtime.ComponentID, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	boolVariable := &BoolVariable{
		policyReadAPI:  policyReadAPI,
		variableProto:  variableProto,
		constantOutput: variableProto.GetConstantOutput(),
	}

	return boolVariable, fx.Options(), nil
}

// Execute executes the BoolVariable component and emits the current boolean value.
func (v *BoolVariable) Execute(inPortReadings runtime.PortToReading, circuitAPI runtime.CircuitAPI) (runtime.PortToReading, error) {
	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewBoolReading(v.constantOutput)},
	}, nil
}

// DynamicConfigUpdate updates the boolean value when the configuration changes.
func (v *BoolVariable) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	key := v.variableProto.GetConfigKey()
	v.constantOutput = config.GetBoolValue(unmarshaller, key, v.variableProto.GetConstantOutput())
}
