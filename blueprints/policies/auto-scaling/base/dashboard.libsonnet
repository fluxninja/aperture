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

local signalAveragePanel(title, datasource, signalName, policyName) =
  local query = |||
    increase(signal_reading_sum{policy_name="%(policy_name)s",signal_name="%(signal_name)s"}[$__rate_interval])
    /
    increase(signal_reading_count{policy_name="%(policy_name)s",signal_name="%(signal_name)s"}[$__rate_interval])
  ||| % { policy_name: policyName, signal_name: signalName };
  newGraphPanel(title, datasource, query);

local signalFrequencyPanel(title, datasource, signalName, policyName) =
  local query = |||
    avg by (valid) (rate(signal_reading_count{policy_name="%(policy_name)s",signal_name="%(signal_name)s"}[$__rate_interval]))
  ||| % { policy_name: policyName, signal_name: signalName };
  newGraphPanel(title, datasource, query);

local dashboardWithPanels(dashboardParams, policyName) =
  local datasource = dashboardParams.datasource;
  local dsName = datasource.name;

  local actualScaleAverage = signalAveragePanel('Actual Scale Average', dsName, 'ACTUAL_SCALE', policyName);
  local configuredScaleAverage = signalAveragePanel('Configured Scale Average', dsName, 'CONFIGURED_SCALE', policyName);
  local desiredScaleAverage = signalAveragePanel('Desired Scale Average', dsName, 'DESIRED_SCALE', policyName);

  local actualScaleFrequency = signalFrequencyPanel('Actual Scale Validity (Frequency)', dsName, 'ACTUAL_SCALE', policyName);
  local configuredScaleFrequency = signalFrequencyPanel('Configured Scale Validity (Frequency)', dsName, 'CONFIGURED_SCALE', policyName);
  local desiredScaleFrequency = signalFrequencyPanel('Desired Scale Validity (Frequency)', dsName, 'DESIRED_SCALE', policyName);

  dashboard.new(
    title='Jsonnet / FluxNinja',
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
  .addPanel(actualScaleAverage, gridPos={ h: 10, w: 24, x: 0, y: 0 })
  .addPanel(actualScaleFrequency, gridPos={ h: 10, w: 24, x: 0, y: 10 })
  .addPanel(configuredScaleAverage, gridPos={ h: 10, w: 24, x: 0, y: 20 })
  .addPanel(configuredScaleFrequency, gridPos={ h: 10, w: 24, x: 0, y: 30 })
  .addPanel(desiredScaleAverage, gridPos={ h: 10, w: 24, x: 0, y: 40 })
  .addPanel(desiredScaleFrequency, gridPos={ h: 10, w: 24, x: 0, y: 50 });


function(cfg) {
  local params = config + cfg,
  local policyName = params.policy.policy_name,

  local dashboardDef = dashboardWithPanels(params.dashboard, policyName),

  dashboard: dashboardDef,
}
