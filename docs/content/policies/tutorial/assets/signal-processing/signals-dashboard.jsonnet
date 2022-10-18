local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local signalsDashboard = aperture.blueprints.dashboards.Signals({
  policyName: 'signal-processing',
  datasource+: {
    name: 'controller-prometheus',
  },
});

signalsDashboard.dashboard
