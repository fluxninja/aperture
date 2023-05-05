{
  new():: {
  },
  withAutoScaler(auto_scaler):: {
    auto_scaler: auto_scaler,
  },
  withAutoScalerMixin(auto_scaler):: {
    auto_scaler+: auto_scaler,
  },
  withPodScaler(pod_scaler):: {
    pod_scaler: pod_scaler,
  },
  withPodScalerMixin(pod_scaler):: {
    pod_scaler+: pod_scaler,
  },
}
