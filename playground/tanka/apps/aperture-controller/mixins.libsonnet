local aperture = import '../../../../blueprints/lib/1.0/main.libsonnet';

local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';

local latencyGradientPolicy = aperture.blueprints.policies.LatencyGradient;

local WorkloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local LabelMatcher = aperture.spec.v1.LabelMatcher;
local Workload = aperture.spec.v1.SchedulerWorkload;

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

local fluxMeterSelector = selector.new()
                          + selector.withServiceSelector(
                            serviceSelector.new()
                            + serviceSelector.withAgentGroup('default')
                            + serviceSelector.withService('service3-demo-app.demoapp.svc.cluster.local')
                          )
                          + selector.withFlowSelector(
                            flowSelector.new()
                            + flowSelector.withControlPoint(controlPoint.new()
                                                            + controlPoint.withTraffic('ingress'))
                          );

local concurrencyLimiterSelector = selector.new()
                                   + selector.withServiceSelector(
                                     serviceSelector.new()
                                     + serviceSelector.withAgentGroup('default')
                                     + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                                   )
                                   + selector.withFlowSelector(
                                     flowSelector.new()
                                     + flowSelector.withControlPoint(controlPoint.new()
                                                                     + controlPoint.withTraffic('ingress'))
                                   );

local rateLimiterSelector = selector.new()
                            + selector.withServiceSelector(
                              serviceSelector.new()
                              + serviceSelector.withAgentGroup('default')
                              + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                            )
                            + selector.withFlowSelector(
                              flowSelector.new()
                              + flowSelector.withControlPoint(controlPoint.new()
                                                              + controlPoint.withTraffic('ingress'))
                              + flowSelector.withLabelMatcher(
                                LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'bot' })
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
            level: 'debug',
          },
          etcd+: {
            endpoints: ['http://controller-etcd.aperture-controller.svc.cluster.local:2379'],
          },
          prometheus+: {
            address: 'http://controller-prometheus-server.aperture-controller.svc.cluster.local:80',
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

local policy = latencyGradientPolicy({
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
      Workload.new()
      + Workload.withWorkloadParameters(WorkloadParameters.withPriority(50))
      // match the label extracted by classifier
      + Workload.withLabelMatcher(LabelMatcher.withMatchLabels({ user_type: 'guest' })),
      Workload.new()
      + Workload.withWorkloadParameters(WorkloadParameters.withPriority(200))
      // match the http header directly
      + Workload.withLabelMatcher(LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
    ],
  },
  components: [
    component.new()
    + component.withDecider(
      decider.new()
      + decider.withOperator('gt')
      + decider.withInPorts({ lhs: port.withSignalName('LSF'), rhs: port.withConstantValue(0.1) })
      + decider.withOutPorts({ output: port.withSignalName('IS_BOT_ESCALATION') })
      + decider.withTrueFor('30s')
    ),
    component.new()
    + component.withSwitcher(
      switcher.new()
      + switcher.withInPorts({ switch: port.withSignalName('IS_BOT_ESCALATION'), on_true: port.withConstantValue(0.0), on_false: port.withConstantValue(10) })
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
}).policy;

local policyMixin = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'service1-demo-app',
    labels: {
      'fluxninja.com/validate': 'true',
    },
  },
  spec: policy,
};

{
  latencyGradientPolicy: policyMixin,
  controller: apertureControllerMixin,
}
