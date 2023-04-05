local podautoscalerouts = import './podautoscalerouts.libsonnet';
{
  new():: {
  },
  outPorts:: podautoscalerouts,
  withCooldownOverridePercentage(cooldown_override_percentage):: {
    cooldown_override_percentage: cooldown_override_percentage,
  },
  withCooldownOverridePercentageMixin(cooldown_override_percentage):: {
    cooldown_override_percentage+: cooldown_override_percentage,
  },
  withMaxReplicas(max_replicas):: {
    max_replicas: max_replicas,
  },
  withMaxReplicasMixin(max_replicas):: {
    max_replicas+: max_replicas,
  },
  withMaxScaleInPercentage(max_scale_in_percentage):: {
    max_scale_in_percentage: max_scale_in_percentage,
  },
  withMaxScaleInPercentageMixin(max_scale_in_percentage):: {
    max_scale_in_percentage+: max_scale_in_percentage,
  },
  withMaxScaleOutPercentage(max_scale_out_percentage):: {
    max_scale_out_percentage: max_scale_out_percentage,
  },
  withMaxScaleOutPercentageMixin(max_scale_out_percentage):: {
    max_scale_out_percentage+: max_scale_out_percentage,
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
  withPodScaler(pod_scaler):: {
    pod_scaler: pod_scaler,
  },
  withPodScalerMixin(pod_scaler):: {
    pod_scaler+: pod_scaler,
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
}
