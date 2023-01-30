local aperture = import '../../lib/1.0/main.libsonnet';
local bundle = aperture.policies.LatencyAIMDConcurrencyLimiting.bundle;

local workloadParameters = aperture.spec.v1.SchedulerWorkloadParameters;
local labelMatcher = aperture.spec.v1.LabelMatcher;
local workload = aperture.spec.v1.SchedulerWorkload;
local classifier = aperture.spec.v1.Classifier;
local fluxMeter = aperture.spec.v1.FluxMeter;
local extractor = aperture.spec.v1.Extractor;
local rule = aperture.spec.v1.Rule;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local controlPoint = aperture.spec.v1.ControlPoint;

local svcSelector = flowSelector.new()
                    + flowSelector.withServiceSelector(
                      serviceSelector.new()
                      + serviceSelector.withAgentGroup('default')
                      + serviceSelector.withService('service1-demo-app.demoapp.svc.cluster.local')
                    )
                    + flowSelector.withFlowMatcher(
                      flowMatcher.new()
                      + flowMatcher.withControlPoint('ingress')
                    );


local config = {
  common+: {
    policy_name: 'example',
  },
  policy+: {
    flux_meter: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector),
    classifiers: [
      classifier.new()
      + classifier.withFlowSelector(svcSelector)
      + classifier.withRules({
        user_type: rule.new()
                   + rule.withExtractor(extractor.new()
                                        + extractor.withFrom('request.http.headers.user-type')),
      }),
    ],
    concurrency_controller+: {
      flow_selector: svcSelector,
      default_workload_parameters: {
        priority: 20,
      },
      scheduler+: {
        workloads: [
          workload.new()
          + workload.withParameters(workloadParameters.withPriority(50))
          // match the label extracted by classifier
          + workload.withLabelMatcher(labelMatcher.withMatchLabels({ user_type: 'guest' })),
          workload.new()
          + workload.withParameters(workloadParameters.withPriority(200))
          // match the http header directly
          + workload.withLabelMatcher(labelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' })),
        ],
      },
    },
  },
};

bundle { _config+:: config }
