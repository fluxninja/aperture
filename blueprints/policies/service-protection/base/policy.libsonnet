local spec = import '../../../spec.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg) {
  local params = config.common + config.policy + cfg,

  local policyName = params.policy_name,

  local flux_meters = params.flux_meters,

  local addOverloadConfirmation = function(confirmationAccumulator, confirmation) {
    local evaluationInterval = params.evaluation_interval,
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
          spec.v1.Port.withConstantSignal(0),
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
    (if std.objectHas(params, 'overload_confirmations') then params.overload_confirmations else []),
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
      local adaptiveLoadScheduler = params.service_protection_core.adaptive_load_scheduler;
      spec.v1.AdaptiveLoadScheduler.new()
      + spec.v1.AdaptiveLoadScheduler.withSelectors(adaptiveLoadScheduler.selectors)
      + spec.v1.AdaptiveLoadScheduler.withSchedulerParameters(adaptiveLoadScheduler.scheduler)
      + spec.v1.AdaptiveLoadScheduler.withGradientParameters(adaptiveLoadScheduler.gradient)
      + spec.v1.AdaptiveLoadScheduler.withMaxLoadMultiplier(adaptiveLoadScheduler.max_load_multiplier)
      + spec.v1.AdaptiveLoadScheduler.withLoadMultiplierLinearIncrement(adaptiveLoadScheduler.load_multiplier_linear_increment)
      + spec.v1.AdaptiveLoadScheduler.withAlerterParameters(adaptiveLoadScheduler.alerter)
      + spec.v1.AdaptiveLoadScheduler.withDynamicConfigKey('load_scheduler')
      + spec.v1.AdaptiveLoadScheduler.withDefaultConfig(adaptiveLoadScheduler.default_config)
      + spec.v1.AdaptiveLoadScheduler.withInPorts({
        overload_confirmation: (if isConfirmationCriteria then spec.v1.Port.withSignalName('OVERLOAD_CONFIRMATION') else spec.v1.Port.withConstantSignal(1)),
        signal: spec.v1.Port.withSignalName('SIGNAL'),
        setpoint: spec.v1.Port.withSignalName('SETPOINT'),
      })
      + spec.v1.AdaptiveLoadScheduler.withOutPorts({
        desired_load_multiplier: spec.v1.Port.withSignalName('DESIRED_LOAD_MULTIPLIER'),
        observed_load_multiplier: spec.v1.Port.withSignalName('OBSERVED_LOAD_MULTIPLIER'),
        accepted_token_rate: spec.v1.Port.withSignalName('ACCEPTED_CONCURRENCY'),
        incoming_token_rate: spec.v1.Port.withSignalName('INCOMING_CONCURRENCY'),
      }),
    ),
  ),

  local policyDef =
    spec.v1.Policy.new()
    + spec.v1.Policy.withResources(params.resources)
    + spec.v1.Policy.withCircuit(
      spec.v1.Circuit.new()
      + spec.v1.Circuit.withEvaluationInterval(evaluation_interval=params.evaluation_interval)
      + spec.v1.Circuit.withComponents(
        confirmationAccumulator.components
        + (if isConfirmationCriteria then [overloadConfirmationAnd] else [])
        + [
          adaptiveLoadSchedulerComponent,
        ]
        + params.components,
      ),
    ),

  policyDef: policyDef,
}
