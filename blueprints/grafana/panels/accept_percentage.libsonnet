local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local componentID = std.get(component.component, 'load_scheduler_component_id', default=component.component_id),
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID }),
  local acceptPercentage_filters = utils.dictToPrometheusFilter(stringFilters { signal_name: 'ACCEPT_PERCENTAGE' }),

  local acceptPercentage = timeSeriesPanel('Accept Percentage',
                                           datasourceName,
                                           'increase(signal_reading_sum{%(filters)s}[$__rate_interval]) / increase(signal_reading_count{%(filters)s}[$__rate_interval])',
                                           acceptPercentage_filters),

  panel: acceptPercentage.panel,
}
