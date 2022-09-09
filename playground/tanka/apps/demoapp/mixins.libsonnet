local grafanaOperator = import 'github.com/jsonnet-libs/grafana-operator-libsonnet/4.3/main.libsonnet';
local k = import 'github.com/jsonnet-libs/k8s-libsonnet/1.22/main.libsonnet';

local demoApp = import 'apps/demoapp/main.libsonnet';
local latencyGradientPolicy = import 'github.com/fluxninja/aperture-blueprints/lib/1.0/policies/latency-gradient.libsonnet';
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local Workload = aperture.v1.SchedulerWorkload;
local LabelMatcher = aperture.v1.LabelMatcher;
local WorkloadWithLabelMatcher = aperture.v1.SchedulerWorkloadAndLabelMatcher;

local classifier = aperture.v1.policylanguagev1Classifier;
local extractor = aperture.v1.Extractor;
local rule = aperture.v1.Rule;
local selector = aperture.v1.Selector;
local controlPoint = aperture.v1.ControlPoint;

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
  policyName: 'service1-demo-app',
  fluxMeterSelector: svcSelector,
  concurrencyLimiterSelector: svcSelector,
  classifiers: [
    classifier.new()
    + classifier.withSelector(svcSelector)
    + classifier.withRules({
      user_type: rule.new()
                 + rule.withExtractor(extractor.new()
                                      + extractor.withFrom('request.http.headers.user_type')),
    }),
  ],
  concurrencyLimiter+: {
    timeoutFactor: 0.1,
    defaultWorkload: {
      priority: 20,
    },
    workloads: [
      WorkloadWithLabelMatcher.new(
        workload=Workload.withPriority(50),
        // match the label extracted by classifier
        label_matcher=LabelMatcher.withMatchLabels({ user_type: 'guest' })
      ),
      WorkloadWithLabelMatcher.new(
        workload=Workload.withPriority(200),
        // match the http header directly
        label_matcher=LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })
      ),
    ],
  },
}).policy;

{
  configMap:
    k.core.v1.configMap.new('policies')
    + k.core.v1.configMap.metadata.withLabels({ 'fluxninja.com/validate': 'true' })
    + k.core.v1.configMap.metadata.withNamespace('aperture-controller')
    + k.core.v1.configMap.withData({
      'service1-demo-app.yaml': std.manifestYamlDoc(policy, quote_keys=false),
    }),
  demoapp: demoappMixin,
}
