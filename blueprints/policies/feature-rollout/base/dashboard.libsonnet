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

  local throughputPanel =
    graphPanel.new(
      title='Throughput - Accept/Reject',
      datasource=dsName,
    )
    .addTarget(
      prometheus.target(
        expr='rate(regulator_counter{policy_name="%(policy_name)s"}[$__rate_interval])' % {
          policy_name: policyName,
        },
        intervalFactor=1,
      ),
    ),


  local acceptPercentagePanel =
    graphPanel.new(
      title='Accept Percentage',
      datasource=dsName,
    )
    .addTarget(
      prometheus.target(
        expr=(
          'increase(signal_reading_sum{policy_name="' + policyName + '",signal_name="ACCEPT_PERCENTAGE"}[$__rate_interval])' +
          '/' +
          'increase(signal_reading_count{policy_name="' + policyName + '",signal_name="ACCEPT_PERCENTAGE"}[$__rate_interval])'
        ),
        intervalFactor=1,
      ),
    ),


  local dashboardDef =
    dashboard.new(
      title='Feature Rollout',
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
    .addPanel(
      panel=throughputPanel,
      gridPos={ x: 0, y: 0, w: 24, h: 8 },
    )
    .addPanel(
      panel=acceptPercentagePanel,
      gridPos={ x: 0, y: 15, w: 24, h: 8 },
    ),

  dashboard: dashboardDef,
}
