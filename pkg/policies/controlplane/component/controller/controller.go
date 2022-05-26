package controller

import (
	"aperture.tech/aperture/pkg/policies/controlplane/reading"
	"aperture.tech/aperture/pkg/policies/controlplane/runtime"
)

// Controller is the interface for the controller.
type Controller interface {
	// ComputeOutput: Compute the output given the current and previous signal readings. Refer to last values of other readings via ControlLoopReadAPI.
	ComputeOutput(signal, setpoint, controlVariable reading.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) (reading.Reading, error)
	// WindOutput: Wind the output given the previous and target output readings. Refer to last values of other readings via ControlLoopReadAPI.
	WindOutput(currentOutput, targetOutput reading.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) (reading.Reading, error)
	// MaintainOutput: Try to maintain the output on setpoint change. Refer to last values of other readings via ControlLoopReadAPI.
	MaintainOutput(prevSetpoint, currentSetpoint reading.Reading, controllerStateReadAPI ControllerStateReadAPI, tickInfo runtime.TickInfo) error
}
