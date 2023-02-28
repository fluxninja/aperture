local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local signalsDashboard = aperture.dashboards.Signals.dashboard({
  policy_name: 'signal-processing',
  datasource+: {
    name: 'controller-prometheus',
  },
});

signalsDashboard.dashboard
