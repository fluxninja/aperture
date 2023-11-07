local statFn = import './stat-panels.libsonnet';
local timeSeriesFn = import './time-series-panels.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters)
  local statPanels = statFn(policyName, infraMeterName, datasource, extraFilters);
  local timePanels = timeSeriesFn(policyName, infraMeterName, datasource, extraFilters);

  local clusterTitleBar = g.panel.text.new('')
                          + g.panel.text.options.withMode('html')
                          + g.panel.text.options.withContent('<div style="text-align: center; width: 100%;"><h2>Cluster Level Stats</h2></div>');

  local nodeTitleBar = g.panel.text.new('')
                       + g.panel.text.options.withMode('html')
                       + g.panel.text.options.withContent('<div style="text-align: center; width: 100%;"><h2>Node Level Stats</h2></div>');

  [
    // Cluster Level Stats
    clusterTitleBar
    + g.panel.text.gridPos.withX(0)
    + g.panel.text.gridPos.withY(0)
    + g.panel.text.gridPos.withH(3)
    + g.panel.text.gridPos.withW(24),
    statPanels.nodes
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.dataNodes
    + g.panel.stat.gridPos.withX(4)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.activeShards
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.activePrimaryShards
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.initializingShards
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.unassignedShards
    + g.panel.stat.gridPos.withX(20)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),

    statPanels.delayedUnassignedShards
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.relocatingShards
    + g.panel.stat.gridPos.withX(4)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.numberOfPendingTasks
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.greenHealth
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.yellowHealth
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(10)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.redHealth
    + g.panel.barGauge.gridPos.withX(20)
    + g.panel.barGauge.gridPos.withY(10)
    + g.panel.barGauge.gridPos.withH(4)
    + g.panel.barGauge.gridPos.withW(4),

    timePanels.cpuUsage
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(15),
    timePanels.breakersTripped
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(25)
    + g.panel.timeSeries.gridPos.withW(24),
    timePanels.indexingPressure
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.indexingRejections
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),

    // Node Level Stats
    nodeTitleBar
    + g.panel.text.gridPos.withX(0)
    + g.panel.text.gridPos.withY(45)
    + g.panel.text.gridPos.withH(3)
    + g.panel.text.gridPos.withW(24),
    timePanels.totalDocs
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(50)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.currentDocs
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(50)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.operationsData
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(60)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.docsData
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(60)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.activeThreads
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(70)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.httpConnections
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(70)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.avgDiskUsage
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(80)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.translogOperations
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(80)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.tasksQueued
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(90)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.tasksThroughput
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(90)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.gcCollectionsPerMinute
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(100)
    + g.panel.timeSeries.gridPos.withW(24),
    timePanels.gcCollectionsPercentTime
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(110)
    + g.panel.timeSeries.gridPos.withW(24),
    timePanels.jvmHeapMemoryUsage
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(120)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.jvmMemoryUsage
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(120)
    + g.panel.timeSeries.gridPos.withW(12),
  ]
