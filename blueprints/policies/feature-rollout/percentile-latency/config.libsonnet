local featureRolloutConfig = import '../base/config.libsonnet';


local rolloutPolicyDefaults = featureRolloutConfig.rollout_policy_base {

  drivers: {
    percentile_latency_drivers: [
      featureRolloutConfig.percentile_latency_driver,
    ],
  },

};


{
  /**
  * @param (common.policy_name: string required) Name of the policy.
  */
  common: {
    policy_name: '__REQUIRED_FIELD__',
  },
  /**
  * @param (policy: policies/feature-rollout/base:schema:rollout_policy required) Parameters for the Feature Rollout policy.
  */
  policy: rolloutPolicyDefaults,
  /**
  * @param (dashboard: policies/feature-rollout/base:param:dashboard) Configuration for the Grafana dashboard accompanying this policy.
  */
  dashboard: featureRolloutConfig.dashboard,
}
