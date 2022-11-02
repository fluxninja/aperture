local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local tablePanel = grafana.tablePanel;
local barGaugePanel = grafana.barGaugePanel;
local statPanel = grafana.statPanel;
local annotation = grafana.annotation;
local timeSeriesPanel = lib.TimeSeriesPanel;

local newTimeSeriesPanel(title, datasource, query, axisLabel='', unit='') =
  local thresholds =
    {
      mode: 'absolute',
      steps: [
        { color: 'green', value: null },
        { color: 'red', value: 80 },
      ],
    };
  local target =
    prometheus.target(query, intervalFactor=1)
    + { range: true, editorMode: 'code' };
  aperture.timeSeriesPanel.new(title, datasource, axisLabel, unit)
  + timeSeriesPanel.withTarget(target)
  + timeSeriesPanel.defaults.withThresholds(thresholds)
  + timeSeriesPanel.withFieldConfigMixin(
    timeSeriesPanel.fieldConfig.withDefaultsMixin(
      timeSeriesPanel.fieldConfig.defaults.withThresholds(thresholds)
    )
  ) + {
    interval: '1s',
  };

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

function(params) {
  _config:: config.common + config.dashboard + params,

  local p = 'service_latency',
  local ds = $._config.datasource.name,

  local fluxMeterPanel =
    newTimeSeriesPanel('FluxMeter',
                       ds,
                       |||
                         sum(increase(flux_meter_sum{valid="true", response_status="OK", flux_meter_name="%(policyName)s"}[$__rate_interval]))
                         / sum(increase(flux_meter_count{valid="true", response_status="OK", flux_meter_name="%(policyName)s"}[$__rate_interval]))
                       ||| % { policyName: $._config.policyName },
                       'Latency (ms)',
                       'ms'),
  local WFQSchedulerFlows =
    newBarGaugePanel('WFQ Scheduler Flows', ds, 'avg(wfq_flows_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }),

  local WFQSchedulerHeapRequests =
    newBarGaugePanel('WFQ Scheduler Heap Requests', ds, 'avg(wfq_requests_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }),

  local TotalBucketLoadSchedFactor =
    newStatPanel('Average Load Multiplier', ds, 'avg(token_bucket_lm_ratio{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }),

  local TokenBucketBucketCapacity =
    newStatPanel('Token Bucket Bucket Capacity', ds, 'avg(token_bucket_capacity_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName })
    + {
      options+: {
        orientation: 'auto',
      },
    },

  local TokenBucketBucketFillRate =
    newStatPanel('Token Bucket Bucket FillRate', ds, 'avg(token_bucket_fill_rate{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }) +
    {
      options+: {
        orientation: 'auto',
      },
    },

  local TokenBucketAvailableTokens =
    newStatPanel('Token Bucket Available Tokens', ds, 'avg(token_bucket_available_tokens_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }) +
    {
      options+: {
        orientation: 'auto',
      },
    },

  local IncomingConcurrency =
    newTimeSeriesPanel('Incoming Concurrency', ds, 'sum(rate(incoming_concurrency_ms{policy_name="%(policyName)s"}[$__rate_interval]))' % { policyName: $._config.policyName }, 'Concurrency', 'ms'),

  local AcceptedConcurrency =
    newTimeSeriesPanel('Accepted Concurrency', ds, 'sum(rate(accepted_concurrency_ms{policy_name="%(policyName)s"}[$__rate_interval]))' % { policyName: $._config.policyName }, 'Concurrency', 'ms'),

  local WorkloadDecisions =
    newTimeSeriesPanel('Workload Decisions', ds, 'sum by(workload_index, decision_type) (rate(workload_requests_total{policy_name="%(policyName)s"}[$__rate_interval]))' % { policyName: $._config.policyName }, 'Decisions', 'reqps'),

  local WorkloadLatency =
    newTimeSeriesPanel('Workload Latency (Auto Tokens)', ds, '(sum by (workload_index) (increase(workload_latency_ms_sum{policy_name="%(policyName)s"}[$__rate_interval])))/(sum by (workload_index) (increase(workload_latency_ms_count{policy_name="%(policyName)s"}[$__rate_interval])))' % { policyName: $._config.policyName }, 'Latency', 'ms'),


  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja',
      editable=true,
      schemaVersion=18,
      refresh=$._config.refreshInterval,
      time_from='now-5m',
      time_to='now'
    )
    .addTemplate(
      {
        current: {
          text: 'default',
          value: $._config.datasource.name,
        },
        hide: 0,
        label: 'Data Source',
        name: 'datasource',
        options: [],
        query: 'prometheus',
        refres: 1,
        regex: $._config.datasource.filterRegex,
        type: 'datasource',
      }
    )
    .addPanel(fluxMeterPanel, gridPos={ h: 10, w: 24, x: 0, y: 0 })
    .addPanel(WorkloadDecisions, gridPos={ h: 10, w: 24, x: 0, y: 10 })
    .addPanel(WorkloadLatency, gridPos={ h: 10, w: 24, x: 0, y: 20 })
    .addPanel(IncomingConcurrency, gridPos={ h: 8, w: 12, x: 0, y: 30 })
    .addPanel(AcceptedConcurrency, gridPos={ h: 8, w: 12, x: 12, y: 30 })
    .addPanel(WFQSchedulerFlows, gridPos={ h: 3, w: 8, x: 0, y: 40 })
    .addPanel(TotalBucketLoadSchedFactor, gridPos={ h: 6, w: 4, x: 8, y: 40 })
    .addPanel(TokenBucketBucketCapacity, gridPos={ h: 6, w: 4, x: 12, y: 40 })
    .addPanel(TokenBucketBucketFillRate, gridPos={ h: 6, w: 4, x: 16, y: 40 })
    .addPanel(TokenBucketAvailableTokens, gridPos={ h: 6, w: 4, x: 20, y: 40 })
    .addPanel(WFQSchedulerHeapRequests, gridPos={ h: 3, w: 8, x: 0, y: 40 }),

  dashboard: dashboardDef,
}
