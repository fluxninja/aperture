local utils = import '../../../policies/policy-utils.libsonnet';
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
  local base_filters = params.dashboard.extra_filters { policy_name: policyName },
  local throughput_filters = utils.dictToPrometheusFilter(base_filters),
  local acceptPercentage_filters = utils.dictToPrometheusFilter(base_filters { signal_name: 'ACCEPT_PERCENTAGE' }),

  local throughputPanel =
    graphPanel.new(
      title='Throughput - Accept/Reject',
      datasource=dsName,
      interval='30',
    )
    .addTarget(
      prometheus.target(
        expr='rate(sampler_counter_total{%(filters)s}[$__rate_interval])' % {
          filters: throughput_filters,
        },
        intervalFactor=1,
      ),
    ),


  local acceptPercentagePanel =
    graphPanel.new(
      title='Accept Percentage',
      datasource=dsName,
      interval='30s',
    )
    .addTarget(
      prometheus.target(
        expr=|||
          increase(signal_reading_sum{%(filters)s}[$__rate_interval])
          /
          increase(signal_reading_count{%(filters)s}[$__rate_interval])
        ||| % { filters: acceptPercentage_filters },
        intervalFactor=1,
      ),
    ),


  local dashboardDef =
    dashboard.new(
      title=params.dashboard.title,
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
