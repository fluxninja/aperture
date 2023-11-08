local barFn = import './bar-gauge-panels.libsonnet';
local statFn = import './stat-panels.libsonnet';
local timeSeriesFn = import './time-series-panels.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters)

  local statPanels = statFn(policyName, infraMeterName, datasource, extraFilters);
  local timePanels = timeSeriesFn(policyName, infraMeterName, datasource, extraFilters);
  local barPanels = barFn(policyName, infraMeterName, datasource, extraFilters);

  [
    statPanels.dbCount
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.dbSize
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.tbCount
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.tbSize
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.currentInsert
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.currentIndex
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.activeConn
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    statPanels.maxConnections
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    timePanels.commitVsRollback
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(20)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    timePanels.operations
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(20)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    timePanels.blockReads
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(30)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    timePanels.checkpointComparison
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(30)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    timePanels.bufferWrites
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(40)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    barPanels.topLiveRowsTables
    + g.panel.barGauge.gridPos.withX(0)
    + g.panel.barGauge.gridPos.withY(50)
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.options.withValueMode('text'),
    barPanels.topDiskUsageTables
    + g.panel.barGauge.gridPos.withX(12)
    + g.panel.barGauge.gridPos.withY(50)
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.options.withValueMode('text'),
  ]
