local barGaugePanel = import '../utils/bar_gauge_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local requestsDuration = barGaugePanel('Requests in Queue Duration',
                                         cfg.dashboard.datasource.name,
                                         'topk(10, (sum by(workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__range])) ) / (sum by(workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__range])) ))',
                                         stringFilters,
                                         unit='ms',
                                         instantQuery=true,
                                         range=false),

  panel: requestsDuration.panel,
}
