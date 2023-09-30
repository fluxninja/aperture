local barGauge = import '../pgsql/pgsql_bar_gauge.libsonnet';
local statPanels = import '../pgsql/pgsql_static.libsonnet';
local timeSeries = import '../pgsql/pgsql_time_series.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {

  local static = statPanels(policyName, infraMeterName, datasource, extraFilters),
  local time = timeSeries(policyName, infraMeterName, datasource, extraFilters),
  local bar = barGauge(policyName, infraMeterName, datasource, extraFilters),
  panels: [
    static.dbCount
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.dbSize
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.tbCount
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.tbSize
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.currentInsert
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.currentIndex
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.activeConn
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    static.maxConnections
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    time.commitVsRollback
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(20)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    time.operations
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(20)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    time.blockReads
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(30)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    time.checkpointComparison
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(30)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    time.bufferWrites
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(40)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    bar.topLiveRowsTables
    + g.panel.barGauge.gridPos.withX(0)
    + g.panel.barGauge.gridPos.withY(50)
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.options.withValueMode('text'),
    bar.topDiskUsageTables
    + g.panel.barGauge.gridPos.withX(12)
    + g.panel.barGauge.gridPos.withY(50)
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.options.withValueMode('text'),
  ],
}
