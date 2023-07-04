local autoScalingDefaults = import '../common/config-defaults.libsonnet';

autoScalingDefaults {
  policy+: {
    scaling_backend: {
      kubernetes_replicas: '__REQUIRED_FIELD__',
    },
  },
}
