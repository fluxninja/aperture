local dynamicConfig = import '../base/dynamic-config.libsonnet';

/**
* @param (dry_run: bool) Dynamic configuration for setting dry run mode at runtime without restarting this policy. Dry run mode ensures that no scaling is invoked by this auto scaler. This is useful for observing the behavior of auto scaler without disrupting any real deployment.
*/

dynamicConfig
