local featureRolloutConfig = import '../base/config.libsonnet';


local rollout_policy_defaults = featureRolloutConfig.rollout_policy_base {

  drivers: {
    average_latency_drivers: [
      featureRolloutConfig.average_latency_driver,
    ],
  },

};


{
  /**
  * @param (policy: policies/feature-rollout/base:schema:rollout_policy required) Configuration for the Feature Rollout policy.
  */
  policy: rollout_policy_defaults,
  /**
  * @param (dashboard: policies/feature-rollout/base:param:dashboard) Configuration for the Grafana dashboard accompanying this policy.
  */
  dashboard: featureRolloutConfig.dashboard,
}
