local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local signalsDashboard = aperture.dashboards.SignalsDashboard.dashboard({
  policy_name: 'signal-processing',
  datasource+: {
    name: 'controller-prometheus',
  },
});

signalsDashboard.dashboard
