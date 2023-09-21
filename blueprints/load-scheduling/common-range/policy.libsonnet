local spec = import '../../spec.libsonnet';
local commonPolicyFn = import '../common/policy.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg, params={}, metadata={}) {
  local updatedConfig = config + cfg,
  local commonPolicy = commonPolicyFn(updatedConfig),

  local policyDef = commonPolicy.policyDef,
  local isConfirmationCriteria = commonPolicy.isConfirmationCriteria,

  local rangeDrivenLoadSchedulerComponent = spec.v1.Component.withFlowControl(
    spec.v1.FlowControl.withRangeDrivenLoadScheduler(
      local range = updatedConfig.policy.service_protection_core.range_driven_load_scheduler;

      spec.v1.RangeDrivenLoadScheduler.new()
      + spec.v1.RangeDrivenLoadScheduler.withInPorts({
        signal: spec.v1.Port.withSignalName('SIGNAL'),
        overload_confirmation: (if isConfirmationCriteria then spec.v1.Port.withSignalName('OVERLOAD_CONFIRMATION') else spec.v1.Port.withConstantSignal(1)),
      })
      + spec.v1.RangeDrivenLoadScheduler.withOutPorts({
        desired_load_multiplier: spec.v1.Port.withSignalName('DESIRED_LOAD_MULTIPLIER'),
        observed_load_multiplier: spec.v1.Port.withSignalName('OBSERVED_LOAD_MULTIPLIER'),
      })
      + spec.v1.RangeDrivenLoadScheduler.withParameters(range)
      + spec.v1.RangeDrivenLoadScheduler.withDryRunConfigKey('dry_run')
      + spec.v1.RangeDrivenLoadScheduler.withDryRun(updatedConfig.policy.service_protection_core.dry_run)
    ),
  ),

  local updatePolicyDef = policyDef
                          + spec.v1.Policy.new()
                          + spec.v1.Policy.withResources(policyDef.resources)
                          + spec.v1.Policy.withCircuit(
                            policyDef.circuit
                            + spec.v1.Circuit.withComponents(
                              policyDef.circuit.components
                              + [rangeDrivenLoadSchedulerComponent],
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
