local autoscalerincreasinggradientins = import './autoscalerincreasinggradientins.libsonnet';
{
  new():: {
  },
  inPorts:: autoscalerincreasinggradientins,
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
