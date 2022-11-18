package alertmgrclient

import (
	"net/http"

	"github.com/prometheus/alertmanager/api/v2/client"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
)

var (
	// AlertMgrClientConfigKey is the key used to store the AlertManagerClientConfig in the config.
	configKey           = "alertmanagers"
	AlertMgrClientFxKey = config.GroupTag("alertmgr_clients")
)

// ProvideNamedAlertManagerClient provides a AlertManagerClient with named http_client.
func ProvideNamedClient(clientName string, httpConfig *commonhttp.HTTPClientConfig) fx.Option {
	return fx.Options(
		commonhttp.ClientConstructor{Name: clientName, ProvidedConfig: httpConfig}.Annotate(),
		fx.Provide(
			fx.Annotate(
				ProvideAlertMgrClient,
				fx.ParamTags(config.NameTag(clientName)),
				fx.ResultTags(AlertMgrClientFxKey),
			)),
	)
}

// AlertManagerClient provides an interface for alert manager client.
type AlertManagerClient interface {
}

// RealAlertManagerClient implements AlertManagerClient interface.
type RealAlertManagerClient struct {
	promAlertClient *client.Alertmanager
}

// ProvideAlertMgrClient provides a new alertmanager client and sets logger.
func ProvideAlertMgrClient(httpClient *http.Client) AlertManagerClient {
	transportCfg := &client.TransportConfig{}
	promClient := client.NewHTTPClientWithConfig(nil, transportCfg)

	alertMgrClient := &RealAlertManagerClient{
		promAlertClient: promClient,
	}
	return alertMgrClient
}
