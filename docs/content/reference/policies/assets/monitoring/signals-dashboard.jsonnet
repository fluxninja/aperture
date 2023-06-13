local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local signalsDashboard = aperture.dashboards.Signals.dashboard({
  policy+: {
    policy_name: 'signal-processing',
  },
  dashboard+: {
    datasource+: {
      name: 'controller-prometheus',
    },
  },
});

signalsDashboard.dashboard
