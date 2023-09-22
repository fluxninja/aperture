local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local dbSize = statPanel('PGSQL Instances Size (GB)',
                           datasource.name,
                           'sum(postgresql_db_size_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName },
                           stringFilters),
  panel: dbSize.panel,
}
