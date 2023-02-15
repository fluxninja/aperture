package main

import (
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
)

func main() {
	info.Service = "aperturectl"

	logger := log.NewLogger(log.GetPrettyConsoleWriter(), "info")
	log.SetGlobalLogger(logger)
	cmd.Execute()
}
