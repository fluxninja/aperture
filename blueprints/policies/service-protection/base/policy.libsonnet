local spec = import '../../../spec.libsonnet';
local utils = import '../../../utils/utils.libsonnet';
local baseAutoScalingPolicyFn = import '../../auto-scaling/base/policy.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg, params={}, metadata={}) {
  local updatedConfig = config + cfg,

  local addOverloadConfirmation = function(confirmationAccumulator, confirmation) {
    local evaluationInterval = updatedConfig.policy.evaluation_interval,

    local promQLSignalName = 'PROMQL_' + std.toString(confirmationAccumulator.overload_confirmation_signals_count),

    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(
      spec.v1.PromQL.withQueryString(confirmation.query_string)
      + spec.v1.PromQL.withEvaluationInterval(evaluationInterval)
      + spec.v1.PromQL.withOutPorts({
        output: spec.v1.Port.withSignalName(promQLSignalName),
      })
    )),

    local confirmationSignal = 'CONFIRMATION_SIGNAL_' + std.toString(confirmationAccumulator.overload_confirmation_signals_count),

    local confirmationDecider = spec.v1.Component.withDecider(
      spec.v1.Decider.withOperator(confirmation.operator)
      + spec.v1.Decider.withInPorts({
        lhs: spec.v1.Port.withSignalName(promQLSignalName),
        rhs: spec.v1.Port.withConstantSignal(confirmation.threshold),
      })
      + spec.v1.Decider.withOutPorts({
        output: spec.v1.Port.withSignalName(confirmationSignal),
      })
    ),

    local overloadConfirmationSignal = 'OVERLOAD_CONFIRMATION_' + std.toString(confirmationAccumulator.overload_confirmation_signals_count),

    local firstValidComponent = spec.v1.Component.withFirstValid(
      spec.v1.FirstValid.withInPorts({
        inputs: [
          spec.v1.Port.withSignalName(confirmationSignal),
          spec.v1.Port.withConstantSignal(0),  // overload confirmation is assumed false if no confirmation signal is received
        ],
      })
      + spec.v1.FirstValid.withOutPorts({
        output: spec.v1.Port.withSignalName(overloadConfirmationSignal),
      }),
    ),

    overload_confirmation_signals: confirmationAccumulator.overload_confirmation_signals + [overloadConfirmationSignal],
    overload_confirmation_signals_count: confirmationAccumulator.overload_confirmation_signals_count + 1,
    components: confirmationAccumulator.components + [promQLComponent, confirmationDecider, firstValidComponent],
  },

  local confirmationAccumulatorInitial = {
    overload_confirmation_signals: [],
    overload_confirmation_signals_count: 0,
    components: [],
  },

  local confirmationAccumulator = std.foldl(
    addOverloadConfirmation,
    (if std.objectHas(updatedConfig.policy.service_protection_core, 'overload_confirmations') then updatedConfig.policy.service_protection_core.overload_confirmations else []),
    confirmationAccumulatorInitial
  ),

  local overloadConfirmationAnd = spec.v1.Component.withAnd(
    spec.v1.And.withInPorts({
      inputs: [
        spec.v1.Port.withSignalName(signal)
        for signal in confirmationAccumulator.overload_confirmation_signals
      ],
    })
    + spec.v1.And.withOutPorts({
      output: spec.v1.Port.withSignalName('OVERLOAD_CONFIRMATION'),
    }),
  ),

  local isConfirmationCriteria = std.length(confirmationAccumulator.overload_confirmation_signals) > 0,

  local adaptiveLoadSchedulerComponent = spec.v1.Component.withFlowControl(
    spec.v1.FlowControl.withAdaptiveLoadScheduler(
      local adaptiveLoadScheduler = updatedConfig.policy.service_protection_core.adaptive_load_scheduler;
      spec.v1.AdaptiveLoadScheduler.new()
      + spec.v1.AdaptiveLoadScheduler.withParameters(adaptiveLoadScheduler)
      + spec.v1.AdaptiveLoadScheduler.withDryRunConfigKey('dry_run')
      + spec.v1.AdaptiveLoadScheduler.withDryRun(updatedConfig.policy.service_protection_core.dry_run)
      + spec.v1.AdaptiveLoadScheduler.withInPorts({
        overload_confirmation: (if isConfirmationCriteria then spec.v1.Port.withSignalName('OVERLOAD_CONFIRMATION') else spec.v1.Port.withConstantSignal(1)),
        signal: spec.v1.Port.withSignalName('SIGNAL'),
        setpoint: spec.v1.Port.withSignalName('SETPOINT'),
      })
      + spec.v1.AdaptiveLoadScheduler.withOutPorts({
        desired_load_multiplier: spec.v1.Port.withSignalName('DESIRED_LOAD_MULTIPLIER'),
        observed_load_multiplier: spec.v1.Port.withSignalName('OBSERVED_LOAD_MULTIPLIER'),
      }),
    ),
  ),

  /** Auto scale escalation **/

  local scaleInControllers =
    (if std.objectHas(params, 'policy') &&
        std.objectHas(params.policy, 'auto_scaling') &&
        std.objectHas(updatedConfig.policy.auto_scaling, 'periodic_decrease')
     then
       [
         spec.v1.ScaleInController.new()
         + spec.v1.ScaleInController.withAlerter(
           spec.v1.AlerterParameters.new()
           + spec.v1.AlerterParameters.withAlertName('Periodic scale in intended')
         )
         + spec.v1.ScaleInController.withController(
           spec.v1.ScaleInControllerController.new()
           + spec.v1.ScaleInControllerController.withPeriodic(updatedConfig.policy.auto_scaling.periodic_decrease)
         ),
       ]
     else []),

  local scaleOutControllers =
    (if std.objectHas(params, 'policy') &&
        std.objectHas(params.policy, 'auto_scaling') then [
       spec.v1.ScaleOutController.new()
       + spec.v1.ScaleOutController.withAlerter(
         spec.v1.AlerterParameters.new()
         + spec.v1.AlerterParameters.withAlertName('Load based scale out intended')
       )
       + spec.v1.ScaleOutController.withController(
         spec.v1.ScaleOutControllerController.new()
         + spec.v1.ScaleOutControllerController.withGradient(
           spec.v1.IncreasingGradient.new()
           + spec.v1.IncreasingGradient.withInPorts(
             spec.v1.IncreasingGradientIns.new()
             + spec.v1.IncreasingGradientIns.withSignal(spec.v1.Port.withSignalName('DESIRED_LOAD_MULTIPLIER'))
             + spec.v1.IncreasingGradientIns.withSetpoint(spec.v1.Port.withConstantSignal(1.0))
           )
           + spec.v1.IncreasingGradient.withParameters(
             spec.v1.IncreasingGradientParameters.new()
             + spec.v1.IncreasingGradientParameters.withSlope(-1.0)
           )
         )
       ),
     ] else []),

  local policyDef =
    spec.v1.Policy.new()
    + spec.v1.Policy.withResources(utils.resources(updatedConfig.policy.resources).updatedResources)
    + spec.v1.Policy.withCircuit(
      spec.v1.Circuit.new()
      + spec.v1.Circuit.withEvaluationInterval(evaluation_interval=updatedConfig.policy.evaluation_interval)
      + spec.v1.Circuit.withComponents(
        confirmationAccumulator.components
        + (if isConfirmationCriteria then [overloadConfirmationAnd] else [])
        + [
          adaptiveLoadSchedulerComponent,
        ]
        + updatedConfig.policy.components,
      ),
    ) +
    (
      if std.objectHas(params, 'policy') &&
         std.objectHas(params.policy, 'auto_scaling') then
        local autoScalingUpdatedConfig = {
          policy+: updatedConfig.policy.auto_scaling {
            policy_name: updatedConfig.policy.policy_name,
            // Set empty defaults for promql_scale_out_controllers and promql_scale_in_controllers
            promql_scale_out_controllers: if std.objectHas(updatedConfig.policy.auto_scaling, 'promql_scale_out_controllers') then updatedConfig.policy.auto_scaling.promql_scale_out_controllers else [],
            promql_scale_in_controllers: if std.objectHas(updatedConfig.policy.auto_scaling, 'promql_scale_in_controllers') then updatedConfig.policy.auto_scaling.promql_scale_in_controllers else [],
          },
        };

        local baseAutoScalingPolicy = baseAutoScalingPolicyFn(autoScalingUpdatedConfig).policyDef;
        {
          circuit+: {
            components+: std.map(
              function(component) if std.objectHas(component, 'auto_scale') then
                component {
                  auto_scale+: {
                    auto_scaler+: {
                      scale_out_controllers+: scaleOutControllers,
                      scale_in_controllers+: scaleInControllers,
                    },
                  },
                }
              else component,
              baseAutoScalingPolicy.circuit.components
            ),
          },
        } else {}
    ),
  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: updatedConfig.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
      annotations: {
        [if std.objectHas(metadata, 'values') then 'fluxninja.com/values']: metadata.values,
        [if std.objectHas(metadata, 'blueprints_uri') then 'fluxninja.com/blueprints-uri']: metadata.blueprints_uri,
        [if std.objectHas(metadata, 'blueprint_name') then 'fluxninja.com/blueprint-name']: metadata.blueprint_name,
      },
    },
    spec: policyDef,
  },
  policyDef: policyDef,
  policyResource: policyResource,
}
