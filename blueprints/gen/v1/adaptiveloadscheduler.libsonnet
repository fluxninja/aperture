local adaptiveloadschedulerins = import './adaptiveloadschedulerins.libsonnet';
local adaptiveloadschedulerouts = import './adaptiveloadschedulerouts.libsonnet';
{
  new():: {
  },
  inPorts:: adaptiveloadschedulerins,
  outPorts:: adaptiveloadschedulerouts,
  withAimdThrottlingStrategy(aimd_throttling_strategy):: {
    aimd_throttling_strategy: aimd_throttling_strategy,
  },
  withAimdThrottlingStrategyMixin(aimd_throttling_strategy):: {
    aimd_throttling_strategy+: aimd_throttling_strategy,
  },
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
