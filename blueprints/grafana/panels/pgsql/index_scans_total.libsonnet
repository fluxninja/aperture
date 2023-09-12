local barGaugePanel = import '../../utils/bar_gauge_panel.libsonnet';
local utils = import '../../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local indexScansTotal = barGaugePanel('Index Scans',
                                        cfg.dashboard.datasource.name,
                                        'avg(postgresql_index_scans_total {%(filters)s})',
                                        stringFilters),

  panel: indexScansTotal.panel,
}
