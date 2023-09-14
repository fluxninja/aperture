local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg, title) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local targets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'avg(rate(postgresql_operations_total{%(filters)s}[1m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Insert'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'avg(rate(postgresql_operations_total{%(filters)s}[1m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Delete'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'avg(rate(postgresql_operations_total{%(filters)s}[1m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Update'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'avg(rate(postgresql_operations_total{%(filters)s}[1m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Hot Update'),
  ],

  local operations = timeSeriesPanel(title, cfg.dashboard.datasource.name, '', stringFilters, targets),
  panel: operations.panel,
}
