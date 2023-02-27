local unaryoperatorins = import './unaryoperatorins.libsonnet';
local unaryoperatorouts = import './unaryoperatorouts.libsonnet';
{
  new():: {
    in_ports: {
      input: error 'Port input is missing',
    },
    out_ports: {
      output: error 'Port output is missing',
    },
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
