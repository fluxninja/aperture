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

local signalAveragePanel(title, datasource, filtersDict) =
  local filters = utils.dictToPrometheusFilter(filtersDict);
  local query = |||
    increase(signal_reading_sum{%(filters)s}[$__rate_interval])
    /
    increase(signal_reading_count{%(filters)s}[$__rate_interval])
  ||| % { filters: filters };
  newGraphPanel(title, datasource, query);

local signalFrequencyPanel(title, datasource, filtersDict) =
  local filters = utils.dictToPrometheusFilter(filtersDict);
  local query = |||
    avg by (valid) (rate(signal_reading_count{%(filters)s}[$__rate_interval]))
  ||| % { filters: filters };
  newGraphPanel(title, datasource, query);

local dashboardWithPanels(dashboardParams, policyName, extra_filters) =
  local datasource = dashboardParams.datasource;
  local dsName = datasource.name;

  local baseFilters = extra_filters { policy_name: policyName };
  local actualScaleFilters = baseFilters { signal_name: 'ACTUAL_SCALE' };
  local configuredScaleFilters = baseFilters { signal_name: 'CONFIGURED_SCALE' };
  local desiredScaleFilters = baseFilters { signal_name: 'DESIRED_SCALE' };

  local actualScaleAverage = signalAveragePanel('Actual Scale Average', dsName, actualScaleFilters);
  local configuredScaleAverage = signalAveragePanel('Configured Scale Average', dsName, configuredScaleFilters);
  local desiredScaleAverage = signalAveragePanel('Desired Scale Average', dsName, desiredScaleFilters);

  local actualScaleFrequency = signalFrequencyPanel('Actual Scale Validity (Frequency)', dsName, actualScaleFilters);
  local configuredScaleFrequency = signalFrequencyPanel('Configured Scale Validity (Frequency)', dsName, configuredScaleFilters);
  local desiredScaleFrequency = signalFrequencyPanel('Desired Scale Validity (Frequency)', dsName, desiredScaleFilters);

  dashboard.new(
    title='Aperture Auto-scale',
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
      refresh: 1,
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

  local dashboardDef = dashboardWithPanels(params.dashboard, policyName, params.dashboard.extra_filters),

  dashboard: dashboardDef,
}
