local average_load_multiplier = import '../average_load_multiplier.libsonnet';
local token_bucket_capacity = import '../token_bucket_capacity.libsonnet';
local token_bucket_fillrate = import '../token_bucket_fillrate.libsonnet';
local token_bucket_available_tokens = import '../token_bucket_available_tokens.libsonnet';

function(cfg) {
  panels: [
    average_load_multiplier(cfg).panel,
    token_bucket_capacity(cfg).panel,
    token_bucket_fillrate(cfg).panel,
    token_bucket_available_tokens(cfg).panel,
  ],
}
