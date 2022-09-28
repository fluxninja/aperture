{
  /**
  * @section Latency Gradient Policy
  *
  * @param (policy.policyName: string required) A name of the policy, used within PromQL queries for fluxmeter metrics.
  * @param (policy.evaluationInterval: string) How often should policy be re-evaluated.
  */
  policy: {
    policyName: error 'policyName is not set',
    evaluationInterval: '0.1s',
    /**
    * @section Latency Gradient Policy
    * @subsection Flux Meter Selector
    *
    * @param (policy.fluxMeterSelector: aperture.spec.v1.Selector required) Flux Meter selector.
    */
    fluxMeterSelector: error 'fluxMeterSelector is not set',
    /**
    * @section Latency Gradient Policy
    * @subsection Flux Meters
    *
    * @param (policy.fluxMeters: map[string]aperture.spec.v1.FluxMeter) Mappings of fluxMeterName to fluxMeter.
    * @param (policy.fluxMeters[policyName].attributeKey: string) Key of the attribute in access log or span.
    * @param (policy.fluxMeters[policyName].histogramBuckets: aperture.spec.v1.FluxMeterStaticBuckets) Flux Meter static histogram buckets.
    */
    fluxMeters: {
      'service1-latency-gradient': {
        attributeKey: 'workload_duration_ms',
        staticBuckets: {
          buckets: [5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0],
        },
      },
    },
    /**
    * @section Latency Gradient Policy
    * @subsection Concurrency Limiter Selector
    *
    * @param (policy.concurrencyLimiterSelector: aperture.spec.v1.Selector required) Concurrency Limiter selector.
    */
    concurrencyLimiterSelector: error 'concurrencyLimiterSelector is not set',
    /**
    * @section Latency Gradient Policy
    * @subsection Classification Rules
    *
    * @param (policy.classifiers: string) List of classification rules.
    */
    classifiers: [],
    constants: {
      emaLimitMultiplier: '2.0',
      tolerance: '1.1',
      concurrencyLimitMultiplier: '2.0',
      minConcurrency: '10.0',
      linearConcurrencyIncrement: '5.0',
      concurrencyIncrementOverload: '10.0',
      sqrtScale: '0.5',
    },
    /**
    * @section Latency Gradient Policy
    * @subsection Exponential Moving Average configuration
    *
    * @param (policy.ema.window: string) How far back to look when calculating moving average
    * @param (policy.ema.warmUpWindow: string) How much time to give circuit before we start calculating averages
    */
    ema: {
      window: '1500s',
      warmUpWindow: '10s',
      correctionFactor: '0.95',
    },
    gradient: {
      slope: '-1',
      minGradient: '0.1',
      maxGradient: '1.0',
    },
    /**
    * @section Latency Gradient Policy
    * @subsection Concurrency Limiter
    *
    * @param (policy.concurrencyLimiter.defaultWorkloadParameters.priority: int) Workload parameters to use in case none of the configured workloads match.
    * @param (policy.concurrencyLimiter.workloads: []aperture.spec.v1.SchedulerWorkload) A list of additional workloads for the scheduler
    */
    concurrencyLimiter: {
      autoTokens: true,
      timeoutFactor: '0.5',
      defaultWorkloadParameters: {
        priority: 20,
      },
      workloads: [],
    },
  },
  dashboard: {
    /**
    * @section FluxMeter Dashboard
    *
    * @param (dashboard.policyName: string required) A name of the policy used as a promQL query filter for flux meter metrics
    */
    policyName: error 'policyName is not set',
  },
}
