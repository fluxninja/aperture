{
  new():: {
  },
  withKubernetesObjectSelector(kubernetes_object_selector):: {
    kubernetes_object_selector: kubernetes_object_selector,
  },
  withKubernetesObjectSelectorMixin(kubernetes_object_selector):: {
    kubernetes_object_selector+: kubernetes_object_selector,
  },
  withScaleActuator(scale_actuator):: {
    scale_actuator: scale_actuator,
  },
  withScaleActuatorMixin(scale_actuator):: {
    scale_actuator+: scale_actuator,
  },
  withScaleReporter(scale_reporter):: {
    scale_reporter: scale_reporter,
  },
  withScaleReporterMixin(scale_reporter):: {
    scale_reporter+: scale_reporter,
  },
}
