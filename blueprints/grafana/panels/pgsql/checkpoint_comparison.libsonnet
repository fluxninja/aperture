local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg, title) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local targets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s,type="requested"}[5m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Requested'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'sum(rate(postgresql_bgwriter_checkpoint_count_total{%(filters)s,type="scheduled"}[5m]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Scheduled'),
  ],

  local checkpointComparison = timeSeriesPanel(title, cfg.dashboard.datasource.name, '', stringFilters, targets=targets),
  panel: checkpointComparison.panel,
}
