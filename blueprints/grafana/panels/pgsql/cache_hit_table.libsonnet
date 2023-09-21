local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local currentInsert = statPanel('Table cache hit rate %',
                                  datasource.name,
                                  'sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_hit"}) / (sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_hit"}) + sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_read"}))' % { filters: stringFilters, infra_meter: infraMeterName },
                                  stringFilters),

  panel: currentInsert.panel,
}
