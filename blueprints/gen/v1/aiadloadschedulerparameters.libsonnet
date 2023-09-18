{
  new():: {
  },
  withAlerter(alerter):: {
    alerter: alerter,
  },
  withAlerterMixin(alerter):: {
    alerter+: alerter,
  },
  withLoadMultiplierLinearDecrement(load_multiplier_linear_decrement):: {
    load_multiplier_linear_decrement: load_multiplier_linear_decrement,
  },
  withLoadMultiplierLinearDecrementMixin(load_multiplier_linear_decrement):: {
    load_multiplier_linear_decrement+: load_multiplier_linear_decrement,
  },
  withLoadMultiplierLinearIncrement(load_multiplier_linear_increment):: {
    load_multiplier_linear_increment: load_multiplier_linear_increment,
  },
  withLoadMultiplierLinearIncrementMixin(load_multiplier_linear_increment):: {
    load_multiplier_linear_increment+: load_multiplier_linear_increment,
  },
  withLoadScheduler(load_scheduler):: {
    load_scheduler: load_scheduler,
  },
  withLoadSchedulerMixin(load_scheduler):: {
    load_scheduler+: load_scheduler,
  },
  withMaxLoadMultiplier(max_load_multiplier):: {
    max_load_multiplier: max_load_multiplier,
  },
  withMaxLoadMultiplierMixin(max_load_multiplier):: {
    max_load_multiplier+: max_load_multiplier,
  },
  withMinLoadMultiplier(min_load_multiplier):: {
    min_load_multiplier: min_load_multiplier,
  },
  withMinLoadMultiplierMixin(min_load_multiplier):: {
    min_load_multiplier+: min_load_multiplier,
  },
  withOverloadCondition(overload_condition):: {
    overload_condition: overload_condition,
  },
  withOverloadConditionMixin(overload_condition):: {
    overload_condition+: overload_condition,
  },
}
