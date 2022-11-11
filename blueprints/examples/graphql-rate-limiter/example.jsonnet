local aperture = import '../../lib/1.0/main.libsonnet';
local bundle = aperture.blueprints.RateLimiter.bundle;

local LazySync = aperture.spec.v1.RateLimiterLazySync;
local selector = aperture.spec.v1.Selector;
local serviceSelector = aperture.spec.v1.ServiceSelector;
local flowSelector = aperture.spec.v1.FlowSelector;
local fluxMeter = aperture.spec.v1.FluxMeter;
local controlPoint = aperture.spec.v1.ControlPoint;
local classifier = aperture.spec.v1.Classifier;
local rule = aperture.spec.v1.Rule;
local rego = aperture.spec.v1.RuleRego;

local svcSelector = selector.new()
                    + selector.withServiceSelector(
                      serviceSelector.new()
                      + serviceSelector.withAgentGroup('default')
                      + serviceSelector.withService('service-graphql-demo-app.demoapp.svc.cluster.local')
                    )
                    + selector.withFlowSelector(
                      flowSelector.new()
                      + flowSelector.withControlPoint(controlPoint.new()
                                                      + controlPoint.withTraffic('ingress'))
                    );

local config = {
  common+: {
    policyName: 'example',
  },
  policy+: {
    rateLimiterSelector: svcSelector,
    rateLimit: '20.0',
    labelKey: 'user_id',
    limitResetInterval: '1s',
    classifiers: [
      classifier.new()
      + classifier.withSelector(svcSelector)
      + classifier.withRules({
        user_id: rule.new()
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
    ],
  },
};

bundle { _config+:: config }
