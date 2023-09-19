local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local tbCount = statPanel('Number of tables',
                            cfg.dashboard.datasource.name,
                            'sum(postgresql_table_count{%(filters)s,infra_meter_name="postgresql"})',
                            stringFilters),
  panel: tbCount.panel,
}
