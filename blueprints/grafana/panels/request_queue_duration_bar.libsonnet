local barChartPanel = import '../utils/bar_chart_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),

  local requestsDuration = barChartPanel('Requests in Queue Duration',
                                         datasourceName,
                                         'topk(10, (sum by(workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__range])) ) / ((sum by(workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__range])) )) != 0)',
                                         stringFilters,
                                         instantQuery=true,
                                         range=false,
                                         unit='ms'),

  panel: requestsDuration.panel,
}
