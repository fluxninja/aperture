local autoscalerouts = import './autoscalerouts.libsonnet';
{
  new():: {
  },
  outPorts:: autoscalerouts,
  withDryRun(dry_run):: {
    dry_run: dry_run,
  },
  withDryRunMixin(dry_run):: {
    dry_run+: dry_run,
  },
  withDryRunConfigKey(dry_run_config_key):: {
    dry_run_config_key: dry_run_config_key,
  },
  withDryRunConfigKeyMixin(dry_run_config_key):: {
    dry_run_config_key+: dry_run_config_key,
  },
  withOutPorts(out_ports):: {
    out_ports: out_ports,
  },
  withOutPortsMixin(out_ports):: {
    out_ports+: out_ports,
  },
  withScaleInControllers(scale_in_controllers):: {
    scale_in_controllers:
      if std.isArray(scale_in_controllers)
      then scale_in_controllers
      else [scale_in_controllers],
  },
  withScaleInControllersMixin(scale_in_controllers):: {
    scale_in_controllers+: scale_in_controllers,
  },
  withScaleOutControllers(scale_out_controllers):: {
    scale_out_controllers:
      if std.isArray(scale_out_controllers)
      then scale_out_controllers
      else [scale_out_controllers],
  },
  withScaleOutControllersMixin(scale_out_controllers):: {
    scale_out_controllers+: scale_out_controllers,
  },
  withScalingBackend(scaling_backend):: {
    scaling_backend: scaling_backend,
  },
  withScalingBackendMixin(scaling_backend):: {
    scaling_backend+: scaling_backend,
  },
  withScalingParameters(scaling_parameters):: {
    scaling_parameters: scaling_parameters,
  },
  withScalingParametersMixin(scaling_parameters):: {
    scaling_parameters+: scaling_parameters,
  },
}
