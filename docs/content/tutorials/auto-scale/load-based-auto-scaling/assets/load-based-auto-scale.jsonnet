local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local latencyAIMDPolicy = aperture.policies.LatencyAIMDConcurrencyLimiting.policy;

local flowSelector = aperture.spec.v1.FlowSelector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local controlPoint = aperture.spec.v1.ControlPoint;
local component = aperture.spec.v1.Component;
local autoScale = aperture.spec.v1.AutoScale;
local autoScaler = aperture.spec.v1.AutoScaler;
local autoScalerKubernetesReplicas = aperture.spec.v1.KubernetesReplicas;
local kubeObjectSelector = aperture.spec.v1.KubernetesObjectSelector;
local increasingGradient = aperture.spec.v1.IncreasingGradient;
local decreasingGradient = aperture.spec.v1.DecreasingGradient;
local port = aperture.spec.v1.Port;
local scaler = aperture.spec.v1.AutoScalerScaler;
local query = aperture.spec.v1.Query;
local promQL = aperture.spec.v1.PromQL;

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
  policy_name: 'load-based-auto-scale',
  flux_meter: fluxMeter.new() + fluxMeter.withFlowSelector(svcSelector),
  concurrency_controller+: {
    flow_selector: svcSelector,
  },
  components: [
    component.withQuery(
      query.new()
      + query.withPromql(
        local q = 'avg(avg_over_time(k8s_pod_cpu_utilization{k8s_deployment_name="service1-demo-app"}[30s]))';
        promQL.new()
        + promQL.withQueryString(q)
        + promQL.withEvaluationInterval('10s')
        + promQL.withOutPorts({ output: port.withSignalName('AVERAGE_CPU') }),
      )
    ),
    component.withAutoScale(
      autoScale.new()
      + autoScale.withAutoScaler(
        autoScaler.new()
        + autoScaler.withScaler(
          scaler.withKubernetesReplicas(
            autoScalerKubernetesReplicas.new()
            + autoScalerKubernetesReplicas.withKubernetesObjectSelector(
              kubeObjectSelector.new()
              + kubeObjectSelector.withNamespace('demoapp')
              + kubeObjectSelector.withApiVersion('apps/v1')
              + kubeObjectSelector.withKind('Deployment')
              + kubeObjectSelector.withName('service1-demo-app')
            )
          )
        )
        + autoScaler.withMinScale('1')
        + autoScaler.withMaxScale('10')
        + autoScaler.withScaleInCooldown('40s')
        + autoScaler.withScaleOutCooldown('30s')
        + autoScaler.withScaleOutControllers(
          [
            {
              controller: {
                gradient: increasingGradient.new()
                          + increasingGradient.withInPorts(
                            {
                              signal: port.withSignalName('OBSERVED_LOAD_MULTIPLIER'),
                              setpoint: port.withConstantSignal(1.0),
                            }
                          )
                          + increasingGradient.withParameters({
                            slope: -1.0,
                          }),
              },
            },
          ]
        )
        + autoScaler.withScaleInControllers(
          [
            {
              controller: {
                gradient: decreasingGradient.new()
                          + decreasingGradient.withInPorts(
                            {
                              signal: port.withSignalName('AVERAGE_CPU'),
                              setpoint: port.withConstantSignal(0.5),
                            }
                          )
                          + decreasingGradient.withParameters({
                            slope: 1.0,
                          }),
              },
            },
          ]
        )
      )
    ),
  ],
}).policyResource;

policyResource
