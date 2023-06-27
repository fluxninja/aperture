{
  spec: import 'spec.libsonnet',
  policies: import 'policies/policies.libsonnet',
  signalsDashboardCreator: import 'grafana/signals_dashboard.libsonnet',
  dashboardPanelLibrary: import 'grafana/panel_library.libsonnet',
  dashboardCreator: import 'grafana/creator.libsonnet',
}
