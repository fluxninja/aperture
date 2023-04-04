local spec = import '../../spec.libsonnet';
local config = import './config.libsonnet';

local policy = spec.v1.Policy;
local circuit = spec.v1.Circuit;
local query = spec.v1.Query;
local component = spec.v1.Component;
local promQL = spec.v1.PromQL;
local port = spec.v1.Port;
local gradient = spec.v1.GradientController;
local gradientInPort = spec.v1.GradientControllerIns;
local gradientOutPort = spec.v1.GradientControllerOuts;
local port = spec.v1.Port;
local outPorts = spec.v1.OutPort;
local autoScale = spec.v1.AutoScale;
local podAutoScaler = spec.v1.PodAutoScaler;
local podScaler = spec.v1.PodScaler;
local constant = spec.v1.ConstantSignal;
local kubernetesObjectSelector = spec.v1.KubernetesObjectSelector;
local scaleInController = spec.v1.ScaleInController;
local scaleInControllerController = spec.v1.ScaleInControllerController;
local scaleOutController = spec.v1.ScaleOutController;
local scaleOutControllerController = spec.v1.ScaleOutControllerController;
local alerterParameters = spec.v1.AlerterParameters;

function(cfg) {
  local params = config.common + config.policy + cfg,

  local policyDef =
    policy.new()
    + policy.withCircuit(
      circuit.new()
      + circuit.withEvaluationInterval(evaluation_interval='0.5s')
      + circuit.withComponents(
        [
          component.withAutoScale(
            autoScale.new()
            + autoScale.withPodAutoScaler(
              podAutoScaler.new()
              + podAutoScaler.withCooldownOverridePercentage(params.cooldown_override_percentage)
              + podAutoScaler.withMaxScaleInPercentage(params.max_scale_in_percentage)
              + podAutoScaler.withMaxScaleOutPercentage(params.max_scale_out_percentage)
              + podAutoScaler.withMinReplicas(params.min_replicas)
              + podAutoScaler.withMaxReplicas(params.max_replicas)
              + podAutoScaler.withScaleInCooldown(params.scale_in_cooldown)
              + podAutoScaler.withScaleOutCooldown(params.scale_out_cooldown)
              + podAutoScaler.withScaleInAlerterParameters(params.scale_in_alerter_parameters)
              + podAutoScaler.withScaleOutAlerterParameters(params.scale_out_alerter_parameters)
              + podAutoScaler.withOutPorts(
                {
                  actual_replicas: port.withSignalName('ACTUAL_REPLICAS'),
                  configured_replicas: port.withSignalName('CONFIGURED_REPLICAS'),
                  desired_replicas: port.withSignalName('DESIRED_REPLICAS'),
                }
              )
              + podAutoScaler.withPodScaler(
                podScaler.new()
                + podScaler.withKubernetesObjectSelector(
                  kubernetesObjectSelector.new()
                  + kubernetesObjectSelector.withNamespace(params.kubernetes_object_selector.namespace)
                  + kubernetesObjectSelector.withName(params.kubernetes_object_selector.name)
                  + kubernetesObjectSelector.withApiVersion(params.kubernetes_object_selector.api_version)
                  + kubernetesObjectSelector.withKind(params.kubernetes_object_selector.kind)
                )
              )
              + podAutoScaler.withScaleInControllers([
                scaleInController.new()
                + scaleInController.withController(
                  scaleInControllerController.new()
                  + scaleInControllerController.withGradient(
                    gradient.new()
                    + gradient.withInPorts(
                      gradientInPort.new()
                      + gradientInPort.withSignal(port.withSignalName(criteria.query.promql.out_ports.output.signal_name))
                      + gradientInPort.withSetpoint(port.withConstantSignal(criteria.set_point))
                    )
                    + gradient.withParameters(criteria.parameters)
                  )
                )
                for criteria in params.scale_in_criteria
              ])
              + podAutoScaler.withScaleOutControllers([
                scaleInController.new()
                + scaleInController.withController(
                  scaleInControllerController.new()
                  + scaleInControllerController.withGradient(
                    gradient.new()
                    + gradient.withInPorts(
                      gradientInPort.new()
                      + gradientInPort.withSignal(port.withSignalName(criteria.query.promql.out_ports.output.signal_name))
                      + gradientInPort.withSetpoint(port.withConstantSignal(criteria.set_point))
                    )
                    + gradient.withParameters(criteria.parameters)
                  )
                )
                for criteria in params.scale_out_criteria
              ])
            )
          ),
        ] + [
          component.withQuery(
            query.new()
            + query.withPromql(
              local q = criteria.query.promql.query_string;
              promQL.new()
              + promQL.withQueryString(q)
              + promQL.withEvaluationInterval(criteria.query.promql.evaluation_interval)
              + promQL.withOutPorts({ output: port.withSignalName(criteria.query.promql.out_ports.output.signal_name) }),
            ),
          )
          for criteria in params.scale_out_criteria + params.scale_in_criteria
        ] + params.components,
      ),
    ),

  local policyResource = {
    kind: 'Policy',
    apiVersion: 'fluxninja.com/v1alpha1',
    metadata: {
      name: params.policy_name,
      labels: {
        'fluxninja.com/validate': 'true',
      },
    },
    spec: policyDef,
  },

  policyResource: policyResource,
  policyDef: policyDef,
}
