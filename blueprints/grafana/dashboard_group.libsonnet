local creator = import 'creator.libsonnet';
local infraMeterDashboard = import 'infra_meter_creator.libsonnet';
local signals = import 'signals_dashboard.libsonnet';

function(policyJSON, cfg) {
  local policyName = cfg.policy.policy_name,
  local dashboards = creator(policyJSON, cfg),
  local mainDashboard = dashboards.dashboard,
  local receiverDashboards = dashboards.receiverDashboards,

  local signalsDashboard = signals({
    policy+: {
      policy_name: policyName,
    },
    dashboard+: {
      title: 'Aperture Signals - %s' % policyName,
    },
  }).dashboard,

  mainDashboard: mainDashboard,
  signalsDashboard: signalsDashboard,
  receiverDashboards: receiverDashboards,
}
