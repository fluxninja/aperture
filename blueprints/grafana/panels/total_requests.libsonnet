local utils = import '../utils/policy_utils.libsonnet';
local statPanel = import '../utils/stat_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),

  local totalRequests = statPanel('Total Requests',
                                  datasourceName,
                                  'sum(increase(workload_requests_total{%(filters)s}[$__range]))',
                                  stringFilters,
                                  h=10,
                                  w=8,
                                  panelColor='blue',
                                  graphMode='area',
                                  unit='short'),

  panel: totalRequests.panel,
}
