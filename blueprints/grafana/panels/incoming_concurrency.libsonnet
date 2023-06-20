local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local incomingConcurrency = timeSeriesPanel('Incoming Token Rate',
                                              cfg.dashboard.datasource.name,
                                              'sum(rate(incoming_tokens_total{%(filters)s}[$__rate_interval]))',
                                              stringFilters,
                                              'Concurrency',
                                              '',
                                              8,
                                              12),

  panel: incomingConcurrency.panel,
}
