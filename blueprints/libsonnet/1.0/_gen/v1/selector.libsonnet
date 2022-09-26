{
  new():: {
  },
  withFlowSelector(flow_selector):: {
    flow_selector: flow_selector,
  },
  withFlowSelectorMixin(flow_selector):: {
    flow_selector+: flow_selector,
  },
  withServiceSelector(service_selector):: {
    service_selector: service_selector,
  },
  withServiceSelectorMixin(service_selector):: {
    service_selector+: service_selector,
  },
}
