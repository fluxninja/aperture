local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),
  local acceptPercentage_filters = utils.dictToPrometheusFilter(stringFilters { signal_name: 'ACCEPT_PERCENTAGE' }),

  local acceptPercentage = timeSeriesPanel('Accept Percentage',
                                           cfg.dashboard.datasource.name,
                                           'increase(signal_reading_sum{%(filters)s}[$__rate_interval]) / increase(signal_reading_count{%(filters)s}[$__rate_interval])',
                                           acceptPercentage_filters),

  panel: acceptPercentage.panel,
}
