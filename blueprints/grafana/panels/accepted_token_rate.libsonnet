local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(datasourceName, policyName, component, extraFilters={}) {
  local stringFilters = utils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: component.component_id }),

  local acceptedTokenRate = timeSeriesPanel('Accepted Token Rate',
                                            datasourceName,
                                            'sum(rate(accepted_tokens_total{%(filters)s}[$__rate_interval]))',
                                            stringFilters,
                                            'Token Rate',
                                            '',
                                            8,
                                            12),

  panel: acceptedTokenRate.panel,
}
