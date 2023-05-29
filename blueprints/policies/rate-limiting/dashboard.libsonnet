local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local utils = import '../policy-utils.libsonnet';
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

  local stringFilters = utils.dictToPrometheusFilter(params.dashboard.extra_filters { policy_name: policyName }),

  local rateLimiterPanel =
    graphPanel.new(
      title='Aperture Rate Limiter',
      datasource=dsName,
      labelY1='Decisions',
      formatY1='reqps',
    )
    .addTarget(
      prometheus.target(
        expr=|||
          sum by(decision_type) (
            rate(
              rate_limiter_counter_total{
                %(filters)s
              }[$__rate_interval]
            )
          )
        ||| % { filters: stringFilters },
        intervalFactor=1,
      )
    ),

  local dashboardDef =
    dashboard.new(
      title='Aperture Rate Limiter',
      editable=true,
      schemaVersion=18,
      refresh=params.dashboard.refresh_interval,
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
