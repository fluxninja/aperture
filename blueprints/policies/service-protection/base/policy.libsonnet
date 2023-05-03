local spec = import '../../../spec.libsonnet';
local config = import './config-defaults.libsonnet';

function(cfg) {
  local params = config.common + config.policy + cfg,

  local policyName = params.policy_name,

  local flux_meters = params.flux_meters,

  local addOverloadConfirmation = function(confirmationAccumulator, confirmation) {
    local evaluationInterval = params.evaluation_interval,
    local promQLSignalName = 'PROMQL_' + std.toString(confirmationAccumulator.enabled_signals_count),
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(
      spec.v1.PromQL.withQueryString(confirmation.query_string)
      + spec.v1.PromQL.withEvaluationInterval(evaluationInterval)
      + spec.v1.PromQL.withOutPorts({
        output: spec.v1.Port.withSignalName(promQLSignalName),
      })
    )),
    local confirmationSignal = 'CONFIRMATION_SIGNAL_' + std.toString(confirmationAccumulator.enabled_signals_count),
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

    local enabledSignal = 'ENABLED_' + std.toString(confirmationAccumulator.enabled_signals_count),
    local firstValidComponent = spec.v1.Component.withFirstValid(
      spec.v1.FirstValid.withInPorts({
        inputs: [
          spec.v1.Port.withSignalName(confirmationSignal),
          spec.v1.Port.withConstantSignal(0),
        ],
      })
      + spec.v1.FirstValid.withOutPorts({
        output: spec.v1.Port.withSignalName(enabledSignal),
      }),
    ),


    enabled_signals: confirmationAccumulator.enabled_signals + [enabledSignal],
    enabled_signals_count: confirmationAccumulator.enabled_signals_count + 1,
    components: confirmationAccumulator.components + [promQLComponent, confirmationDecider, firstValidComponent],
  },

  local confirmationAccumulatorInitial = {
    enabled_signals: [],
    enabled_signals_count: 0,
    components: [],
  },

  local confirmationAccumulator = std.foldl(
    addOverloadConfirmation,
    (if std.objectHas(params, 'overload_confirmations') then params.overload_confirmations else []),
    confirmationAccumulatorInitial
  ),

  local enabledAnd = spec.v1.Component.withAnd(
    spec.v1.And.withInPorts({
      inputs: [
        spec.v1.Port.withSignalName(signal)
        for signal in confirmationAccumulator.enabled_signals
      ],
    })
    + spec.v1.Or.withOutPorts({
      output: spec.v1.Port.withSignalName('ENABLED'),
    }),
  ),

  local isConfirmationCriteria = std.length(confirmationAccumulator.enabled_signals) > 0,
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
        enabled: (if isConfirmationCriteria then spec.v1.Port.withSignalName('ENABLED') else spec.v1.Port.withConstantSignal(1)),
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
        + (if isConfirmationCriteria then [enabledAnd] else [])
        + [
          adaptiveLoadSchedulerComponent,
        ]
        + params.components,
      ),
    ),

  policyDef: policyDef,
}
