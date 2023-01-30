local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

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

local policyResource = latencyAIMDPolicy({
  policy_name: 'service1-demo-app',
  flux_meter: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector),
  concurrency_controller+: {
    flow_selector: svcSelector,
    dynamicConfig: {
      dryRun: false,
    },
    // highlight-start
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
  // highlight-end
}).policyResource;

policyResource
