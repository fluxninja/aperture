local nestedsignalingressouts = import './nestedsignalingressouts.libsonnet';
{
  new():: {
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
