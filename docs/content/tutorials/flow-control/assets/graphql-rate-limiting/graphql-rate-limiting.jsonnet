local aperture = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local resources = aperture.spec.v1.Resources;
local component = aperture.spec.v1.Component;
local rateLimiter = aperture.spec.v1.RateLimiter;
local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local circuit = aperture.spec.v1.Circuit;
local port = aperture.spec.v1.Port;
local classifier = aperture.spec.v1.Classifier;
local rule = aperture.spec.v1.Rule;
local rego = aperture.spec.v1.RuleRego;

local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withAgentGroup('default')
    + serviceSelector.withService('service-graphql-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint({ traffic: 'ingress' })
  );

local policyDef =
  policy.new()
  + policy.withResources(
    resources.new()
    + resources.withClassifiers(
      classifier.new()
      + classifier.withFlowSelector(svcSelector)
      + classifier.withRules({
        user_id: rule.new()
                 + rule.withTelemetry(true)
                 + rule.withRego(
                   local source = |||
                     package graphql_example

                     query_ast := graphql.parse_query(input.parsed_body.query)

                     createTodoByUserIds[name] := value {
                       some operation
                       walk(query_ast, [_, operation])
                       operation.Name == "createTodo"
                       count(operation.SelectionSet) > 0
                       some selection
                       walk(operation.SelectionSet, [_, selection])
                       selection.Name == "createTodo"
                       count(selection.Arguments) > 0
                       argument := selection.Arguments[_]
                       argument.Name == "input"
                       count(argument.Value.Children) > 0
                       child := argument.Value.Children[_]
                       child.Name == "userId"
                       name := child.Name
                       value := child.Value.Raw
                     }

                     userID := createTodoByUserIds.userId
                   |||;
                   rego.new()
                   + rego.withQuery('data.graphql_example.userID')
                   + rego.withSource(source)
                 ),
      }),
    )
  )
  + policy.withCircuit(
    circuit.new()
    + circuit.withEvaluationInterval('300s')
    + circuit.withComponents([
      component.withRateLimiter(
        rateLimiter.new()
        + rateLimiter.withInPorts({ limit: port.withConstantValue(10) })
        + rateLimiter.withFlowSelector(svcSelector)
        + rateLimiter.withLimitResetInterval('1s')
        + rateLimiter.withLabelKey('user_id')
        + rateLimiter.withLazySync({ enabled: false, num_sync: 5 })
      ),
    ]),
  );

local policyResource = {
  kind: 'Policy',
  apiVersion: 'fluxninja.com/v1alpha1',
  metadata: {
    name: 'graphql-static-rate-limiting',
    labels: {
      'fluxninja.com/validate': 'true',
    },
  },
  spec: policyDef,
};

policyResource
