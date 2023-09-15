local barGaugePanel = import '../../utils/bar_gauge_panel.libsonnet';
local utils = import '../../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local topTables = barGaugePanel('Tables with most disk usage',
                                  cfg.dashboard.datasource.name,
                                  'topk(5, sum by (postgresql_table_name) (postgresql_table_size_bytes{%(filters)s}))',
                                  stringFilters),
  panel: topTables.panel,
}
