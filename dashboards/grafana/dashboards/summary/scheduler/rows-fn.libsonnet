local barChartPanel = import '../../../panels/bar-chart.libsonnet';
local barGaugePanel = import '../../../panels/bar-gauge.libsonnet';
local statPanel = import '../../../panels/stat.libsonnet';
local timeSeriesPanel = import '../../../panels/time-series.libsonnet';
local promUtils = import '../../../utils/prometheus.libsonnet';
local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v10.1.0/main.libsonnet';

function(datasourceName, policyName, componentID, scheduler, extraFilters={})
  local stringFilters = promUtils.dictToPrometheusFilter(extraFilters { policy_name: policyName, component_id: componentID });

  local row1 = [
    timeSeriesPanel(
      'Workload Decisions',
      datasourceName,
      query='sum by(decision_type) (rate(workload_requests_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters },
      axisLabel='Decisions',
      unit='reqps'
    ),
  ];

  local row2 = [
    timeSeriesPanel(
      'Workload Decisions (accepted)',
      datasourceName,
      query='sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_ACCEPTED"}[$__rate_interval]))' % { filters: stringFilters },
      axisLabel='Decisions',
      unit='reqps'
    ),
  ];

  local row3 = [
    timeSeriesPanel(
      'Workload Decisions (rejected)',
      datasourceName,
      query='sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_REJECTED"}[$__rate_interval]))' % { filters: stringFilters },
      axisLabel='Decisions',
      unit='reqps'
    ),
  ];

  local row4 = [
    statPanel(
      'Total Requests',
      datasourceName,
      'sum(increase(workload_requests_total{%(filters)s}[$__range]))' % { filters: stringFilters },
      x=0,
      h=10,
      w=8,
      panelColor='blue',
      graphMode='area',
      unit='short'
    ),
    statPanel(
      'Total Accepted Requests',
      datasourceName,
      'sum(increase(workload_requests_total{%(filters)s, decision_type="DECISION_TYPE_ACCEPTED"}[$__range]))' % { filters: stringFilters },
      x=8,
      h=10,
      w=8,
      graphMode='area',
      unit='short'
    ),
    statPanel(
      'Total Rejected Requests',
      datasourceName,
      'sum(increase(workload_requests_total{%(filters)s, decision_type="DECISION_TYPE_REJECTED"}[$__range]))' % { filters: stringFilters },
      x=16,
      h=10,
      w=8,
      panelColor='red',
      graphMode='area',
      noValue='No rejected requests',
      unit='short'
    ),
  ];

  local row5 = [
    timeSeriesPanel(
      'Workload Latency',
      datasourceName,
      query='(sum by (workload_index) (increase(workload_latency_ms_sum{%(filters)s}[$__rate_interval])))/(sum by (workload_index) (increase(workload_latency_ms_count{%(filters)s}[$__rate_interval])))' % { filters: stringFilters },
      axisLabel='Latency',
      unit='ms'
    ),
  ];

  local row6 = [
    timeSeriesPanel(
      'Request in Queue Duration',
      datasourceName,
      query='(sum by (workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__rate_interval])))/ ((sum by (workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__rate_interval]))) != 0)' % { filters: stringFilters },
      axisLabel='Wait Time',
      unit='ms',
      x=0,
      w=12
    ),
    barChartPanel(
      'Request in Queue Duration',
      datasourceName,
      query='topk(10, (sum by(workload_index) (increase(request_in_queue_duration_ms_sum{%(filters)s}[$__range])) ) / ((sum by(workload_index) (increase(request_in_queue_duration_ms_count{%(filters)s}[$__range])) )) != 0)' % { filters: stringFilters },
      x=12,
      w=12,
      instantQuery=true,
      range=false,
      unit='ms'
    ),
  ];

  local row7 = [
    timeSeriesPanel(
      'Average preemption per request by Workload',
      datasourceName,
      query='( sum by (workload_index) (rate(workload_preempted_tokens_sum{%(filters)s}[$__rate_interval])) - sum by (workload_index) (rate(workload_delayed_tokens_sum{%(filters)s}[$__rate_interval])) ) / (( sum by (workload_index) (rate(workload_preempted_tokens_count{%(filters)s}[$__rate_interval])) + sum by (workload_index) (rate(workload_delayed_tokens_count{%(filters)s}[$__rate_interval])) + sum by (workload_index) (rate(workload_on_time_total{%(filters)s}[$__rate_interval])) ) != 0)' % { filters: stringFilters },
      axisLabel='Tokens',
      description='Preemption measures the average number of tokens by which a request from a specific workload was preempted, compared to a purely FIFO (first-in-first-out) ordering of requests. A negative preemption value indicates that the workload was delayed rather than preempted.',
      x=0,
      w=12
    ),
    barChartPanel(
      'Average Preemption Tokens per Request by Workload',
      datasourceName,
      query='( sum by (workload_index) (rate(workload_preempted_tokens_sum{%(filters)s}[$__range])) - sum by (workload_index) (rate(workload_delayed_tokens_sum{%(filters)s}[$__range])) ) / (( sum by (workload_index) (rate(workload_preempted_tokens_count{%(filters)s}[$__range])) + sum by (workload_index) (rate(workload_delayed_tokens_count{%(filters)s}[$__range])) + sum by (workload_index) (rate(workload_on_time_total{%(filters)s}[$__range])) ) != 0)' % { filters: stringFilters },
      description='Preemption measures the average number of tokens by which a request from a specific workload was preempted, compared to a purely FIFO (first-in-first-out) ordering of requests. A negative preemption value indicates that the workload was delayed rather than preempted.',
      instantQuery=true,
      range=false,
      x=12,
      w=12
    ),
  ];

  local acceptedVsRejectedTargets = [
    g.query.prometheus.new(datasourceName, 'sum(rate(accepted_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Accepted Token Rate'),

    g.query.prometheus.new(datasourceName, 'sum(rate(rejected_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters })
    + g.query.prometheus.withIntervalFactor(1)
    + g.query.prometheus.withLegendFormat('Rejected Token Rate'),
  ];
  local row8 = [
    timeSeriesPanel(
      'Incoming Token Rate',
      datasourceName,
      query='sum(rate(incoming_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: stringFilters },
      axisLabel='Token Rate',
      x=0,
      h=8,
      w=12
    ),
    timeSeriesPanel(
      'Accepted Token Rate vs Rejected Token Rate',
      datasourceName,
      targets=acceptedVsRejectedTargets,
      axisLabel='Token Rate',
      x=12,
      h=8,
      w=12
    ),
  ];

  local row9 = [
    statPanel(
      'Total Incoming Tokens',
      datasourceName,
      'sum(increase(incoming_tokens_total{%(filters)s}[$__range]))' % { filters: stringFilters },
      x=0,
      h=10,
      w=8,
      panelColor='blue',
      graphMode='area',
      unit='short'
    ),
    statPanel(
      'Total Accepted Tokens',
      datasourceName,
      'sum(increase(accepted_tokens_total{%(filters)s}[$__range]))' % { filters: stringFilters },
      x=8,
      h=10,
      w=8,
      graphMode='area',
      unit='short'
    ),
    statPanel(
      'Total Rejected Tokens',
      datasourceName,
      'sum(increase(rejected_tokens_total{%(filters)s}[$__range]))' % { filters: stringFilters },
      x=16,
      h=10,
      w=8,
      panelColor='red',
      graphMode='area',
      unit='short'
    ),
  ];

  if 'fairness_label_key' in scheduler && scheduler.fairness_label_key != ''
  then
    local row10 = [
      timeSeriesPanel(
        'Tokens adjusted per second to maintain fairness within a workload',
        datasourceName,
        query='sum by (workload_index) (rate(fairness_delayed_tokens_sum{%(filters)s}[$__rate_interval]))' % { filters: stringFilters },
        axisLabel='Token Rate',
      ),
    ];
    [row1, row2, row3, row4, row5, row6, row7, row8, row9, row10]
  else
    [row1, row2, row3, row4, row5, row6, row7, row8, row9]
