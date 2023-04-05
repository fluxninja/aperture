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

local signalAveragePanel(title, datasource, signalName, policyName) =
  local query = |||
    increase(signal_reading_sum{policy_name="%(policy_name)s",signal_name="%(signal_name)s",valid="true"}[$__rate_interval])
    /
    increase(signal_reading_count{policy_name="%(policy_name)s",signal_name="%(signal_name)s",valid="true"}[$__rate_interval])
  ||| % { policy_name: policyName, signal_name: signalName };
  newTimeSeriesPanel(title, datasource, query);

local signalFrequencyPanel(title, datasource, signalName, policyName) =
  local query = |||
    avg by (valid) (rate(signal_reading_count{policy_name="%(policy_name)s",signal_name="%(signal_name)s"}[$__rate_interval]))
  ||| % { policy_name: policyName, signal_name: signalName };
  newTimeSeriesPanel(title, datasource, query);

function(cfg) {
  local p = 'service_latency',
  local params = config.common + config.dashboard + cfg,
  local policyName = params.policy_name,
  local ds = params.datasource,
  local dsName = ds.name,

  local actualScaleAverage = signalAveragePanel('Actual Scale Average', dsName, 'ACTUAL_REPLICAS', policyName),
  local configuredScaleAverage = signalAveragePanel('Configured Scale Average', dsName, 'CONFIGURED_REPLICAS', policyName),
  local desiredScaleAverage = signalAveragePanel('Desired Scale Average', dsName, 'DESIRED_REPLICAS', policyName),

  local actualScaleFrequency = signalFrequencyPanel('Actual Scale Validity (Frequency)', dsName, 'ACTUAL_REPLICAS', policyName),
  local configuredScaleFrequency = signalFrequencyPanel('Configured Scale Validity (Frequency)', dsName, 'CONFIGURED_REPLICAS', policyName),
  local desiredScaleFrequency = signalFrequencyPanel('Desired Scale Validity (Frequency)', dsName, 'DESIRED_REPLICAS', policyName),

  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja - Kubernetes Auto Scaler',
      editable=true,
      schemaVersion=18,
      refresh=params.refresh_interval,
      time_from=params.time_from,
      time_to=params.time_to
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
    .addPanel(actualScaleAverage, gridPos={ h: 10, w: 24, x: 0, y: 0 })
    .addPanel(actualScaleFrequency, gridPos={ h: 10, w: 24, x: 0, y: 10 })
    .addPanel(configuredScaleAverage, gridPos={ h: 10, w: 24, x: 0, y: 20 })
    .addPanel(configuredScaleFrequency, gridPos={ h: 10, w: 24, x: 0, y: 30 })
    .addPanel(desiredScaleAverage, gridPos={ h: 10, w: 24, x: 0, y: 40 })
    .addPanel(desiredScaleFrequency, gridPos={ h: 10, w: 24, x: 0, y: 50 }),

  dashboard: dashboardDef,
}
