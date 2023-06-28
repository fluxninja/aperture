local featureRolloutConfig = import '../base/config-defaults.libsonnet';

featureRolloutConfig {
  policy+: {
    drivers+: {
      /**
      * @param (policy.drivers.average_latency_drivers: []average_latency_driver) List of drivers that compare average latency against forward, backward and reset thresholds.
      * @schema (average_latency_driver.selectors: []aperture.spec.v1.Selector) Identify the service and flows whose latency needs to be measured.
      * @schema (average_latency_driver.criteria: criteria) The criteria for average latency comparison.
      */
      average_latency_drivers: [
        {
          selectors: featureRolloutConfig.selectors_defaults,
          criteria: featureRolloutConfig.criteria_defaults,
        },
      ],
    },
  },
}
