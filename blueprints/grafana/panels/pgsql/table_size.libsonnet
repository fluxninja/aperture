local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local tbSize = statPanel('Tables Size (MB) ',
                           cfg.dashboard.datasource.name,
                           'sum(postgresql_table_size_bytes{%(filters)s,infra_meter_name="postgresql"}) / 1024 / 1024',
                           stringFilters),
  panel: tbSize.panel,
}
