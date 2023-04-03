local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local StaticRateLimiting = aperture.policies.StaticRateLimiting.policy;

local flowSelector = aperture.spec.v1.FlowSelector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowMatcher = aperture.spec.v1.FlowMatcher;
local classifier = aperture.spec.v1.Classifier;
local rego = aperture.spec.v1.Rego;

local svcSelector =
  flowSelector.new()
  + flowSelector.withServiceSelector(
    serviceSelector.new()
    + serviceSelector.withService('service-graphql-demo-app.demoapp.svc.cluster.local')
  )
  + flowSelector.withFlowMatcher(
    flowMatcher.new()
    + flowMatcher.withControlPoint('ingress')
  );

local policyResource = StaticRateLimiting({
  policy_name: 'graphql-static-rate-limiting',
  rate_limiter+: {
    flow_selector: svcSelector,
    rate_limit: 10.0,
    parameters+: {
      label_key: 'user_id',
      limit_reset_interval: '1s',
    },
  },

  classifiers: [
    classifier.new()
    + classifier.withFlowSelector(svcSelector)
    + classifier.withRego(
      local module = |||
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
        user_id := u if {
          queryIsCreateTodo
          u := claims.userID
        }
      |||;
      rego.new()
      + rego.withLabels({
        user_id: {
          telemetry: true,
        },
      })
      + rego.withModule(module),
    ),
  ],
}).policyResource;

policyResource
