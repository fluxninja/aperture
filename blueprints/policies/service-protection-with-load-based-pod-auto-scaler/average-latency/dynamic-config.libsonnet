local dynamicConfig = import '../../service-protection/average-latency/dynamic-config.libsonnet';
/**
* @param (dry_run: bool) Dynamic configuration for setting dry run mode at runtime without restarting this policy. In dry run mode the scheduler acts as pass through to all flow and does not queue flows. The Auto Scaler does not perform any scaling in dry mode. This mode is useful for observing the behavior of load scheduler and auto scaler without disrupting any real deployment or traffic.
*/
dynamicConfig
