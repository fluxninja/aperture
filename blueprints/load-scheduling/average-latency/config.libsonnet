local serviceProtectionDefaults = import '../common/config-defaults.libsonnet';

serviceProtectionDefaults {
  policy+: {
    latency_baseliner: {
      /**
      * @param (policy.latency_baseliner.flux_meter: aperture.spec.v1.FluxMeter) Flux Meter defines the scope of latency measurements.
      */
      flux_meter: {
        selectors: serviceProtectionDefaults.selectors_defaults,
      },
      /**
      * @param (policy.latency_baseliner.long_term_query_interval: string) Interval for long-term latency query, i.e., how far back in time the query is run. The value should be a string representing the duration in seconds.
      * @param (policy.latency_baseliner.long_term_query_periodic_interval: string) Periodic interval for long-term latency query, i.e., how often the query is run. The value should be a string representing the duration in seconds.
      * @param (policy.latency_baseliner.latency_tolerance_multiplier: float64) Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if the long-term average of latency is L and if the tolerance is T, then the service is considered to be in an overloaded state if the short-term average of latency is more than L*T.
      */
      long_term_query_interval: '1800s',
      long_term_query_periodic_interval: '30s',
      latency_tolerance_multiplier: 1.25,
    },
  },

  dashboard+: {
    variant_name: 'Average Latency',
  },
}
