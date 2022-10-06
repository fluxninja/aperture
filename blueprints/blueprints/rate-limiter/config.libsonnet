{
  policy: {
    /**
    * @section Rate Limiter Policy
    *
    * #### Rate Limiter Overrides
    *
    * To tweak rate limiter behaviour based on specific label values, a list of `RateLimierOverride` objects
    * can be added to the policy:
    *
    * ```jsonnet
    * local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';
    *
    * local Override = aperture.spec.v1.RateLimiterOverride;
    *
    * {
    *   policy+: {
    *     overrides: [
    *       Override.new() + Override.withLabelValue('gold') + Override.withLimitScaleFactor(1)
    *     ]
    *   }
    * }
    *
    * ```
    * This allows us to prioritize some incoming requests over others.
    *
    * @param (policy.policyName: string required) An unique name for the policy created by this blueprint
    * @param (policy.evaluationInterval: string) How often should the policy be re-evaluated
    * @param (policy.rateLimit: string required) How many requests per `policy.limitResetInterval` to accept
    * @param (policy.rateLimiterSelector: aperture.spec.v1.Selector required) A selector to match requests against
    * @param (policy.limitResetInterval: string) The window for `policy.rateLimit`
    * @param (policy.labelKey: string required) What flow label to use for rate limiting
    * @param (policy.overrides: []aperture.spec.v1.RateLimiterOverride) A list of overrides for the rate limiter
    */
    policyName: error 'policy.policyName is required',
    evaluationInterval: '0.5s',
    rateLimit: error 'policy.rateLimit must be set',
    rateLimiterSelector: error 'rateLimiterSelector must be set',
    limitResetInterval: '1s',
    labelKey: error 'policy.labelKey is required',
    overrides: [],
    lazySync: {
      /**
      * @section Rate Limiter Policy
      * @subsection Rate Limiter Lazy Sync
      *
      *
      * @param (policy.lazySync.enabled: boolean) TODO document what happens when lazy sync is disabled
      * @param (policy.lazySync.numSync: integer) Number of times to lazy sync within the `policy.limitResetInterval`.
      */
      enabled: true,
      numSync: 5,
    },
  },
}
