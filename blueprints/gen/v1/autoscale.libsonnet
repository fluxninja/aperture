{
  new():: {
  },
  withAutoscaler(autoscaler):: {
    autoscaler: autoscaler,
  },
  withAutoscalerMixin(autoscaler):: {
    autoscaler+: autoscaler,
  },
  withPodScaler(pod_scaler):: {
    pod_scaler: pod_scaler,
  },
  withPodScalerMixin(pod_scaler):: {
    pod_scaler+: pod_scaler,
  },
}
