local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local tbSize = statPanel('Tables Size (MB) ',
                           datasource.name,
                           'sum(postgresql_table_size_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}) / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName },
                           stringFilters),
  panel: tbSize.panel,
}
