local baseDashboardFn = import '../base/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

function(cfg) {
  local params = config + cfg,
  local variantName = params.dashboard.variant_name,
  local query = params.policy.promql_query,

  local baseDashboard = baseDashboardFn(params),

  local queryPanel =
    graphPanel.new(
      title=variantName + ' Query',
      datasource=params.dashboard.datasource.name,
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
