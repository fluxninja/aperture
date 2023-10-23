local signal_average = import '../signal_average.libsonnet';
local signal_frequency = import '../signal_frequency.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  panels: [
    signal_average(datasourceName, policyName, component, extraFilters, 'Signal Average - ${signal_name} (${sub_circuit_id})', true).panel,
    signal_frequency(datasourceName, policyName, component, extraFilters, 'Signal Validity (Frequency) - ${signal_name} (${sub_circuit_id})', true).panel,
  ],
}
