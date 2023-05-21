local aperture = import '../../grafana/aperture.libsonnet';
local lib = import '../../grafana/grafana.libsonnet';
local baseDashboardFn = import '../service-protection/base/dashboard.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local graphPanel = grafana.graphPanel;

function(cfg) {
  local params = config + cfg,

  local baseDashboard = baseDashboardFn(params),

  local policyName = params.policy.policy_name,
  local ds = params.dashboard.datasource,
  local dsName = ds.name,

  local quotaSchedulerPanel =
    graphPanel.new(
      title='Quota Checks',
      datasource=dsName,
      labelY1='Decisions',
      formatY1='reqps',
    )
    .addTarget(
      prometheus.target(
        expr='sum by(decision_type) (rate(quota_check_counter_total{policy_name="%(policy_name)s"}[$__rate_interval]))' % { policy_name: policyName },
        intervalFactor=1,
      )
    ),

  local quotaSchedulerPanelWithID =
    quotaSchedulerPanel {
      id: '0',
      gridPos: { x: 0, y: 0, w: 24, h: 10 },
    },

  // extend the base dashboard to add the panels
  local extendedDashboard =
    baseDashboard.dashboard {
      title: 'Aperture Quota Scheduler',
      panels+: [
        quotaSchedulerPanelWithID,
      ],
    },

  dashboard: extendedDashboard,
}
