package controller

import "github.com/fluxninja/aperture/pkg/policies/controlplane/runtime"

// ControllerStateReadAPI is the interface to the Controller state.
type ControllerStateReadAPI interface {
	GetSignal() runtime.Reading
	GetSetpoint() runtime.Reading
	GetControlVariable() runtime.Reading
	GetControllerOutput() runtime.Reading
}
