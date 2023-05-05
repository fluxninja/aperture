local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
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

function(cfg) {
  local params = config.common + config.dashboard + cfg,
  local policyName = params.policy_name,
  local ds = params.datasource,
  local dsName = ds.name,

  local rateLimiterPanel =
    newTimeSeriesPanel('Rate Limiter',
                       dsName,
                       'sum by(decision_type) (rate(rate_limiter_counter{policy_name="%(policy_name)s"}[$__rate_interval]))' % { policy_name: policyName },
                       'Decisions',
                       'reqps'),

  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja - Rate Limiter',
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
    .addPanel(rateLimiterPanel, gridPos={ h: 10, w: 24, x: 0, y: 0 }),

  dashboard: dashboardDef,
}
