package components

import (
	"go.uber.org/fx"

	policylangv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/notifiers"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/iface"
	"github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"
)

// Variable is a dynamically configurable variable signal.
type Variable struct {
	constantSignal *runtime.ConstantSignal
	policyReadAPI  iface.Policy
	variableProto  *policylangv1.Variable
}

// Name implements runtime.Component.
func (*Variable) Name() string { return "Variable" }

// Type implements runtime.Component.
func (*Variable) Type() runtime.ComponentType { return runtime.ComponentTypeSource }

// ShortDescription implements runtime.Component.
func (v *Variable) ShortDescription() string {
	return v.constantSignal.Description()
}

// IsActuator implements runtime.Component.
func (*Variable) IsActuator() bool { return false }

// NewConstantSignal creates a variable component with a value that is always valid.
func NewConstantSignal(value float64) runtime.Component {
	return &Variable{
		constantSignal: &runtime.ConstantSignal{
			Value: value,
		},
	}
}

// NewVariableAndOptions creates a variable components and its fx options.
func NewVariableAndOptions(variableProto *policylangv1.Variable, _ string, policyReadAPI iface.Policy) (runtime.Component, fx.Option, error) {
	variable := &Variable{
		policyReadAPI: policyReadAPI,
		variableProto: variableProto,
	}
	variable.constantSignal = runtime.ConstantSignalFromProto(variableProto.GetDefaultConfig().ConstantSignal)

	return variable, fx.Options(), nil
}

// Execute implements runtime.Component.Execute.
func (v *Variable) Execute(inPortReadings runtime.PortToReading, tickInfo runtime.TickInfo) (runtime.PortToReading, error) {
	// Always emit the value.
	return runtime.PortToReading{
		"output": []runtime.Reading{runtime.NewReading(v.constantSignal.Float())},
	}, nil
}

// DynamicConfigUpdate finds the dynamic config and syncs the constant value.
func (v *Variable) DynamicConfigUpdate(event notifiers.Event, unmarshaller config.Unmarshaller) {
	logger := v.policyReadAPI.GetStatusRegistry().GetLogger()
	key := v.variableProto.GetDynamicConfigKey()
	// read dynamic config
	if unmarshaller.IsSet(key) {
		dynamicConfig := &policylangv1.Variable_DynamicConfig{}
		if err := unmarshaller.UnmarshalKey(key, dynamicConfig); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal dynamic config")
			return
		}
		v.constantSignal = runtime.ConstantSignalFromProto(dynamicConfig.ConstantSignal)
	} else {
		v.constantSignal = runtime.ConstantSignalFromProto(v.variableProto.GetDefaultConfig().ConstantSignal)
	}
}
