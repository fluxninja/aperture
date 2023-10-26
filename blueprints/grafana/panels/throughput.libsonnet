local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),

  local throughput = timeSeriesPanel('Throughput - Accept/Reject',
                                     datasourceName,
                                     'sum(rate(sampler_counter_total{%(filters)s}[$__rate_interval]))',
                                     stringFilters),

  panel: throughput.panel,
}
