local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  local dashboardDef =
    g.dashboard.new(cfg.dashboard.title)
    + g.dashboard.time.withFrom('now-15m')
    + g.dashboard.withTimezone('browser')
    + g.dashboard.withRefresh(cfg.dashboard.refresh_interval),

  baseDashboard: dashboardDef,
}
