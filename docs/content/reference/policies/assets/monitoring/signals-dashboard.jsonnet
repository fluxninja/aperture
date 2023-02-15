local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local signalsDashboard = aperture.dashboards.SignalsDashboard.dashboard({
  policy_name: 'signal-processing',
  datasource+: {
    name: 'controller-prometheus',
  },
});

signalsDashboard.dashboard
