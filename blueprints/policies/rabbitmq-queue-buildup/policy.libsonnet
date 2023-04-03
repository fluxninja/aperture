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
                             + flowControlResources.withClassifiers(params.classifiers)
                           ))
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval('1s')
      + circuit.withComponents(
        [
          component.withQuery(
            query.new()
            + query.withPromql(
              local q = 'sum(rabbitmq_message_current{rabbitmq_queue_name="service3-demo-app.demoapp.svc.cluster.local",state="ready"})';
              promQL.new()
              + promQL.withQueryString(q)
              + promQL.withEvaluationInterval('1s')
              + promQL.withOutPorts({ output: port.withSignalName('RABBITMQ_QUEUE_READY') }),
            )
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
                signal: port.withSignalName('RABBITMQ_QUEUE_READY'),
                setpoint: port.withConstantSignal(1000),
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
