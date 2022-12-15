{
  /**
  * @section Common
  *
  * @param (common.policyName: string required) Name of the policy.
  */
  common: {
    policyName: error 'policyName is required',
  },
  /**
  * @section Policy
  *
  * @param (policy.evaluationInterval: string) How often should the policy be re-evaluated
  * @param (policy.rateLimit: float64 required) How many requests per `policy.limitResetInterval` to accept
  * @param (policy.rateLimiterFlowSelector: aperture.spec.v1.FlowSelector required) A flow selector to match requests against
  * @param (policy.limitResetInterval: string) The window for `policy.rateLimit`
  * @param (policy.labelKey: string required) What flow label to use for rate limiting
  */
  policy: {
    evaluationInterval: '300s',
    rateLimit: error 'policy.rateLimit must be set',
    rateLimiterFlowSelector: error 'rateLimiterFlowSelector must be set',
    limitResetInterval: '1s',
    labelKey: error 'policy.labelKey is required',
    classifiers: [],
    /**
    * @section Policy
    * @subsection Overrides
    *
    * @param (policy.overrides: []aperture.spec.v1.RateLimiterOverride) A list of limit overrides for the rate limiter.
    *
    */
    overrides: [],
    /**
    * @section Policy
    * @subsection Lazy Sync
    *
    *
    * @param (policy.lazySync.enabled: boolean) Enable lazy syncing.
    * @param (policy.lazySync.numSync: integer) Number of times to lazy sync within the `policy.limitResetInterval`.
    */
    lazySync: {
      enabled: true,
      numSync: 5,
    },
  },
}
