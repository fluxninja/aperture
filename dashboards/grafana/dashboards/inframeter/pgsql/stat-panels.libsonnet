local statPanel = import '../../../panels/stat.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, infra_meter_name: infraMeterName }),

  local activeConn = statPanel('Active Connections %',
                               datasource,
                               'sum(postgresql_backends{%(filters)s"}) / sum(postgresql_connection_max{%(filters)s"}) * 100' % { filters: stringFilters },
                               instantQuery=true,
                               range=false,
                               unit='percent'),

  local currentIndex = statPanel('Index cache hit rate %',
                                 datasource,
                                 'sum(rate(postgresql_blocks_read_total{%(filters)s",source="idx_hit"}[$__range])) / (sum(rate(postgresql_blocks_read_total{%(filters)s",source="idx_hit"}[$__range])) + sum(rate(postgresql_blocks_read_total{%(filters)s",source="idx_read"}[$__range]))) * 100' % { filters: stringFilters },
                                 instantQuery=true,
                                 range=true,
                                 unit='percent'),

  local currentInsert = statPanel('Table cache hit rate %',
                                  datasource,
                                  'sum(rate(postgresql_blocks_read_total{%(filters)s",source="heap_hit"}[$__range])) / (sum(rate(postgresql_blocks_read_total{%(filters)s",source="heap_hit"}[$__range])) + sum(rate(postgresql_blocks_read_total{%(filters)s",source="heap_read"}[$__range]))) * 100' % { filters: stringFilters },
                                  instantQuery=true,
                                  range=true,
                                  unit='percent'),

  local dbCount = statPanel('PGSQL Instances',
                            datasource,
                            'sum(postgresql_database_count{%(filters)s"})' % { filters: stringFilters },
                            instantQuery=true,
                            range=false,
                            panelColor='blue'),

  local dbSize = statPanel('PGSQL Instances Size (GB)',
                           datasource,
                           'sum(postgresql_db_size_bytes{%(filters)s"})' % { filters: stringFilters },
                           instantQuery=true,
                           range=false,
                           unit='bytes'),

  local maxConnections = statPanel('Max Connections',
                                   datasource,
                                   'sum(postgresql_connection_max{%(filters)s"})' % { filters: stringFilters },
                                   instantQuery=true,
                                   range=false,
                                   panelColor='blue'),

  local tbCount = statPanel('Number of tables',
                            datasource,
                            'sum(postgresql_table_count{%(filters)s"})' % { filters: stringFilters },
                            instantQuery=true,
                            range=false,
                            panelColor='blue'),

  local tbSize = statPanel('Tables Size (MB) ',
                           datasource,
                           'sum(postgresql_table_size_bytes{%(filters)s"})' % { filters: stringFilters },
                           instantQuery=true,
                           range=false,
                           unit='bytes'),

  activeConn: activeConn,
  currentIndex: currentIndex,
  currentInsert: currentInsert,
  dbCount: dbCount,
  dbSize: dbSize,
  maxConnections: maxConnections,
  tbCount: tbCount,
  tbSize: tbSize,
}
