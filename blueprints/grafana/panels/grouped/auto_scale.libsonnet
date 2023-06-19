local signal_average = import '../signal_average.libsonnet';
local signal_frequency = import '../signal_frequency.libsonnet';

function(cfg) {
  panels: [
    signal_average(cfg, 'Actual Scale Average', false, 'ACTUAL_SCALE').panel,
    signal_frequency(cfg, 'Actual Scale Validity (Frequency)', false, 'ACTUAL_SCALE').panel,

    signal_average(cfg, 'Configured Scale Average', false, 'CONFIGURED_SCALE').panel,
    signal_frequency(cfg, 'Configured Scale Validity (Frequency)', false, 'CONFIGURED_SCALE').panel,

    signal_average(cfg, 'Desired Scale Average', false, 'DESIRED_SCALE').panel,
    signal_frequency(cfg, 'Desired Scale Validity (Frequency)', false, 'DESIRED_SCALE').panel,
  ],
}
