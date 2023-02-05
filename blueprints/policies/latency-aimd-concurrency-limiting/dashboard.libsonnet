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

function(cfg) {
  local p = 'service_latency',
  local params = config.common + config.dashboard + cfg,
  local policyName = params.policy_name,
  local ds = params.datasource,
  local dsName = ds.name,

  local fluxMeterPanel =
    newTimeSeriesPanel('FluxMeter',
                       dsName,
                       |||
                         sum(increase(flux_meter_sum{valid="true", flow_status="OK", flux_meter_name="%(policy_name)s"}[$__rate_interval]))
                         / sum(increase(flux_meter_count{valid="true", flow_status="OK", flux_meter_name="%(policy_name)s"}[$__rate_interval]))
                       ||| % { policy_name: policyName },
                       'Latency (ms)',
                       'ms'),
  local WFQSchedulerFlows =
    newBarGaugePanel('WFQ Scheduler Flows', dsName, 'avg(wfq_flows_total{policy_name="%(policy_name)s"})' % { policy_name: policyName }),

  local WFQSchedulerHeapRequests =
    newBarGaugePanel('WFQ Scheduler Heap Requests', dsName, 'avg(wfq_requests_total{policy_name="%(policy_name)s"})' % { policy_name: policyName }),

  local TotalBucketLoadSchedFactor =
    newStatPanel('Average Load Multiplier', dsName, 'avg(token_bucket_lm_ratio{policy_name="%(policy_name)s"})' % { policy_name: policyName }),

  local TokenBucketBucketCapacity =
    newStatPanel('Token Bucket Bucket Capacity', dsName, 'avg(token_bucket_capacity_total{policy_name="%(policy_name)s"})' % { policy_name: policyName })
    + {
      options+: {
        orientation: 'auto',
      },
    },

  local TokenBucketBucketFillRate =
    newStatPanel('Token Bucket Bucket FillRate', dsName, 'avg(token_bucket_fill_rate{policy_name="%(policy_name)s"})' % { policy_name: policyName }) +
    {
      options+: {
        orientation: 'auto',
      },
    },

  local TokenBucketAvailableTokens =
    newStatPanel('Token Bucket Available Tokens', dsName, 'avg(token_bucket_available_tokens_total{policy_name="%(policy_name)s"})' % { policy_name: policyName }) +
    {
      options+: {
        orientation: 'auto',
      },
    },

  local IncomingConcurrency =
    newTimeSeriesPanel('Incoming Concurrency', dsName, 'sum(rate(incoming_work_seconds_total{policy_name="%(policy_name)s"}[$__rate_interval]))' % { policy_name: policyName }, 'Concurrency', 'none'),

  local AcceptedConcurrency =
    newTimeSeriesPanel('Accepted Concurrency', dsName, 'sum(rate(accepted_work_seconds_total{policy_name="%(policy_name)s"}[$__rate_interval]))' % { policy_name: policyName }, 'Concurrency', 'none'),

  local WorkloadDecisions =
    newTimeSeriesPanel('Workload Decisions', dsName, 'sum by(workload_index, decision_type) (rate(workload_requests_total{policy_name="%(policy_name)s"}[$__rate_interval]))' % { policy_name: policyName }, 'Decisions', 'reqps'),

  local WorkloadLatency =
    newTimeSeriesPanel('Workload Latency (Auto Tokens)', dsName, '(sum by (workload_index) (increase(workload_latency_ms_sum{policy_name="%(policy_name)s"}[$__rate_interval])))/(sum by (workload_index) (increase(workload_latency_ms_count{policy_name="%(policy_name)s"}[$__rate_interval])))' % { policy_name: policyName }, 'Latency', 'ms'),


  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja',
      editable=true,
      schemaVersion=18,
      refresh=params.refresh_interval,
      time_from='now-5m',
      time_to='now'
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
        regex: ds.filter_regex,
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
