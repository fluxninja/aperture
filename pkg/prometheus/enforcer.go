package prometheus

import (
	plp "github.com/prometheus-community/prom-label-proxy/injectproxy"
	promlabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"go.uber.org/fx"

	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/log"
	promconfig "github.com/fluxninja/aperture/v2/pkg/prometheus/config"
)

// EnforcerIn holds fields, parameters, to provide Prometheus Label Injector .
type EnforcerIn struct {
	fx.In
	Unmarshaller config.Unmarshaller
}

// PrometheusEnforcer is responsible for enforcing extra set of labels that
// should be present in every PromQL query executed.
type PrometheusEnforcer struct {
	enforcer *plp.Enforcer
}

// EnforceLabels transforms given query, making sure that all the required
// labels are in place.
func (e *PrometheusEnforcer) EnforceLabels(query string) (string, error) {
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return "", err
	}

	if e == nil {
		log.Error().Msg("e is not initialized?")
		return query, nil
	}

	if e.enforcer == nil {
		log.Error().Msg("e.enforcer is not initialized?")
		return query, nil
	}

	if err := e.enforcer.EnforceNode(expr); err != nil {
		return "", err
	}

	log.Debug().Str("query", expr.String()).Msg("Enforcing additional PromQL labels")

	return expr.String(), nil
}

func providePrometheusEnforcer(in EnforcerIn) (*PrometheusEnforcer, error) {
	var config promconfig.PrometheusConfig
	if err := in.Unmarshaller.UnmarshalKey(prometheusConfigKey, &config); err != nil {
		log.Error().Err(err).Msg("unable to deserialize")
		return nil, err
	}

	labels := []*promlabels.Matcher{}
	for _, label := range config.Labels {
		labels = append(labels, &promlabels.Matcher{
			Name:  label.Name,
			Type:  promlabels.MatchEqual,
			Value: label.Value,
		})
	}

	enforcer := plp.NewEnforcer(false, labels...)

	log.Info().Msg("Initializing prometheus labels exporter")

	return &PrometheusEnforcer{enforcer: enforcer}, nil
}