local utils = import '../../policy-utils.libsonnet';
local config = import './config-defaults.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local barGaugePanel = grafana.barGaugePanel;
local statPanel = grafana.statPanel;
local graphPanel = grafana.graphPanel;

local newGraphPanel(title, datasource, query, axisLabel='', unit='') =
  graphPanel.new(
    title=title,
    datasource=datasource,
    labelY1=axisLabel,
    formatY1=unit,
  )
  .addTarget(
    prometheus.target(
      expr=query,
      intervalFactor=1,
    )
  );

local newBarGaugePanel(graphTitle, datasource, graphQuery) =
  local target =
    prometheus.target(graphQuery) +
    {
      legendFormat: '{{ instance }} - {{ policy_name }}',
      format: 'time_series',
      instant: false,
      range: true,
    };

  barGaugePanel.new(
    title=graphTitle,
    datasource=datasource,
  ).addTarget(target) +
  {
    fieldConfig: {
      defaults: {
        color: {
          mode: 'thresholds',
        },
        mappings: [],
        thresholds: {
          mode: 'absolute',
          steps: [
            { color: 'green', value: null },
          ],
        },
      },
      overrides: [],
    },
    options: {
      displayMode: 'gradient',
      minVizHeight: 10,
      minVizWidth: 0,
      orientation: 'horizontal',
      reduceOptions: {
        calcs: ['lastNotNull'],
        fields: '',
        values: false,
      },
      showUnfilled: true,
    },
  };

local newStatPanel(graphTitle, datasource, graphQuery) =
  local target =
    prometheus.target(graphQuery) +
    {
      legendFormat: '{{ instance }} - {{ policy_name }}',
      editorMode: 'code',
      range: true,
    };

  statPanel.new(
    title=graphTitle,
    datasource=datasource,
  ).addTarget(target) +
  {
    fieldConfig: {
      defaults: {
        color: {
          mode: 'thresholds',
        },
        mappings: [],
        thresholds: {
          mode: 'absolute',
          steps: [
            { color: 'green', value: null },
          ],
        },
      },
      overrides: [],
    },
    options: {
      colorMode: 'value',
      graphMode: 'none',
      justifyMode: 'center',
      orientation: 'horizontal',
      reduceOptions: {
        calcs: ['lastNotNull'],
        fields: '',
        values: false,
      },
      textMode: 'auto',
    },
  };

local dashboardWithPanels(dashboardParams, filters) =
  local datasource = dashboardParams.datasource;
  local dsName = datasource.name;

  local WFQSchedulerFlows =
    newBarGaugePanel('WFQ Scheduler Flows', dsName, 'avg(wfq_flows_total{%(filters)s})' % { filters: filters });

  local WFQSchedulerHeapRequests =
    newBarGaugePanel('WFQ Scheduler Heap Requests', dsName, 'avg(wfq_requests_total{%(filters)s})' % { filters: filters });

  local TotalBucketLoadSchedFactor =
    newStatPanel('Average Load Multiplier', dsName, 'avg(token_bucket_lm_ratio{%(filters)s})' % { filters: filters });

  local TokenBucketBucketCapacity =
    newStatPanel('Token Bucket Bucket Capacity', dsName, 'avg(token_bucket_capacity_total{%(filters)s})' % { filters: filters })
    + {
      options+: {
        orientation: 'auto',
      },
    };

  local TokenBucketBucketFillRate =
    newStatPanel('Token Bucket Bucket FillRate', dsName, 'avg(token_bucket_fill_rate{%(filters)s})' % { filters: filters }) +
    {
      options+: {
        orientation: 'auto',
      },
    };

  local TokenBucketAvailableTokens =
    newStatPanel('Token Bucket Available Tokens', dsName, 'avg(token_bucket_available_tokens_total{%(filters)s})' % { filters: filters }) +
    {
      options+: {
        orientation: 'auto',
      },
    };

  local IncomingConcurrency =
    newGraphPanel('Incoming Token Rate', dsName, 'sum(rate(incoming_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: filters }, 'Concurrency', 'none');

  local AcceptedConcurrency =
    newGraphPanel('Accepted Token Rate', dsName, 'sum(rate(accepted_tokens_total{%(filters)s}[$__rate_interval]))' % { filters: filters }, 'Concurrency', 'none');

  local WorkloadDecisionsAccepted =
    newGraphPanel('Workload Decisions (accepted)', dsName, 'sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_ACCEPTED"}[$__rate_interval]))' % { filters: filters }, 'Decisions', 'reqps');

  local WorkloadDecisionsRejected =
    newGraphPanel('Workload Decisions (rejected)', dsName, 'sum by(workload_index, decision_type) (rate(workload_requests_total{%(filters)s,decision_type="DECISION_TYPE_REJECTED"}[$__rate_interval]))' % { filters: filters }, 'Decisions', 'reqps');

  local WorkloadLatency =
    newGraphPanel('Workload Latency', dsName, '(sum by (workload_index) (increase(workload_latency_ms_sum{%(filters)s}[$__rate_interval])))/(sum by (workload_index) (increase(workload_latency_ms_count{%(filters)s}[$__rate_interval])))' % { filters: filters }, 'Latency', 'ms');

  dashboard.new(
    title='Aperture Service Protection',
    editable=true,
    schemaVersion=18,
    refresh=dashboardParams.refresh_interval,
    time_from=dashboardParams.time_from,
    time_to=dashboardParams.time_to
  )
  .addTemplate(
    {
      current: {
        text: 'default',
        value: dsName,
      },
      hide: 0,
      label: 'Data Source',
      name: 'datasource',
      options: [],
      query: 'prometheus',
      refres: 1,
      regex: datasource.filter_regex,
      type: 'datasource',
    }
  )
  .addPanel(WorkloadDecisionsAccepted, gridPos={ h: 10, w: 24, x: 0, y: 10 })
  .addPanel(WorkloadDecisionsRejected, gridPos={ h: 10, w: 24, x: 0, y: 20 })
  .addPanel(WorkloadLatency, gridPos={ h: 10, w: 24, x: 0, y: 30 })
  .addPanel(IncomingConcurrency, gridPos={ h: 8, w: 12, x: 0, y: 40 })
  .addPanel(AcceptedConcurrency, gridPos={ h: 8, w: 12, x: 12, y: 40 })
  .addPanel(WFQSchedulerFlows, gridPos={ h: 3, w: 8, x: 0, y: 50 })
  .addPanel(TotalBucketLoadSchedFactor, gridPos={ h: 6, w: 4, x: 8, y: 50 })
  .addPanel(TokenBucketBucketCapacity, gridPos={ h: 6, w: 4, x: 12, y: 50 })
  .addPanel(TokenBucketBucketFillRate, gridPos={ h: 6, w: 4, x: 16, y: 50 })
  .addPanel(TokenBucketAvailableTokens, gridPos={ h: 6, w: 4, x: 20, y: 50 })
  .addPanel(WFQSchedulerHeapRequests, gridPos={ h: 3, w: 8, x: 0, y: 50 });


function(cfg) {
  local params = config + cfg,
  local policyName = params.policy.policy_name,
  local filters = utils.dictToPrometheusFilter(params.dashboard.extra_filters { policy_name: policyName }),

  local dashboardDef = dashboardWithPanels(params.dashboard, filters),

  dashboard: dashboardDef,
}
