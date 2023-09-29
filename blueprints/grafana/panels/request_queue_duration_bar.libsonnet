local barChartPanel = import '../utils/bar_chart_panel.libsonnet';
local utils = import '../utils/policy_utils.libsonnet';

function(cfg) {
  local stringFilters = utils.dictToPrometheusFilter(cfg.dashboard.extra_filters { policy_name: cfg.policy.policy_name }),

  local requestsDuration = barChartPanel('Requests Duration',
                                         cfg.dashboard.datasource.name,
                                         'topk(10, (sum by(workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__range])) ) / (sum by(workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__range])) ))',
                                         stringFilters,
                                         legendFormat='workload_index',
                                         queryFormat='table',
                                         instantQuery=true,
                                         range=false,
                                         labelSpacing=100,
                                         axisGridshow=false,
                                         axisPlacement='hidden',
                                         mode='multi',
                                         sort='asc'),

  panel: requestsDuration.panel,
}
