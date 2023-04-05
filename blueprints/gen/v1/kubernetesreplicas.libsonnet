{
  new():: {
  },
  withDefaultConfig(default_config):: {
    default_config: default_config,
  },
  withDefaultConfigMixin(default_config):: {
    default_config+: default_config,
  },
  withDynamicConfigKey(dynamic_config_key):: {
    dynamic_config_key: dynamic_config_key,
  },
  withDynamicConfigKeyMixin(dynamic_config_key):: {
    dynamic_config_key+: dynamic_config_key,
  },
  withKubernetesObjectSelector(kubernetes_object_selector):: {
    kubernetes_object_selector: kubernetes_object_selector,
  },
  withKubernetesObjectSelectorMixin(kubernetes_object_selector):: {
    kubernetes_object_selector+: kubernetes_object_selector,
  },
}
