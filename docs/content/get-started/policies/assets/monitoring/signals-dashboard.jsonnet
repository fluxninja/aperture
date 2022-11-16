local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local signalsDashboard = aperture.blueprints.SignalsDashboard.dashboard({
  policyName: 'signal-processing',
  datasource+: {
    name: 'controller-prometheus',
  },
});

signalsDashboard.dashboard
