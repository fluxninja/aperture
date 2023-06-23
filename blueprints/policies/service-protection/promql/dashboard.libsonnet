local baseDashboardFn = import '../base/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

function(cfg, params={}) {
  local updatedConfig = config + cfg,
  local variantName = updatedConfig.dashboard.variant_name,
  local query = updatedConfig.policy.promql_query,

  local baseDashboard = baseDashboardFn(updatedConfig, params),

  local queryPanel =
    graphPanel.new(
      title=variantName + ' Query',
      datasource=updatedConfig.dashboard.datasource.name,
      labelY1='Messages',
    ).addTarget(
      prometheus.target(
        expr=query,
        intervalFactor=1,
      )
    ),

  local queryPanelWithID =
    queryPanel {
      id: '0',
      gridPos: { x: 0, y: 0, w: 24, h: 10 },
    },


  // extend the base dashboard to add the panels
  local extendedDashboard =
    baseDashboard.dashboard {
      panels+: [
        queryPanelWithID,
      ],
    },

  dashboard: extendedDashboard,
}
