local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local StaticRateLimiting = aperture.policies.StaticRateLimiting.policy;

local selector = aperture.spec.v1.Selector;
local classifier = aperture.spec.v1.Classifier;
local rego = aperture.spec.v1.Rego;

local svcSelectors = [
  selector.new()
  + selector.withControlPoint('ingress')
  + selector.withService('service-graphql-demo-app.demoapp.svc.cluster.local'),
];

local policyResource = StaticRateLimiting({
  policy+: {
    policy_name: 'graphql-static-rate-limiting',
    rate_limiter+: {
      selectors: svcSelectors,
      rate_limit: 10.0,
      parameters+: {
        label_key: 'user_id',
        limit_reset_interval: '1s',
      },
    },
    classifiers: [
      classifier.new()
      + classifier.withSelectors(svcSelectors)
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
  },
}).policyResource;

policyResource
