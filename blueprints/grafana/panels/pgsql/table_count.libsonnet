local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local tbCount = statPanel('Number of tables',
                            datasource.name,
                            'sum(postgresql_table_count{%(filters)s,infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                            stringFilters),
  panel: tbCount.panel,
}
