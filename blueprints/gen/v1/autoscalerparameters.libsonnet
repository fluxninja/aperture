{
  new():: {
  },
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
  withMinScale(min_scale):: {
    min_scale: min_scale,
  },
  withMinScaleMixin(min_scale):: {
    min_scale+: min_scale,
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
  withScaler(scaler):: {
    scaler: scaler,
  },
  withScalerMixin(scaler):: {
    scaler+: scaler,
  },
}
