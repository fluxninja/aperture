local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, infra_meter_name: infraMeterName }),

  // CPU Usage
  local cpuQuery = g.query.prometheus.new(datasource, 'avg by (elasticsearch_node_name) (elasticsearch_os_cpu_usage_percent{%(filters)s})' % { filters: stringFilters })
                   + g.query.prometheus.withIntervalFactor(1)
                   + g.query.prometheus.withLegendFormat('{{ elasticsearch_node_name }}'),

  local cpuUsage = timeSeriesPanel('Average CPU Usage per Node', datasource, axisLabel='Percentage', unit='percent', targets=[cpuQuery]),

  // Breakers Tripped
  local breakerTrippedQueries = [
    g.query.prometheus.new(datasource, 'avg by (name) (rate(elasticsearch_breaker_tripped_total{name="fielddata",%(filters)s}[$__rate_interval])) * 60' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ name }}'),
  ],

  local breakersTripped = timeSeriesPanel('Average number of Breaker Trips per Minute', datasource, axisLabel='Count per minute', unit='cpm', targets=breakerTrippedQueries),

  // Indexing Pressure Metrics
  local indexingPressureQueries = [
    g.query.prometheus.new(datasource, 'avg by (stage) (elasticsearch_memory_indexing_pressure_bytes{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ stage }}'),
  ],

  local indexingPressure = timeSeriesPanel('Average Indexing Pressure by Stage', datasource, axisLabel='Memory', unit='bytes', targets=indexingPressureQueries),

  // Indexing Rejection Rate Metrics
  local indexingRejectionsQueries = [
    g.query.prometheus.new(datasource, 'avg by (elasticsearch_node_name) (rate(elasticsearch_indexing_pressure_memory_primary_rejections_total{%(filters)s}[$__rate_interval])) * 60' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Primary'),

    g.query.prometheus.new(datasource, 'avg by (elasticsearch_node_name) (rate(elasticsearch_indexing_pressure_memory_replica_rejections_total{%(filters)s}[$__rate_interval])) * 60' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Replica'),
  ],

  local indexingRejections = timeSeriesPanel('Indexing Rejection Rate per Node', datasource, axisLabel='Count per minute', unit='cpm', targets=indexingRejectionsQueries),

  // Operations (query, get, merge, fetch, flush, refresh, scroll, suggest, delete, warmer)
  local operationsQueries = [
    g.query.prometheus.new(datasource, 'avg by (operation) (rate(elasticsearch_node_operations_completed_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1),
  ],

  local operationsData = timeSeriesPanel('Operations Completed per Second', datasource, axisLabel='Operations', targets=operationsQueries),

  // Documents State (active and deleted)
  local DocsQueries = [
    g.query.prometheus.new(datasource, 'avg by (elasticsearch_node_name) (elasticsearch_node_documents{state="active", %(filters)s}[$__rate_interval])' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1),
  ],

  local docsData = timeSeriesPanel('Active Documents per Node', datasource, axisLabel='Count', unit='short', targets=DocsQueries),

  // Total Docs
  local totalDocsQuery =
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_ingest_documents_total{ %(filters)s}[$__range])) * 60' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Total'),

  local totalDocs = timeSeriesPanel('Documents Ingested per Minute across the entire cluster', datasource, axisLabel='Count per minute', unit='cpm', targets=[totalDocsQuery]),

  // Current Docs
  local currentDocsQuery =
    g.query.prometheus.new(datasource, 'sum(elasticsearch_node_ingest_documents_current{%(filters)s})' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Total'),

  local currentDocs = timeSeriesPanel('Documents Currently Ingesting across the entire cluster', datasource, axisLabel='Count', unit='short', targets=[currentDocsQuery]),

  local translogOperationsQuery =
    g.query.prometheus.new(datasource, 'avg(rate(elasticsearch_node_translog_operations_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Translog Operations'),

  local translogOperations = timeSeriesPanel('Average Transaction Log Operations per Second', datasource, axisLabel='Operations per second', unit='ops', targets=[translogOperationsQuery]),

  // Active Threads for All Thread Pools
  local activeThreadsQueries =
    g.query.prometheus.new(datasource, 'avg by (thread_pool_name) (elasticsearch_node_thread_pool_threads{state="active", %(filters)s})' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1),

  local activeThreads = timeSeriesPanel('Active Threads per Pool', datasource, axisLabel='Count', unit='short', targets=[activeThreadsQueries]),

  local httpConnectionsQuery = g.query.prometheus.new(datasource, 'sum(elasticsearch_node_http_connections{%(filters)s})' % { filters: stringFilters })
                               + g.query.prometheus.withIntervalFactor(1)
                               + g.query.prometheus.withLegendFormat('Total Connections'),

  local httpConnections = timeSeriesPanel('Total HTTP Connections across the Cluster', datasource, axisLabel='Count', unit='short', targets=[httpConnectionsQuery]),

  local avgDiskUsageQuery = g.query.prometheus.new(
                              datasource,
                              |||
                                avg by (elasticsearch_node_name) (
                                  (
                                    1 - (
                                      elasticsearch_node_fs_disk_free_bytes{%(filters)s}
                                      /
                                      elasticsearch_node_fs_disk_total_bytes{%(filters)s}
                                    )
                                  ) * 100
                                )
                              |||
                              % { filters: stringFilters }
                            )
                            + g.query.prometheus.withIntervalFactor(1),


  local avgDiskUsage = timeSeriesPanel('Disk Usage Percentage', datasource, axisLabel='Percentage', unit='percent', targets=[avgDiskUsageQuery]),

  // Tasks Queued for All Thread Pools
  local tasksQueuedQueries =
    g.query.prometheus.new(datasource, 'avg by (thread_pool_name) (elasticsearch_node_thread_pool_tasks_queued{%(filters)s })' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1),

  local tasksQueued = timeSeriesPanel('Tasks Queued per Pool', datasource, axisLabel='Count', unit='short', targets=[tasksQueuedQueries]),

  // Tasks Throughput for All Thread Pools
  local tasksThroughputQueries =
    g.query.prometheus.new(datasource, 'sum by (state) (rate(elasticsearch_node_thread_pool_tasks_finished_total{%(filters)s }[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{ state }}'),

  local tasksThroughput = timeSeriesPanel('Task Throughput across the Cluster', datasource, axisLabel='Count per second', unit='cps', targets=[tasksThroughputQueries]),

  // GC Collections (young and old)
  local gcCollectionsPerMinuteQuery = g.query.prometheus.new(datasource, 'avg by (name) (rate(jvm_gc_collections_count_total{%(filters)s}[$__rate_interval])) * 60' % { filters: stringFilters })
                                      + g.query.prometheus.withIntervalFactor(1)
                                      + g.query.prometheus.withLegendFormat('{{ name }}'),

  local gcCollectionsPerMinute = timeSeriesPanel('Average GC Collections per Minute', datasource, axisLabel='Count per minute', unit='cps', targets=[gcCollectionsPerMinuteQuery]),

  // GC Collections Elapsed (young and old)
  local gcCollectionsPercentTimeQuery = g.query.prometheus.new(datasource, 'avg by (name) (rate(jvm_gc_collections_elapsed_milliseconds_total{%(filters)s}[$__rate_interval])) / 10' % { filters: stringFilters })
                                        + g.query.prometheus.withIntervalFactor(1)
                                        + g.query.prometheus.withLegendFormat('{{ name }}'),

  local gcCollectionsPercentTime = timeSeriesPanel('Percentage of time spent in GC Collection', datasource, 'percent', targets=[gcCollectionsPercentTimeQuery]),


  // JVM Heap Usage
  local jvmQueries = [
    g.query.prometheus.new(datasource, 'avg by (elasticsearch_node_name) (jvm_memory_heap_utilization{%(filters)s}) * 100' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1),
  ],

  local jvmHeapMemoryUsage = timeSeriesPanel('Percentage Heap Used per Node', datasource, axisLabel='Percentage', unit='percent', targets=jvmQueries),

  // JVM Memory Heap and Non-Heap Usage
  local memoryUsageQueries = [
    g.query.prometheus.new(datasource, 'avg(jvm_memory_nonheap_used_bytes{%(filters)s}) ' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Non-Heap Used'),

    g.query.prometheus.new(datasource, 'avg(jvm_memory_heap_used_bytes{%(filters)s}) ' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Heap Used'),
  ],

  local jvmMemoryUsage = timeSeriesPanel('Average Memory Usage', datasource, axisLabel='Memory', unit='bytes', targets=memoryUsageQueries),

  cpuUsage: cpuUsage,
  jvmHeapMemoryUsage: jvmHeapMemoryUsage,
  gcCollectionsPerMinute: gcCollectionsPerMinute,
  gcCollectionsPercentTime: gcCollectionsPercentTime,
  breakersTripped: breakersTripped,
  httpConnections: httpConnections,
  docsData: docsData,
  operationsData: operationsData,
  jvmMemoryUsage: jvmMemoryUsage,
  indexingPressure: indexingPressure,
  indexingRejections: indexingRejections,
  activeThreads: activeThreads,
  tasksQueued: tasksQueued,
  tasksThroughput: tasksThroughput,
  avgDiskUsage: avgDiskUsage,
  totalDocs: totalDocs,
  currentDocs: currentDocs,
  translogOperations: translogOperations,
}
