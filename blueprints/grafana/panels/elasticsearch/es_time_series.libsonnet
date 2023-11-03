local utils = import '../../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../../utils/time_series_panel.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName }),

  // CPU Usage
  local cpuQuery = g.query.prometheus.new(datasource, 'avg(elasticsearch_os_cpu_usage_percent{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName })
                   + g.query.prometheus.withIntervalFactor(1)
                   + g.query.prometheus.withLegendFormat('Avg CPU Usage'),

  local cpuUsage = timeSeriesPanel('Average CPU Usage', datasource, '', stringFilters, '%', targets=[cpuQuery]),

  // Breakers (limit vs estimated)
  local breakerQueries = [
    g.query.prometheus.new(datasource, 'sum(elasticsearch_breaker_memory_limit_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Limit'),
    g.query.prometheus.new(datasource, 'sum(elasticsearch_breaker_memory_estimated_bytes{%(filters)s,infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Estimated'),

  ],

  local breakersComparison = timeSeriesPanel('Breakers Limit vs Estimated', datasource, '', stringFilters, 'Gigabytes', targets=breakerQueries),

  // Breakers Tripped
  local breakerTrippedQueries = [
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_breaker_tripped_total{name="fielddata",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Field Data'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_breaker_tripped_total{name="parent",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Parent'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_breaker_tripped_total{name="inflight_requests",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('In-flight Requests'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_breaker_tripped_total{name="eql_sequence",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('EQL Sequence'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_breaker_tripped_total{name="model_inference",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Model Inference'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_breaker_tripped_total{name="request",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Request'),
  ],

  local breakersTripped = timeSeriesPanel('Breakers Tripped', datasource, '', stringFilters, 'Short', targets=breakerTrippedQueries),

  // Indexing Pressure Metrics
  local indexingPressureQueries = [
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_memory_indexing_pressure_bytes{stage="coordinating", %(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Coordinating'),
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_memory_indexing_pressure_bytes{stage="primary", %(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Primary'),
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_memory_indexing_pressure_bytes{stage="replica", %(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Replica'),
  ],

  local indexingPressure = timeSeriesPanel('Indexing Pressure', datasource, '', stringFilters, 'bytes', targets=indexingPressureQueries),

  // Indexing Rejection Rate Metrics
  local indexingRejectionsQueries = [
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_indexing_pressure_memory_primary_rejections_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Primary'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_indexing_pressure_memory_replica_rejections_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Replica'),
  ],

  local indexingRejections = timeSeriesPanel('Indexing Rejection Rate', datasource, '', stringFilters, 'rejections', targets=indexingRejectionsQueries),

  // Operations (query, get, merge, fetch, flush, refresh, scroll, suggest, delete, warmer)
  local operationsQueries = [
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="query"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Query'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="get"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Get'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="merge"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Merge'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="fetch"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Fetch'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="flush"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Flush'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="refresh"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Refresh'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="scroll"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Scroll'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="suggest"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Suggest'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="delete"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Delete'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_operations_completed_total{%(filters)s, infra_meter_name="%(infra_meter)s", operation="warmer"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Warmer'),
  ],

  local operationsData = timeSeriesPanel('Node Operations', datasource, '', stringFilters, 'ops', targets=operationsQueries),

  // Documents State (active and deleted)
  local DocsQueries = [
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_documents{state="active", %(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Active Docs'),

    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_documents{state="deleted", %(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Deleted Docs'),
  ],

  local docsData = timeSeriesPanel('Documents State', datasource, '', stringFilters, 'short', targets=DocsQueries),

  // Total Docs
  local totalDocsQuery =
    g.query.prometheus.new(datasource, 'sum(increase(elasticsearch_node_ingest_documents_total{ %(filters)s, infra_meter_name="%(infra_meter)s"}[$__range]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Total Docs'),

  local totalDocs = timeSeriesPanel('Total Ingested Documents', datasource, '', stringFilters, 'short', targets=[totalDocsQuery]),

  // Current Docs
  local currentDocsQuery =
    g.query.prometheus.new(datasource, 'sum(increase(elasticsearch_node_ingest_documents_current{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__range]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Current Docs'),

  local currentDocs = timeSeriesPanel('Current Ingested Documents', datasource, '', stringFilters, 'short', targets=[currentDocsQuery]),

  local translogOperationsQuery =
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_translog_operations_total{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Translog Operations'),

  local translogOperations = timeSeriesPanel('Translog Operations', datasource, '', stringFilters, 'short', targets=[translogOperationsQuery]),

  // Active Threads for All Thread Pools
  local activeThreadsQueries =
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_thread_pool_threads{%(filters)s, infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{threadpool}} Active Threads'),

  local activeThreads = timeSeriesPanel('Active Threads', datasource, '', stringFilters, 'short', targets=[activeThreadsQueries]),

  local httpConnectionsQuery = g.query.prometheus.new(datasource, 'sum(elasticsearch_node_http_connections{%(filters)s, infra_meter_name="%(infra_meter)s"})' % { filters: stringFilters, infra_meter: infraMeterName })
                               + g.query.prometheus.withIntervalFactor(1)
                               + g.query.prometheus.withLegendFormat('HTTP Connections'),

  local httpConnections = timeSeriesPanel('HTTP Connections', datasource, '', stringFilters, 'short', targets=[httpConnectionsQuery]),

  local avgDiskUsageQuery = g.query.prometheus.new(
                              datasource,
                              |||
                                avg(
                                  (
                                    1 - (
                                      elasticsearch_node_fs_disk_free_bytes{%(filters)s, infra_meter_name="%(infra_meter)s"}
                                      /
                                      elasticsearch_node_fs_disk_total_bytes{%(filters)s, infra_meter_name="%(infra_meter)s"}
                                    )
                                  ) * 100
                                )
                              |||
                              % {
                                filters: stringFilters,
                                infra_meter: infraMeterName,
                              }
                            )
                            + g.query.prometheus.withIntervalFactor(1)
                            + g.query.prometheus.withLegendFormat('Avg Disk Usage Percentage'),


  local avgDiskUsage = timeSeriesPanel('Average Disk Usage', datasource, '', stringFilters, '%', targets=[avgDiskUsageQuery]),

  // Tasks Queued for All Thread Pools
  local tasksQueuedQueries =
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_thread_pool_tasks_queued{%(filters)s , infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{threadpool}} Tasks Queued'),

  local tasksQueued = timeSeriesPanel('Tasks Queued', datasource, '', stringFilters, 'short', targets=[tasksQueuedQueries]),

  // Tasks Throughput for All Thread Pools
  local tasksThroughputQueries =
    g.query.prometheus.new(datasource, 'sum(rate(elasticsearch_node_thread_pool_tasks_finished_total{%(filters)s , infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('{{threadpool}} Tasks Throughput'),

  local tasksThroughput = timeSeriesPanel('Tasks Throughput', datasource, '', stringFilters, 'ops', targets=[tasksThroughputQueries]),

  // GC Collections (young and old)
  local youngGcCollectionsQuery = g.query.prometheus.new(datasource, 'sum(rate(jvm_gc_collections_count_total{name="young",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
                                  + g.query.prometheus.withIntervalFactor(1)
                                  + g.query.prometheus.withLegendFormat('Young GC Collections'),

  local youngGcCollections = timeSeriesPanel('Young GC Collections', datasource, '', stringFilters, 'cps', targets=[youngGcCollectionsQuery]),

  local oldGcCollectionsQuery = g.query.prometheus.new(datasource, 'sum(rate(jvm_gc_collections_count_total{name="old",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
                                + g.query.prometheus.withIntervalFactor(1)
                                + g.query.prometheus.withLegendFormat('Old GC Collections'),

  local oldGcCollections = timeSeriesPanel('Old GC Collections', datasource, '', stringFilters, 'cps', targets=[oldGcCollectionsQuery]),

  // GC Collections Elapsed (young and old)
  local youngGcCollectionsElapsedQuery = g.query.prometheus.new(datasource, 'sum(rate(jvm_gc_collections_elapsed_milliseconds_total{name="young",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
                                         + g.query.prometheus.withIntervalFactor(1)
                                         + g.query.prometheus.withLegendFormat('Young GC Collections Elapsed'),

  local youngGcCollectionsElapsed = timeSeriesPanel('Young GC Collections Elapsed', datasource, '', stringFilters, 'msps', targets=[youngGcCollectionsElapsedQuery]),

  local oldGcCollectionsElapsedQuery = g.query.prometheus.new(datasource, 'sum(rate(jvm_gc_collections_elapsed_milliseconds_total{name="old",%(filters)s,infra_meter_name="%(infra_meter)s"}[$__rate_interval]))' % { filters: stringFilters, infra_meter: infraMeterName })
                                       + g.query.prometheus.withIntervalFactor(1)
                                       + g.query.prometheus.withLegendFormat('Old GC Collections Elapsed'),

  local oldGcCollectionsElapsed = timeSeriesPanel('Old GC Collections Elapsed', datasource, '', stringFilters, 'msps', targets=[oldGcCollectionsElapsedQuery]),


  // JVM Usage
  local jvmQueries = [
    g.query.prometheus.new(datasource, 'sum(jvm_memory_heap_used_bytes{%(filters)s, infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Heap Used'),

    g.query.prometheus.new(datasource, 'sum(jvm_memory_heap_max_bytes{%(filters)s, infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Heap Max'),
  ],

  local jvmMemoryUsage = timeSeriesPanel('JVM Memory Usage', datasource, '', stringFilters, 'Gigabytes', targets=jvmQueries),
  // JVM Non-Heap Memory (Used and Committed)
  local nonHeapQueries = [
    g.query.prometheus.new(datasource, 'sum(jvm_memory_nonheap_used_bytes{%(filters)s, infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Non-Heap Used'),

    g.query.prometheus.new(datasource, 'sum(jvm_memory_nonheap_committed_bytes{%(filters)s, infra_meter_name="%(infra_meter)s"}) / 1024 / 1024 / 1024' % { filters: stringFilters, infra_meter: infraMeterName })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Non-Heap Committed'),
  ],

  local jvmNonHeap = timeSeriesPanel('JVM Non-Heap Memory Usage', datasource, '', stringFilters, 'Gigabytes', targets=nonHeapQueries),

  cpuUsage: cpuUsage.panel,
  jvmMemoryUsage: jvmMemoryUsage.panel,
  youngGcCollections: youngGcCollections.panel,
  oldGcCollections: oldGcCollections.panel,
  youngGcCollectionsElapsed: youngGcCollectionsElapsed.panel,
  oldGcCollectionsElapsed: oldGcCollectionsElapsed.panel,
  breakersComparison: breakersComparison.panel,
  breakersTripped: breakersTripped.panel,
  httpConnections: httpConnections.panel,
  docsData: docsData.panel,
  operationsData: operationsData.panel,
  jvmNonHeap: jvmNonHeap.panel,
  indexingPressure: indexingPressure.panel,
  indexingRejections: indexingRejections.panel,
  activeThreads: activeThreads.panel,
  tasksQueued: tasksQueued.panel,
  tasksThroughput: tasksThroughput.panel,
  avgDiskUsage: avgDiskUsage.panel,
  totalDocs: totalDocs.panel,
  currentDocs: currentDocs.panel,
  translogOperations: translogOperations.panel,
}
