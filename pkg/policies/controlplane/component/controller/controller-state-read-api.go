package controller

import "github.com/FluxNinja/aperture/pkg/policies/controlplane/reading"

// ControllerStateReadAPI is the interface to the Controller state.
type ControllerStateReadAPI interface {
	GetSignal() reading.Reading
	GetSetpoint() reading.Reading
	GetControlVariable() reading.Reading
	GetControllerOutput() reading.Reading
}
