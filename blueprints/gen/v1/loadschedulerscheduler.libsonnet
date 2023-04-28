local loadschedulerschedulerouts = import './loadschedulerschedulerouts.libsonnet';
{
  new():: {
  },
  outPorts:: loadschedulerschedulerouts,
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
}
