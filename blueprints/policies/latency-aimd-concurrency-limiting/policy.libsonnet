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
local flowControlResources = spec.v1.FlowControlResources;
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

function(cfg) {
  local params = config.common + config.policy + cfg,

  local policyDef =
    policy.new()
    + policy.withResources(resources.new()
                           + resources.withFlowControl(
                             flowControlResources.new()
                             + flowControlResources.withFluxMetersMixin({ [params.policy_name]: params.flux_meter })
                             + flowControlResources.withClassifiers(params.classifiers)
                           ))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval='0.5s')
      + circuit.withComponents(
        [
          component.withQuery(
            query.new()
            + query.withPromql(
              local q = 'sum(increase(flux_meter_sum{flow_status="OK", flux_meter_name="%(policy_name)s"}[5s]))/sum(increase(flux_meter_count{flow_status="OK", flux_meter_name="%(policy_name)s"}[5s]))' % { policy_name: params.policy_name };
              promQL.new()
              + promQL.withQueryString(q)
              + promQL.withEvaluationInterval('1s')
              + promQL.withOutPorts({ output: port.withSignalName('LATENCY') }),
            ),
          ),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY'),
                                                            port.withConstantSignal(params.latency_baseliner.latency_ema_limit_multiplier),
                                                            output=port.withSignalName('MAX_EMA'))),
          component.withArithmeticCombinator(combinator.mul(port.withSignalName('LATENCY_EMA'),
                                                            port.withConstantSignal(params.latency_baseliner.latency_tolerance_multiplier),
                                                            output=port.withSignalName('LATENCY_SETPOINT'))),
          component.withEma(
            ema.withParameters(params.latency_baseliner.ema)
            + ema.withInPortsMixin(
              ema.inPorts.withInput(port.withSignalName('LATENCY'))
              + ema.inPorts.withMaxEnvelope(port.withSignalName('MAX_EMA'))
            )
            + ema.withOutPortsMixin(ema.outPorts.withOutput(port.withSignalName('LATENCY_EMA')))
          ),
          component.withFlowControl(
            flowControl.new()
            + flowControl.withAimdConcurrencyController(
              local cc = params.concurrency_controller;
              aimdConcurrencyController.new()
              + aimdConcurrencyController.withFlowSelector(cc.flow_selector)
              + aimdConcurrencyController.withSchedulerParameters(cc.scheduler)
              + aimdConcurrencyController.withGradientParameters(cc.gradient)
              + aimdConcurrencyController.withMaxLoadMultiplier(cc.max_load_multiplier)
              + aimdConcurrencyController.withLoadMultiplierLinearIncrement(cc.load_multiplier_linear_increment)
              + aimdConcurrencyController.withAlerterParameters(cc.alerter)
              + aimdConcurrencyController.withDynamicConfigKey('concurrency_controller')
              + aimdConcurrencyController.withDefaultConfig(params.concurrency_controller.default_config)
              + aimdConcurrencyController.withInPorts({
                signal: port.withSignalName('LATENCY'),
                setpoint: port.withSignalName('LATENCY_SETPOINT'),
              })
              + aimdConcurrencyController.withOutPorts({
                is_overload: port.withSignalName('IS_OVERLOAD'),
                desired_load_multiplier: port.withSignalName('DESIRED_LOAD_MULTIPLIER'),
                observed_load_multiplier: port.withSignalName('OBSERVED_LOAD_MULTIPLIER'),
                accepted_concurrency: port.withSignalName('ACCEPTED_CONCURRENCY'),
                incoming_concurrency: port.withSignalName('INCOMING_CONCURRENCY'),
              }),
            ),
          ),
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
