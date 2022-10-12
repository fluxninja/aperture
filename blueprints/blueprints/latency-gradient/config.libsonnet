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
    * @subsection Flux Meter
    *
    * @param (policy.fluxMeter: aperture.v1.FluxMeter required) Flux Meter selector.
    */
    fluxMeter: error 'fluxMeter is not set',
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
    * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
    */
    classifiers: [],
    /**
    * @section Latency Gradient Policy
    * @subsection Additional Circuit Components
    *
    * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
    */
    components: [],
    constants: {
      emaLimitMultiplier: '2.0',
      tolerance: '1.1',
      concurrencyLimitMultiplier: '2.0',
      minConcurrency: '10.0',
      linearConcurrencyIncrement: '5.0',
      concurrencyIncrementOverload: '10.0',
      sqrtScale: '0.1',
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
  /**
  * @section FluxMeter Dashboard
  *
  * @param (dashboard.policyName: string required) Dashboard Configuration.
  */
  dashboard: {
    policyName: error 'policyName is not set',
  },
  /**
  * @section Signals Dashboard
  *
  * @param (signalsDashboard.policyName: string required) Dashboard Configuration.
  */
  signalsDashboard: {
    policyName: error 'policyName is not set',
  },
}
