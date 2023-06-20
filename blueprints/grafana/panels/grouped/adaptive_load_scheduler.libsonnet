local average_load_multiplier = import '../average_load_multiplier.libsonnet';
local token_bucket_available_tokens = import '../token_bucket_available_tokens.libsonnet';
local token_bucket_capacity = import '../token_bucket_capacity.libsonnet';
local token_bucket_fillrate = import '../token_bucket_fillrate.libsonnet';

local g = import 'github.com/grafana/grafonnet/gen/grafonnet-v9.4.0/main.libsonnet';

function(cfg) {
  panels: [
    average_load_multiplier(cfg).panel
    + g.panel.stat.gridPos.withY(50),
    token_bucket_capacity(cfg).panel
    + g.panel.stat.gridPos.withX(6)
    + g.panel.stat.gridPos.withY(50),
    token_bucket_fillrate(cfg).panel
    + g.panel.stat.gridPos.withX(12)
    + g.panel.stat.gridPos.withY(50),
    token_bucket_available_tokens(cfg).panel
    + g.panel.stat.gridPos.withX(18)
    + g.panel.stat.gridPos.withY(50),
  ],
}
