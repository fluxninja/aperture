local autoScalingDefaults = import '../base/config-defaults.libsonnet';

autoScalingDefaults {
  policy+: {
    scaling_backend: {
      kubernetes_replicas: '__REQUIRED_FIELD__',
    },
  },
}
