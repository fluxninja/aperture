local signal_average = import '../signal_average.libsonnet';
local throughput = import '../throughput.libsonnet';

function(cfg) {
  panels: [
    signal_average(cfg, 'Accept Percentage', false, 'ACCEPT_PERCENTAGE').panel,
    throughput(cfg).panel,
  ],
}
