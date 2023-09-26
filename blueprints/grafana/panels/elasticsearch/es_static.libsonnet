local utils = import '../../utils/policy_utils.libsonnet';
local statPanel = import '../../utils/stat_panel.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  local nodes = statPanel('Nodes',
                          datasource.name,
                          'sum(elasticsearch_cluster_nodes{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                          stringFilters),

  local dataNodes = statPanel('Data Nodes',
                              datasource.name,
                              'sum(elasticsearch_cluster_data_nodes{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                              stringFilters),

  local activeShards = statPanel('Active Shards',
                                 datasource.name,
                                 'sum(elasticsearch_cluster_shards{%(filters)s, infra_meter_name="%(infra_meter)s", state="active"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                 stringFilters),

  local activePrimaryShards = statPanel('Active Pri. Shards',
                                        datasource.name,
                                        'sum(elasticsearch_cluster_shards{%(filters)s, infra_meter_name="%(infra_meter)s", state="active_primary"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                        stringFilters),

  local initializingShards = statPanel('Initializing Shards',
                                       datasource.name,
                                       'sum(elasticsearch_cluster_shards{%(filters)s, infra_meter_name="%(infra_meter)s", state="initializing"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                       stringFilters),

  local unassignedShards = statPanel('Unassigned Shards',
                                     datasource.name,
                                     'sum(elasticsearch_cluster_shards{%(filters)s, infra_meter_name="%(infra_meter)s", state="unassigned"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                     stringFilters),

  local delayedUnassignedShards = statPanel('Delayed Un. Shards',
                                            datasource.name,
                                            'sum(elasticsearch_cluster_shards{%(filters)s, infra_meter_name="%(infra_meter)s", state="unassigned_delayed"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                            stringFilters),

  local relocatingShards = statPanel('Relocating Shards',
                                     datasource.name,
                                     'sum(elasticsearch_cluster_shards{%(filters)s, infra_meter_name="%(infra_meter)s", state="relocating"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                     stringFilters),

  local numberOfPendingTasks = statPanel('Pending Tasks',
                                         datasource.name,
                                         'sum(elasticsearch_cluster_pending_tasks{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                         stringFilters),
  local greenHealth = statPanel('Green Health',
                                datasource.name,
                                'sum(elasticsearch_cluster_health{%(filters)s, infra_meter_name="%(infra_meter)s",status="green"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                stringFilters),

  local yellowHealth = statPanel('Yellow Health',
                                 datasource.name,
                                 'sum(elasticsearch_cluster_health{%(filters)s, infra_meter_name="%(infra_meter)s",status="yellow"})' % { filters: stringFilters, infra_meter: infraMeterName },
                                 stringFilters),

  local redHealth = statPanel('Red Health',
                              datasource.name,
                              'sum(elasticsearch_cluster_health{%(filters)s, infra_meter_name="%(infra_meter)s",status="red"})' % { filters: stringFilters, infra_meter: infraMeterName },
                              stringFilters),

  nodes: nodes.panel,
  dataNodes: dataNodes.panel,
  activeShards: activeShards.panel,
  activePrimaryShards: activePrimaryShards.panel,
  initializingShards: initializingShards.panel,
  unassignedShards: unassignedShards.panel,
  delayedUnassignedShards: delayedUnassignedShards.panel,
  relocatingShards: relocatingShards.panel,
  numberOfPendingTasks: numberOfPendingTasks.panel,
  greenHealth: greenHealth.panel,
  yellowHealth: yellowHealth.panel,
  redHealth: redHealth.panel,
}
