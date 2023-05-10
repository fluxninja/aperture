local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local graphPanel = grafana.graphPanel;

function(cfg) {
  local params = config + cfg,
  local policyName = params.policy.policy_name,
  local ds = params.dashboard.datasource,
  local dsName = ds.name,
  local refresh = params.dashboard.refresh_interval,
  local time_from = params.dashboard.time_from,
  local time_to = params.dashboard.time_to,

  local signalAveragePanel =
    graphPanel.new(
      title='Signal Average',
      datasource=dsName,
    )
    .addTarget(
      prometheus.target(
        expr=(
          'increase(signal_reading_sum{policy_name="' + policyName + '",signal_name="${signal_name}",sub_circuit_id="${sub_circuit_id}"}[$__rate_interval])' +
          '/' +
          'increase(signal_reading_count{policy_name="' + policyName + '",signal_name="${signal_name}",sub_circuit_id="${sub_circuit_id}"}[$__rate_interval])'
        ),
        intervalFactor=1,
      ),
    ),

  local InvalidFrequencyPanel =
    graphPanel.new(
      title='Signal Validity (Frequency)',
      datasource=dsName,
      stack=true,
      bars=true,
    )
    .addTarget(
      prometheus.target(
        expr='avg(rate(signal_reading_count{policy_name="%(policy_name)s",signal_name="${signal_name}",sub_circuit_id="${sub_circuit_id}"}[$__rate_interval]))' % { policy_name: policyName },
        intervalFactor=1,
        legendFormat='Valid',
      ),
    )
    .addTarget(
      prometheus.target(
        expr='sum(rate(invalid_signal_readings_total{policy_name="%(policy_name)s",signal_name="${signal_name}",sub_circuit_id="${sub_circuit_id}"}[$__rate_interval]))' % { policy_name: policyName },
        intervalFactor=1,
        legendFormat='Invalid',
      ),
    ),

  local dashboardDef =
    dashboard.new(
      title='Signals',
      schemaVersion=36,
      editable=true,
      refresh=refresh,
      time_from=time_from,
      time_to=time_to,
    )
    .addTemplate(
      grafana.template.datasource(
        name='datasource',
        query='prometheus',
        label='Data Source',
        current=dsName,
        hide=0,
        regex=ds.filter_regex,
      )
    )
    .addTemplate({
      datasource: {
        type: 'prometheus',
        uid: '${datasource}',
      },
      query: 'label_values(signal_reading{policy_name="%(policy_name)s"}, signal_name)' % { policy_name: policyName },
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
      label: 'Signal Name',
    })
    .addTemplate({
      datasource: {
        type: 'prometheus',
        uid: '${datasource}',
      },
      query: 'label_values(signal_reading{policy_name="%(policy_name)s",signal_name="${signal_name}"}, sub_circuit_id)' % { policy_name: policyName },
      hide: 0,
      includeAll: false,
      multi: false,
      name: 'sub_circuit_id',
      options: [],
      refresh: 1,
      regex: '',
      skipUrlSync: false,
      sort: 0,
      type: 'query',
      label: 'Sub Circuit ID',
    })
    .addPanel(
      panel=signalAveragePanel,
      gridPos={ x: 0, y: 0, w: 24, h: 10 },
    )
    .addPanel(
      panel=InvalidFrequencyPanel,
      gridPos={ x: 0, y: 15, w: 24, h: 10 },
    ),

  dashboard: dashboardDef,
}
