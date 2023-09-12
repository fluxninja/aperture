local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local tbVacuumTot = statPanel('Table Vacuum Total',
                                cfg.dashboard.datasource.name,
                                'postgresql_table_vacuum_total{%(filters)s}',
                                stringFilters),

  panel: tbVacuumTot.panel,
}
