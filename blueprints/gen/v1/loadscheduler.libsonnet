{
  new():: {
  },
  withActuator(actuator):: {
    actuator: actuator,
  },
  withActuatorMixin(actuator):: {
    actuator+: actuator,
  },
  withFlowSelector(flow_selector):: {
    flow_selector: flow_selector,
  },
  withFlowSelectorMixin(flow_selector):: {
    flow_selector+: flow_selector,
  },
  withScheduler(scheduler):: {
    scheduler: scheduler,
  },
  withSchedulerMixin(scheduler):: {
    scheduler+: scheduler,
  },
  withSelectors(selectors):: {
    selectors: selectors,
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
}
