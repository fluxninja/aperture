# Aperture Controller Configuration Reference

## Table of contents

### CLASSIFICATION CONFIGURATION

| Key | Reference                                  |
| --- | ------------------------------------------ |
|     | [ClassificationRule](#classification-rule) |

### Object Index

- [MatchExpressionList](#match-expression-list) – List of MatchExpressions that is used for all/any matching.
  eg. {any: {of: [expr…
- [RuleRego](#rule-rego) – Raw rego rules are compiled 1:1 to rego queries.
  High-level extractor-based rule…
- [v1AddressExtractor](#v1-address-extractor) – Display an [Address][ext-authz-address] as a single string, eg. `<ip>:<port>`.
  I…
- [v1Classifier](#v1-classifier) – Set of classification rules sharing a common selector.

Example:

````yaml
selecto…
* [v1ControlPoint](#v1-control-point) – Identifies control point within a service that the rule or policy should apply t…
* [v1EqualsMatchExpression](#v1-equals-match-expression) – Label selector expression of the equal form "label == value".
* [v1Extractor](#v1-extractor) – Defines a high-level way to specify how to extract a flow label given http request metadata, without a need to write regod code.
There are multiple variants of extractor, specify exactly one:
- JSON Extractor
- Address Extractor
- JWT Extractor
* [v1JSONExtractor](#v1-json-extractor) – Deserialize a json, and extract one of the fields.

Example:
```yaml
from: reque…
* [v1JWTExtractor](#v1-j-w-t-extractor) – Parse the attribute as JWT and read the payload.
Specify a field to be extracted…
* [v1K8sLabelMatcherRequirement](#v1-k8s-label-matcher-requirement) – Label selector requirement which is a selector that contains values, a key, and …
* [v1LabelMatcher](#v1-label-matcher) – Allows to define rules whether a map of labels should be considered a match or not.
It provides three ways to define requirements:
- matchLabels
- matchExpressions
- arbitrary expression
* [v1MatchExpression](#v1-match-expression) – Defines a [map<string, string> → bool] expression to be evaluated on labels.
…
* [v1MatchesMatchExpression](#v1-matches-match-expression) – Label selector expression of the matches form "label matches regex".
* [v1PathTemplateMatcher](#v1-path-template-matcher) – Matches HTTP Path to given path templates.
HTTP path will be matched against giv…
* [v1Rule](#v1-rule) – Rule describes a single Flow Classification Rule.
Flow classification rule extra…
* [v1Selector](#v1-selector) – Describes where a rule or actuation component should apply to.

Example:
```yaml…


## Reference

### <span id="classification-rule"></span> *ClassificationRule*







#### Members

<dl>

<dt></dt>
<dd>


Type: [V1Classifier](#v1-classifier)
</dd>
</dl>

## Objects

### <span id="match-expression-list"></span> MatchExpressionList


List of MatchExpressions that is used for all/any matching.
eg. {any: {of: [expr1, expr2]}}.



#### Properties
<dl>
<dt>of</dt>
<dd>

([[]V1MatchExpression](#v1-match-expression)) List of subexpressions of the match expression.

</dd>
</dl>

### <span id="rule-rego"></span> RuleRego


Raw rego rules are compiled 1:1 to rego queries.
High-level extractor-based rules are compiled into a single rego query.



#### Properties
<dl>
<dt>query</dt>
<dd>

(string) Query string to extract a value (eg. `data.<mymodulename>.<variablename>`).

Note: The module name must match the package name from the "source".

</dd>
</dl>
<dl>
<dt>source</dt>
<dd>

(string) Source code of the rego module.

Note: Must include a "package" declaration.

</dd>
</dl>

### <span id="v1-address-extractor"></span> v1AddressExtractor


Display an [Address][ext-authz-address] as a single string, eg. `<ip>:<port>`.
IP addresses in attribute context are defined as objects with separate ip and port fields.
This is a helper to display an address as a single string.

Note: Use with care, as it might accidentally introduce a high-cardinality flow label values.

[ext-authz-address]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/address.proto#config-core-v3-address

Example:
```yaml
from: "source.address # or dstination.address"
````

#### Properties

<dl>
<dt>from</dt>
<dd>

(string, `required`) Attribute path pointing to some string - eg. "source.address".

</dd>
</dl>

### <span id="v1-classifier"></span> v1Classifier

Set of classification rules sharing a common selector.

Example:

```yaml
selector:
  namespace: default
  service: service1
  control_point:
    traffic: ingress
rules:
  user:
    extractor:
      from: request.http.headers.user
```

#### Properties

<dl>
<dt>rules</dt>
<dd>

(map of [V1Rule](#v1-rule)) A map of {key, value} pairs mapping from flow label names to rules that define how to extract and propagate them.

</dd>
</dl>
<dl>
<dt>selector</dt>
<dd>

([V1Selector](#v1-selector)) Defines where to apply the flow classification rule.

</dd>
</dl>

### <span id="v1-control-point"></span> v1ControlPoint

Identifies control point within a service that the rule or policy should apply to.
Controlpoint is either a library feature name or one of ingress/egress traffic control point.

#### Properties

<dl>
<dt>feature</dt>
<dd>

(string, `required`) Name of FlunxNinja library's feature.
Feature corresponds to a block of code that can be "switched off" which usually is a "named opentelemetry's Span".

Note: Flowcontrol only.

</dd>
</dl>
<dl>
<dt>traffic</dt>
<dd>

(string, `required,oneof=ingress egress`) Type of traffic service, either "ingress" or "egress".
Apply the policy to the whole incoming/outgoing traffic of a service.
Usually powered by integration with a proxy (like envoy) or a web framework.

- Flowcontrol: Blockable atom here is a single HTTP-transaction.
- Classification: Apply the classification rules to every incoming/outgoing request and attach the resulting flow labels to baggage and telemetry.

</dd>
</dl>

### <span id="v1-equals-match-expression"></span> v1EqualsMatchExpression

Label selector expression of the equal form "label == value".

#### Properties

<dl>
<dt>label</dt>
<dd>

(string, `required`) Name of the label to equal match the value.

</dd>
</dl>
<dl>
<dt>value</dt>
<dd>

(string) Exact value that the label should be equal to.

</dd>
</dl>

### <span id="v1-extractor"></span> v1Extractor

Defines a high-level way to specify how to extract a flow label given http request metadata, without a need to write regod code.
There are multiple variants of extractor, specify exactly one:

- JSON Extractor
- Address Extractor
- JWT Extractor

#### Properties

<dl>
<dt>address</dt>
<dd>

([V1AddressExtractor](#v1-address-extractor)) Display an address as a single string - `<ip>:<port>`.

</dd>
</dl>
<dl>
<dt>from</dt>
<dd>

(string) Use an attribute with no convertion.
Attribute path is a dot-separated path to attribute.

Should be either:

- one of the fields of [Attribute Context][attribute-context], or
- a special "request.http.bearer" pseudo-attribute.
  Eg. "request.http.method" or "request.http.header.user-agent"

Note: The same attribute path syntax is shared by other extractor variants,
wherever attribute path is needed in their "from" syntax.

Example:

```yaml
from: request.http.headers.user-agent
```

[attribute-context]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto"

</dd>
</dl>
<dl>
<dt>json</dt>
<dd>

([V1JSONExtractor](#v1-json-extractor)) Deserialize a json, and extract one of the fields.

</dd>
</dl>
<dl>
<dt>jwt</dt>
<dd>

([V1JWTExtractor](#v1-j-w-t-extractor)) Parse the attribute as JWT and read the payload.

</dd>
</dl>
<dl>
<dt>path_templates</dt>
<dd>

([V1PathTemplateMatcher](#v1-path-template-matcher)) Match HTTP Path to given path templates.

</dd>
</dl>

### <span id="v1-json-extractor"></span> v1JSONExtractor

Deserialize a json, and extract one of the fields.

Example:

```yaml
from: request.http.body
pointer: /user/name
```

#### Properties

<dl>
<dt>from</dt>
<dd>

(string, `required`) Attribute path pointing to some strings - eg. "request.http.body".

</dd>
</dl>
<dl>
<dt>pointer</dt>
<dd>

(string) Json pointer represents a parsed json pointer which allows to select a specified field from the json payload.

Note: Uses [json pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
eg. `/foo/bar`. If the pointer points into an object, it'd be stringified.

</dd>
</dl>

### <span id="v1-j-w-t-extractor"></span> v1JWTExtractor

Parse the attribute as JWT and read the payload.
Specify a field to be extracted from payload using "json_pointer".

Note: The signature is not verified against the secret (we're assuming there's some
other parts of the system that handles such verification).

Example:

```yaml
from: request.http.bearer
json_pointer: /user/email
```

#### Properties

<dl>
<dt>from</dt>
<dd>

(string, `required`) Jwt token can be pulled from any input attribute, but most likely you'd want to use "request.http.bearer".

</dd>
</dl>
<dl>
<dt>json_pointer</dt>
<dd>

(string) Json pointer allowing to select a specified field from the json payload.

Note: Uses [json pointer](https://datatracker.ietf.org/doc/html/rfc6901) syntax,
eg. `/foo/bar`. If the pointer points into an object, it'd be stringified.

</dd>
</dl>

### <span id="v1-k8s-label-matcher-requirement"></span> v1K8sLabelMatcherRequirement

Label selector requirement which is a selector that contains values, a key, and an operator that relates the key and values.

#### Properties

<dl>
<dt>key</dt>
<dd>

(string, `required`) Label key that the selector applies to.

</dd>
</dl>
<dl>
<dt>operator</dt>
<dd>

(string, `oneof=In NotIn Exists DoesNotExists`) Logical operator which represents a key's relationship to a set of values.
Valid operators are In, NotIn, Exists and DoesNotExist.

</dd>
</dl>
<dl>
<dt>values</dt>
<dd>

([]string) An array of string values that relates to the key by an operator.
If the operator is In or NotIn, the values array must be non-empty.
If the operator is Exists or DoesNotExist, the values array must be empty.

</dd>
</dl>

### <span id="v1-label-matcher"></span> v1LabelMatcher

Allows to define rules whether a map of labels should be considered a match or not.
It provides three ways to define requirements:

- matchLabels
- matchExpressions
- arbitrary expression

If multiple requirements are set, they are all ANDed.
An empty label matcher always matches.

#### Properties

<dl>
<dt>expression</dt>
<dd>

([V1MatchExpression](#v1-match-expression)) An arbitrary expression to be evaluated on the labels.

</dd>
</dl>
<dl>
<dt>match_expressions</dt>
<dd>

([[]V1K8sLabelMatcherRequirement](#v1-k8s-label-matcher-requirement)) List of k8s-style label matcher requirements.

Note: The requirements are ANDed.

</dd>
</dl>
<dl>
<dt>match_labels</dt>
<dd>

(map of string) A map of {key,value} pairs representing labels to be matched.
A single {key,value} in the matchLabels requires that the label "key" is present and equal to "value".

Note: The requirements are ANDed.

</dd>
</dl>

### <span id="v1-match-expression"></span> v1MatchExpression

Defines a [map<string, string> → bool] expression to be evaluated on labels.
MatchExpression has multiple variants, exactly one should be set.

Example:

```yaml
all:
  of:
    - label_exists: foo
    - label_equals: { label = app, value = frobnicator }
```

#### Properties

<dl>
<dt>all</dt>
<dd>

([MatchExpressionList](#match-expression-list)) The expression is true when all subexpressions are true.

</dd>
</dl>
<dl>
<dt>any</dt>
<dd>

([MatchExpressionList](#match-expression-list)) The expression is true when any subexpression is true.

</dd>
</dl>
<dl>
<dt>label_equals</dt>
<dd>

([V1EqualsMatchExpression](#v1-equals-match-expression)) The expression is true when label value equals given value.

</dd>
</dl>
<dl>
<dt>label_exists</dt>
<dd>

(string, `required`) The expression is true when label with given name exists.

</dd>
</dl>
<dl>
<dt>label_matches</dt>
<dd>

([V1MatchesMatchExpression](#v1-matches-match-expression)) The expression is true when label matches given regex.

</dd>
</dl>
<dl>
<dt>not</dt>
<dd>

([V1MatchExpression](#v1-match-expression)) The expression negates the result of subexpression.

</dd>
</dl>

### <span id="v1-matches-match-expression"></span> v1MatchesMatchExpression

Label selector expression of the matches form "label matches regex".

#### Properties

<dl>
<dt>label</dt>
<dd>

(string, `required`) Name of the label to match the regular expression.

</dd>
</dl>
<dl>
<dt>regex</dt>
<dd>

(string, `required`) Regular expression that should match the label value.
It uses [golang's regular expression syntax](https://github.com/google/re2/wiki/Syntax).

</dd>
</dl>

### <span id="v1-path-template-matcher"></span> v1PathTemplateMatcher

Matches HTTP Path to given path templates.
HTTP path will be matched against given path templates.
If a match occurs, the value associated with the path template will be treated as a result.
In case of multiple path templates matching, the most specific one will be chosen.

#### Properties

<dl>
<dt>template_values</dt>
<dd>

(map of string) Template value keys are OpenAPI-inspired path templates.

Examples:

```
/register
/users/{user_id}
/static/*
```

- Static path segment `/foo` matches a path segment exactly.
- `/{param}` matches arbitrary path segment.
  (The param name is ignored and can be omitted (`{}`))
- The parameter must cover whole segment.
- Additionally, path template can end with `/*` wildcard to match
  arbitrary number of trailing segments (0 or more).
- Multiple consecutive `/` are ignored, as well as trailing `/`.
- Parametrized path segments must come after static segments.
- `*`, if present, must come last.
- Most specific template \"wins\" (`/foo` over `/{}` and `/{}` over `/*`).

See also <https://swagger.io/specification/#path-templating-matching>"

</dd>
</dl>

### <span id="v1-rule"></span> v1Rule

Rule describes a single Flow Classification Rule.
Flow classification rule extracts a value from request metadata.
More specifically, from `input`, which has the same spec as [Envoy's External Authorization Attribute Context][attribute-context].
See <https://play.openpolicyagent.org/p/gU7vcLkc70> for an example input.
There are two ways to define a flow classification rule:

- Using a declarative extractor – suitable from simple cases, such as directly reading a value from header or a field from json body.
- Rego expression.

Performance note: It's recommended to use declarative extractors where possible, as they may be slightly performant than Rego expressions.
[attribute-context](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto)

Example:

```yaml
Example of Declarative JSON extractor:
  yaml:
    extractor:
      json:
        from: request.http.body
        pointer: /user/name
    propagate: true
    hidden: false
Example of Rego module:
  yaml:
    rego:
      query: data.user_from_cookie.user
      source:
        package: user_from_cookie
        cookies: "split(input.attributes.request.http.headers.cookie, ';')"
        cookie: "cookies[_]"
        cookie.startswith: "('session=')"
        session: "substring(cookie, count('session='), -1)"
        parts: "split(session, '.')"
        object: "json.unmarshal(base64url.decode(parts[0]))"
        user: object.user
    propagate: false
    hidden: true
```

#### Properties

<dl>
<dt>extractor</dt>
<dd>

([V1Extractor](#v1-extractor)) High-level flow label declarative extractor.
Rego extractor extracts a value from the rego module.

</dd>
</dl>
<dl>
<dt>hidden</dt>
<dd>

(bool) Decides if the created flow label should be hidden from the telemetry.

</dd>
</dl>
<dl>
<dt>propagate</dt>
<dd>

(bool) Decides if the created label should be applied to the whole flow (propagated in baggage) (default=true).

</dd>
</dl>
<dl>
<dt>rego</dt>
<dd>

([RuleRego](#rule-rego)) Rego module to extract a value from the rego module.

</dd>
</dl>

### <span id="v1-selector"></span> v1Selector

Describes where a rule or actuation component should apply to.

Example:

```yaml
selector:
  namespace: default
  service: service1
  control_point:
    traffic: ingress # Allowed values are `ingress` and `egress`.
  label_matcher:
    match_labels:
      user_tier: gold
    match_expressions:
      - key: query
        operator: In
        values:
          - insert
          - delete
      - label: user_agent
        regex: ^(?!.*Chrome).*Safari
```

#### Properties

<dl>
<dt>agent_group</dt>
<dd>

(string, default: `default`) Describes where this selector applies to.

</dd>
</dl>
<dl>
<dt>control_point</dt>
<dd>

([V1ControlPoint](#v1-control-point), `required`) Describes control point within the entity where the policy should apply to.

</dd>
</dl>
<dl>
<dt>label_matcher</dt>
<dd>

([V1LabelMatcher](#v1-label-matcher)) Allows to add _additional_ condition on labels that must also be satisfied (in addition to namespace+service+control point matching).
The label matcher allows to match on infra labels, flow labels and request labels.
Arbitrary label matcher can be used to match infra labels.
For flowcontrol policies, the matcher can be used to match flow labels.

Note: For classification we can only match flow labels that were created at some **previous** control point.

In case of k8s, infra labels are labels on entities (note: there might exist some additional labels).
Flow label names are always prefixed with `flow_`.
Request labels are always prefixed with `request_`.
Available request labels are `id` (available as `request_id`), `method`, `path`, `host`, `scheme`, `size`, `protocol`.
(mapped from fields of [HttpRequest](https://github.com/envoyproxy/envoy/blob/637a92a56e2739b5f78441c337171968f18b46ee/api/envoy/service/auth/v3/attribute_context.proto#L102)).
Also, (non-pseudo) headers are available as `request_header_<headername>`.

Note: Request headers are only available for "traffic" control points.

</dd>
</dl>
<dl>
<dt>service</dt>
<dd>

(string) The service (name) of the entities.
In k8s, this is the FQDN of the Service object.

Note: Entity may belong to multiple services.

</dd>
</dl>

<!---
Generated File Ends
-->
