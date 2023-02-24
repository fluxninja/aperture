{
  /**
  * @section Common
  *
  * @param (common.policy_name: string required) Name of the policy.
  */
  common: {
    policy_name: '__REQUIRED_FIELD__',
  },
  /**
  * @section Policy
  *
  * @param (policy.classifiers: []aperture.spec.v1.Classifier) List of classification rules.
  */
  policy: {
    classifiers: [],
    /**
    * @section Policy
    * @subsection Rate Limiter
    *
    * @param (policy.rate_limiter.rate_limit: float64 required) Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept
    * @param (policy.rate_limiter.flow_selector: aperture.spec.v1.FlowSelector) A flow selector to match requests against
    * @param (policy.rate_limiter.flow_selector.service_selector.service: string required) Service Name.
    * @param (policy.rate_limiter.flow_selector.flow_matcher.control_point: string required) Control Point Name.
    * @param (policy.rate_limiter.parameters: aperture.spec.v1.RateLimiterParameters) Parameters.
    * @param (policy.rate_limiter.parameters.limit_reset_interval: string required) Time after which the limit for a given label value will be reset.
    * @param (policy.rate_limiter.parameters.label_key: string required) Flow label to use for rate limiting.
    * @param (policy.rate_limiter.default_config: aperture.spec.v1.RateLimiterDynamicConfig) Default configuration for rate limiter that can be updated at the runtime without shutting down the policy.
    * @param (policy.rate_limiter.default_config.overrides: []aperture.spec.v1.RateLimiterOverride) Allows to specify different limits for particular label values.
    */
    rate_limiter: {
      rate_limit: '__REQUIRED_FIELD__',
      flow_selector: {
        service_selector: {
          service: '__REQUIRED_FIELD__',
        },
        flow_matcher: {
          control_point: '__REQUIRED_FIELD__',
        },
      },
      parameters: {
        limit_reset_interval: '__REQUIRED_FIELD__',
        label_key: '__REQUIRED_FIELD__',
      },
      default_config: {
        overrides: [],
      },
    },
  },
  /**
  * @section Dashboard
  *
  * @param (dashboard.refresh_interval: string) Refresh interval for dashboard panels.
  */
  dashboard: {
    refresh_interval: '10s',
    /**
    * @section Dashboard
    * @subsection Datasource
    *
    * @param (dashboard.datasource.name: string) Datasource name.
    * @param (dashboard.datasource.filter_regex: string) Datasource filter regex.
    */
    datasource: {
      name: '$datasource',
      filter_regex: '',
    },
  },
}
