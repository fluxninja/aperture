local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local activeConn = statPanel('Active Connections %',
                               datasource.name,
                               'sum(postgresql_backends{%(filters)s,infra_meter_name="%(infra_meter)s"}) / sum(postgresql_connection_max{%(filters)s,infra_meter_name="%(infra_meter)s"}) * 100' % { filters: stringFilters, infra_meter: infraMeterName },
                               stringFilters,
                               instantQuery=true,
                               range=false),

  local currentIndex = statPanel('Index cache hit rate %',
                                 datasource.name,
                                 'sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="idx_hit"}) / (sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="idx_hit"}) + sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="idx_read"}))' % { filters: stringFilters, infra_meter: infraMeterName },
                                 stringFilters,
                                 instantQuery=true,
                                 range=false),

  local currentInsert = statPanel('Table cache hit rate %',
                                  datasource.name,
                                  'sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_hit"}) / (sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_hit"}) + sum(postgresql_blocks_read_total{%(filters)s,infra_meter_name="%(infra_meter)s",source="heap_read"}))' % { filters: stringFilters, infra_meter: infraMeterName },
                                  stringFilters,
                                  instantQuery=true,
                                  range=false),

  local dbCount = statPanel('PGSQL Instances',
                            datasource.name,
                            'count(postgresql_database_count{%(filters)s,infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                            stringFilters,
                            instantQuery=true,
                            range=false,
                            panelColor='blue'),

  local dbSize = statPanel('PGSQL Instances Size (GB)',
                           datasource.name,
                           'sum(postgresql_db_size_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName },
                           stringFilters,
                           instantQuery=true,
                           range=false),

  local maxConnections = statPanel('Max Connections',
                                   datasource.name,
                                   'postgresql_connection_max{%(filters)s,infra_meter_name="%(infra_meter)s"} / postgresql_database_count' % { filters: stringFilters, infra_meter: infraMeterName },
                                   stringFilters,
                                   instantQuery=true,
                                   range=false,
                                   panelColor='blue'),

  local tbCount = statPanel('Number of tables',
                            datasource.name,
                            'sum(postgresql_table_count{%(filters)s,infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                            stringFilters,
                            instantQuery=true,
                            range=false,
                            panelColor='blue'),

  local tbSize = statPanel('Tables Size (MB) ',
                           datasource.name,
                           'sum(postgresql_table_size_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}) / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName },
                           stringFilters,
                           instantQuery=true,
                           range=false),

  activeConn: activeConn.panel,
  currentIndex: currentIndex.panel,
  currentInsert: currentInsert.panel,
  dbCount: dbCount.panel,
  dbSize: dbSize.panel,
  maxConnections: maxConnections.panel,
  tbCount: tbCount.panel,
  tbSize: tbSize.panel,
}
