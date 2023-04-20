package main

import (
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd"
	"github.com/fluxninja/aperture/pkg/info"
)

func main() {
	info.Service = "aperturectl"
	cmd.Execute()
}
