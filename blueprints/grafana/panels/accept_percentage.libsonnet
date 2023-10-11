local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),
  local acceptPercentage_filters = utils.dictToPrometheusFilter(stringFilters { signal_name: 'ACCEPT_PERCENTAGE' }),

  local acceptPercentage = timeSeriesPanel('Accept Percentage',
                                           datasourceName,
                                           'increase(signal_reading_sum{%(filters)s}[$__rate_interval]) / increase(signal_reading_count{%(filters)s}[$__rate_interval])',
                                           acceptPercentage_filters),

  panel: acceptPercentage.panel,
}
