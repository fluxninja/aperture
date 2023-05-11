local autoScalingDynamicConfig = import '../../auto-scaling/base/dynamic-config.libsonnet';
local serviceProtectionDynamicConfig = import '../../service-protection/base/dynamic-config.libsonnet';

/**
* @param (load_scheduler: aperture.spec.v1.LoadSchedulerDynamicConfig required) Default configuration for load scheduler that can be updated at the runtime without shutting down the policy.
* @param (auto_scaling: bool required) Dry run mode ensures that no scaling is invoked by this auto scaler.
*/

serviceProtectionDynamicConfig + autoScalingDynamicConfig
