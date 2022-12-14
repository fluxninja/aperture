{
  new():: {
  },
  withFlowSelector(flow_selector):: {
    flow_selector: flow_selector,
  },
  withFlowSelectorMixin(flow_selector):: {
    flow_selector+: flow_selector,
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
}
