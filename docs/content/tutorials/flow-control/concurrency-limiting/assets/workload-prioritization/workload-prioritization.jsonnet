local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local latencyAIMDPolicy = aperture.policies.LatencyAIMDConcurrencyLimiting.policy;

local selector = aperture.spec.v1.Selector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local classifier = aperture.spec.v1.Classifier;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local workloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local labelMatcher = aperture.spec.v1.LabelMatcher;
local workload = aperture.spec.v1.SchedulerWorkload;


local svcSelectors = [
  selector.new()
  + selector.withControlPoint('ingress')
  + selector.withService('service1-demo-app.demoapp.svc.cluster.local')
  + selector.withAgentGroup('default'),
];

local policyResource = latencyAIMDPolicy({
  policy_name: 'service1-demo-app',
  flux_meter: fluxMeter.new() + fluxMeter.withSelectors(svcSelectors),
  concurrency_controller+: {
    selectors: svcSelectors,
    default_config: {
      dry_run: false,
    },
    // highlight-start
    scheduler+: {
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
    + classifier.withSelectors(svcSelectors)
    + classifier.withRules({
      user_type: rule.new()
                 + rule.withExtractor(extractor.new()
                                      + extractor.withFrom('request.http.headers.user-type')),
    }),
  ],
  // highlight-end
}).policyResource;

policyResource
