local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local timeSeriesPanel = lib.TimeSeriesPanel;


function(params) {
  _config:: config.common + config.dashboard + params,

  local ds = $._config.datasource.name,

  local SignalAveragePanel =
    local query = |||
      increase(signal_reading_sum{policy_name="%(policyName)s",signal_name="${signal_name}",valid="true"}[$__rate_interval])
      /
      increase(signal_reading_count{policy_name="%(policyName)s",signal_name="${signal_name}",valid="true"}[$__rate_interval])
    ||| % { policyName: $._config.policyName };
    local target =
      grafana.prometheus.target(query) +
      {
        legendFormat: 'Avg',
        editorMode: 'code',
        range: true,
      };
    local thresholds =
      {
        mode: 'absolute',
        steps: [
          { color: 'green', value: null },
          { color: 'red', value: 80 },
        ],
      };
    aperture.timeSeriesPanel.new('Signal Average', ds)
    + timeSeriesPanel.withTarget(target)
    + timeSeriesPanel.defaults.withThresholds(thresholds)
    + timeSeriesPanel.withFieldConfigMixin(
      timeSeriesPanel.fieldConfig.withDefaultsMixin(
        timeSeriesPanel.fieldConfig.defaults.withThresholds(thresholds)
      )
    ) + {
      interval: '1s',
    },

  local InvalidFrequencyPanel =
    local query = |||
      avg by (valid) (rate(signal_reading_count{policy_name="%(policyName)s",signal_name="${signal_name}"}[$__rate_interval]))
    ||| % { policyName: $._config.policyName };
    local target =
      grafana.prometheus.target(query) +
      {
        editorMode: 'code',
        range: true,
      };
    local thresholds =
      {
        mode: 'absolute',
        steps: [
          { color: 'green', value: null },
          { color: 'red', value: 80 },
        ],
      };
    aperture.timeSeriesPanel.new('Signal Validity (Frequency)', ds)
    + timeSeriesPanel.withTarget(target)
    + timeSeriesPanel.defaults.withThresholds(thresholds)
    + timeSeriesPanel.withFieldConfigMixin(
      timeSeriesPanel.fieldConfig.withDefaultsMixin(
        timeSeriesPanel.fieldConfig.defaults.withThresholds(thresholds)
      )
    ) + {
      interval: '1s',
    },

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
        uid: '${datasource}',
      },
      query: 'label_values(signal_reading{policy_name="%(policyName)s"}, signal_name)' % { policyName: $._config.policyName },
      hide: 0,
      includeAll: false,
      multi: false,
      name: 'signal_name',
      options: [],
      refresh: 1,
      regex: '',
      skipUrlSync: false,
      sort: 0,
      type: 'query',
    })
    .addPanel(SignalAveragePanel, gridPos={ h: 10, w: 24, x: 0, y: 0 })
    .addPanel(InvalidFrequencyPanel, gridPos={ h: 10, w: 24, x: 0, y: 10 }),
}
