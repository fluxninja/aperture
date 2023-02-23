local podautoscalerouts = import './podautoscalerouts.libsonnet';
{
  new():: {
    out_ports: {
      actual_replicas: error 'Port actual_replicas is missing',
      configured_replicas: error 'Port configured_replicas is missing',
      desired_replicas: error 'Port desired_replicas is missing',
    },
  },
  outPorts:: podautoscalerouts,
  withCooldownOverridePercentage(cooldown_override_percentage):: {
    cooldown_override_percentage: cooldown_override_percentage,
  },
  withCooldownOverridePercentageMixin(cooldown_override_percentage):: {
    cooldown_override_percentage+: cooldown_override_percentage,
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
  withMaxReplicas(max_replicas):: {
    max_replicas: max_replicas,
  },
  withMaxReplicasMixin(max_replicas):: {
    max_replicas+: max_replicas,
  },
  withMinReplicas(min_replicas):: {
    min_replicas: min_replicas,
  },
  withMinReplicasMixin(min_replicas):: {
    min_replicas+: min_replicas,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withScaleInAlerterParameters(scale_in_alerter_parameters):: {
    scale_in_alerter_parameters: scale_in_alerter_parameters,
  },
  withScaleInAlerterParametersMixin(scale_in_alerter_parameters):: {
    scale_in_alerter_parameters+: scale_in_alerter_parameters,
  },
  withScaleInControllers(scale_in_controllers):: {
    scale_in_controllers:
      if std.isArray(scale_in_controllers)
      then scale_in_controllers
      else [scale_in_controllers],
  },
  withScaleInControllersMixin(scale_in_controllers):: {
    scale_in_controllers+: scale_in_controllers,
  },
  withScaleInCooldown(scale_in_cooldown):: {
    scale_in_cooldown: scale_in_cooldown,
  },
  withScaleInCooldownMixin(scale_in_cooldown):: {
    scale_in_cooldown+: scale_in_cooldown,
  },
  withScaleInMaxPercentage(scale_in_max_percentage):: {
    scale_in_max_percentage: scale_in_max_percentage,
  },
  withScaleInMaxPercentageMixin(scale_in_max_percentage):: {
    scale_in_max_percentage+: scale_in_max_percentage,
  },
  withScaleOutAlerterParameters(scale_out_alerter_parameters):: {
    scale_out_alerter_parameters: scale_out_alerter_parameters,
  },
  withScaleOutAlerterParametersMixin(scale_out_alerter_parameters):: {
    scale_out_alerter_parameters+: scale_out_alerter_parameters,
  },
  withScaleOutControllers(scale_out_controllers):: {
    scale_out_controllers:
      if std.isArray(scale_out_controllers)
      then scale_out_controllers
      else [scale_out_controllers],
  },
  withScaleOutControllersMixin(scale_out_controllers):: {
    scale_out_controllers+: scale_out_controllers,
  },
  withScaleOutCooldown(scale_out_cooldown):: {
    scale_out_cooldown: scale_out_cooldown,
  },
  withScaleOutCooldownMixin(scale_out_cooldown):: {
    scale_out_cooldown+: scale_out_cooldown,
  },
  withScaleOutMaxPercentage(scale_out_max_percentage):: {
    scale_out_max_percentage: scale_out_max_percentage,
  },
  withScaleOutMaxPercentageMixin(scale_out_max_percentage):: {
    scale_out_max_percentage+: scale_out_max_percentage,
  },
}
