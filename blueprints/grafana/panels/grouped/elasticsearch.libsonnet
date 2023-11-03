local static_panels = import '../elasticsearch/es_static.libsonnet';
local time_series_panels = import '../elasticsearch/es_time_series.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {
  local static = static_panels(policyName, infraMeterName, datasource, extraFilters),
  local timeSeries = time_series_panels(policyName, infraMeterName, datasource, extraFilters),

  local clusterTitleBar = g.panel.text.new('')
                          + g.panel.text.options.withMode('html')
                          + g.panel.text.options.withContent('<div style="text-align: center; width: 100%;"><h2>Cluster Level Stats</h2></div>'),

  local nodeTitleBar = g.panel.text.new('')
                       + g.panel.text.options.withMode('html')
                       + g.panel.text.options.withContent('<div style="text-align: center; width: 100%;"><h2>Node Level Stats</h2></div>'),

  panels: [

    // Cluster Level Stats
    clusterTitleBar
    + g.panel.text.gridPos.withX(0)
    + g.panel.text.gridPos.withY(0)
    + g.panel.text.gridPos.withH(3)
    + g.panel.text.gridPos.withW(24),
    static.nodes
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.dataNodes
    + g.panel.stat.gridPos.withX(4)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.activeShards
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.activePrimaryShards
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.initializingShards
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.unassignedShards
    + g.panel.stat.gridPos.withX(20)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),

    static.delayedUnassignedShards
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.relocatingShards
    + g.panel.stat.gridPos.withX(4)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.numberOfPendingTasks
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.greenHealth
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.yellowHealth
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.redHealth
    + g.panel.barGauge.gridPos.withX(20)
    + g.panel.barGauge.gridPos.withY(10)
    + g.panel.barGauge.gridPos.withH(4)
    + g.panel.barGauge.gridPos.withW(4),

    timeSeries.cpuUsage
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(15),
    timeSeries.breakersComparison
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(25)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.breakersTripped
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(25)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.indexingPressure
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.indexingRejections
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),

    // Node Level Stats
    nodeTitleBar
    + g.panel.text.gridPos.withX(0)
    + g.panel.text.gridPos.withY(45)
    + g.panel.text.gridPos.withH(3)
    + g.panel.text.gridPos.withW(24),
    timeSeries.totalDocs
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(50)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.currentDocs
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(50)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.operationsData
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(60)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.docsData
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(60)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.activeThreads
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(70)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.httpConnections
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(70)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.avgDiskUsage
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(80)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.translogOperations
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(80)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.tasksQueued
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(90)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.tasksThroughput
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(90)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.youngGcCollections
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(100)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.youngGcCollectionsElapsed
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(100)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.oldGcCollections
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(110)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.oldGcCollectionsElapsed
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(110)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.jvmMemoryUsage
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(120)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.jvmNonHeap
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(120)
    + g.panel.timeSeries.gridPos.withW(12),
  ],
}
