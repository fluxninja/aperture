local spec = import '../../spec.libsonnet';
local commonPolicyFn = import '../common/policy.libsonnet';

function(cfg, params={}) {
  local updatedConfig = cfg,
  local commonPolicy = commonPolicyFn(updatedConfig),

  local policyDef = commonPolicy.policyDef,
  local isConfirmationCriteria = commonPolicy.isConfirmationCriteria,

  local aiadLoadSchedulerComponent = spec.v1.Component.withFlowControl(
    spec.v1.FlowControl.withAiadLoadScheduler(
      local aiad = updatedConfig.policy.load_scheduling_core.aiad_load_scheduler;

      spec.v1.AIADLoadScheduler.new()
      + spec.v1.AIADLoadScheduler.withOverloadCondition(updatedConfig.policy.overload_condition)
      + spec.v1.AIADLoadScheduler.withInPorts({
        signal: spec.v1.Port.withSignalName('SIGNAL'),
        setpoint: spec.v1.Port.withSignalName('SETPOINT'),
        overload_confirmation: (if isConfirmationCriteria then spec.v1.Port.withSignalName('OVERLOAD_CONFIRMATION') else spec.v1.Port.withConstantSignal(1)),
      })
      + spec.v1.AIADLoadScheduler.withOutPorts({
        desired_load_multiplier: spec.v1.Port.withSignalName('DESIRED_LOAD_MULTIPLIER'),
        observed_load_multiplier: spec.v1.Port.withSignalName('OBSERVED_LOAD_MULTIPLIER'),
      })
      + spec.v1.AIADLoadScheduler.withParameters(aiad)
      + spec.v1.AIADLoadScheduler.withDryRunConfigKey('dry_run')
      + spec.v1.AIADLoadScheduler.withDryRun(updatedConfig.policy.load_scheduling_core.dry_run)
    ),
  ),

  local updatePolicyDef = policyDef
                          + spec.v1.Policy.new()
                          + spec.v1.Policy.withResources(policyDef.resources)
                          + spec.v1.Policy.withCircuit(
                            policyDef.circuit
                            + spec.v1.Circuit.withComponents(
                              policyDef.circuit.components
                              + [aiadLoadSchedulerComponent],
                            ),
                          ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: updatedConfig.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: updatePolicyDef,
  },

  policyDef: updatePolicyDef,
  policyResource: policyResource,
}
