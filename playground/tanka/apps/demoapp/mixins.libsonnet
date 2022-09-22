local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local k = import 'github.com/jsonnet-libs/k8s-libsonnet/1.22/main.libsonnet';

local demoApp = import 'apps/demoapp/main.libsonnet';
local latencyGradientPolicy = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/policies/latency-gradient.libsonnet';
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local Workload = aperture.v1.SchedulerWorkload;
local LabelMatcher = aperture.v1.LabelMatcher;
local WorkloadWithLabelMatcher = aperture.v1.SchedulerWorkloadAndLabelMatcher;

local classifier = aperture.v1.Classifier;
local fluxMeter = aperture.v1.FluxMeter;
local extractor = aperture.v1.Extractor;
local rule = aperture.v1.Rule;
local selector = aperture.v1.Selector;
local controlPoint = aperture.v1.ControlPoint;
local staticBuckets = aperture.v1.FluxMeterStaticBuckets;
local linearBuckets = aperture.v1.FluxMeterLinearBuckets;
local exponentialBuckets = aperture.v1.ExponentialBuckets;
local exponentialBucketsRange = aperture.v1.ExponentialBucketsRange;

local svcSelector = selector.new()
                    + selector.withAgentGroup('default')
                    + selector.withService('service1-demo-app.demoapp.svc.cluster.local')
                    + selector.withControlPoint(controlPoint.new()
                                                + controlPoint.withTraffic('ingress'));
local demoappMixin =
  demoApp {
    values+: {
      replicaCount: 2,
      simplesrv+: {
        image: {
          repository: 'docker.io/fluxninja/demo-app',
          tag: 'test',
        },
      },
      resources+: {
        limits+: {
          cpu: '100m',
          memory: '128Mi',
        },
        requests+: {
          cpu: '100m',
          memory: '128Mi',
        },
      },
    },
  };

local policy = latencyGradientPolicy({
  policyName: 'service1-demoapp',
  fluxMeterSelector: svcSelector,
  fluxMeters: {
    'service1-demoapp':
      fluxMeter.new()
      + fluxMeter.withSelector(svcSelector)
      + fluxMeter.withAttributeKey('workload_duration_ms')
      + fluxMeter.withStaticBuckets(
        staticBuckets.new()
        + staticBuckets.withBuckets([5.0, 10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0, 10000.0])
      )
      + fluxMeter.withLinearBuckets(
        linearBuckets.new()
        + linearBuckets.withStart(1.0)
        + linearBuckets.withWidth(300.0)
        + linearBuckets.withCount(50)
      )
      + fluxMeter.withExponentialBuckets(
        exponentialBuckets.new()
        + exponentialBuckets.withStart(2.0)
        + exponentialBuckets.withFactor(1.5)
        + exponentialBuckets.withCount(50)
      )
      + fluxMeter.withExponentialBucketsRange(
        exponentialBucketsRange.new()
        + exponentialBucketsRange.withMin(1.0)
        + exponentialBucketsRange.withMax(1000.0)
        + exponentialBucketsRange.withCount(50)
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
      + WorkloadWithLabelMatcher.withWorkload(Workload.withPriority(50)),
      // match the label extracted by classifier
      +WorkloadWithLabelMatcher.withLabelMatcher(LabelMatcher.withMatchLabels({ user_type: 'guest' })),
      WorkloadWithLabelMatcher.new()
      + WorkloadWithLabelMatcher.withWorkload(Workload.withPriority(200)),
      // match the http header directly
      +WorkloadWithLabelMatcher.withLabelMatcher(LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
    ],
  },
}).policy;

local policMixin = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'service1',
  },
  spec: policy,
};

{
  policy: policMixin,
  demoapp: demoappMixin,
}
