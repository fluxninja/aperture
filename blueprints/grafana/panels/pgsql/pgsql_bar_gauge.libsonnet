local barGaugePanel = import '../../utils/bar_gauge_panel.libsonnet';
local utils = import '../../utils/policy_utils.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local topDiskUsageTables = barGaugePanel('Tables with most disk usage',
                                           datasource,
                                           'topk(5, sum by (postgresql_table_name,postgresql_database_name) (postgresql_table_size_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}))' % { filters: stringFilters, infra_meter: infraMeterName },
                                           stringFilters),

  local topLiveRowsTables = barGaugePanel('Tables with most live rows',
                                          datasource,
                                          'topk(5, sum by (postgresql_table_name,postgresql_database_name) (postgresql_rows{%(filters)s,infra_meter_name="%(infra_meter)s",state="live"}))' % { filters: stringFilters, infra_meter: infraMeterName },
                                          stringFilters),

  topDiskUsageTables: topDiskUsageTables.panel,
  topLiveRowsTables: topLiveRowsTables.panel,
}
