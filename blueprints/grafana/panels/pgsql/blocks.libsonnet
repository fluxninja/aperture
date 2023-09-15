local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg, title) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local targets = [
    g.query.prometheus.new(cfg.dashboard.datasource.name, 'rate(postgresql_blocks_read_total{%(filters)s, source="heap_read"}[5m])' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Heap Read'),

    g.query.prometheus.new(cfg.dashboard.datasource.name, 'rate(postgresql_blocks_read_total{%(filters)s, source="idx_read"}[5m])' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Index Read'),
  ],

  local blockReads = timeSeriesPanel(title, cfg.dashboard.datasource.name, '', stringFilters, targets=targets),
  panel: blockReads.panel,
}
