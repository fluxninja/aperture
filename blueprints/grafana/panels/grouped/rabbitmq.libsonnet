local static_panels = import '../rabbitmq/rmq_static.libsonnet';
local time_series_panels = import '../rabbitmq/rmq_time_series.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters) {

  local static = static_panels(policyName, infraMeterName, datasource, extraFilters),
  local timeSeries = time_series_panels(policyName, infraMeterName, datasource, extraFilters),

  panels: [
    //Static Panels
    static.ready
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.unacknowledged
    + g.panel.stat.gridPos.withX(4)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.incoming
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.outgoing
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.consumers
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    static.queues
    + g.panel.stat.gridPos.withX(20)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),

    //Time Series Panels
    timeSeries.ready
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(5)
    + g.panel.timeSeries.gridPos.withW(8),
    timeSeries.acknowledged
    + g.panel.timeSeries.gridPos.withX(8)
    + g.panel.timeSeries.gridPos.withY(5)
    + g.panel.timeSeries.gridPos.withW(8),
    timeSeries.unacknowledged
    + g.panel.timeSeries.gridPos.withX(16)
    + g.panel.timeSeries.gridPos.withY(5)
    + g.panel.timeSeries.gridPos.withW(8),
    timeSeries.readyPerVhost
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(15)
    + g.panel.timeSeries.gridPos.withW(8),
    timeSeries.acknowledgedPerVhost
    + g.panel.timeSeries.gridPos.withX(8)
    + g.panel.timeSeries.gridPos.withY(15)
    + g.panel.timeSeries.gridPos.withW(8),
    timeSeries.unacknowledgedPerVhost
    + g.panel.timeSeries.gridPos.withX(16)
    + g.panel.timeSeries.gridPos.withY(15)
    + g.panel.timeSeries.gridPos.withW(8),
    timeSeries.queuesGrowth
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(25),
    timeSeries.published
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),
    timeSeries.delivered
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),

  ],
}
