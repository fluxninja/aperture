local creator = import 'creator.libsonnet';
local signals = import 'signals_dashboard.libsonnet';

function(policyJSON, cfg) {
  local policyName = cfg.policy.policy_name,
  local mainDashboard = creator(policyJSON, cfg).dashboard,
  local signalsDashboard = signals({
    policy+: {
      policy_name: policyName,
    },
    dashboard+: {
      title: 'Aperture Signals - %s' % policyName,
    },
  }).dashboard,
  local additionalDashboard = creator(policyJSON, cfg).additionalDashboard,

  mainDashboard: mainDashboard,
  signalsDashboard: signalsDashboard,
  additionalDashboard: additionalDashboard,
}
