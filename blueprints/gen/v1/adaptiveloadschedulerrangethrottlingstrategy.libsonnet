local adaptiveloadschedulerrangethrottlingstrategyins = import './adaptiveloadschedulerrangethrottlingstrategyins.libsonnet';
{
  new():: {
  },
  inPorts:: adaptiveloadschedulerrangethrottlingstrategyins,
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
