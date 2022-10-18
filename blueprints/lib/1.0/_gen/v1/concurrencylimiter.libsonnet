{
  new():: {
  },
  withLoadActuator(load_actuator):: {
    load_actuator: load_actuator,
  },
  withLoadActuatorMixin(load_actuator):: {
    load_actuator+: load_actuator,
  },
  withScheduler(scheduler):: {
    scheduler: scheduler,
  },
  withSchedulerMixin(scheduler):: {
    scheduler+: scheduler,
  },
  withSelector(selector):: {
    selector: selector,
  },
  withSelectorMixin(selector):: {
    selector+: selector,
  },
}
