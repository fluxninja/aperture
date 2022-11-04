package fluxmeter

import "go.uber.org/fx"

// Module returns the fx options for dataplane side pieces of flux meter.
func Module() fx.Option {
	return fluxMeterModule()
}
