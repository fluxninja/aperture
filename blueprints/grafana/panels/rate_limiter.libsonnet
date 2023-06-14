local utils = import '../utils/policy_utils.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local rateLimiterPanel =
    g.panel.timeSeries.new('Aperture Rate Limiter')
    + g.panel.timeSeries.datasource.withType('prometheus')
    + g.panel.timeSeries.datasource.withUid(cfg.dashboard.datasource.name)
    + g.panel.timeSeries.withTargets([
      g.query.prometheus.new(cfg.dashboard.datasource.name,
                             |||
                               sum by(decision_type) (
                               rate(rate_limiter_counter_total{ %(filters)s}[$__rate_interval]))
                             ||| % { filters: stringFilters })
      + g.query.prometheus.withIntervalFactor(1),
    ])
    + g.panel.timeSeries.standardOptions.withUnit('reqps')
    + g.panel.timeSeries.fieldConfig.defaults.custom.withAxisLabel('Decisions')
    + g.panel.timeSeries.gridPos.withH(10)
    + g.panel.timeSeries.gridPos.withW(24),

  panel: rateLimiterPanel,
}
