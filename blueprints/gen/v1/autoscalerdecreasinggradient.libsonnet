local autoscalerdecreasinggradientins = import './autoscalerdecreasinggradientins.libsonnet';
{
  new():: {
  },
  inPorts:: autoscalerdecreasinggradientins,
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
