local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local dbSize = statPanel('PGSQL Instances Size (GB)',
                           cfg.dashboard.datasource.name,
                           'postgresql_db_size_bytes{%(filters)s} / 1024 / 1024 / 1024',
                           stringFilters),
  panel: dbSize.panel,
}
