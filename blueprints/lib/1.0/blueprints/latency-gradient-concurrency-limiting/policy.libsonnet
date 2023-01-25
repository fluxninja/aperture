local spec = import '../../spec.libsonnet';
local config = import './config.libsonnet';

local policy = spec.v1.Policy;
local resources = spec.v1.Resources;
local circuit = spec.v1.Circuit;
local classifier = spec.v1.Classifier;
local flowSelector = spec.v1.FlowSelector;
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
local loadActuatorDynamicConfig = spec.v1.LoadActuatorDynamicConfig;
local alerterConfig = spec.v1.AlerterConfig;
local max = spec.v1.Max;
local min = spec.v1.Min;
local sqrt = spec.v1.Sqrt;
local firstValid = spec.v1.FirstValid;
local extrapolator = spec.v1.Extrapolator;
local integrator = spec.v1.Integrator;
local constantSignal = spec.v1.ConstantSignal;

function(params) {
  _config:: config.common + config.policy + params,

  local c = $._config.constants,

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFluxMetersMixin({ [$._config.policyName]: $._config.fluxMeter })
                           + resources.withClassifiers($._config.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval='0.5s')
      + circuit.withComponents(
        [
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY'),
                                                            port.withConstantSignal(c.latencyEMALimitMultiplier),
                                                            output=port.withSignalName('MAX_EMA'))),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                            port.withConstantSignal(c.latencyToleranceMultiplier),
                                                            output=port.withSignalName('LATENCY_SETPOINT'))),
          component.withArithmeticCombinator(combinator.div(port.withSignalName('DESIRED_CONCURRENCY'),
                                                            port.withSignalName('INCOMING_CONCURRENCY'),
                                                            output=port.withSignalName('DESIRED_CONCURRENCY_RATIO'))),
          component.withArithmeticCombinator(combinator.mul(port.withConstantSignal(c.concurrencyLimitMultiplier),
                                                            port.withSignalName('ACCEPTED_CONCURRENCY'),
                                                            output=port.withSignalName('NORMAL_CONCURRENCY_LIMIT'))),
          component.withArithmeticCombinator(combinator.add(port.withConstantSignal(c.concurrencyLinearIncrement),
                                                            port.withSignalName('SQRT_CONCURRENCY_INCREMENT'),
                                                            output=port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK'))),
          component.withIntegrator(
            integrator.new()
            + integrator.withInPorts({
              input: port.withSignalName('CONCURRENCY_INCREMENT_SINGLE_TICK'),
              max: port.withSignalName('NORMAL_CONCURRENCY_LIMIT'),
              reset: port.withSignalName('IS_OVERLOAD'),
            })
            + integrator.withOutPorts({ output: port.withSignalName('CONCURRENCY_INCREMENT') })
          ),
          component.withSqrt(
            sqrt.new()
            + sqrt.withInPorts({ input: port.withSignalName('ACCEPTED_CONCURRENCY') })
            + sqrt.withOutPorts({ output: port.withSignalName('SQRT_CONCURRENCY_INCREMENT') })
            + sqrt.withScale($._config.constants.concurrencySQRTIncrementMultiplier),
          ),
          component.withPromql(
            local q = 'sum(increase(flux_meter_sum{valid="true", flow_status="OK", flux_meter_name="%(policyName)s"}[5s]))/sum(increase(flux_meter_count{valid="true", flow_status="OK", flux_meter_name="%(policyName)s"}[5s]))' % { policyName: $._config.policyName };
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
              max: port.withSignalName('NORMAL_CONCURRENCY_LIMIT'),
              control_variable: port.withSignalName('ACCEPTED_CONCURRENCY'),
              optimize: port.withSignalName('CONCURRENCY_INCREMENT'),
            })
            + gradient.withOutPortsMixin({
              output: port.withSignalName('DESIRED_CONCURRENCY'),
            })
          ),
          component.withExtrapolator(
            extrapolator.new()
            + extrapolator.withMaxExtrapolationInterval('5s')
            + extrapolator.withInPortsMixin({
              input: port.withSignalName('DESIRED_CONCURRENCY_RATIO'),
            })
            + extrapolator.withOutPortsMixin({
              output: port.withSignalName('LOAD_MULTIPLIER'),
            })
          ),
          component.withConcurrencyLimiter(
            local c = $._config.concurrencyLimiter;
            concurrencyLimiter.new()
            + concurrencyLimiter.withFlowSelector($._config.concurrencyLimiterFlowSelector)
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
              + loadActuator.withDynamicConfigKey('concurrency_limiter')
              + loadActuator.withAlerterConfig(
                alerterConfig.new()
                + alerterConfig.withAlertName(c.alerterName)
                + alerterConfig.withAlertChannels(c.alerterChannels)
                + alerterConfig.withResolveTimeout(c.alerterResolveTimeout)
              )
            )
          ),
          component.withDecider(
            decider.new()
            + decider.withOperator('gt')
            + decider.withInPortsMixin(
              decider.inPorts.withLhs(port.withSignalName('LATENCY'))
              + decider.inPorts.withRhs(port.withSignalName('LATENCY_SETPOINT'))
            )
            + decider.withOutPortsMixin(decider.outPorts.withOutput(port.withSignalName('IS_OVERLOAD')))
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
    dynamicConfig: {
      concurrency_limiter: loadActuatorDynamicConfig.withDryRun($._config.dynamicConfig.dryRun),
    },
  },

  policyResource: policyResource,
}
