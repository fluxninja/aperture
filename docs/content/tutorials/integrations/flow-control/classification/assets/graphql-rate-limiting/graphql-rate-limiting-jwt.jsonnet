local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local policy = aperture.spec.v1.Policy;
local resources = aperture.spec.v1.Resources;
local component = aperture.spec.v1.Component;
local flowControl = aperture.spec.v1.FlowControl;
local rateLimiter = aperture.spec.v1.RateLimiter;
local rateLimiterParameters = aperture.spec.v1.RateLimiterParameters;
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
    + flowMatcher.withControlPoint('ingress')
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
                     import future.keywords.if
                     query_ast := graphql.parse_query(input.parsed_body.query)
                     claims := payload if {
                       io.jwt.verify_hs256(bearer_token, "secret")
                       [_, payload, _] := io.jwt.decode(bearer_token)
                     }
                     bearer_token := t if {
                       v := input.attributes.request.http.headers.authorization
                       startswith(v, "Bearer ")
                       t := substring(v, count("Bearer "), -1)
                     }
                     queryIsCreateTodo if {
                       some operation
                       walk(query_ast, [_, operation])
                       operation.Name == "createTodo"
                       count(operation.SelectionSet) > 0
                       some selection
                       walk(operation.SelectionSet, [_, selection])
                       selection.Name == "createTodo"
                     }
                     userID := u if {
                       queryIsCreateTodo
                       u := claims.userID
                     }
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
    + circuit.withEvaluationInterval('0.5s')
    + circuit.withComponents([
      component.withFlowControl(
        flowControl.new()
        + flowControl.withRateLimiter(
          rateLimiter.new()
          + rateLimiter.withInPorts({ limit: port.withConstantSignal(10.0) })
          + rateLimiter.withFlowSelector(svcSelector)
          + rateLimiter.withParameters(
            rateLimiterParameters.new()
            + rateLimiterParameters.withLimitResetInterval('1s')
            + rateLimiterParameters.withLabelKey('user_id')
            + rateLimiterParameters.withLazySync({ enabled: false, num_sync: 5 })
          )
        ),
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
