{
  new():: {
  },
  withPodAutoscaler(pod_autoscaler):: {
    pod_autoscaler: pod_autoscaler,
  },
  withPodAutoscalerMixin(pod_autoscaler):: {
    pod_autoscaler+: pod_autoscaler,
  },
  withPodScaler(pod_scaler):: {
    pod_scaler: pod_scaler,
  },
  withPodScalerMixin(pod_scaler):: {
    pod_scaler+: pod_scaler,
  },
}
