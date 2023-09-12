local backends = import '../pgsql/backends.libsonnet';
local bgwriter_allocated = import '../pgsql/bgwriter_allocated.libsonnet';
local bgwriter_writes = import '../pgsql/bgwriter_writes.libsonnet';
local blocks = import '../pgsql/blocks.libsonnet';
local db_count = import '../pgsql/db_count.libsonnet';
local index_scans_total = import '../pgsql/index_scans_total.libsonnet';
local max_connections = import '../pgsql/max_connections.libsonnet';
local table_count = import '../pgsql/table_count.libsonnet';
local total_commits = import '../pgsql/total_commits.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    db_count(cfg).panel
    + g.panel.stat.gridPos.withY(5),
    table_count(cfg).panel
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(5),
    max_connections(cfg).panel
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(5),
    bgwriter_writes(cfg).panel
    + g.panel.stat.gridPos.withY(10),
    bgwriter_allocated(cfg).panel
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(10),
    total_commits(cfg).panel
    + g.panel.timeSeries.gridPos.withY(20),
    blocks(cfg).panel
    + g.panel.stat.gridPos.withY(40),
    backends(cfg).panel
    + g.panel.timeSeries.gridPos.withY(60),
    index_scans_total(cfg).panel
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.gridPos.withX(24)
    + g.panel.timeSeries.gridPos.withY(80),
  ],
}
