local spec = import '../../spec.libsonnet';

local defaults = {
  policyName: error 'policyName must be set',
  evaluationInterval: '0.5s',
  fluxMeterSelector: error 'fluxMeterSelector must be set',
  fluxMeters: {},
  concurrencyLimiterSelector: error 'concurrencyLimiterSelector must be set',
  classifiers: [],
  constants: {
    emaLimitMultiplier: '2.0',
    tolerance: '1.1',
    concurrencyLimitMultiplier: '2.0',
    minConcurrency: '10.0',
    linearConcurrencyIncrement: '5.0',
    concurrencyIncrementOverload: '10.0',
    sqrtScale: '0.5',
  },
  ema: {
    window: '1500s',
    warmUpWindow: '10s',
    correctionFactor: '0.95',
  },
  gradient: {
    slope: '-1',
    minGradient: '0.1',
    maxGradient: '1.0',
  },
  concurrencyLimiter: {
    autoTokens: true,
    timeoutFactor: '0.5',
    defaultWorkloadParameters: {
      priority: 20,
    },
    workloads: [],
  },
};


local policy = spec.v1.Policy;
local resources = spec.v1.Resources;
local circuit = spec.v1.Circuit;
local fluxMeter = spec.v1.FluxMeter;
local classifier = spec.v1.Classifier;
local selector = spec.v1.Selector;
local component = spec.v1.Component;
local promQL = spec.v1.PromQL;
local port = spec.v1.Port;
local constant = spec.v1.Constant;
local combinator = spec.v1.ArithmeticCombinator;
local ema = spec.v1.EMA;
local gradient = spec.v1.GradientController;
local limiter = spec.v1.ConcurrencyLimiter;
local scheduler = spec.v1.Scheduler;
local decider = spec.v1.Decider;
local switcher = spec.v1.Switcher;
local loadShed = spec.v1.LoadShedActuator;
local max = spec.v1.Max;
local min = spec.v1.Min;
local sqrt = spec.v1.Sqrt;

local latencyPort = port.new() + port.withSignalName('LATENCY');

// constant ports
local emaLimitMultiplierPort = port.new() + port.withSignalName('EMA_LIMIT_MULTIPLIER');
local tolerancePort = port.new() + port.withSignalName('TOLERANCE');
local concurrencyLimitMultiplierPort = port.new() + port.withSignalName('CONCURRENCY_LIMIT_MULTIPLIER');
local minConcurrencyPort = port.new() + port.withSignalName('MIN_CONCURRENCY');
local linearConcurrencyIncrementPort = port.new() + port.withSignalName('LINEAR_CONCURRENCY_INCREMENT');
local sqrtConcurrencyIncrementPort = port.new() + port.withSignalName('SQRT_CONCURRENCY_INCREMENT');
local zeroPort = port.new() + port.withSignalName('ZERO');

local concurrencyIncrementOverloadPort = port.new() + port.withSignalName('CONCURRENCY_INCREMENT_OVERLOAD');
local concurrencyIncrementSingleTickPort = port.new() + port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK');
local concurrencyIncrementFeedbackPort = port.new() + port.withSignalName('CONCURRENCY_INCREMENT_FEEDBACK');
local isOverloadSwitchPort = port.new() + port.withSignalName('IS_OVERLOAD_SWITCH');
local concurrencyIncrementIntegralPort = port.new() + port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL');
local concurrencyIncrementNormalPort = port.new() + port.withSignalName('CONCURRENCY_INCREMENT_NORMAL');

local maxEmaPort = port.new() + port.withSignalName('MAX_EMA');
local latencySetpointPort = port.new() + port.withSignalName('LATENCY_SETPOINT');
local latencyEmaPort = port.new() + port.withSignalName('LATENCY_EMA');
local LSFPort = port.new() + port.withSignalName('LSF');
local upperConcurrencyLimitPort = port.new() + port.withSignalName('UPPER_CONCURRENCY_LIMIT');
local latencyOverloadPort = port.new() + port.withSignalName('LATENCY_OVERLOAD');

