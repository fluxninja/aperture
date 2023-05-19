local leakybucketratelimiterins = import './leakybucketratelimiterins.libsonnet';
{
  new():: {
  },
  inPorts:: leakybucketratelimiterins,
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
  withSelectors(selectors):: {
    selectors:
      if std.isArray(selectors)
      then selectors
      else [selectors],
  },
  withSelectorsMixin(selectors):: {
    selectors+: selectors,
  },
}
