local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local maxConnections = statPanel('Max Connections',
                                   cfg.dashboard.datasource.name,
                                   'postgresql_connection_max{%(filters)s,infra_meter_name="postgresql"} / postgresql_database_count',
                                   stringFilters),
  panel: maxConnections.panel,
}
