local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local acceptedVsRejectedTargets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(accepted_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Accepted Token Rate'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(rejected_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rejected Token Rate'),
  ],

  local acceptedVsRejected = timeSeriesPanel('Accepted Token Rate vs Rejected Token Rate', cfg.dashboard.datasource.name, 'Token Rate', stringFilters, targets=acceptedVsRejectedTargets, h=8, w=12),

  panel: acceptedVsRejected.panel,
}
