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
local staticBuckets = aperture.spec.v1.FluxMeterStaticBuckets;

local svcSelector = selector.new()
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
  fluxMeterSelector: svcSelector,
  fluxMeters: {
    'service1-demo-app':
      fluxMeter.new()
      + fluxMeter.withSelector(svcSelector)
      + fluxMeter.withAttributeKey('workload_duration_ms')
      + fluxMeter.withStaticBuckets(
        staticBuckets.new()
        + staticBuckets.withBuckets([5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0])
      ),
  },

  concurrencyLimiterSelector: svcSelector,
  classifiers: [
    classifier.new()
    + classifier.withSelector(svcSelector)
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
  policy: policyMixin,
  controller: apertureControllerMixin,
}
