{
  /**
  * @section Common
  *
  * @param (common.policyName: string required) Name of the policy.
  */
  common: {
    policyName: error 'policyName is not set',
  },
  /**
  * @section Policy
  *
  * @param (policy.evaluationInterval: string) How often should policy be re-evaluated.
  * @param (policy.fluxMeter: aperture.v1.FluxMeter required) Flux Meter selector.
  * @param (policy.concurrencyLimiterSelector: aperture.spec.v1.Selector required) Concurrency Limiter selector.
  * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  */
  policy: {
    evaluationInterval: '0.1s',
    fluxMeter: error 'fluxMeter is not set',
    concurrencyLimiterSelector: error 'concurrencyLimiterSelector is not set',
    classifiers: [],
    components: [],
    /**
    * @section Policy
    * @subsection Concurrency Limiter
    *
    * @param (policy.concurrencyLimiter.defaultWorkloadParameters.priority: int) Workload parameters to use in case none of the configured workloads match.
    * @param (policy.concurrencyLimiter.workloads: []aperture.spec.v1.SchedulerWorkload) A list of additional workloads for the scheduler.
    */
    concurrencyLimiter: {
      autoTokens: true,
      timeoutFactor: '0.5',
      defaultWorkloadParameters: {
        priority: 20,
      },
      workloads: [],
    },
    /**
    * @section Policy
    * @subsection Constants
    *
    * @param (policy.constants.tolerance: float64) Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.
    * @param (policy.constants.emaLimitMultiplier: float64) Current latency value is multiplied with this factor to calculate max EMA envelope.
    * @param (policy.constants.concurrencyLimitMultiplier: float64) Current accepted concurrency is multiplied with this number to calculate upper concurrency limit that can be allowed at the scheduler. This prevents from system to be protected from sudden spikes while the controller catches up on the observability data.
    * @param (policy.constants.minConcurrency: float64) Minimum concurrency allowed in the system during no overload state.
    * @param (policy.constants.linearConcurrencyIncrement: float64) Linear increment to concurrency in each execution tick when the system is not in overloaded state.
    * @param (policy.constants.sqrtScale: float64) Scale factor to multiply square root of current accepted concurrrency. This, along with linearConcurrencyIncrement helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during non-overload state (integral effect).
    * @param (policy.constants.concurrencyIncrementOverload: float64) Concurrent increment to apply during overload state that is still applied despite lowering the overall concurrency in the gradient controller. This is the minimum concurrency that will still be allowed during overload state.
    */
    constants: {
      tolerance: '1.1',
      emaLimitMultiplier: '2.0',
      concurrencyLimitMultiplier: '2.0',
      minConcurrency: '10.0',
      linearConcurrencyIncrement: '5.0',
      sqrtScale: '0.1',
      concurrencyIncrementOverload: '10.0',
    },
    /**
    * @section Policy
    * @subsection EMA
    *
    * @param (policy.ema.window: string) How far back to look when calculating moving average.
    * @param (policy.ema.warmUpWindow: string) How much time to give circuit to learn the average value before we start emitting EMA values.
    * @param (policy.ema.correctionFactor: string) Factor that is applied to the EMA value when it's above the maximum envelope.
    */
    ema: {
      window: '1500s',
      warmUpWindow: '10s',
      correctionFactor: '0.95',
    },
    /**
    * @section Policy
    * @subsection Gradient Controller
    *
    * @param (policy.gradient.slope: float64) Gradient that adjusts the response of the controller based on current latency and setpoint latency.
    * @param (policy.gradient.minGradient: float64) Minimum gradient cap.
    * @param (policy.gradient.maxGradient: float64) Maximum gradient cap.
    */
    gradient: {
      slope: '-1',
      minGradient: '0.1',
      maxGradient: '1.0',
    },
  },
  /**
  * @section Dashboard
  *
  * @param (dashboard.refreshInterval: string) Refresh interval for dashboard panels.
  */
  dashboard: {
    refreshInterval: '10s',
    /**
    * @section Dashboard
    * @subsection Datasource
    *
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filterRegex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filterRegex: '',
    },
  },
}
