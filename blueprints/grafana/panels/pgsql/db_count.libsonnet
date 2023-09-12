local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local dbCount = statPanel('PGSQL Instances',
                            cfg.dashboard.datasource.name,
                            'postgresql_database_count{%(filters)s}',
                            stringFilters),
  panel: dbCount.panel,
}
