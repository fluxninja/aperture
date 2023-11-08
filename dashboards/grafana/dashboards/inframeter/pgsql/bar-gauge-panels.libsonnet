local barGaugePanel = import '../../../panels/bar-gauge.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, infra_meter_name: infraMeterName }),

  local topDiskUsageTables = barGaugePanel('Tables with most disk usage',
                                           datasource,
                                           'topk(5, sum by (postgresql_table_name,postgresql_database_name) (postgresql_table_size_bytes{%(filters)s}))' % { filters: stringFilters }),

  local topLiveRowsTables = barGaugePanel('Tables with most live rows',
                                          datasource,
                                          'topk(5, sum by (postgresql_table_name,postgresql_database_name) (postgresql_rows{%(filters)s,state="live"}))' % { filters: stringFilters }),

  topDiskUsageTables: topDiskUsageTables,
  topLiveRowsTables: topLiveRowsTables,
}
