local samplerins = import './samplerins.libsonnet';
{
  new():: {
  },
  inPorts:: samplerins,
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
  withPassThroughLabelValues(pass_through_label_values):: {
    pass_through_label_values:
      if std.isArray(pass_through_label_values)
      then pass_through_label_values
      else [pass_through_label_values],
  },
  withPassThroughLabelValuesMixin(pass_through_label_values):: {
    pass_through_label_values+: pass_through_label_values,
  },
  withPassThroughLabelValuesConfigKey(pass_through_label_values_config_key):: {
    pass_through_label_values_config_key: pass_through_label_values_config_key,
  },
  withPassThroughLabelValuesConfigKeyMixin(pass_through_label_values_config_key):: {
    pass_through_label_values_config_key+: pass_through_label_values_config_key,
  },
}
