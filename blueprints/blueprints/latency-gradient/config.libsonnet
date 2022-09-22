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
    * @param (policy.fluxMeterSelector: aperture.v1.Selector required) Flux Meter selector.
    */
    fluxMeterSelector: error 'fluxMeterSelector is not set',
    /**
    * @section Latency Gradient Policy
    * @subsection Concurrency Limiter Selector
    *
    * @param (policy.concurrencyLimiterSelector: aperture.v1.Selector required) Concurrency Limiter selector.
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
    * Concurrency limiter is responsible for.... TODO
    *
    * `concurrencyLimiter.workloads` is a list of `aperture.v1.ServiceWorkloadAndLabelMatcher` objects, which can be generated
    * using aperture libsonnet library (which allows us to do some static checks for validity):
    *
    * ```jsonnet
    * local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';
    *
    * local Workload = aperture.v1.SchedulerWorkload;
    * local LabelMatcher = aperture.v1.LabelMatcher;
    * local WorkloadMatcher = aperture.v1.ServiceWorkloadAndLabelMatcher; // Make a local typedef for quicker access
    *
    * {
    *   concurrencyLimiter+: {
    *     workloads: [
    *       WorkloadMatcher.new(
    *         workload=Workload.new() + Workload.withPriority(50) + Workload.withTimeout('0.005s')
    *         label_matcher=LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'guest' })),
    *       WorkloadWithLabelMatcher.new(
    *         workload=Workload.new() + Workload.withPriority(200) + Workload.withTimeout('0.005s'),
    *         label_matcher=LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' }))
    *     ]
    *   }
    * }
    * ```
    *
    * Or it can be passed as a list of jsonnet objects directly:
    *
    * ```jsonnet
    * {
    *   concurrencyLimiter+: {
    *     workloads: [
    *       {
    *         label_matcher: {
    *           match_labels: {
    *             "http.request.header.user_type": "guest"
    *           }
    *         },
    *         workload: { priority: 50, timeout: "0.005s" }
    *       },
    *       {
    *         label_matcher: {
    *           match_labels: {
    *             "http.request.header.user_type": "subscriber"
    *           }
    *         },
    *         workload: { priority: 200, timeout: "0.005s" }
    *       }
    *     ]
    *   }
    * }
    * ```
    *
    * @param (policy.concurrencyLimiter.defaultWorkload.priority: int) TODO
    * @param (policy.concurrencyLimiter.workloads: []aperture.v1.SchedulerWorkloadAndLabelMatcher) A list of additional workloads for the scheduler
    */
    concurrencyLimiter: {
      autoTokens: true,
      timeoutFactor: '0.5',
      defaultWorkload: {
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