local desiredConcurrencyPort = port.new() + port.withSignalName('DESIRED_CONCURRENCY');
local maxConcurrencyPort = port.new() + port.withSignalName('MAX_CONCURRENCY');
local acceptedConcurrencyPort = port.new() + port.withSignalName('ACCEPTED_CONCURRENCY');
local concurrencyIncrementPort = port.new() + port.withSignalName('CONCURRENCY_INCREMENT');
local incomingConcurrencyPort = port.new() + port.withSignalName('INCOMING_CONCURRENCY');
local deltaConcurrencyPort = port.new() + port.withSignalName('DELTA_CONCURRENCY');


function(params) {
  _config:: defaults + params,

  local c = $._config.constants,

  local constants = [
    component.withConstant(constant.new() + constant.withValue(c.emaLimitMultiplier) + constant.withOutPorts({ output: emaLimitMultiplierPort })),
    component.withConstant(constant.new() + constant.withValue(c.concurrencyLimitMultiplier) + constant.withOutPorts({ output: concurrencyLimitMultiplierPort })),
    component.withConstant(constant.new() + constant.withValue(c.minConcurrency) + constant.withOutPorts({ output: minConcurrencyPort })),
    component.withConstant(constant.new() + constant.withValue(c.linearConcurrencyIncrement) + constant.withOutPorts({ output: linearConcurrencyIncrementPort })),
    component.withConstant(constant.new() + constant.withValue(c.concurrencyIncrementOverload) + constant.withOutPorts({ output: concurrencyIncrementOverloadPort })),
    component.withConstant(constant.new() + constant.withValue(c.tolerance) + constant.withOutPorts({ output: tolerancePort })),
    component.withConstant(constant.new() + constant.withValue(0) + constant.withOutPorts({ output: zeroPort })),
  ],

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFluxMeters($._config.fluxMeters)
                           + resources.withClassifiers($._config.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval=$._config.evaluationInterval)
      + circuit.withComponents(constants + [
        component.withArithmeticCombinator(combinator.mul(latencyPort,
                                                          emaLimitMultiplierPort,
                                                          output=maxEmaPort)),
        component.withArithmeticCombinator(combinator.mul(latencyEmaPort,
                                                          tolerancePort,
                                                          output=latencySetpointPort)),
        component.withArithmeticCombinator(combinator.sub(incomingConcurrencyPort,
                                                          desiredConcurrencyPort,
                                                          output=deltaConcurrencyPort)),
        component.withArithmeticCombinator(combinator.div(deltaConcurrencyPort,
                                                          incomingConcurrencyPort,
                                                          output=LSFPort)),
        component.withArithmeticCombinator(combinator.mul(concurrencyLimitMultiplierPort,
                                                          acceptedConcurrencyPort,
                                                          output=upperConcurrencyLimitPort)),
        component.withArithmeticCombinator(combinator.mul(latencyEmaPort,
                                                          tolerancePort,
                                                          output=latencyOverloadPort)),
        component.withArithmeticCombinator(combinator.add(linearConcurrencyIncrementPort,
                                                          sqrtConcurrencyIncrementPort,
                                                          output=concurrencyIncrementSingleTickPort)),
        component.withArithmeticCombinator(combinator.add(concurrencyIncrementSingleTickPort,
                                                          concurrencyIncrementFeedbackPort,
                                                          output=concurrencyIncrementIntegralPort)),
        component.withMin(
          min.new()
          + min.withInPorts(min.inPorts.withInputs([concurrencyIncrementIntegralPort, acceptedConcurrencyPort]))
          + min.withOutPorts(min.outPorts.withOutput(concurrencyIncrementNormalPort)),
        ),
        component.withMax(
          max.new()
          + max.withInPorts(max.inPorts.withInputs([upperConcurrencyLimitPort, minConcurrencyPort]))
          + max.withOutPorts(max.outPorts.withOutput(maxConcurrencyPort)),
        ),
        component.withSqrt(
          sqrt.new()
          + sqrt.withInPorts({ input: acceptedConcurrencyPort })
          + sqrt.withOutPorts({ output: sqrtConcurrencyIncrementPort })
          + sqrt.withScale($._config.constants.sqrtScale),
        ),
        component.withPromql(
          local q = 'sum(increase(flux_meter_sum{decision_type!="DECISION_TYPE_REJECTED", flux_meter_name="%(policyName)s"}[5s]))/sum(increase(flux_meter_count{decision_type!="DECISION_TYPE_REJECTED", flux_meter_name="%(policyName)s"}[5s]))' % { policyName: $._config.policyName };
          promQL.new()
          + promQL.withQueryString(q)
          + promQL.withEvaluationInterval('1s')
          + promQL.withOutPorts({ output: latencyPort }),
        ),
        component.withEma(
          local e = $._config.ema;
          ema.withEmaWindow(e.window)
          + ema.withWarmUpWindow(e.warmUpWindow)
          + ema.withCorrectionFactorOnMaxEnvelopeViolation(e.correctionFactor)
          + ema.withInPortsMixin(
            ema.inPorts.withInput(latencyPort)
            + ema.inPorts.withMaxEnvelope(maxEmaPort)
          )
          + ema.withOutPortsMixin(ema.outPorts.withOutput(latencyEmaPort))
        ),
        component.withGradientController(
          local g = $._config.gradient;
          gradient.new()
          + gradient.withSlope(g.slope)
          + gradient.withMinGradient(g.minGradient)
          + gradient.withMaxGradient(g.maxGradient)
          + gradient.withInPorts({
            signal: latencyPort,
            setpoint: latencySetpointPort,
            max: maxConcurrencyPort,
            control_variable: acceptedConcurrencyPort,
            optimize: concurrencyIncrementPort,
          })
          + gradient.withOutPortsMixin({
            output: desiredConcurrencyPort,
          })
        ),
        component.withConcurrencyLimiter(
          local c = $._config.concurrencyLimiter;
          limiter.new()
          + limiter.withScheduler(
            scheduler.new()
            + scheduler.withSelector($._config.concurrencyLimiterSelector)
            + scheduler.withAutoTokens(c.autoTokens)
            + scheduler.withTimeoutFactor(c.timeoutFactor)
            + scheduler.withDefaultWorkloadParameters(c.defaultWorkloadParameters)
            + scheduler.withWorkloads(c.workloads)
            + scheduler.withOutPortsMixin({
              accepted_concurrency: acceptedConcurrencyPort,
              incoming_concurrency: incomingConcurrencyPort,
            })
          )
          + limiter.withLoadShedActuator(
            loadShed.withInPortsMixin({ load_shed_factor: LSFPort })
          )
        ),
        component.withDecider(
          decider.new()
          + decider.withOperator('gt')
          + decider.withInPortsMixin(
            decider.inPorts.withLhs(latencyPort)
            + decider.inPorts.withRhs(latencyOverloadPort)
          )
          + decider.withOutPortsMixin(decider.outPorts.withOutput(isOverloadSwitchPort))
        ),
        component.withSwitcher(
          switcher.new()
          + switcher.withInPortsMixin(
            switcher.inPorts.withOnTrue(concurrencyIncrementOverloadPort)
            + switcher.inPorts.withOnFalse(concurrencyIncrementNormalPort)
            + switcher.inPorts.withSwitch(isOverloadSwitchPort)
          )
          + switcher.withOutPortsMixin(switcher.outPorts.withOutput(concurrencyIncrementPort))
        ),
        component.withSwitcher(
          switcher.new()
          + switcher.withInPortsMixin(
            switcher.inPorts.withOnTrue(zeroPort)
            + switcher.inPorts.withOnFalse(concurrencyIncrementNormalPort)
            + switcher.inPorts.withSwitch(isOverloadSwitchPort)
          )
          + switcher.withOutPortsMixin(switcher.outPorts.withOutput(concurrencyIncrementFeedbackPort))
        ),
      ]),
    ),
  policy: policyDef,
}
