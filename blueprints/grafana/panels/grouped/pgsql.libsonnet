local active_connections = import '../pgsql/active_connections.libsonnet';
local backends = import '../pgsql/backends.libsonnet';
local bgwriter_signals = import '../pgsql/bgwriter_signals.libsonnet';
local blocks = import '../pgsql/blocks.libsonnet';
local cache_hit_idx = import '../pgsql/cache_hit_idx.libsonnet';
local cache_hit_table = import '../pgsql/cache_hit_table.libsonnet';
local db_count = import '../pgsql/db_count.libsonnet';
local db_size = import '../pgsql/db_size.libsonnet';
local disk_blocks = import '../pgsql/disk_blocks.libsonnet';
local max_connections = import '../pgsql/max_connections.libsonnet';
local operations = import '../pgsql/operations.libsonnet';
local table_count = import '../pgsql/table_count.libsonnet';
local table_size = import '../pgsql/table_size.libsonnet';
local top_tables_disk = import '../pgsql/top_tables_disk.libsonnet';
local top_tables_rows = import '../pgsql/top_tables_rows.libsonnet';
local total_commits = import '../pgsql/total_commits.libsonnet';


local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    db_count(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    db_size(cfg).panel
    + g.panel.stat.gridPos.withX(5)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    table_count(cfg).panel
    + g.panel.stat.gridPos.withX(10)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    table_size(cfg).panel
    + g.panel.stat.gridPos.withX(15)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    cache_hit_table(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    cache_hit_idx(cfg).panel
    + g.panel.stat.gridPos.withX(5)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    active_connections(cfg).panel
    + g.panel.stat.gridPos.withX(10)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    max_connections(cfg).panel
    + g.panel.stat.gridPos.withX(15)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    total_commits(cfg, 'Commits per minute').panel
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(10),
    operations(cfg, 'Operations').panel
    + g.panel.timeSeries.gridPos.withX(10)
    + g.panel.timeSeries.gridPos.withY(10)
    + g.panel.timeSeries.standardOptions.withMax(true)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    blocks(cfg, 'Heap Read vs Index Read').panel
    + g.panel.timeSeries.gridPos.withX(10)
    + g.panel.timeSeries.gridPos.withY(15),
    top_tables_rows(cfg).panel
    + g.panel.barGauge.standardOptions.withDisplayName('postgresql_table_name')
    + g.panel.barGauge.gridPos.withX(0)
    + g.panel.barGauge.gridPos.withY(20)
    + g.panel.barGauge.gridPos.withH(10)
    + g.panel.barGauge.gridPos.withW(5),
    top_tables_disk(cfg).panel
    + g.panel.barGauge.standardOptions.withDisplayName('postgresql_table_name')
    + g.panel.barGauge.gridPos.withX(10)
    + g.panel.barGauge.gridPos.withY(20)
    + g.panel.barGauge.gridPos.withH(10)
    + g.panel.barGauge.gridPos.withW(5),
  ],
}
