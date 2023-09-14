local barGaugePanel = import '../../utils/bar_gauge_panel.libsonnet';
local utils = import '../../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local indexSize = barGaugePanel('Index Size (MB)',
                                  cfg.dashboard.datasource.name,
                                  'sum(postgresql_index_size_bytes{%(filters)s}) / (1024 * 1024)',
                                  stringFilters),

  panel: indexSize.panel,
}
