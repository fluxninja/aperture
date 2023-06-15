{
  /**
  * @param (pass_through_label_values: []string) Specify certain label values to be always accepted by the _Sampler_ regardless of accept percentage. This configuration can be updated at the runtime without shutting down the policy.
  * @param (rollout: bool) Start feature rollout. This setting can be updated at runtime without shutting down the policy. The feature rollout gets paused if this flag is set to false in the middle of a feature rollout.
  * @param (reset: bool) Reset feature rollout to the first step. This setting can be updated at the runtime without shutting down the policy.
  */
  pass_through_label_values: ['__REQUIRED_FIELD__'],
  rollout: false,
  reset: false,
}
