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
  * @param (policy.fluxMeter: aperture.v1.FluxMeter required) Flux Meter selector.
  * @param (policy.concurrencyLimiterSelector: aperture.spec.v1.Selector required) Concurrency Limiter selector.
  * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
  * @param (policy.components: []aperture.spec.v1.Component) List of additional circuit components.
  */
  policy: {
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
    * @param (policy.constants.latencyToleranceMultiplier: float64) Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.
    * @param (policy.constants.latencyEMALimitMultiplier: float64) Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.
    * @param (policy.constants.concurrencyLimitMultiplier: float64) Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.
    * @param (policy.constants.concurrencyLinearIncrement: float64) Linear increment to concurrency in each execution tick when the system is not in overloaded state.
    * @param (policy.constants.concurrencySQRTIncrementMultiplier: float64) Scale factor to multiply square root of current accepted concurrrency. This, along with concurrencyLinearIncrement helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect).
    */
    constants: {
      latencyToleranceMultiplier: '1.1',
      latencyEMALimitMultiplier: '2.0',
      concurrencyLimitMultiplier: '2.0',
      concurrencyLinearIncrement: '5.0',
      concurrencySQRTIncrementMultiplier: '0.1',
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
