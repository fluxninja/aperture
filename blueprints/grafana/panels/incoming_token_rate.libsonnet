local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),

  local incomingTokenRate = timeSeriesPanel('Incoming Token Rate',
                                            datasourceName,
                                            'sum(rate(incoming_tokens_total{%(filters)s}[$__rate_interval]))',
                                            stringFilters,
                                            'Token Rate',
                                            '',
                                            8,
                                            12),

  panel: incomingTokenRate.panel,
}
