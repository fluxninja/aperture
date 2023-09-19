// Stat Panels
local active_connections = import '../pgsql/active_connections.libsonnet';
local cache_hit_idx = import '../pgsql/cache_hit_idx.libsonnet';
local cache_hit_table = import '../pgsql/cache_hit_table.libsonnet';
local db_count = import '../pgsql/db_count.libsonnet';
local db_size = import '../pgsql/db_size.libsonnet';
local max_connections = import '../pgsql/max_connections.libsonnet';
local table_count = import '../pgsql/table_count.libsonnet';
local table_size = import '../pgsql/table_size.libsonnet';

// Time Series and Bar Gauge Panels
local checkpoint_comparison = import '../pgsql/checkpoint_comparison.libsonnet';
local commits_vs_rollbacks = import '../pgsql/commits_vs_rollbacks.libsonnet';
local heap_vs_index = import '../pgsql/heap_vs_index.libsonnet';
local operations = import '../pgsql/operations.libsonnet';
local top_tables_disk = import '../pgsql/top_tables_disk.libsonnet';
local top_tables_rows = import '../pgsql/top_tables_rows.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    db_count(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    db_size(cfg).panel
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    table_count(cfg).panel
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    table_size(cfg).panel
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(0)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    cache_hit_table(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    cache_hit_idx(cfg).panel
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    active_connections(cfg).panel
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    max_connections(cfg).panel
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(3)
    + g.panel.stat.gridPos.withW(6),
    commits_vs_rollbacks(cfg, 'Commits vs Rollbacks').panel
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(20)
    + g.panel.timeSeries.gridPos.withH(6)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    operations(cfg, 'Operations').panel
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(20)
    + g.panel.timeSeries.gridPos.withH(6)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    heap_vs_index(cfg, 'Heap Read vs Index Read').panel
    + g.panel.timeSeries.gridPos.withX(0)
    + g.panel.timeSeries.gridPos.withY(30)
    + g.panel.timeSeries.gridPos.withH(6)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    checkpoint_comparison(cfg, 'Checkpoint Scheduled vs Requested').panel
    + g.panel.timeSeries.gridPos.withX(12)
    + g.panel.timeSeries.gridPos.withY(30)
    + g.panel.timeSeries.gridPos.withH(6)
    + g.panel.timeSeries.gridPos.withW(12)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
    top_tables_rows(cfg).panel
    + g.panel.barGauge.gridPos.withX(0)
    + g.panel.barGauge.gridPos.withY(40)
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.options.withValueMode('text'),
    top_tables_disk(cfg).panel
    + g.panel.barGauge.gridPos.withX(12)
    + g.panel.barGauge.gridPos.withY(40)
    + g.panel.barGauge.gridPos.withH(6)
    + g.panel.barGauge.gridPos.withW(12)
    + g.panel.barGauge.options.withValueMode('text'),
  ],
}
