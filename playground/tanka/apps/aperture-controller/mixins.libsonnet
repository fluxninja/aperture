local latencyGradientPolicy = import '../../../../blueprints/lib/1.0/policies/latency-gradient.libsonnet';
local aperture = import '../../../../blueprints/libsonnet/1.0/main.libsonnet';
local apertureControllerApp = import 'apps/aperture-controller/main.libsonnet';
local demoApp = import 'apps/demoapp/main.libsonnet';

local Workload = aperture.v1.SchedulerWorkload;
local LabelMatcher = aperture.v1.LabelMatcher;
local WorkloadWithLabelMatcher = aperture.v1.SchedulerWorkloadAndLabelMatcher;

local classifier = aperture.v1.Classifier;
local fluxMeter = aperture.v1.FluxMeter;
local extractor = aperture.v1.Extractor;
local rule = aperture.v1.Rule;
local selector = aperture.v1.Selector;
local serviceSelector = aperture.v1.ServiceSelector;
local flowSelector = aperture.v1.FlowSelector;
local controlPoint = aperture.v1.ControlPoint;
local staticBuckets = aperture.v1.FluxMeterStaticBuckets;

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
    defaultWorkload: {
      priority: 20,
    },
    workloads: [
      WorkloadWithLabelMatcher.new()
      + WorkloadWithLabelMatcher.withWorkload(Workload.withPriority(50))
      // match the label extracted by classifier
      + WorkloadWithLabelMatcher.withLabelMatcher(LabelMatcher.withMatchLabels({ user_type: 'guest' })),
      WorkloadWithLabelMatcher.new()
      + WorkloadWithLabelMatcher.withWorkload(Workload.withPriority(200))
      // match the http header directly
      + WorkloadWithLabelMatcher.withLabelMatcher(LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
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
