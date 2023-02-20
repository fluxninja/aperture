{
  new():: {
  },
  withPodScaler(pod_scaler):: {
    pod_scaler: pod_scaler,
  },
  withPodScalerMixin(pod_scaler):: {
    pod_scaler+: pod_scaler,
  },
}
