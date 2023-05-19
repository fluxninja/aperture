local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
  },
  inPorts:: ratelimiterins,
  withCustomLimits(custom_limits):: {
    custom_limits:
      if std.isArray(custom_limits)
      then custom_limits
      else [custom_limits],
  },
  withCustomLimitsMixin(custom_limits):: {
    custom_limits+: custom_limits,
  },
  withCustomLimitsConfigKey(custom_limits_config_key):: {
    custom_limits_config_key: custom_limits_config_key,
  },
  withCustomLimitsConfigKeyMixin(custom_limits_config_key):: {
    custom_limits_config_key+: custom_limits_config_key,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withParameters(parameters):: {
    parameters: parameters,
  },
  withParametersMixin(parameters):: {
    parameters+: parameters,
  },
}
