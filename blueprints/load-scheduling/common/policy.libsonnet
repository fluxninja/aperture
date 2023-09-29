local consts = import '../../consts.libsonnet';
local spec = import '../../spec.libsonnet';
local utils = import '../../utils/utils.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg, params={}) {
  local updatedConfig = config + cfg,

  local addOverloadConfirmation = function(confirmationAccumulator, confirmation) {
    local promQLSignalName = 'PROMQL_' + std.toString(confirmationAccumulator.overload_confirmation_signals_count),

    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(
      spec.v1.PromQL.withQueryString(confirmation.query_string)
      + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=consts.metricScrapeInterval)
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

  local policyDef =
    spec.v1.Policy.new()
    + spec.v1.Policy.withResources(utils.resources(updatedConfig.policy.resources).updatedResources)
    + spec.v1.Policy.withCircuit(
      spec.v1.Circuit.new()
      + spec.v1.Circuit.withEvaluationInterval(evaluation_interval=consts.circuitEvaluationInterval)
      + spec.v1.Circuit.withComponents(
        confirmationAccumulator.components
        + (if isConfirmationCriteria then [overloadConfirmationAnd] else [])
        + updatedConfig.policy.components,
      ),
    ),

  policyDef: policyDef,
  isConfirmationCriteria: isConfirmationCriteria,
}
