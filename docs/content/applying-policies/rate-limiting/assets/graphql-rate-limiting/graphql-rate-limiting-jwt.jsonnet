local aperture = import 'github.com/fluxninja/aperture/blueprints/main.libsonnet';

local RateLimiting = aperture.policies.RateLimiting.policy;

local selector = aperture.spec.v1.Selector;
local classifier = aperture.spec.v1.Classifier;
local rego = aperture.spec.v1.Rego;

local svcSelectors = [
  selector.new()
  + selector.withControlPoint('ingress')
  + selector.withService('service-graphql-demo-app.demoapp.svc.cluster.local'),
];

local policyResource = RateLimiting({
  policy+: {
    policy_name: 'graphql-rate-limiting',
    rate_limiter+: {
      selectors: svcSelectors,
      bucket_capacity: 40,
      fill_amount: 2,
      parameters+: {
        label_key: 'user_id',
        interval: '1s',
      },
    },
    resources: {
      flow_control: {
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
    },
  },
}).policyResource;

policyResource
