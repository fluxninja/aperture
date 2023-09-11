local integratorins = import './integratorins.libsonnet';
local integratorouts = import './integratorouts.libsonnet';
{
  new():: {
  },
  inPorts:: integratorins,
  outPorts:: integratorouts,
  withEvaluationInterval(evaluation_interval):: {
    evaluation_interval: evaluation_interval,
  },
  withEvaluationIntervalMixin(evaluation_interval):: {
    evaluation_interval+: evaluation_interval,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
  withInitialValue(initial_value):: {
    initial_value: initial_value,
  },
  withInitialValueMixin(initial_value):: {
    initial_value+: initial_value,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
}
