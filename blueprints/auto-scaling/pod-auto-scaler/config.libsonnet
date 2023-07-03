local autoScalingDefaults = import '../common/config-defaults.libsonnet';

/**
* @schema (scaling_driver.scale_out.enabled: bool) Enables the driver to do scale out of the resource.
* @schema (scaling_driver.scale_out.threshold: float64) Threshold for the driver.
* @schema (scaling_driver.scale_in.enabled: bool) Enables the Driver to do scale in of the resource.
* @schema (scaling_driver.scale_in.threshold: float64) Threshold for the driver.
*/
local scaling_driver_defaults = {
  scale_out: {
    enabled: '__REQUIRED_FIELD__',
    threshold: '__REQUIRED_FIELD__',
  },
  scale_in: {
    enabled: '__REQUIRED_FIELD__',
    threshold: '__REQUIRED_FIELD__',
  },
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
}
