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
local gradientParameters = spec.v1.GradientParameters;
local concurrencyLimiter = spec.v1.ConcurrencyLimiter;
local scheduler = spec.v1.Scheduler;
local schedulerParameters = spec.v1.SchedulerParameters;
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
local aimdConcurrencyController = spec.v1.AIMDConcurrencyController;

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
          component.withPromql(
            local q = 'sum(increase(flux_meter_sum{valid="true", flow_status="OK", flux_meter_name="%(policyName)s"}[5s]))/sum(increase(flux_meter_count{valid="true", flow_status="OK", flux_meter_name="%(policyName)s"}[5s]))' % { policyName: $._config.policyName };
            promQL.new()
            + promQL.withQueryString(q)
            + promQL.withEvaluationInterval('1s')
            + promQL.withOutPorts({ output: port.withSignalName('LATENCY') }),
          ),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY'),
                                                            port.withConstantSignal(c.latencyEMALimitMultiplier),
                                                            output=port.withSignalName('MAX_EMA'))),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                            port.withConstantSignal(c.latencyToleranceMultiplier),
                                                            output=port.withSignalName('LATENCY_SETPOINT'))),
          component.withEma(
            local e = $._config.ema;
            ema.withEmaWindow(e.window)
            + ema.withWarmupWindow(e.warmupWindow)
            + ema.withCorrectionFactorOnMaxEnvelopeViolation(e.correctionFactor)
            + ema.withInPortsMixin(
              ema.inPorts.withInput(port.withSignalName('LATENCY'))
              + ema.inPorts.withMaxEnvelope(port.withSignalName('MAX_EMA'))
            )
            + ema.withOutPortsMixin(ema.outPorts.withOutput(port.withSignalName('LATENCY_EMA')))
          ),
          component.withAimdConcurrencyController(
            aimdConcurrencyController.new()
            + aimdConcurrencyController.withFlowSelector($._config.concurrencyLimiterFlowSelector)
            + aimdConcurrencyController.withSchedulerParameters(
              schedulerParameters.new()
              + schedulerParameters.withAutoTokens($._config.concurrencyLimiter.autoTokens)
              + schedulerParameters.withTimeoutFactor($._config.concurrencyLimiter.timeoutFactor)
              + schedulerParameters.withDefaultWorkloadParameters($._config.concurrencyLimiter.defaultWorkloadParameters)
              + schedulerParameters.withWorkloads($._config.concurrencyLimiter.workloads)
            )
            + aimdConcurrencyController.withGradientParameters(
              local g = $._config.gradient;
              gradientParameters.new()
              + gradientParameters.withSlope(g.slope)
              + gradientParameters.withMinGradient(g.minGradient)
              + gradientParameters.withMaxGradient(g.maxGradient)
            )
            + aimdConcurrencyController.withConcurrencyLimitMultiplier(c.concurrencyLimitMultiplier)
            + aimdConcurrencyController.withConcurrencyLinearIncrement(c.concurrencyLinearIncrement)
            + aimdConcurrencyController.withConcurrencySqrtIncrementMultiplier(c.concurrencySQRTIncrementMultiplier)
            + aimdConcurrencyController.withAlerterConfig(
              alerterConfig.new()
              + alerterConfig.withAlertName($._config.concurrencyLimiter.alerterName)
              + alerterConfig.withAlertChannels($._config.concurrencyLimiter.alerterChannels)
              + alerterConfig.withResolveTimeout($._config.concurrencyLimiter.alerterResolveTimeout)
            )
            + aimdConcurrencyController.withDryRunDynamicConfigKey('concurrency_limiter')
            + aimdConcurrencyController.withInPorts({
              signal: port.withSignalName('LATENCY'),
              setpoint: port.withSignalName('LATENCY_SETPOINT'),
            })
            + aimdConcurrencyController.withOutPorts({
              accepted_concurrency: port.withSignalName('ACCEPTED_CONCURRENCY'),
              incoming_concurrency: port.withSignalName('INCOMING_CONCURRENCY'),
              desired_concurrency: port.withSignalName('DESIRED_CONCURRENCY'),
              is_overload: port.withSignalName('IS_OVERLOAD'),
              load_multiplier: port.withSignalName('LOAD_MULTIPLIER'),
            }),
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
