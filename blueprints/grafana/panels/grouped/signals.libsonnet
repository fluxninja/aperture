local signal_average = import '../signal_average.libsonnet';
local signal_frequency = import '../signal_frequency.libsonnet';

function(cfg) {
  panels: [
    signal_average(cfg, 'Signal Average - ${signal_name} (${sub_circuit_id})', true).panel,
    signal_frequency(cfg, 'Signal Validity (Frequency) - ${signal_name} (${sub_circuit_id})', true).panel,
  ],
}
