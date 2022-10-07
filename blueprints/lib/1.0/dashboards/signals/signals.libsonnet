local lib = import '../../grafana/grafana.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;


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

  dashboard:
    dashboard.new(
      title='Signals',
      editable=true,
      schema=18,
      refresh=$._config.refreshInterval,
      time_from=$._config.timeFrom,
      time_to=$._config.timeTo
    )
    .addTemplate(
      {
        current: {
          text: 'default',
          value: $._config.dataSource.name,
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
    ),
}
