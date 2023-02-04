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

function(params) {
  _config:: config.common + config.dashboard + params,

  local ds = $._config.datasource.name,

  local rateLimiterPanel =
    newTimeSeriesPanel('Rate Limiter',
                       ds,
                       'sum by(decision_type) (rate(rate_limiter_counter{policy_name="%(policy_name)s"}[$__rate_interval]))' % { policy_name: $._config.policy_name },
                       'Decisions',
                       'reqps'),

  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja - Rate Limiter',
      editable=true,
      schemaVersion=18,
      refresh=$._config.refresh_interval,
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
        regex: $._config.datasource.filter_regex,
        type: 'datasource',
      }
    )
    .addPanel(rateLimiterPanel, gridPos={ h: 10, w: 24, x: 0, y: 0 }),

  dashboard: dashboardDef,
}
