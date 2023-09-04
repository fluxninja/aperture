local spec = import '../../spec.libsonnet';
local utils = import '../../utils/utils.libsonnet';
local config = import './config.libsonnet';

function(cfg, metadata={}) {
  local params = config + cfg,

  local policyName = params.policy.policy_name,

  local addPromQLDriver = function(driverAccumulator, driver) {
    local promQLSignalName = 'PROMQL_' + std.toString(driverAccumulator.promql_driver_count),
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(driver.query_string)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(promQLSignalName),
                                                                                 }))),
    local criteria = (if std.objectHas(driver, 'criteria') then driver.criteria else {}),
    local isForward = (if std.objectHas(criteria, 'forward') then criteria.forward != null else false),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator(criteria.forward.operator)
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(promQLSignalName),
                                                           rhs: spec.v1.Port.withConstantSignal(criteria.forward.threshold),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(criteria, 'backward') then criteria.backward != null else false),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator(criteria.backward.operator)
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(promQLSignalName),
                                                            rhs: spec.v1.Port.withConstantSignal(criteria.backward.threshold),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(criteria, 'reset') then criteria.reset != null else false),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator(criteria.reset.operator)
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(promQLSignalName),
                                                         rhs: spec.v1.Port.withConstantSignal(criteria.reset.threshold),
                                                       })
                                                       + spec.v1.Decider.withOutPorts({
                                                         output: spec.v1.Port.withSignalName(resetSignal),
                                                       })),

    flux_meters: driverAccumulator.flux_meters,
    forward_signals: driverAccumulator.forward_signals + (if isForward then [forwardSignal] else []),
    backward_signals: driverAccumulator.backward_signals + (if isBackward then [backwardSignal] else []),
    reset_signals: driverAccumulator.reset_signals + (if isReset then [resetSignal] else []),
    forward_signals_count: driverAccumulator.forward_signals_count + if isForward then 1 else 0,
    backward_signals_count: driverAccumulator.backward_signals_count + if isBackward then 1 else 0,
    reset_signals_count: driverAccumulator.reset_signals_count + if isReset then 1 else 0,
    components: driverAccumulator.components + [promQLComponent] + (if isForward then [forwardDecider] else []) + (if isBackward then [backwardDecider] else []) + (if isReset then [resetDecider] else []),
    promql_driver_count: driverAccumulator.promql_driver_count + 1,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count,
  },

  local addAverageLatencyDriver = function(driverAccumulator, driver) {
    local flux_meter_name = policyName + '/average_latency/' + std.toString(driverAccumulator.average_latency_driver_count),
    local averageLatencySignalName = 'AVERAGE_LATENCY_' + std.toString(driverAccumulator.average_latency_driver_count),
    local q = 'sum(increase(flux_meter_sum{flow_status="OK", flux_meter_name="%(flux_meter_name)s", policy_name="%(policy_name)s"}[30s]))/sum(increase(flux_meter_count{flow_status="OK", flux_meter_name="%(flux_meter_name)s", policy_name="%(policy_name)s"}[30s]))' % { flux_meter_name: flux_meter_name, policy_name: policyName },
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(q)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                                                 }))),

    local criteria = (if std.objectHas(driver, 'criteria') then driver.criteria else {}),
    local isForward = (if std.objectHas(criteria, 'forward') then criteria.forward != null else false),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('lt')
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                           rhs: spec.v1.Port.withConstantSignal(criteria.forward.threshold),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(criteria, 'backward') then criteria.backward != null else false),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                            rhs: spec.v1.Port.withConstantSignal(criteria.backward.threshold),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(criteria, 'reset') then criteria.reset != null else false),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                         rhs: spec.v1.Port.withConstantSignal(criteria.reset.threshold),
                                                       })
                                                       + spec.v1.Decider.withOutPorts({
                                                         output: spec.v1.Port.withSignalName(resetSignal),
                                                       })),
    local newFluxMeters = {
      [flux_meter_name]: spec.v1.FluxMeter.withSelectors(driver.selectors),
    },

    flux_meters: driverAccumulator.flux_meters + newFluxMeters,
    forward_signals: driverAccumulator.forward_signals + (if isForward then [forwardSignal] else []),
    backward_signals: driverAccumulator.backward_signals + (if isBackward then [backwardSignal] else []),
    reset_signals: driverAccumulator.reset_signals + (if isReset then [resetSignal] else []),
    forward_signals_count: driverAccumulator.forward_signals_count + if isForward then 1 else 0,
    backward_signals_count: driverAccumulator.backward_signals_count + if isBackward then 1 else 0,
    reset_signals_count: driverAccumulator.reset_signals_count + if isReset then 1 else 0,
    components: driverAccumulator.components + [promQLComponent] + (if isForward then [forwardDecider] else []) + (if isBackward then [backwardDecider] else []) + (if isReset then [resetDecider] else []),
    promql_driver_count: driverAccumulator.promql_driver_count,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count + 1,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count,
  },

  local addPercentileLatencyDriver = function(driverAccumulator, driver) {
    local flux_meter_name = policyName + '/percentile_latency/' + std.toString(driverAccumulator.percentile_latency_driver_count),
    local percentileLatencySignalName = 'PERCENTILE_LATENCY_' + std.toString(driverAccumulator.percentile_latency_driver_count),
    local q = 'histogram_quantile(%(percentile)f, sum(rate(flux_meter_bucket{flow_status="OK", flux_meter_name="%(flux_meter_name)s", policy_name="%(policy_name)s"}[30s])) by (le))' % { percentile: driver.percentile, flux_meter_name: flux_meter_name, policy_name: policyName },
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(q)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                                                 }))),

    local criteria = (if std.objectHas(driver, 'criteria') then driver.criteria else {}),
    local isForward = (if std.objectHas(criteria, 'forward') then criteria.forward != null else false),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('lt')
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                           rhs: spec.v1.Port.withConstantSignal(criteria.forward.threshold),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(criteria, 'backward') then criteria.backward != null else false),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                            rhs: spec.v1.Port.withConstantSignal(criteria.backward.threshold),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(criteria, 'reset') then criteria.reset != null else false),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                         rhs: spec.v1.Port.withConstantSignal(criteria.reset.threshold),
                                                       })
                                                       + spec.v1.Decider.withOutPorts({
                                                         output: spec.v1.Port.withSignalName(resetSignal),
                                                       })),

    local newFluxMeters = {
      [flux_meter_name]: driver.flux_meter,
    },

    flux_meters: driverAccumulator.flux_meters + newFluxMeters,
    forward_signals: driverAccumulator.forward_signals + (if isForward then [forwardSignal] else []),
    backward_signals: driverAccumulator.backward_signals + (if isBackward then [backwardSignal] else []),
    reset_signals: driverAccumulator.reset_signals + (if isReset then [resetSignal] else []),
    forward_signals_count: driverAccumulator.forward_signals_count + if isForward then 1 else 0,
    backward_signals_count: driverAccumulator.backward_signals_count + if isBackward then 1 else 0,
    reset_signals_count: driverAccumulator.reset_signals_count + if isReset then 1 else 0,
    components: driverAccumulator.components + [promQLComponent] + (if isForward then [forwardDecider] else []) + (if isBackward then [backwardDecider] else []) + (if isReset then [resetDecider] else []),
    promql_driver_count: driverAccumulator.promql_driver_count,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count + 1,
  },

  local driverAccumulatorInitial = {
    flux_meters: {},
    forward_signals: [],
    backward_signals: [],
    reset_signals: [],
    forward_signals_count: 0,
    backward_signals_count: 0,
    reset_signals_count: 0,
    components: [],
    promql_driver_count: 0,
    average_latency_driver_count: 0,
    percentile_latency_driver_count: 0,
  },


  local driverAccumulatorStep1 = std.foldl(
    addPromQLDriver,
    (if std.objectHas(params.policy.drivers, 'promql_drivers') then params.policy.drivers.promql_drivers else []),
    driverAccumulatorInitial
  ),
  local driverAccumulatorStep2 = std.foldl(
    addAverageLatencyDriver,
    (if std.objectHas(params.policy.drivers, 'average_latency_drivers') then params.policy.drivers.average_latency_drivers else []),
    driverAccumulatorStep1
  ),
  local driverAccumulator = std.foldl(
    addPercentileLatencyDriver,
    (if std.objectHas(params.policy.drivers, 'percentile_latency_drivers') then params.policy.drivers.percentile_latency_drivers else []),
    driverAccumulatorStep2
  ),

  local userStartControlComponent = spec.v1.Component.withBoolVariable(
    spec.v1.BoolVariable.withConstantOutput(params.policy.start) +
    spec.v1.BoolVariable.withConfigKey('start') +
    spec.v1.BoolVariable.withOutPorts({
      output: spec.v1.Port.withSignalName('USER_START_CONTROL'),
    }),
  ),

  local userResetControlComponent = spec.v1.Component.withBoolVariable(
    spec.v1.BoolVariable.withConstantOutput(false) +
    spec.v1.BoolVariable.withConfigKey('reset') +
    spec.v1.BoolVariable.withOutPorts({
      output: spec.v1.Port.withSignalName('USER_RESET_CONTROL'),
    }),
  ),

  local alwaysForward =
    if std.length(driverAccumulator.forward_signals) > 0 then
      false
    else
      true,

  local forwardIntentComponent = spec.v1.Component.withOr(
    spec.v1.Or.withInPorts({
      inputs: [
        spec.v1.Port.withSignalName(signal)
        for signal in driverAccumulator.forward_signals
      ],
    })
    + spec.v1.Or.withOutPorts({
      output: spec.v1.Port.withSignalName('FORWARD_INTENT'),
    }),
  ),

  local backwardIntentComponent = spec.v1.Component.withOr(
    spec.v1.Or.withInPorts({
      inputs: [
        spec.v1.Port.withSignalName(signal)
        for signal in driverAccumulator.backward_signals
      ],
    })
    + spec.v1.Or.withOutPorts({
      output: spec.v1.Port.withSignalName('BACKWARD_INTENT'),
    }),
  ),

  local resetIntentComponent = spec.v1.Component.withOr(
    spec.v1.Or.withInPorts({
      inputs: [
        spec.v1.Port.withSignalName(signal)
        for signal in driverAccumulator.reset_signals
      ] + [
        spec.v1.Port.withSignalName('USER_RESET_CONTROL'),
      ],
    })
    + spec.v1.Or.withOutPorts({
      output: spec.v1.Port.withSignalName('RESET'),
    }),
  ),

  local createNotIntentComponents = function(inSignal, outSignal) [
    spec.v1.Component.withInverter(
      spec.v1.Inverter.withInPorts({
        input: spec.v1.Port.withSignalName(inSignal),
      })
      + spec.v1.Inverter.withOutPorts({
        output: spec.v1.Port.withSignalName('INVERTED_' + inSignal),
      }),
    ),
    spec.v1.Component.withFirstValid(
      spec.v1.FirstValid.withInPorts({
        inputs: [
          spec.v1.Port.withSignalName('INVERTED_' + inSignal),
          spec.v1.Port.withConstantSignal(1),
        ],
      })
      + spec.v1.FirstValid.withOutPorts({
        output: spec.v1.Port.withSignalName(outSignal),
      }),
    ),
  ],

  local notBackwardIntentComponents = createNotIntentComponents('BACKWARD_INTENT', 'NOT_BACKWARD'),
  local notResetComponents = createNotIntentComponents('RESET', 'NOT_RESET'),

  local forwardCriteriaComponent = spec.v1.Component.withAnd(
    spec.v1.And.withInPorts({
      inputs:
        [
          spec.v1.Port.withSignalName('NOT_BACKWARD'),
          spec.v1.Port.withSignalName('NOT_RESET'),
          spec.v1.Port.withSignalName('USER_START_CONTROL'),
        ]
        +
        (
          if alwaysForward then
            []
          else
            [spec.v1.Port.withSignalName('FORWARD_INTENT')]
        ),
    })
    + spec.v1.And.withOutPorts({
      output: spec.v1.Port.withSignalName('FORWARD'),
    }),
  ),

  local backwardCriteriaComponent = spec.v1.Component.withAnd(
    spec.v1.And.withInPorts({
      inputs: [
        spec.v1.Port.withSignalName('BACKWARD_INTENT'),
        spec.v1.Port.withSignalName('NOT_RESET'),
      ],
    })
    + spec.v1.And.withOutPorts({
      output: spec.v1.Port.withSignalName('BACKWARD'),
    }),
  ),

  local loadRamp = spec.v1.Component.withFlowControl(
    spec.v1.FlowControl.withLoadRamp(
      spec.v1.LoadRamp.withInPorts({
        forward: spec.v1.Port.withSignalName('FORWARD'),
        backward: spec.v1.Port.withSignalName('BACKWARD'),
        reset: spec.v1.Port.withSignalName('RESET'),
      })
      + spec.v1.LoadRamp.withParameters(params.policy.load_ramp)
      + spec.v1.LoadRamp.withPassThroughLabelValuesConfigKey('pass_through_label_values'),
    ),
  ),


  local policyDef =
    spec.v1.Policy.new()
    + spec.v1.Policy.withResources(
      utils.resources(params.policy.resources).updatedResources +
      spec.v1.Resources.new()
      + spec.v1.Resources.withFlowControl(
        spec.v1.FlowControlResources.new()
        + spec.v1.FlowControlResources.withFluxMeters((
          if std.objectHas(params.policy.resources.flow_control, 'flux_meters') then
            params.policy.resources.flow_control.flux_meters else {}
        ) + driverAccumulator.flux_meters)
      )
    )
    + spec.v1.Policy.withCircuit(
      spec.v1.Circuit.new()
      + spec.v1.Circuit.withEvaluationInterval(evaluation_interval=params.policy.evaluation_interval)
      + spec.v1.Circuit.withComponents(
        driverAccumulator.components + [
          userStartControlComponent,
          userResetControlComponent,
          backwardIntentComponent,
          resetIntentComponent,
        ] +
        (if alwaysForward
         then []
         else [forwardIntentComponent]) +
        notBackwardIntentComponents + notResetComponents + [
          forwardCriteriaComponent,
          backwardCriteriaComponent,
          loadRamp,
        ] + params.policy.components,
      ),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
