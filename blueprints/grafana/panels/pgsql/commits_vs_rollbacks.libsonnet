local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg, title) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local targets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'rate(postgresql_commits_total{%(filters)s, infra_meter_name="postgresql"}[$__rate_interval])' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Commits'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'rate(postgresql_rollbacks_total{%(filters)s,infra_meter_name="postgresql"}[$__rate_interval])' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rollbacks'),
  ],

  local commitVsRollback = timeSeriesPanel(title, cfg.dashboard.datasource.name, '', stringFilters, h=8, w=10, targets=targets),
  panel: commitVsRollback.panel,
}
