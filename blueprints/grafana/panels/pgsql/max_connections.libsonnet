local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local maxConnections = statPanel('Max Connections',
                                   datasource.name,
                                   'postgresql_connection_max{%(filters)s,infra_meter_name="%(infra_meter)s"} / postgresql_database_count' % { filters: stringFilters, infra_meter: infraMeterName },
                                   stringFilters),
  panel: maxConnections.panel,
}
