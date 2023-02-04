local horizontalpodscalerscaleactuatorins = import './horizontalpodscalerscaleactuatorins.libsonnet';
{
  new():: {
    in_ports: {
      desired_replicas: error 'Port desired_replicas is missing',
    },
  },
  inPorts:: horizontalpodscalerscaleactuatorins,
  withDefaultConfig(default_config):: {
    default_config: default_config,
  },
  withDefaultConfigMixin(default_config):: {
    default_config+: default_config,
  },
  withDynamicConfigKey(dynamic_config_key):: {
    dynamic_config_key: dynamic_config_key,
  },
  withDynamicConfigKeyMixin(dynamic_config_key):: {
    dynamic_config_key+: dynamic_config_key,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
}
