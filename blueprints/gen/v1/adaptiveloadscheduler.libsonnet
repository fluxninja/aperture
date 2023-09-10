local adaptiveloadschedulerins = import './adaptiveloadschedulerins.libsonnet';
local adaptiveloadschedulerouts = import './adaptiveloadschedulerouts.libsonnet';
{
  new():: {
  },
  inPorts:: adaptiveloadschedulerins,
  outPorts:: adaptiveloadschedulerouts,
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
  withGradientThrottlingStrategy(gradient_throttling_strategy):: {
    gradient_throttling_strategy: gradient_throttling_strategy,
  },
  withGradientThrottlingStrategyMixin(gradient_throttling_strategy):: {
    gradient_throttling_strategy+: gradient_throttling_strategy,
  },
  withInPorts(in_ports):: {
    in_ports: in_ports,
  },
  withInPortsMixin(in_ports):: {
    in_ports+: in_ports,
  },
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
  withRangeThrottlingStrategy(range_throttling_strategy):: {
    range_throttling_strategy: range_throttling_strategy,
  },
  withRangeThrottlingStrategyMixin(range_throttling_strategy):: {
    range_throttling_strategy+: range_throttling_strategy,
  },
}
