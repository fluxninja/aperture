local aperture = import '../../../../blueprints/lib/1.0/main.libsonnet';

local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';

local latencyGradientPolicy = aperture.blueprints.LatencyGradient.policy;

local workloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local labelMatcher = aperture.spec.v1.LabelMatcher;
local workload = aperture.spec.v1.SchedulerWorkload;

local classifier = aperture.spec.v1.Classifier;
local fluxMeter = aperture.spec.v1.FluxMeter;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local selector = aperture.spec.v1.Selector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local controlPoint = aperture.spec.v1.ControlPoint;
local component = aperture.spec.v1.Component;
local rateLimiter = aperture.spec.v1.RateLimiter;
local decider = aperture.spec.v1.Decider;
local switcher = aperture.spec.v1.Switcher;
local port = aperture.spec.v1.Port;
local alerter = aperture.spec.v1.Alerter;

local fluxMeterSelector = selector.new()
                          + selector.withServiceSelector(
                            serviceSelector.new()
                            + serviceSelector.withAgentGroup('default')
                            + serviceSelector.withService('service3-demo-app.demoapp.svc.cluster.local')
                          )
                          + selector.withFlowSelector(
                            flowSelector.new()
                            + flowSelector.withControlPoint('ingress')
                          );

local concurrencyLimiterSelector = selector.new()
                                   + selector.withServiceSelector(
                                     serviceSelector.new()
                                     + serviceSelector.withAgentGroup('default')
                                     + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                                   )
                                   + selector.withFlowSelector(
                                     flowSelector.new()
                                     + flowSelector.withControlPoint('ingress')
                                   );

// Restrict this selector to only bot traffic
local rateLimiterSelector = selector.new()
                            + selector.withServiceSelector(
                              serviceSelector.new()
                              + serviceSelector.withAgentGroup('default')
                              + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                            )
                            + selector.withFlowSelector(
                              flowSelector.new()
                              + flowSelector.withControlPoint('ingress')
                              + flowSelector.withLabelMatcher(
                                labelMatcher.withMatchLabels({ 'http.request.header.user_type': 'bot' })
                              )
                            );


local apertureControllerMixin =
  apertureControllerApp {
    values+:: {
      operator+: {
        image: {
          registry: 'docker.io/fluxninja',
          repository: 'aperture-operator',
          tag: 'latest',
          pullPolicy: 'IfNotPresent',
        },
      },
      controller+: {
        createUninstallHook: false,
        config+: {
          plugins+: {
            disabled_plugins: [
              'aperture-plugin-fluxninja',
            ],
          },
          log+: {
            pretty_console: true,
            non_blocking: true,
            level: 'info',
          },
          etcd+: {
            endpoints: ['http://controller-etcd.aperture-controller.svc.cluster.local:2379'],
          },
          prometheus+: {
            address: 'http://controller-prometheus-server.aperture-controller.svc.cluster.local:80',
          },
          alertmanagers+: {
            clients: [
              {
                name: 'test1',
                address: 'http://ingestion-service.cloud.svc.cluster.local:80',
              },
              {
                name: 'test2',
                address: 'http://ingestion-service.cloud.svc.cluster.local:80',
              },
            ],
          },
        },
        image: {
          registry: '',
          repository: 'docker.io/fluxninja/aperture-controller',
          tag: 'latest',
        },
      },
    },
  };

local policyResource = latencyGradientPolicy({
  policyName: 'service1-demo-app',
  fluxMeter: fluxMeter.new() + fluxMeter.withSelector(fluxMeterSelector),
  concurrencyLimiterSelector: concurrencyLimiterSelector,
  classifiers: [
    classifier.new()
    + classifier.withSelector(concurrencyLimiterSelector)
    + classifier.withRules({
      user_type: rule.new()
                 + rule.withExtractor(extractor.new()
                                      + extractor.withFrom('request.http.headers.user-type')),
    }),
  ],
  concurrencyLimiter+: {
    timeoutFactor: 0.5,
    defaultWorkloadParameters: {
      priority: 20,
    },
    workloads: [
      workload.new()
      + workload.withWorkloadParameters(workloadParameters.withPriority(50))
      // match the label extracted by classifier
      + workload.withLabelMatcher(labelMatcher.withMatchLabels({ user_type: 'guest' })),
      workload.new()
      + workload.withWorkloadParameters(workloadParameters.withPriority(200))
      // match the http header directly
      + workload.withLabelMatcher(labelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
    ],
  },
  components: [
    component.new()
    + component.withDecider(
      decider.new()
      + decider.withOperator('lt')
      + decider.withInPorts({ lhs: port.withSignalName('LOAD_MULTIPLIER'), rhs: port.withConstantValue(1.0) })
      + decider.withOutPorts({ output: port.withSignalName('IS_BOT_ESCALATION') })
      + decider.withTrueFor('30s')
    ),
    component.new()
    + component.withSwitcher(
      switcher.new()
      + switcher.withInPorts({
        switch: port.withSignalName('IS_BOT_ESCALATION'),
        on_true: port.withConstantValue(0.0),
        on_false: port.withConstantValue(10),
      })
      + switcher.withOutPorts({ output: port.withSignalName('RATE_LIMIT') })
    ),
    component.new()
    + component.withRateLimiter(
      rateLimiter.new()
      + rateLimiter.withSelector(rateLimiterSelector)
      + rateLimiter.withInPorts({ limit: port.withSignalName('RATE_LIMIT') })
      + rateLimiter.withLimitResetInterval('1s')
      + rateLimiter.withLabelKey('http.request.header.user_id')
      + rateLimiter.withDynamicConfigKey('rate_limiter'),
    ),
  ],
}).policyResource;

{
  latencyGradientPolicy: policyResource,
  controller: apertureControllerMixin,
}
