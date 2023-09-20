local ratelimiterins = import './ratelimiterins.libsonnet';
{
  new():: {
  },
  inPorts:: ratelimiterins,
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
  withRequestParameters(request_parameters):: {
    request_parameters: request_parameters,
  },
  withRequestParametersMixin(request_parameters):: {
    request_parameters+: request_parameters,
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
