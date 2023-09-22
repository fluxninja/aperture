local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local activeConn = statPanel('Active Connections %',
                               datasource.name,
                               'sum(postgresql_backends{%(filters)s,infra_meter_name="%(infra_meter)s"}) / sum(postgresql_connection_max{%(filters)s,infra_meter_name="%(infra_meter)s"}) * 100' % { filters: stringFilters, infra_meter: infraMeterName },
                               stringFilters),
  panel: activeConn.panel,
}
