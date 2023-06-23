local utils = import '../../policy-utils.libsonnet';
local baseDashboardFn = import '../base/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

function(cfg, params={}) {
  local updatedConfig = config + cfg,
  local policyName = updatedConfig.policy.policy_name,
  local variantName = updatedConfig.dashboard.variant_name,
  local filters = utils.dictToPrometheusFilter(updatedConfig.dashboard.extra_filters { flux_meter_name: policyName, flow_status: 'OK' }),

  local baseDashboard = baseDashboardFn(updatedConfig, params),

  local fluxMeterPanel =
    graphPanel.new(
      title=variantName + ' Query',
      datasource=updatedConfig.dashboard.datasource.name,
      labelY1='Latency (ms)',
      formatY1='ms'
    ).addTarget(
      prometheus.target(
        expr=|||
          sum(increase(flux_meter_sum{%(filters)s}[$__rate_interval]))
          / sum(increase(flux_meter_count{%(filters)s}[$__rate_interval]))
        ||| % { filters: filters },
        intervalFactor=1,
      )
    ),

  local fluxMeterPanelWithID =
    fluxMeterPanel {
      id: '0',
      gridPos: { x: 0, y: 0, w: 24, h: 10 },
    },


  // extend the base dashboard to add the panels
  local extendedDashboard =
    baseDashboard.dashboard {
      panels+: [
        fluxMeterPanelWithID,
      ],
    },

  dashboard: extendedDashboard,
}
