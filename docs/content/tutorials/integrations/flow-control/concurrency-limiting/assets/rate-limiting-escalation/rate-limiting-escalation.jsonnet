local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local latencyAIMDPolicy = aperture.policies.LatencyAIMDConcurrencyLimiting.policy;

local flowSelector = aperture.spec.v1.FlowSelector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local controlPoint = aperture.spec.v1.ControlPoint;
local classifier = aperture.spec.v1.Classifier;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local workloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local labelMatcher = aperture.spec.v1.LabelMatcher;
local workload = aperture.spec.v1.SchedulerWorkload;
local component = aperture.spec.v1.Component;
local flowControl = aperture.spec.v1.FlowControl;
local rateLimiter = aperture.spec.v1.RateLimiter;
local rateLimiterParameters = aperture.spec.v1.RateLimiterParameters;
local decider = aperture.spec.v1.Decider;
local switcher = aperture.spec.v1.Switcher;
local port = aperture.spec.v1.Port;


local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint('ingress')
  );

// Restrict this selector to only bot traffic
local rateLimiterSelector = flowSelector.new()
                            + flowSelector.withServiceSelector(
                              serviceSelector.new()
                              + serviceSelector.withAgentGroup('default')
                              + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                            )
                            + flowSelector.withFlowMatcher(
                              flowMatcher.new()
                              + flowMatcher.withControlPoint('ingress')
                              + flowMatcher.withLabelMatcher(
                                labelMatcher.withMatchLabels({ 'http.request.header.user_type': 'bot' })
                              )
                            );

local policyResource = latencyAIMDPolicy({
  policy_name: 'service1-demo-app',
  flux_meter: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector),
  concurrency_controller+: {
    flow_selector: svcSelector,
    dynamic_config: {
      dry_run: false,
    },
    scheduler+: {
      timeout_factor: 0.5,
      default_workload_parameters: {
        priority: 20,
      },
      workloads: [
        workload.new()
        + workload.withParameters(workloadParameters.withPriority(50))
        // match the label extracted by classifier
        + workload.withLabelMatcher(labelMatcher.withMatchLabels({ user_type: 'guest' })),
        workload.new()
        + workload.withParameters(workloadParameters.withPriority(200))
        // alternatively, match the http header directly
        + workload.withLabelMatcher(labelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
      ],
    },
  },
  classifiers: [
    classifier.new()
    + classifier.withFlowSelector(svcSelector)
    + classifier.withRules({
      user_type: rule.new()
                 + rule.withExtractor(extractor.new()
                                      + extractor.withFrom('request.http.headers.user-type')),
    }),
  ],
  // highlight-start
  components: [
    component.new()
    + component.withDecider(
      decider.new()
      + decider.withOperator('lt')
      + decider.withInPorts({ lhs: port.withSignalName('LOAD_MULTIPLIER'), rhs: port.withConstantSignal(1.0) })
      + decider.withOutPorts({ output: port.withSignalName('IS_BOT_ESCALATION') })
      + decider.withTrueFor('30s')
    ),
    component.new()
    + component.withSwitcher(
      switcher.new()
      + switcher.withInPorts({
        switch: port.withSignalName('IS_BOT_ESCALATION'),
        on_true: port.withConstantSignal(0.0),
        on_false: port.withConstantSignal(10.0),
      })
      + switcher.withOutPorts({ output: port.withSignalName('RATE_LIMIT') })
    ),
    component.new()
    + component.withFlowControl(
      flowControl.new()
      + flowControl.withRateLimiter(
        rateLimiter.new()
        + rateLimiter.withFlowSelector(rateLimiterSelector)
        + rateLimiter.withInPorts({ limit: port.withSignalName('RATE_LIMIT') })
        + rateLimiter.withParameters(
          rateLimiterParameters.new()
          + rateLimiterParameters.withLimitResetInterval('1s')
          + rateLimiterParameters.withLabelKey('http.request.header.user_id')
        )
        + rateLimiter.withDynamicConfigKey('rate_limiter'),
      ),
    ),
  ],
  // highlight-end
}).policyResource;

policyResource
