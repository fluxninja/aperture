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
  * @param (policy.ema.window: string) How far back to look when calculating moving average
  * @param (policy.ema.warmUpWindow: string) How much time to give circuit before we start calculating averages
  * @param (policy.concurrencyLimiter.defaultWorkloadParameters.priority: int) Workload parameters to use in case none of the configured workloads match.
  * @param (policy.concurrencyLimiter.workloads: []aperture.spec.v1.SchedulerWorkload) A list of additional workloads for the scheduler
  */
  policy: {
    evaluationInterval: '0.1s',
    fluxMeter: error 'fluxMeter is not set',
    concurrencyLimiterSelector: error 'concurrencyLimiterSelector is not set',
    classifiers: [],
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
  * @section Dashboard
  *
  * @param (dashboard.refreshInterval: string) Refresh interval for dashboard panels.
  * @param (dashboard.datasourceName: string) Datasource name.
  * @param (dashboard.datasourceFilterRegex: string) Datasource filter regex.
  */
  dashboard: {
    refreshInterval: '10s',
    datasourceName: '$datasource',
    datasourceFilterRegex: '',
  },
}
