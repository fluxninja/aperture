local alerterins = import './alerterins.libsonnet';
{
  new():: {
    in_ports: {
      signal: error 'Port signal is missing',
    },
  },
  inPorts:: alerterins,
  withAlerterConfig(alerter_config):: {
    alerter_config: alerter_config,
  },
  withAlerterConfigMixin(alerter_config):: {
    alerter_config+: alerter_config,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
}
