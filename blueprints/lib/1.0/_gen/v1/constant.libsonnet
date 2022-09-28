local constantouts = import './constantouts.libsonnet';
{
  new():: {
    out_ports: {
      output: error 'Port output is missing',
    },
  },
  outPorts:: constantouts,
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withValue(value):: {
    value: value,
  },
  withValueMixin(value):: {
    value+: value,
  },
}
