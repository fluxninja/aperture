local nestedsignalingressouts = import './nestedsignalingressouts.libsonnet';
{
  new():: {
    out_ports: {
      signal: error 'Port signal is missing',
    },
  },
  outPorts:: nestedsignalingressouts,
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withPortName(port_name):: {
    port_name: port_name,
  },
  withPortNameMixin(port_name):: {
    port_name+: port_name,
  },
}
