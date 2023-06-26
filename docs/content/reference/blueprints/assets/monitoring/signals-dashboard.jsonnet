local signalsDashboardCreator = import 'github.com/fluxninja/aperture/blueprints/grafana/signals_dashboard.libsonnet';

local signalsDashboard = signalsDashboardCreator({
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
