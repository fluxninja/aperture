{
  new():: {
  },
  withAutoScaler(auto_scaler):: {
    auto_scaler: auto_scaler,
  },
  withAutoScalerMixin(auto_scaler):: {
    auto_scaler+: auto_scaler,
  },
  withPodAutoScaler(pod_auto_scaler):: {
    pod_auto_scaler: pod_auto_scaler,
  },
  withPodAutoScalerMixin(pod_auto_scaler):: {
    pod_auto_scaler+: pod_auto_scaler,
  },
  withPodScaler(pod_scaler):: {
    pod_scaler: pod_scaler,
  },
  withPodScalerMixin(pod_scaler):: {
    pod_scaler+: pod_scaler,
  },
}
