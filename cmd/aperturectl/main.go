package main

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd"
	"github.com/fluxninja/aperture/v2/pkg/info"
)

func main() {
	info.Service = "aperturectl"
	cmd.Execute()
}
