local utils = import '../utils/policy_utils.libsonnet';
local timeSeriesPanel = import '../utils/time_series_panel.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local throughput = timeSeriesPanel('Throughput - Accept/Reject',
                                     cfg.dashboard.datasource.name,
                                     'rate(sampler_counter_total{%(filters)s}[$__rate_interval])',
                                     stringFilters),

  panel: throughput.panel,
}
