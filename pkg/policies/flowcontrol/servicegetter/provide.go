package servicegetter

import "go.uber.org/fx"

// Module is a set of default providers for servicegetter components.
var Module = fx.Options(
	fx.Provide(ProvideFromEntityCache),
	fx.Provide(NewMetrics),
)
