package alertmanager

import (
	"go.uber.org/fx"

	amclient "github.com/fluxninja/aperture/pkg/alertmanager/client"
)

// Module returns an fx.Option that provides the alertmanager module.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideAlertManager,
			amclient.ProvideNamedAlertManagerClients,
		),
		fx.Invoke(
			setupAlertManager,
		),
	)
}

// AlertManager is a struct that aggregates all of the alert manager clients.
type AlertManager struct {
	Clients []amclient.AlertManagerClient
}

// ProvideAlertManager creates empty AlertManager.
func ProvideAlertManager() *AlertManager {
	return &AlertManager{}
}

func setupAlertManager(clientSlice []amclient.AlertManagerClient, alertMgr *AlertManager) {
	alertMgr.Clients = clientSlice
}
