package prometheus

import (
	"errors"
	"net/http"
	"strings"

	promapi "github.com/prometheus/client_golang/api"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/log"
	commonhttp "github.com/fluxninja/aperture/pkg/net/http"
)

var (
	// swagger:operation POST /prometheus common-configuration Prometheus
	// ---
	// x-fn-config-env: true
	// parameters:
	// - in: body
	//   schema:
	//     $ref: "#/definitions/PrometheusConfig"
	// - name: http_client
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/HTTPClientConfig"

	// PrometheusConfigKey is the key used to store the PrometheusConfig in the config.
	prometheusConfigKey = "prometheus"
	// HttpConfigKey is the key used to store the HTTPClientConfig in the config.
	httpConfigKey = strings.Join([]string{prometheusConfigKey, "http_client"}, ".")
)

// PrometheusConfig holds configuration for Prometheus Server.
// swagger:model
// +kubebuilder:object:generate=true
type PrometheusConfig struct {
	// Address of the prometheus server
	//+kubebuilder:validation:Required
	Address string `json:"address" validate:"required,hostname_port|url|fqdn"`
}

// Module provides a singleton pointer to prometheusv1.API via FX.
func Module() fx.Option {
	return fx.Options(
		fx.Provide(providePrometheusClient),
		commonhttp.ClientConstructor{Name: "prometheus.http-client", Key: httpConfigKey}.Annotate(),
	)
}

// ClientIn holds fields, parameters, to provide Prometheus Client.
type ClientIn struct {
	fx.In
	HTTPClient   *http.Client `name:"prometheus.http-client"`
	Unmarshaller config.Unmarshaller
}

func providePrometheusClient(in ClientIn) (prometheusv1.API, promapi.Client, error) {
	var config PrometheusConfig
	if err := in.Unmarshaller.UnmarshalKey(prometheusConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("unable to deserialize")
		return nil, nil, err
	}
	if config.Address == "" {
		err := errors.New("prometheus address not specified")
		log.Error().Err(err).Msg("")
		return nil, nil, err
	}
	client, err := promapi.NewClient(promapi.Config{
		Address:      config.Address,
		RoundTripper: in.HTTPClient.Transport,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error creating client")
		return nil, nil, err
	}

	return prometheusv1.NewAPI(client), client, nil
}
