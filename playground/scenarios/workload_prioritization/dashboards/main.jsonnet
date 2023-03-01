local apertureDashboard = import '../../../resources/grafana-dashboard/main.libsonnet';

apertureDashboard(std.parseJson(std.extVar('APERTURE_DASHBOARD'))).dashboards
