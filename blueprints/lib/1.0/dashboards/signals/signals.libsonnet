local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local timeSeriesPanel = lib.TimeSeriesPanel;


local defaults = {
  refreshInterval: '30s',
  timeFrom: 'now-30m',
  timeTo: 'now',
  datasource: {
    name: '$datasource',
    filterRegex: '',
  },
};

function(params) {
  _config:: defaults + params,

  local ds = $._config.datasource.name,

  local SignalAveragePanel =
    local thresholds =
      {
        mode: 'absolute',
        steps: [
          { color: 'green', value: null },
          { color: 'red', value: 80 },
        ],
      };
    local prometheusQuery = |||
      increase(signal_reading_sum{signal_name=\"${signal_name}\"}[10s])
      /
      increase(signal_reading_count{signal_name=\"${signal_name}\"}[10s])
    |||;
    aperture.timeSeriesPanel.new('Signal Average', ds, prometheusQuery)
    + timeSeriesPanel.defaults.withThresholds(thresholds)
    + timeSeriesPanel.withFieldConfigMixin(
      timeSeriesPanel.fieldConfig.withDefaultsMixin(
        timeSeriesPanel.fieldConfig.defaults.withThresholds(thresholds)
      )
    ),

  dashboard:
    dashboard.new(
      title='Signals',
      editable=true,
      schemaVersion=18,
      refresh=$._config.refreshInterval,
      time_from=$._config.timeFrom,
      time_to=$._config.timeTo
    )
    .addTemplate(
      {
        current: {
          text: 'default',
          value: $._config.datasource.name,
        },
        regex: $._config.datasource.filterRegex,
        hide: 0,
        label: 'Data Source',
        name: 'datasource',
        options: [],
        query: 'prometheus',
        refresh: 1,
        type: 'datasource',
      }
    )
    .addTemplate({
      current: {
        selected: false,
        text: 'ACCEPTED_CONCURRENCY',
        value: 'ACCEPTED_CONCURRENCY',
      },
      datasource: {
        type: 'prometheus',
      },
      definition: 'label_values(signal_reading, signal_name)',
      hide: 0,
      includeAll: false,
      multi: false,
      name: 'signal_name',
      options: [],
      query: {
        query: 'label_values(signal_reading, signal_name)',
        refId: 'StandardVariableQuery',
      },
      refresh: 1,
      regex: '',
      skipUrlSync: false,
      sort: 0,
      type: 'query',
    })
    .addPanel(SignalAveragePanel, gridPos={ h: 7, w: 23, x: 0, y: 0 }),
}
