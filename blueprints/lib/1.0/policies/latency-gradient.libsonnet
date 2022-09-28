local aperture = import '../../../libsonnet/1.0/main.libsonnet';


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


local policy = aperture.v1.Policy;
local resources = aperture.v1.Resources;
local circuit = aperture.v1.Circuit;
local fluxMeter = aperture.v1.FluxMeter;
local classifier = aperture.v1.Classifier;
local selector = aperture.v1.Selector;
local component = aperture.v1.Component;
local promQL = aperture.v1.PromQL;
local inPort = aperture.v1.InPort;
local outPort = aperture.v1.OutPort;
local combinator = aperture.v1.ArithmeticCombinator;
local ema = aperture.v1.EMA;
local gradient = aperture.v1.GradientController;
local limiter = aperture.v1.ConcurrencyLimiter;
local scheduler = aperture.v1.Scheduler;
local decider = aperture.v1.Decider;
local switcher = aperture.v1.Switcher;
local loadShed = aperture.v1.LoadShedActuator;
local max = aperture.v1.Max;
local min = aperture.v1.Min;
local sqrt = aperture.v1.Sqrt;


// constant ports
local emaLimitMultiplierPort = inPort.new() + inPort.withConstantValue(defaults.constants.emaLimitMultiplier);
local tolerancePort = inPort.new() + inPort.withConstantValue(defaults.constants.tolerance);
local concurrencyLimitMultiplierPort = inPort.new() + inPort.withConstantValue(defaults.constants.concurrencyLimitMultiplier);
local minConcurrencyPort = inPort.new() + inPort.withConstantValue(defaults.constants.minConcurrency);
local linearConcurrencyIncrementPort = inPort.new() + inPort.withConstantValue(defaults.constants.linearConcurrencyIncrement);
local concurrencyIncrementOverloadPort = inPort.new() + inPort.withConstantValue(defaults.constants.concurrencyIncrementOverload);
local zeroPort = inPort.new() + inPort.withConstantValue(0);

// signal input ports
local latencyPort = inPort.new() + inPort.withSignalName('LATENCY');
local sqrtConcurrencyIncrementPort = inPort.new() + inPort.withSignalName('SQRT_CONCURRENCY_INCREMENT');
local concurrencyIncrementSingleTickPort = inPort.new() + inPort.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK');
local concurrencyIncrementFeedbackPort = inPort.new() + inPort.withSignalName('CONCURRENCY_INCREMENT_FEEDBACK');
local isOverloadSwitchPort = inPort.new() + inPort.withSignalName('IS_OVERLOAD_SWITCH');
local concurrencyIncrementIntegralPort = inPort.new() + inPort.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL');
local concurrencyIncrementNormalPort = inPort.new() + inPort.withSignalName('CONCURRENCY_INCREMENT_NORMAL');

local maxEmaPort = inPort.new() + inPort.withSignalName('MAX_EMA');
local latencySetpointPort = inPort.new() + inPort.withSignalName('LATENCY_SETPOINT');
local latencyEmaPort = inPort.new() + inPort.withSignalName('LATENCY_EMA');
local LSFPort = inPort.new() + inPort.withSignalName('LSF');
local upperConcurrencyLimitPort = inPort.new() + inPort.withSignalName('UPPER_CONCURRENCY_LIMIT');
local latencyOverloadPort = inPort.new() + inPort.withSignalName('LATENCY_OVERLOAD');

local desiredConcurrencyPort = inPort.new() + inPort.withSignalName('DESIRED_CONCURRENCY');
local maxConcurrencyPort = inPort.new() + inPort.withSignalName('MAX_CONCURRENCY');
local acceptedConcurrencyPort = inPort.new() + inPort.withSignalName('ACCEPTED_CONCURRENCY');
local concurrencyIncrementPort = inPort.new() + inPort.withSignalName('CONCURRENCY_INCREMENT');
local incomingConcurrencyPort = inPort.new() + inPort.withSignalName('INCOMING_CONCURRENCY');
local deltaConcurrencyPort = inPort.new() + inPort.withSignalName('DELTA_CONCURRENCY');

// signal output ports
local latencyPortOut = outPort.new() + outPort.withSignalName('LATENCY');
local sqrtConcurrencyIncrementPortOut = outPort.new() + outPort.withSignalName('SQRT_CONCURRENCY_INCREMENT');
local concurrencyIncrementSingleTickPortOut = outPort.new() + outPort.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK');
local concurrencyIncrementFeedbackPortOut = outPort.new() + outPort.withSignalName('CONCURRENCY_INCREMENT_FEEDBACK');
local isOverloadSwitchPortOut = outPort.new() + outPort.withSignalName('IS_OVERLOAD_SWITCH');
local concurrencyIncrementIntegralPortOut = outPort.new() + outPort.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL');
local concurrencyIncrementNormalPortOut = outPort.new() + outPort.withSignalName('CONCURRENCY_INCREMENT_NORMAL');

