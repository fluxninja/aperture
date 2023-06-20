local signal_average = import '../signal_average.libsonnet';
local throughput = import '../throughput.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    throughput(cfg).panel
    + g.panel.timeSeries.gridPos.withH(8),
    signal_average(cfg, 'Accept Percentage', false, 'ACCEPT_PERCENTAGE').panel
    + g.panel.timeSeries.gridPos.withH(8),
  ],
}
