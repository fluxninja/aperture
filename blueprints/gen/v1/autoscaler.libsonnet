local autoscalerouts = import './autoscalerouts.libsonnet';
{
  new():: {
    out_ports: {
      actual_scale: error 'Port actual_scale is missing',
      configured_scale: error 'Port configured_scale is missing',
      desired_scale: error 'Port desired_scale is missing',
    },
  },
  outPorts:: autoscalerouts,
  withCooldownOverridePercentage(cooldown_override_percentage):: {
    cooldown_override_percentage: cooldown_override_percentage,
  },
  withCooldownOverridePercentageMixin(cooldown_override_percentage):: {
    cooldown_override_percentage+: cooldown_override_percentage,
  },
  withMaxScale(max_scale):: {
    max_scale: max_scale,
  },
  withMaxScaleMixin(max_scale):: {
    max_scale+: max_scale,
  },
  withMinScale(min_scale):: {
    min_scale: min_scale,
  },
  withMinScaleMixin(min_scale):: {
    min_scale+: min_scale,
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
  withScaler(scaler):: {
    scaler: scaler,
  },
  withScalerMixin(scaler):: {
    scaler+: scaler,
  },
}