local maxEmaPortOut = outPort.new() + outPort.withSignalName('MAX_EMA');
local latencySetpointPortOut = outPort.new() + outPort.withSignalName('LATENCY_SETPOINT');
local latencyEmaPortOut = outPort.new() + outPort.withSignalName('LATENCY_EMA');
local LSFPortOut = outPort.new() + outPort.withSignalName('LSF');
local upperConcurrencyLimitPortOut = outPort.new() + outPort.withSignalName('UPPER_CONCURRENCY_LIMIT');
local latencyOverloadPortOut = outPort.new() + outPort.withSignalName('LATENCY_OVERLOAD');

local desiredConcurrencyPortOut = outPort.new() + outPort.withSignalName('DESIRED_CONCURRENCY');
local maxConcurrencyPortOut = outPort.new() + outPort.withSignalName('MAX_CONCURRENCY');
local acceptedConcurrencyPortOut = outPort.new() + outPort.withSignalName('ACCEPTED_CONCURRENCY');
local concurrencyIncrementPortOut = outPort.new() + outPort.withSignalName('CONCURRENCY_INCREMENT');
local incomingConcurrencyPortOut = outPort.new() + outPort.withSignalName('INCOMING_CONCURRENCY');
local deltaConcurrencyPortOut = outPort.new() + outPort.withSignalName('DELTA_CONCURRENCY');

function(params) {
  _config:: defaults + params,

  local c = $._config.constants,

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFluxMeters($._config.fluxMeters)
                           + resources.withClassifiers($._config.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval=$._config.evaluationInterval)
      + circuit.withComponents([
        component.withArithmeticCombinator(combinator.mul(latencyPort,
                                                          emaLimitMultiplierPort,
                                                          output=maxEmaPortOut)),
        component.withArithmeticCombinator(combinator.mul(latencyEmaPort,
                                                          tolerancePort,
                                                          output=latencySetpointPortOut)),
        component.withArithmeticCombinator(combinator.sub(incomingConcurrencyPort,
                                                          desiredConcurrencyPort,
                                                          output=deltaConcurrencyPortOut)),
        component.withArithmeticCombinator(combinator.div(deltaConcurrencyPort,
                                                          incomingConcurrencyPort,
                                                          output=LSFPortOut)),
        component.withArithmeticCombinator(combinator.mul(concurrencyLimitMultiplierPort,
                                                          acceptedConcurrencyPort,
                                                          output=upperConcurrencyLimitPortOut)),
        component.withArithmeticCombinator(combinator.mul(latencyEmaPort,
                                                          tolerancePort,
                                                          output=latencyOverloadPortOut)),
        component.withArithmeticCombinator(combinator.add(linearConcurrencyIncrementPort,
                                                          sqrtConcurrencyIncrementPort,
                                                          output=concurrencyIncrementSingleTickPortOut)),
        component.withArithmeticCombinator(combinator.add(concurrencyIncrementSingleTickPort,
                                                          concurrencyIncrementFeedbackPort,
                                                          output=concurrencyIncrementIntegralPortOut)),
        component.withMin(
          min.new()
          + min.withInPorts(min.inPorts.withInputs([concurrencyIncrementIntegralPort, acceptedConcurrencyPort]))
          + min.withOutPorts(min.outPorts.withOutput(concurrencyIncrementNormalPortOut)),
        ),
        component.withMax(
          max.new()
          + max.withInPorts(max.inPorts.withInputs([upperConcurrencyLimitPort, minConcurrencyPort]))
          + max.withOutPorts(max.outPorts.withOutput(maxConcurrencyPortOut)),
        ),
        component.withSqrt(
          sqrt.new()
          + sqrt.withInPorts({ input: acceptedConcurrencyPort })
          + sqrt.withOutPorts({ output: sqrtConcurrencyIncrementPortOut })
          + sqrt.withScale($._config.constants.sqrtScale),
        ),
        component.withPromql(
          local q = 'sum(increase(flux_meter_sum{decision_type!="DECISION_TYPE_REJECTED", flux_meter_name="%(policyName)s"}[5s]))/sum(increase(flux_meter_count{decision_type!="DECISION_TYPE_REJECTED", flux_meter_name="%(policyName)s"}[5s]))' % { policyName: $._config.policyName };
          promQL.new()
          + promQL.withQueryString(q)
          + promQL.withEvaluationInterval('1s')
          + promQL.withOutPorts({ output: latencyPortOut }),
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
          + ema.withOutPortsMixin(ema.outPorts.withOutput(latencyEmaPortOut))
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
            output: desiredConcurrencyPortOut,
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
              accepted_concurrency: acceptedConcurrencyPortOut,
              incoming_concurrency: incomingConcurrencyPortOut,
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
          + switcher.withOutPortsMixin(switcher.outPorts.withOutput(concurrencyIncrementPortOut))
        ),
        component.withSwitcher(
          switcher.new()
          + switcher.withInPortsMixin(
            switcher.inPorts.withOnTrue(zeroPort)
            + switcher.inPorts.withOnFalse(concurrencyIncrementNormalPort)
            + switcher.inPorts.withSwitch(isOverloadSwitchPort)
          )
          + switcher.withOutPortsMixin(switcher.outPorts.withOutput(concurrencyIncrementFeedbackPortOut))
        ),
      ]),
    ),
  policy: policyDef,
}
