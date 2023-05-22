package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/notifiers"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/v2/pkg/policies/controlplane/runtime"
)

// Variable is a dynamically configurable variable signal.
type Variable struct {
	constantOutput *runtime.ConstantSignal
	policyReadAPI  iface.Policy
	variableProto  *policylangv1.Variable
}

// Name implements runtime.Component.
func (*Variable) Name() string { return "Variable" }

// Type implements runtime.Component.
func (*Variable) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription implements runtime.Component.
func (v *Variable) ShortDescription() string {
	return v.constantOutput.Description()
}

// IsActuator implements runtime.Component.
func (*Variable) IsActuator() bool { return false }

// NewConstantSignal creates a variable component with a value that is always valid.
func NewConstantSignal(value float64) runtime.Component {
	return &Variable{
		constantOutput: &runtime.ConstantSignal{
			Value: value,
		},
	}
}

// NewVariableAndOptions creates a variable components and its fx options.
func NewVariableAndOptions(variableProto *policylangv1.Variable, _ runtime.ComponentID, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	variable := &Variable{
		policyReadAPI: policyReadAPI,
		variableProto: variableProto,
	}
	variable.constantOutput = runtime.ConstantSignalFromProto(variableProto.GetConstantOutput())

	return variable, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (v *Variable) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	// Always emit the value.
	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewReading(v.constantOutput.Float())},
	}, nil
}

// DynamicConfigUpdate finds the dynamic config and syncs the constant value.
func (v *Variable) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := v.policyReadAPI.GetStatusRegistry().GetLogger()
	key := v.variableProto.GetConfigKey()
	// read dynamic config
	if !unmarshaller.IsSet(key) {
		v.constantOutput = runtime.ConstantSignalFromProto(v.variableProto.GetConstantOutput())
	}
	constantSignalProto := &policylangv1.ConstantSignal{}
	if err := unmarshaller.UnmarshalKey(key, constantSignalProto); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
		v.constantOutput = runtime.ConstantSignalFromProto(v.variableProto.GetConstantOutput())
		return
	}
	v.constantOutput = runtime.ConstantSignalFromProto(constantSignalProto)
}
