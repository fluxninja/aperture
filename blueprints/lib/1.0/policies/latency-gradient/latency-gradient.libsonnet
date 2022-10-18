local spec = import '../../spec.libsonnet';

local defaults = {
  policyName: error 'policyName must be set',
  evaluationInterval: '0.5s',
  fluxMeter: error 'fluxMeter must be set',
  concurrencyLimiterSelector: error 'concurrencyLimiterSelector must be set',
  classifiers: [],
  components: [],
  constants: {
    emaLimitMultiplier: '2.0',
    tolerance: '1.1',
    concurrencyLimitMultiplier: '2.0',
    minConcurrency: '10.0',
    linearConcurrencyIncrement: '5.0',
    concurrencyIncrementOverload: '10.0',
    sqrtScale: '0.1',
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
local classifier = spec.v1.Classifier;
local selector = spec.v1.Selector;
local component = spec.v1.Component;
local promQL = spec.v1.PromQL;
local port = spec.v1.Port;
local combinator = spec.v1.ArithmeticCombinator;
local ema = spec.v1.EMA;
local gradient = spec.v1.GradientController;
local concurrencyLimiter = spec.v1.ConcurrencyLimiter;
local scheduler = spec.v1.Scheduler;
local decider = spec.v1.Decider;
local switcher = spec.v1.Switcher;
local loadActuator = spec.v1.LoadActuator;
local max = spec.v1.Max;
local min = spec.v1.Min;
local sqrt = spec.v1.Sqrt;
local firstValid = spec.v1.FirstValid;

function(params) {
  _config:: defaults + params,

  local c = $._config.constants,

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFluxMetersMixin({ [$._config.policyName]: $._config.fluxMeter })
                           + resources.withClassifiers($._config.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval=$._config.evaluationInterval)
      + circuit.withComponents(
        [
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY'),
                                                            port.withConstantValue(c.emaLimitMultiplier),
                                                            output=port.withSignalName('MAX_EMA'))),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                            port.withConstantValue(c.tolerance),
                                                            output=port.withSignalName('LATENCY_SETPOINT'))),
          component.withArithmeticCombinator(combinator.div(port.withSignalName('DESIRED_CONCURRENCY'),
                                                            port.withSignalName('INCOMING_CONCURRENCY'),
                                                            output=port.withSignalName('LOAD_MULTIPLIER'))),
          component.withArithmeticCombinator(combinator.mul(port.withConstantValue(c.concurrencyLimitMultiplier),
                                                            port.withSignalName('ACCEPTED_CONCURRENCY'),
                                                            output=port.withSignalName('UPPER_CONCURRENCY_LIMIT'))),
          component.withArithmeticCombinator(combinator.add(port.withConstantValue(c.linearConcurrencyIncrement),
                                                            port.withSignalName('SQRT_CONCURRENCY_INCREMENT'),
                                                            output=port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK'))),
          component.withArithmeticCombinator(combinator.add(port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK'),
                                                            port.withSignalName('CONCURRENCY_INCREMENT_FEEDBACK'),
                                                            output=port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL'))),
          component.withMin(
            min.new()
            + min.withInPorts({ inputs: [port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL'), port.withSignalName('ACCEPTED_CONCURRENCY')] })
            + min.withOutPorts({ output: port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL_CAPPED') }),
          ),
          component.withFirstValid(
            firstValid.new()
            + firstValid.withInPorts({ inputs: [port.withSignalName('CONCURRENCY_INCREMENT_INTEGRAL_CAPPED'), port.withConstantValue(0)] })
            + firstValid.withOutPorts({ output: port.withSignalName('CONCURRENCY_INCREMENT_NORMAL') }),
          ),
          component.withMax(
            max.new()
            + max.withInPorts({ inputs: [port.withSignalName('UPPER_CONCURRENCY_LIMIT'), port.withConstantValue(c.minConcurrency)] })
            + max.withOutPorts({ output: port.withSignalName('MAX_CONCURRENCY') }),
          ),
          component.withSqrt(
            sqrt.new()
            + sqrt.withInPorts({ input: port.withSignalName('ACCEPTED_CONCURRENCY') })
            + sqrt.withOutPorts({ output: port.withSignalName('SQRT_CONCURRENCY_INCREMENT') })
            + sqrt.withScale($._config.constants.sqrtScale),
          ),
          component.withPromql(
            local q = 'sum(increase(flux_meter_sum{decision_type!="DECISION_TYPE_REJECTED", response_status="OK", flux_meter_name="%(policyName)s"}[5s]))/sum(increase(flux_meter_count{decision_type!="DECISION_TYPE_REJECTED", response_status="OK", flux_meter_name="%(policyName)s"}[5s]))' % { policyName: $._config.policyName };
            promQL.new()
            + promQL.withQueryString(q)
            + promQL.withEvaluationInterval('1s')
            + promQL.withOutPorts({ output: port.withSignalName('LATENCY') }),
          ),
          component.withEma(
            local e = $._config.ema;
            ema.withEmaWindow(e.window)
            + ema.withWarmUpWindow(e.warmUpWindow)
            + ema.withCorrectionFactorOnMaxEnvelopeViolation(e.correctionFactor)
            + ema.withInPortsMixin(
              ema.inPorts.withInput(port.withSignalName('LATENCY'))
              + ema.inPorts.withMaxEnvelope(port.withSignalName('MAX_EMA'))
            )
            + ema.withOutPortsMixin(ema.outPorts.withOutput(port.withSignalName('LATENCY_EMA')))
          ),
          component.withGradientController(
            local g = $._config.gradient;
            gradient.new()
            + gradient.withSlope(g.slope)
            + gradient.withMinGradient(g.minGradient)
            + gradient.withMaxGradient(g.maxGradient)
            + gradient.withInPorts({
              signal: port.withSignalName('LATENCY'),
              setpoint: port.withSignalName('LATENCY_SETPOINT'),
              max: port.withSignalName('MAX_CONCURRENCY'),
              control_variable: port.withSignalName('ACCEPTED_CONCURRENCY'),
              optimize: port.withSignalName('CONCURRENCY_INCREMENT'),
            })
            + gradient.withOutPortsMixin({
              output: port.withSignalName('DESIRED_CONCURRENCY'),
            })
          ),
          component.withConcurrencyLimiter(
            local c = $._config.concurrencyLimiter;
            concurrencyLimiter.new()
            + concurrencyLimiter.withSelector($._config.concurrencyLimiterSelector)
            + concurrencyLimiter.withScheduler(
              scheduler.new()
              + scheduler.withAutoTokens(c.autoTokens)
              + scheduler.withTimeoutFactor(c.timeoutFactor)
              + scheduler.withDefaultWorkloadParameters(c.defaultWorkloadParameters)
              + scheduler.withWorkloads(c.workloads)
              + scheduler.withOutPortsMixin({
                accepted_concurrency: port.withSignalName('ACCEPTED_CONCURRENCY'),
                incoming_concurrency: port.withSignalName('INCOMING_CONCURRENCY'),
              })
            )
            + concurrencyLimiter.withLoadActuator(
              loadActuator.withInPortsMixin({ load_multiplier: port.withSignalName('LOAD_MULTIPLIER') })
            )
          ),
          component.withDecider(
            decider.new()
            + decider.withOperator('gt')
            + decider.withInPortsMixin(
              decider.inPorts.withLhs(port.withSignalName('LATENCY'))
              + decider.inPorts.withRhs(port.withSignalName('LATENCY_SETPOINT'))
            )
            + decider.withOutPortsMixin(decider.outPorts.withOutput(port.withSignalName('IS_OVERLOAD_SWITCH')))
          ),
          component.withSwitcher(
            switcher.new()
            + switcher.withInPortsMixin(
              switcher.inPorts.withOnTrue(port.withConstantValue(c.concurrencyIncrementOverload))
              + switcher.inPorts.withOnFalse(port.withSignalName('CONCURRENCY_INCREMENT_NORMAL'))
              + switcher.inPorts.withSwitch(port.withSignalName('IS_OVERLOAD_SWITCH'))
            )
            + switcher.withOutPortsMixin(switcher.outPorts.withOutput(port.withSignalName('CONCURRENCY_INCREMENT')))
          ),
          component.withSwitcher(
            switcher.new()
            + switcher.withInPortsMixin(
              switcher.inPorts.withOnTrue(port.withConstantValue(0))
              + switcher.inPorts.withOnFalse(port.withSignalName('CONCURRENCY_INCREMENT_NORMAL'))
              + switcher.inPorts.withSwitch(port.withSignalName('IS_OVERLOAD_SWITCH'))
            )
            + switcher.withOutPortsMixin(switcher.outPorts.withOutput(port.withSignalName('CONCURRENCY_INCREMENT_FEEDBACK')))
          ),
        ] + $._config.components,
      ),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: $._config.policyName,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
}
