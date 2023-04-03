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
local inPort = spec.v1.InPort;
local outPorts = spec.v1.OutPort;
local autoScale = spec.v1.AutoScale;
local podAutoScaler = spec.v1.PodAutoScaler;
local constant = spec.v1.ConstantSignal;

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
              + podAutoScaler.withMinReplicas(params.min_replicas)
              + podAutoScaler.withMaxReplicas(params.max_replicas)
              + podAutoScaler.withScaleInCooldown(params.scale_in_cooldown)
              + podAutoScaler.withScaleOutCooldown(params.scale_out_cooldown)
              + podAutoScaler.withScaleInControllers([
                gradient.new()
                + gradient.withInPorts(
                  gradientInPort.new()
                  + gradientInPort.withSignal(inPort.withSignalName(criteria.query.promql.out_ports.output.signal_name))
                  + gradientInPort.withSetpoint(inPort.withConstantSignal(constant.withValue(criteria.set_point)))
                )
                + gradient.withParameters(criteria.parameters)
                for criteria in params.scale_in_criteria
              ])
              + podAutoScaler.withScaleOutControllers([
                gradient.new()
                + gradient.withInPorts(
                  gradientInPort.new()
                  + gradientInPort.withSignal(inPort.withSignalName(criteria.query.promql.out_ports.output.signal_name))
                  + gradientInPort.withSetpoint(inPort.withConstantSignal(constant.withValue(criteria.set_point)))
                )
                + gradient.withParameters(criteria.parameters)
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
