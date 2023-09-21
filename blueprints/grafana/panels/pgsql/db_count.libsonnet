local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local dbCount = statPanel('PGSQL Instances',
                            datasource,
                            'count(postgresql_database_count{%(filters)s,infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                            policyName),
  panel: dbCount.panel,
}
