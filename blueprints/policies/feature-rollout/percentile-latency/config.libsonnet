local featureRolloutConfig = import '../base/config-defaults.libsonnet';

featureRolloutConfig {
  policy+: {
    drivers+: {
      /**
      * @param (policy.drivers.percentile_latency_drivers: []percentile_latency_driver) List of drivers that compare percentile latency against forward, backward and reset thresholds.
      * @schema (percentile_latency_driver.flux_meter: aperture.spec.v1.FluxMeter) FluxMeter specifies the flows whose latency needs to be measured and parameters for the histogram metrics.
      * @schema (percentile_latency_driver.criteria: criteria) The criteria for percentile latency comparison.
      * @schema (percentile_latency_driver.percentile: float64) The percentile to be used for latency measurement.
      */
      percentile_latency_drivers: [
        {
          flux_meter: {
            selector: featureRolloutConfig.selectors_defaults,
            static_buckets: {
              buckets: [5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0],
            },
          },
          criteria: featureRolloutConfig.criteria_defaults,
          percentile: 95,
        },
      ],
    },
  },
}
