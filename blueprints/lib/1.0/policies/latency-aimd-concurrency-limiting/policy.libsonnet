local spec = import '../../spec.libsonnet';
local config = import './config.libsonnet';

local policy = spec.v1.Policy;
local resources = spec.v1.Resources;
local circuit = spec.v1.Circuit;
local classifier = spec.v1.Classifier;
local flowSelector = spec.v1.FlowSelector;
local query = spec.v1.Query;
local component = spec.v1.Component;
local flowControl = spec.v1.FlowControl;
local promQL = spec.v1.PromQL;
local port = spec.v1.Port;
local combinator = spec.v1.ArithmeticCombinator;
local ema = spec.v1.EMA;
local gradient = spec.v1.GradientController;
local concurrency_controller = spec.v1.ConcurrencyLimiter;
local scheduler = spec.v1.Scheduler;
local schedulerParameters = spec.v1.SchedulerParameters;
local decider = spec.v1.Decider;
local switcher = spec.v1.Switcher;
local loadActuator = spec.v1.LoadActuator;
local alerterConfig = spec.v1.AlerterConfig;
local max = spec.v1.Max;
local min = spec.v1.Min;
local sqrt = spec.v1.Sqrt;
local firstValid = spec.v1.FirstValid;
local extrapolator = spec.v1.Extrapolator;
local constantSignal = spec.v1.ConstantSignal;
local aimdConcurrencyController = spec.v1.AIMDConcurrencyController;

function(params) {
  _config:: config.common + config.policy + params,

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFluxMetersMixin({ [$._config.policy_name]: $._config.flux_meter })
                           + resources.withClassifiers($._config.classifiers))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval='0.5s')
      + circuit.withComponents(
        [
          component.withQuery(
            query.new()
            + query.withPromql(
              local q = 'sum(increase(flux_meter_sum{valid="true", flow_status="OK", flux_meter_name="%(policy_name)s"}[5s]))/sum(increase(flux_meter_count{valid="true", flow_status="OK", flux_meter_name="%(policy_name)s"}[5s]))' % { policy_name: $._config.policy_name };
              promQL.new()
              + promQL.withQueryString(q)
              + promQL.withEvaluationInterval('1s')
              + promQL.withOutPorts({ output: port.withSignalName('LATENCY') }),
            ),
          ),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY'),
                                                            port.withConstantSignal($._config.overload_detection.latency_ema_limit_multiplier),
                                                            output=port.withSignalName('MAX_EMA'))),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                            port.withConstantSignal($._config.overload_detection.latency_tolerance_multiplier),
                                                            output=port.withSignalName('LATENCY_SETPOINT'))),
          component.withEma(
            ema.withParameters($._config.overload_detection.ema)
            + ema.withInPortsMixin(
              ema.inPorts.withInput(port.withSignalName('LATENCY'))
              + ema.inPorts.withMaxEnvelope(port.withSignalName('MAX_EMA'))
            )
            + ema.withOutPortsMixin(ema.outPorts.withOutput(port.withSignalName('LATENCY_EMA')))
          ),
          component.withFlowControl(
            flowControl.new()
            + flowControl.withAimdConcurrencyController(
              local cc = $._config.concurrency_controller;
              aimdConcurrencyController.new()
              + aimdConcurrencyController.withFlowSelector(cc.flow_selector)
              + aimdConcurrencyController.withSchedulerParameters(cc.scheduler)
              + aimdConcurrencyController.withGradientParameters(cc.gradient)
              + aimdConcurrencyController.withConcurrencyLimitMultiplier(cc.concurrency_limit_multiplier)
              + aimdConcurrencyController.withConcurrencyLinearIncrement(cc.concurrency_linear_increment)
              + aimdConcurrencyController.withConcurrencySqrtIncrementMultiplier(cc.concurrency_sqrt_increment_multiplier)
              + aimdConcurrencyController.withAlerterParameters(cc.alerter)
              + aimdConcurrencyController.withDryRunDynamicConfigKey('concurrency_controller')
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
          ),
        ] + $._config.components,
      ),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: $._config.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
    dynamicConfig: {
      concurrency_controller: $._config.concurrency_controller.dynamic_config,
    },
  },

  policyResource: policyResource,
}
