local signal_average = import '../signal_average.libsonnet';
local signal_frequency = import '../signal_frequency.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  panels: [
    signal_average(datasourceName, policyName, component, extraFilters, 'Actual Scale Average', false, 'ACTUAL_SCALE').panel,
    signal_frequency(datasourceName, policyName, component, extraFilters, 'Actual Scale Validity (Frequency)', false, 'ACTUAL_SCALE').panel,

    signal_average(datasourceName, policyName, component, extraFilters, 'Configured Scale Average', false, 'CONFIGURED_SCALE').panel,
    signal_frequency(datasourceName, policyName, component, extraFilters, 'Configured Scale Validity (Frequency)', false, 'CONFIGURED_SCALE').panel,

    signal_average(datasourceName, policyName, component, extraFilters, 'Desired Scale Average', false, 'DESIRED_SCALE').panel,
    signal_frequency(datasourceName, policyName, component, extraFilters, 'Desired Scale Validity (Frequency)', false, 'DESIRED_SCALE').panel,
  ],
}
