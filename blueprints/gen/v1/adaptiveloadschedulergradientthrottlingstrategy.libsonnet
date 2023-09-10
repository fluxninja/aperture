local adaptiveloadschedulergradientthrottlingstrategyins = import './adaptiveloadschedulergradientthrottlingstrategyins.libsonnet';
{
  new():: {
  },
  inPorts:: adaptiveloadschedulergradientthrottlingstrategyins,
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
