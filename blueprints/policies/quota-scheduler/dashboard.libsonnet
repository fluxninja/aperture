local quotaSchedulerDashboard = import '../../../blueprints/dashboards/flow-control/quota-scheduler/dashboard.libsonnet';
local utils = import '../policy-utils.libsonnet';
local config = import './config.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local prometheus = grafana.prometheus;
local graphPanel = grafana.graphPanel;

function(cfg) {
  local params = config + cfg,

  local baseDashboard = quotaSchedulerDashboard(params),

  local policyName = params.policy.policy_name,
  local ds = params.dashboard.datasource,
  local dsName = ds.name,
  local filters = utils.dictToPrometheusFilter(params.dashboard.extra_filters { policy_name: policyName }),

  local quotaSchedulerPanel =
    graphPanel.new(
      title='Quota Checks',
      datasource=dsName,
      labelY1='Decisions',
      formatY1='reqps',
    )
    .addTarget(
      prometheus.target(
        expr='sum by(decision_type) (rate(workload_requests_total{%(filters)s}[$__rate_interval]))' % { filters: filters },
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
