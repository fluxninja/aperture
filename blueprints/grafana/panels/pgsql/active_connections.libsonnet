local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local activeConn = statPanel('Active Connections %',
                               cfg.dashboard.datasource.name,
                               'sum(postgresql_backends{%(filters)s,infra_meter_name="postgresql"}) / sum(postgresql_connection_max{%(filters)s,infra_meter_name="postgresql"}) * 100',
                               stringFilters),
  panel: activeConn.panel,
}
