{
  /**
  * @section Circuit
  *
  * @param (circuit.concurrencyLimiterFlowSelector: aperture.spec.v1.FlowSelector required) Concurrency Limiter flow selector.
  * @param (circuit.components: []aperture.spec.v1.Component) List of additional circuit components.
  */
  circuit: {
    concurrencyLimiterFlowSelector: error 'concurrencyLimiterFlowSelector is not set',
    components: [],
    /**
    * @section Circuit
    * @subsection Concurrency Limiter
    *
    * @param (circuit.concurrencyLimiter.autoTokens: bool) Whether tokens for workloads are computed dynamically or set statically by the user.
    * @param (circuit.concurrencyLimiter.timeoutFactor: float64) The maximum time a request can wait for tokens as a factor of tokens for a flow in a workload.
    * @param (circuit.concurrencyLimiter.defaultWorkloadParameters.priority: int) Workload parameters to use in case none of the configured workloads match.
    * @param (circuit.concurrencyLimiter.workloads: []aperture.spec.v1.SchedulerParametersWorkload) A list of additional workloads for the scheduler.
    * @param (circuit.concurrencyLimiter.alerterName: string) Name of the alert sent on Load Shed Event.
    * @param (circuit.concurrencyLimiter.alerterChannels: []string) A list of alert channels to which the alert will be sent.
    * @param (circuit.concurrencyLimiter.alerterResolveTimeout: string) A timeout after which alert is marked as resolved if alert is not repeated.
    * @param (circuit.concurrencyLimiter.dynamicConfigKey: string) Dynamic configuration key for concurrency limiter.
    */
    concurrencyLimiter: {
      autoTokens: true,
      timeoutFactor: '0.5',
      defaultWorkloadParameters: {
        priority: 20,
      },
      workloads: [],
      alerterName: 'Load Shed Event',
      alerterChannels: [],
      alerterResolveTimeout: '5s',
      dynamicConfigKey: 'concurrency_limiter',
    },
    /**
    * @section Circuit
    * @subsection Constants
    *
    * @param (circuit.constants.concurrencyLimitMultiplier: float64) Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.
    * @param (circuit.constants.concurrencyLinearIncrement: float64) Linear increment to concurrency in each execution tick when the system is not in overloaded state.
    * @param (circuit.constants.concurrencySQRTIncrementMultiplier: float64) Scale factor to multiply square root of current accepted concurrrency. This, along with concurrencyLinearIncrement helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect).
    */
    constants: {
      concurrencyLimitMultiplier: '2.0',
      concurrencyLinearIncrement: '5.0',
      concurrencySQRTIncrementMultiplier: '1',
    },
    /**
    * @section Circuit
    * @subsection Gradient Controller
    *
    * @param (circuit.gradient.slope: float64) Gradient that adjusts the response of the controller based on current latency and setpoint latency.
    * @param (circuit.gradient.minGradient: float64) Minimum gradient cap.
    * @param (circuit.gradient.maxGradient: float64) Maximum gradient cap.
    */
    gradient: {
      slope: '-1',
      minGradient: '0.1',
      maxGradient: '1.0',
    },
  },
}
