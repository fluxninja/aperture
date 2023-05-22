local dynamicConfig = import '../base/dynamic-config.libsonnet';

/**
  * @param (pass_through_label_values: []string) Specify certain label values to be always accepted by the _Regulator_ regardless of accept percentage. This configuration can be updated at the runtime without shutting down the policy.
*/
dynamicConfig
