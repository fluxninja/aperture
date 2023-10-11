local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local throughput = timeSeriesPanel('Throughput - Accept/Reject',
                                     datasourceName,
                                     'sum(rate(sampler_counter_total{%(filters)s}[$__rate_interval]))',
                                     stringFilters),

  panel: throughput.panel,
}
