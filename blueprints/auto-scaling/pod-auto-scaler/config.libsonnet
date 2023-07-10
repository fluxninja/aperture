local autoScalingDefaults = import '../common/config-defaults.libsonnet';

/**
* @schema (scaling_criteria.enabled: bool) Enables the driver to do scale in or out of the resource.
* @schema (scaling_criteria.threshold: float64) Threshold for the driver.
*/
local scaling_criteria_defaults = {
  enabled: '__REQUIRED_FIELD__',
  threshold: '__REQUIRED_FIELD__',
};

/**
* @schema (scaling_driver.scale_out: scaling_criteria) The scale out criteria for the driver.
* @schema (scaling_driver.scale_in: scaling_criteria) The scale in criteria for the driver.
*/
local scaling_driver_defaults = {
  scale_out: {},
  scale_in: {},
};

autoScalingDefaults {
  policy+: {
    scaling_backend: {
      kubernetes_replicas: '__REQUIRED_FIELD__',
    },
    /**
    * @param (policy.pod_cpu: scaling_driver) Driver to do scaling of the resource based on the CPU usage.
    * @param (policy.pod_memory: scaling_driver) Driver to do scaling of the resource based on the Memory usage.
    */
    pod_cpu: {},
    pod_memory: {},
  },

  scaling_driver: scaling_driver_defaults,
  scaling_criteria: scaling_criteria_defaults,
}
