local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg, title) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local targets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(postgresql_commits_total{%(filters)s}[5m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Commits'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(postgresql_backends{%(filters)s}[5m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Backends'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(postgresql_rollbacks_total{%(filters)s}[5m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rollbacks'),
  ],

  local totalCommits = timeSeriesPanel(title, cfg.dashboard.datasource.name, '', stringFilters, h=8, w=10, targets=targets),
  panel: totalCommits.panel,
}
