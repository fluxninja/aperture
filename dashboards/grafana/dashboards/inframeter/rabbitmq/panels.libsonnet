local statFn = import './stat-panels.libsonnet';
local timeSeriesFn = import './time-series-panels.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(policyName, infraMeterName, datasource, extraFilters)

  local statPanels = statFn(policyName, infraMeterName, datasource, extraFilters);
  local timePanels = timeSeriesFn(policyName, infraMeterName, datasource, extraFilters);

  [
    //stat Panels
    statPanels.ready
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.unacknowledged
    + g.panel.stat.gridPos.withX(4)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.incoming
    + g.panel.stat.gridPos.withX(8)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.outgoing
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.consumers
    + g.panel.stat.gridPos.withX(16)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    statPanels.queues
    + g.panel.stat.gridPos.withX(20)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),

    //Time Series Panels
    timePanels.ready
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(5)
    + g.panel.timeSeries.gridPos.withW(8),
    timePanels.acknowledged
    + g.panel.timeSeries.gridPos.withX(8)
    + g.panel.timeSeries.gridPos.withY(5)
    + g.panel.timeSeries.gridPos.withW(8),
    timePanels.unacknowledged
    + g.panel.timeSeries.gridPos.withX(16)
    + g.panel.timeSeries.gridPos.withY(5)
    + g.panel.timeSeries.gridPos.withW(8),
    timePanels.readyPerVhost
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(15)
    + g.panel.timeSeries.gridPos.withW(8),
    timePanels.acknowledgedPerVhost
    + g.panel.timeSeries.gridPos.withX(8)
    + g.panel.timeSeries.gridPos.withY(15)
    + g.panel.timeSeries.gridPos.withW(8),
    timePanels.unacknowledgedPerVhost
    + g.panel.timeSeries.gridPos.withX(16)
    + g.panel.timeSeries.gridPos.withY(15)
    + g.panel.timeSeries.gridPos.withW(8),
    timePanels.queuesGrowth
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(25),
    timePanels.published
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),
    timePanels.delivered
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(35)
    + g.panel.timeSeries.gridPos.withW(12),
  ]
