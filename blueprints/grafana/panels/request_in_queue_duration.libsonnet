local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),

  local workloadLatency = timeSeriesPanel('Request in Queue Duration',
                                          datasourceName,
                                          '(sum by (workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__rate_interval])))/ ((sum by (workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__rate_interval]))) != 0)',
                                          stringFilters,
                                          'Wait Time',
                                          'ms'),

  panel: workloadLatency.panel,
}
