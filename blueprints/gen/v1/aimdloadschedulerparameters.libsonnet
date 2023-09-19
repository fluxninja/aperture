{
  new():: {
  },
  withAlerter(alerter):: {
    alerter: alerter,
  },
  withAlerterMixin(alerter):: {
    alerter+: alerter,
  },
  withGradient(gradient):: {
    gradient: gradient,
  },
  withGradientMixin(gradient):: {
    gradient+: gradient,
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
}
