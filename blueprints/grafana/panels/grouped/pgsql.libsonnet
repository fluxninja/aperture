local backends = import '../pgsql/backends.libsonnet';
local bgwriter_signals = import '../pgsql/bgwriter_signals.libsonnet';
local blocks = import '../pgsql/blocks.libsonnet';
local db_count = import '../pgsql/db_count.libsonnet';
local db_size = import '../pgsql/db_size.libsonnet';
local max_connections = import '../pgsql/max_connections.libsonnet';
local operations = import '../pgsql/operations.libsonnet';
local table_count = import '../pgsql/table_count.libsonnet';
local table_size = import '../pgsql/table_size.libsonnet';
local total_commits = import '../pgsql/total_commits.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    db_count(cfg).panel
    + g.panel.stat.gridPos.withX(0)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    db_size(cfg).panel
    + g.panel.stat.gridPos.withX(5)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    table_count(cfg).panel
    + g.panel.stat.gridPos.withX(10)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    table_size(cfg).panel
    + g.panel.stat.gridPos.withX(15)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    max_connections(cfg).panel
    + g.panel.stat.gridPos.withX(20)
    + g.panel.stat.gridPos.withY(5)
    + g.panel.stat.gridPos.withH(4)
    + g.panel.stat.gridPos.withW(4),
    bgwriter_signals(cfg, 'BgWriter Writes').panel
    + g.panel.timeSeries.gridPos.withY(10)
    + g.panel.timeSeries.options.legend.withAsTable(true),
    //+ g.panel.timeSeries.fieldConfig.defaults.custom.withDrawStyle("bars"),
    total_commits(cfg).panel
    + g.panel.timeSeries.gridPos.withY(20),
    blocks(cfg).panel
    + g.panel.timeSeries.gridPos.withY(30),
    backends(cfg).panel
    + g.panel.barGauge.gridPos.withY(40),
    operations(cfg, 'Operations').panel
    + g.panel.timeSeries.gridPos.withY(50)
    + g.panel.timeSeries.standardOptions.withMax(true)
    + g.panel.timeSeries.options.legend.withPlacement('right'),
  ],
}
