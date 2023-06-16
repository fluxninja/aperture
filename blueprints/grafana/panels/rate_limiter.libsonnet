local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local rateLimiter = timeSeriesPanel('Aperture Rate Limiter',
                                      cfg.dashboard.datasource.name,
                                      'sum by(decision_type) (rate(rate_limiter_counter_total{ %(filters)s}[$__rate_interval]))',
                                      stringFilters,
                                      'Decisions',
                                      'reqps'),

  panel: rateLimiter.panel,
}
