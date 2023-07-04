{
  /**
  * @param (pass_through_label_values: []string) Specify certain label values to be always accepted by the _Sampler_ regardless of accept percentage. This configuration can be updated at the runtime without shutting down the policy.
  * @param (start: bool) Start load ramp. This setting can be updated at runtime without shutting down the policy. The load ramp gets paused if this flag is set to false in the middle of a load ramp.
  * @param (reset: bool) Reset load ramp to the first step. This setting can be updated at the runtime without shutting down the policy.
  */
  pass_through_label_values: ['__REQUIRED_FIELD__'],
  start: false,
  reset: false,
}
