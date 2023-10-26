local signal_average = import '../signal_average.libsonnet';
local throughput = import '../throughput.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  panels: [
    throughput(datasourceName, policyName, component, extraFilters).panel
    + g.panel.timeSeries.gridPos.withH(8),
    signal_average(datasourceName, policyName, component, extraFilters, 'Accept Percentage', false, 'ACCEPT_PERCENTAGE').panel
    + g.panel.timeSeries.gridPos.withH(8),
  ],
}
