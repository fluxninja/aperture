local creator = import 'creator.libsonnet';
local pgsql = import 'pgsql_dashboard.libsonnet';
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
  local pgsqlDashboard = pgsql({
    policy+: {
      policy_name: policyName,
    },
    dashboard+: {
      title: 'pgsql- %s' % policyName,
    },
  }).dashboard,

  mainDashboard: mainDashboard,
  signalsDashboard: signalsDashboard,
  pgsqlDashboard: pgsqlDashboard,
}
