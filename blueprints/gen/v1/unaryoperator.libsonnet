local unaryoperatorins = import './unaryoperatorins.libsonnet';
local unaryoperatorouts = import './unaryoperatorouts.libsonnet';
{
  new():: {
  },
  inPorts:: unaryoperatorins,
  outPorts:: unaryoperatorouts,
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withOperator(operator):: {
    operator: operator,
  },
  withOperatorMixin(operator):: {
    operator+: operator,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
