local spec = import '../../../spec.libsonnet';
local config = import './config.libsonnet';

function(cfg) {
  local params = config.common + config.policy + cfg,

  local policyName = params.policy_name,

  local flux_meters = params.flux_meters,

  local addPromQLDriver = function(driverAccumulator, driver) {
    local evaluationInterval = params.evaluation_interval,
    local promQLSignalName = 'PROMQL_' + std.toString(driverAccumulator.promql_driver_count),
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(driver.query_string)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluationInterval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(promQLSignalName),
                                                                                 }))),
    local isForward = (if std.objectHas(driver, 'forward') then driver.forward != null else false),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator(driver.forward.operator)
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(promQLSignalName),
                                                           rhs: spec.v1.Port.withConstantSignal(driver.forward.threshold),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(driver, 'backward') then driver.backward != null else false),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator(driver.backward.operator)
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(promQLSignalName),
                                                            rhs: spec.v1.Port.withConstantSignal(driver.backward.threshold),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(driver, 'reset') then driver.reset != null else false),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator(driver.reset.operator)
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(promQLSignalName),
                                                         rhs: spec.v1.Port.withConstantSignal(driver.reset.threshold),
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
    reset_signals_count: driverAccumulator.reset_signals_count + if driver.reset != null then 1 else 0,
    components: driverAccumulator.components + [promQLComponent] + (if isForward then [forwardDecider] else []) + (if isBackward then [backwardDecider] else []) + (if isReset then [resetDecider] else []),
    driver_count: driverAccumulator.driver_count + 1,
    promql_driver_count: driverAccumulator.promql_driver_count + 1,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count,
    ema_latency_driver_count: driverAccumulator.ema_latency_driver_count,
  },

  local addAverageLatencyDriver = function(driverAccumulator, driver) {
    local flux_meter_name = policyName + '/average_latency/' + std.toString(driverAccumulator.average_latency_driver_count),
    local evaluationInterval = params.evaluation_interval,
    local averageLatencySignalName = 'AVERAGE_LATENCY_' + std.toString(driverAccumulator.average_latency_driver_count),
    local q = 'sum(increase(flux_meter_sum{valid="true", flow_status="OK", flux_meter_name="%(flux_meter_name)s"}[5s]))/sum(increase(flux_meter_count{valid="true", flow_status="OK", flux_meter_name="%(flux_meter_name)s"}[5s]))' % { flux_meter_name: flux_meter_name },
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(q)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluationInterval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                                                 }))),

    local isForward = (if std.objectHas(driver, 'forward') then driver.forward != null else false),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('lt')
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                           rhs: spec.v1.Port.withConstantSignal(driver.forward.threshold),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(driver, 'backward') then driver.backward != null else false),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                            rhs: spec.v1.Port.withConstantSignal(driver.backward.threshold),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(driver, 'reset') then driver.reset != null else false),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(averageLatencySignalName),
                                                         rhs: spec.v1.Port.withConstantSignal(driver.reset.threshold),
                                                       })
                                                       + spec.v1.Decider.withOutPorts({
                                                         output: spec.v1.Port.withSignalName(resetSignal),
                                                       })),
    local newFluxMeters = {
      [flux_meter_name]: spec.v1.FluxMeter.withFlowSelector(driver.flow_selector),
    },

    flux_meters: driverAccumulator.flux_meters + newFluxMeters,
    forward_signals: driverAccumulator.forward_signals + (if isForward then [forwardSignal] else []),
    backward_signals: driverAccumulator.backward_signals + (if isBackward then [backwardSignal] else []),
    reset_signals: driverAccumulator.reset_signals + (if isReset then [resetSignal] else []),
    forward_signals_count: driverAccumulator.forward_signals_count + if isForward then 1 else 0,
    backward_signals_count: driverAccumulator.backward_signals_count + if isBackward then 1 else 0,
    reset_signals_count: driverAccumulator.reset_signals_count + if driver.reset != null then 1 else 0,
    components: driverAccumulator.components + [promQLComponent] + (if isForward then [forwardDecider] else []) + (if isBackward then [backwardDecider] else []) + (if isReset then [resetDecider] else []),
    driver_count: driverAccumulator.driver_count + 1,
    promql_driver_count: driverAccumulator.promql_driver_count,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count + 1,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count,
    ema_latency_driver_count: driverAccumulator.ema_latency_driver_count,
  },

  local addPercentileLatencyDriver = function(driverAccumulator, driver) {
    local flux_meter_name = policyName + '/percentile_latency/' + std.toString(driverAccumulator.percentile_latency_driver_count),
    local evaluationInterval = params.evaluation_interval,
    local percentileLatencySignalName = 'PERCENTILE_LATENCY_' + std.toString(driverAccumulator.percentile_latency_driver_count),
    local q = 'histogram_quantile(%(percentile)f, sum(rate(flux_meter_bucket{valid="true", flow_status="OK", flux_meter_name="%(flux_meter_name)s"}[5s])) by (le))' % { percentile: driver.percentile, flux_meter_name: flux_meter_name },
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(q)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluationInterval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                                                 }))),

    local isForward = (if std.objectHas(driver, 'forward') then driver.forward != null else false),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('lt')
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                           rhs: spec.v1.Port.withConstantSignal(driver.forward.threshold),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(driver, 'backward') then driver.backward != null else false),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                            rhs: spec.v1.Port.withConstantSignal(driver.backward.threshold),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(driver, 'reset') then driver.reset != null else false),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(percentileLatencySignalName),
                                                         rhs: spec.v1.Port.withConstantSignal(driver.reset.threshold),
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
    reset_signals_count: driverAccumulator.reset_signals_count + if driver.reset != null then 1 else 0,
    components: driverAccumulator.components + [promQLComponent] + (if isForward then [forwardDecider] else []) + (if isBackward then [backwardDecider] else []) + (if isReset then [resetDecider] else []),
    driver_count: driverAccumulator.driver_count + 1,
    promql_driver_count: driverAccumulator.promql_driver_count,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count + 1,
    ema_latency_driver_count: driverAccumulator.ema_latency_driver_count,
  },

  local addEMALatencyDriver = function(driverAccumulator, driver) {
    local flux_meter_name = policyName + '/ema_latency/' + std.toString(driverAccumulator.ema_latency_driver_count),
    local evaluationInterval = params.evaluation_interval,
    local latencySignalName = 'LATENCY_' + std.toString(driverAccumulator.ema_latency_driver_count),
    local q = 'sum(rate(flux_meter_sum{valid="true", flow_status="OK", flux_meter_name="%(flux_meter_name)s"}[5s]))/sum(rate(flux_meter_count{valid="true", flow_status="OK", flux_meter_name="%(flux_meter_name)s"}[5s]))' % { flux_meter_name: flux_meter_name },
    local promQLComponent = spec.v1.Component.withQuery(spec.v1.Query.withPromql(spec.v1.PromQL.withQueryString(q)
                                                                                 + spec.v1.PromQL.withEvaluationInterval(evaluationInterval)
                                                                                 + spec.v1.PromQL.withOutPorts({
                                                                                   output: spec.v1.Port.withSignalName(latencySignalName),
                                                                                 }))),

    local emaLatencySignalName = 'EMA_LATENCY_' + std.toString(driverAccumulator.ema_latency_driver_count),
    local emaComponent = spec.v1.Component.withEma(
      spec.v1.EMA.withParameters(driver.ema)
      + spec.v1.EMA.withInPortsMixin(
        spec.v1.EMA.inPorts.withInput(spec.v1.Port.withSignalName(latencySignalName))
      )
      + spec.v1.EMA.withOutPortsMixin(
        spec.v1.EMA.outPorts.withOutput(spec.v1.Port.withSignalName(emaLatencySignalName))
      )
    ),

    local isForward = (if std.objectHas(driver, 'forward') then driver.forward != null else false),
    local forwardLatencySetpoint = 'FORWARD_LATENCY_SETPOINT_' + std.toString(driverAccumulator.ema_latency_driver_count),
    local forwardSignal = 'FORWARD_' + std.toString(driverAccumulator.forward_signals_count),
    local forwardMultiplier = spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.withOperator('mul')
                                                                         + spec.v1.ArithmeticCombinator.withInPorts({
                                                                           lhs: spec.v1.Port.withSignalName(emaLatencySignalName),
                                                                           rhs: spec.v1.Port.withConstantSignal(driver.forward.latency_tolerance_multiplier),
                                                                         })
                                                                         + spec.v1.ArithmeticCombinator.withOutPorts({
                                                                           output: spec.v1.Port.withSignalName(forwardLatencySetpoint),
                                                                         })),
    local forwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('lt')
                                                         + spec.v1.Decider.withInPorts({
                                                           lhs: spec.v1.Port.withSignalName(latencySignalName),
                                                           rhs: spec.v1.Port.withSignalName(forwardLatencySetpoint),
                                                         })
                                                         + spec.v1.Decider.withOutPorts({
                                                           output: spec.v1.Port.withSignalName(forwardSignal),
                                                         })),
    local isBackward = (if std.objectHas(driver, 'backward') then driver.backward != null else false),
    local backwardLatencySetpoint = 'BACKWARD_LATENCY_SETPOINT_' + std.toString(driverAccumulator.ema_latency_driver_count),
    local backwardSignal = 'BACKWARD_' + std.toString(driverAccumulator.backward_signals_count),
    local backwardMultiplier = spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.withOperator('mul')
                                                                          + spec.v1.ArithmeticCombinator.withInPorts({
                                                                            lhs: spec.v1.Port.withSignalName(emaLatencySignalName),
                                                                            rhs: spec.v1.Port.withConstantSignal(driver.backward.latency_tolerance_multiplier),
                                                                          })
                                                                          + spec.v1.ArithmeticCombinator.withOutPorts({
                                                                            output: spec.v1.Port.withSignalName(backwardLatencySetpoint),
                                                                          })),
    local backwardDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                          + spec.v1.Decider.withInPorts({
                                                            lhs: spec.v1.Port.withSignalName(latencySignalName),
                                                            rhs: spec.v1.Port.withSignalName(backwardLatencySetpoint),
                                                          })
                                                          + spec.v1.Decider.withOutPorts({
                                                            output: spec.v1.Port.withSignalName(backwardSignal),
                                                          })),
    local isReset = (if std.objectHas(driver, 'reset') then driver.reset != null else false),
    local resetLatencySetpoint = 'RESET_LATENCY_SETPOINT_' + std.toString(driverAccumulator.ema_latency_driver_count),
    local resetSignal = 'RESET_' + std.toString(driverAccumulator.reset_signals_count),
    local resetMultiplier = spec.v1.Component.withArithmeticCombinator(spec.v1.ArithmeticCombinator.withOperator('mul')
                                                                       + spec.v1.ArithmeticCombinator.withInPorts({
                                                                         lhs: spec.v1.Port.withSignalName(emaLatencySignalName),
                                                                         rhs: spec.v1.Port.withConstantSignal(driver.reset.latency_tolerance_multiplier),
                                                                       })
                                                                       + spec.v1.ArithmeticCombinator.withOutPorts({
                                                                         output: spec.v1.Port.withSignalName(resetLatencySetpoint),
                                                                       })),
    local resetDecider = spec.v1.Component.withDecider(spec.v1.Decider.withOperator('gt')
                                                       + spec.v1.Decider.withInPorts({
                                                         lhs: spec.v1.Port.withSignalName(latencySignalName),
                                                         rhs: spec.v1.Port.withSignalName(resetLatencySetpoint),
                                                       })
                                                       + spec.v1.Decider.withOutPorts({
                                                         output: spec.v1.Port.withSignalName(resetSignal),
                                                       })),

    local newFluxMeters = {
      [flux_meter_name]: spec.v1.FluxMeter.withFlowSelector(driver.flow_selector),
    },

    flux_meters: driverAccumulator.flux_meters + newFluxMeters,
    forward_signals: driverAccumulator.forward_signals + (if isForward then [forwardSignal] else []),
    backward_signals: driverAccumulator.backward_signals + (if isBackward then [backwardSignal] else []),
    reset_signals: driverAccumulator.reset_signals + (if isReset then [resetSignal] else []),
    forward_signals_count: driverAccumulator.forward_signals_count + (if isForward then 1 else 0),
    backward_signals_count: driverAccumulator.backward_signals_count + (if isBackward then 1 else 0),
    reset_signals_count: driverAccumulator.reset_signals_count + (if isReset then 1 else 0),
    components: driverAccumulator.components + [promQLComponent, emaComponent] + (if isForward then [forwardMultiplier, forwardDecider] else []) + (if isBackward then [backwardMultiplier, backwardDecider] else []) + (if isReset then [resetMultiplier, resetDecider] else []),
    driver_count: driverAccumulator.driver_count + 1,
    promql_driver_count: driverAccumulator.promql_driver_count,
    average_latency_driver_count: driverAccumulator.average_latency_driver_count,
    percentile_latency_driver_count: driverAccumulator.percentile_latency_driver_count,
    ema_latency_driver_count: driverAccumulator.ema_latency_driver_count + 1,
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
    driver_count: 0,
    promql_driver_count: 0,
    average_latency_driver_count: 0,
    percentile_latency_driver_count: 0,
    ema_latency_driver_count: 0,
  },


  local driverAccumulatorStep1 = std.foldl(
    addPromQLDriver,
    (if std.objectHas(params.drivers, 'promql_drivers') then params.drivers.promql_drivers else []),
    driverAccumulatorInitial
  ),
  local driverAccumulatorStep2 = std.foldl(
    addAverageLatencyDriver,
    (if std.objectHas(params.drivers, 'average_latency_drivers') then params.drivers.average_latency_drivers else []),
    driverAccumulatorStep1
  ),
  local driverAccumulatorStep3 = std.foldl(
    addPercentileLatencyDriver,
    (if std.objectHas(params.drivers, 'percentile_latency_drivers') then params.drivers.percentile_latency_drivers else []),
    driverAccumulatorStep2
  ),
  local driverAccumulator = std.foldl(
    addEMALatencyDriver,
    (if std.objectHas(params.drivers, 'ema_latency_drivers') then params.drivers.ema_latency_drivers else []),
    driverAccumulatorStep3
  ),

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
      inputs: [
        spec.v1.Port.withSignalName('FORWARD_INTENT'),
        spec.v1.Port.withSignalName('NOT_BACKWARD'),
        spec.v1.Port.withSignalName('NOT_RESET'),
      ],
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

  local loadShaper = spec.v1.Component.withFlowControl(
    spec.v1.FlowControl.withLoadShaper(
      spec.v1.LoadShaper.withInPorts({
        forward: spec.v1.Port.withSignalName('FORWARD'),
        backward: spec.v1.Port.withSignalName('BACKWARD'),
        reset: spec.v1.Port.withSignalName('RESET'),
      })
      + spec.v1.LoadShaper.withParameters(params.load_shaper)
      + spec.v1.LoadShaper.withDynamicConfigKey('load_shaper'),
    ),
  ),


  local policyDef =
    spec.v1.Policy.new()
    + spec.v1.Policy.withResources(spec.v1.Resources.new()
                                   + spec.v1.Resources.withFlowControl(
                                     spec.v1.FlowControlResources.new()
                                     + spec.v1.FlowControlResources.withFluxMeters((if std.objectHas(params.resources.flow_control, 'flux_meters') then params.resources.flow_control.flux_meters else {}) + driverAccumulator.flux_meters)
                                     + spec.v1.FlowControlResources.withClassifiers(params.resources.flow_control.classifiers)
                                   ))
    + spec.v1.Policy.withCircuit(
      spec.v1.Circuit.new()
      + spec.v1.Circuit.withEvaluationInterval(evaluation_interval=params.evaluation_interval)
      + spec.v1.Circuit.withComponents(
        driverAccumulator.components + [
          forwardIntentComponent,
          backwardIntentComponent,
          resetIntentComponent,
        ] + notBackwardIntentComponents + notResetComponents + [
          forwardCriteriaComponent,
          backwardCriteriaComponent,
          loadShaper,
        ] + params.components,
      ),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
