local ratelimiterins = import './ratelimiterins.libsonnet';
local ratelimiterouts = import './ratelimiterouts.libsonnet';
{
  new():: {
  },
  inPorts:: ratelimiterins,
  outPorts:: ratelimiterouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
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
