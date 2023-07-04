---
title: Classifier
sidebar_position: 7
---

:::info See also

[_Classifier_ reference][reference].

:::

A _Classifier_ can be used to create additional [flow labels][label] based on
request attributes, if the existing flow labels aren't sufficient. _Classifiers_
are defined using [rules](#rules) based on the [Rego][rego] query language.

_Classifiers_ are available only at _HTTP_ [control points][control-point]. For
feature-based control points, developers can provide arbitrary flow labels by
setting baggage or directly as arguments to the `startFlow()` call.

## Defining a Classifier

To define a _Classifier_, it needs to be added as a resource in a
[policy][policies]. Like any other flow control component, a _Classifier_ is
applied only to the flows matching the specified [selectors][selector]

A _Classifier_ defines a set of rules to create new flow labels based on request
attributes. Envoy's [External Authorization][ext-authz] definition is used by
Aperture to describe the request attributes, specifically the
[`AttributeContext`][attr-context].

An example of how the request attributes might look can be seen in the [INPUT
section at this Rego playground][rego-playground].

Example of a Classifier that creates a flow label `user-type` based on the value
of the `user-type` HTTP header:

```yaml
policy:
  policy_name: checkout-service-policy
  resources:
    flow_control:
      classifiers: # resource name
        - selectors: # Selector Defining Classifier's Scope
            - service: checkout.default.svc.cluster.local
              control_point: ingress
          rules: # classification rules
            user_type: # flow label key
              extractor: # Declarative Extractor
                from: request.http.headers.user-type # HTTP header
```

:::caution

Although Classifier is defined as a resource in a [policy][policies], flow
labels aren't isolated in any way and are shared across policies.

:::

Any _Flow Labels_ created through the _Classifier_ become available in
subsequent stages of [flow processing](./flow-lifecycle.md). Additionally, the
_Flow Label_ is injected as baggage, so it will be available as a flow label in
downstream flows too (assuming you have [baggage propagation][baggage]
configured in your system). If [FluxNinja ARC extension][arc] plugin is enabled,
all flow labels including the ones created through classifier are available in
traffic analytics.

:::note

The extracted label's baggage propagation and inclusion in traffic analytics can
be [disabled][rule] in the classifier configuration.

:::

## Live Previewing Requests Attributes {#live-previewing-requests}

Live previewing of request attributes is a feature that allows real-time
examination of attributes flowing through services and control points. This can
be done using the [`aperturectl`][aperturectl] tool, aiding in setting up or
debugging a Classifier. It provides insights into the data a Classifier will
handle, enabling the creation of more effective classification rules.

For example:

```sh
aperturectl flow-control preview --kube checkout.default.svc.cluster.local ingress --http
```

Returns:

```json
{
  "samples": [
    {
      "attributes": {
        "destination": {
          "address": {
            "socketAddress": {
              "address": "10.244.1.20",
              "portValue": 8099
            }
          }
        },
        "metadataContext": {},
        "request": {
          "http": {
            "headers": {
              ":authority": "checkout.default.svc.cluster.local",
              ":method": "POST",
              ":path": "/request",
              ":scheme": "http",
              "content-length": "201",
              "content-type": "application/json",
              "cookie": "session=eyJ1c2VyIjoia2Vub2JpIn0.YbsY4Q.kTaKRTyOIfVlIbNB48d9YH6Q0wo",
              "user-agent": "k6/0.42.0 (https://k6.io/)",
              "user-id": "19",
              "user-type": "guest",
              "x-forwarded-proto": "http",
              "x-request-id": "26f01736-ec45-4b07-a202-bdec8930c7f8"
            },
            "host": "checkout.default.svc.cluster.local",
            "id": "14553976531353216255",
            "method": "POST",
            "path": "/request",
            "protocol": "HTTP/1.1",
            "scheme": "http"
          },
          "time": "2023-01-15T07:07:48.693035Z"
        },
        "source": {
          "address": {
            "socketAddress": {
              "address": "10.244.2.36",
              "portValue": 35388
            }
          }
        }
      },
      "parsed_body": null,
      "parsed_path": ["request"],
      "parsed_query": {},
      "truncated_body": false,
      "version": {
        "encoding": "protojson",
        "ext_authz": "v3"
      }
    }
  ]
}
```

Alternatively, use the
[Introspection API](/reference/api/agent/flow-preview-service-preview-http-requests.api.mdx)
directly on an Aperture agent:

Example:

```sh
curl -X POST localhost:8080/v1/flowcontrol/preview/http_requests/checkout.default.svc.cluster.local/ingress?samples=1
```

## Classification Rules {#rules}

:::note See Also

Rules ([reference][rule])

:::

Each classification rule consists of two main components:

- **Flow Label Key**: This is the identifier for the flow label. It is used to
  reference the flow label in other parts of the system.
- **Extraction Rule**: A rule how to extract the flow label value based on
  request attributes.

There are two ways to specify a classification rule:

- **Declarative extractors**
- **Rego modules**

:::caution Request body availability

The possibility of extracting values from the request body depends on how
[External Authorization in Envoy][ext-authz-extension] was configured. The
default [Istio Configuration][install-istio] does not enable request body
buffering, as it _might_ break some streaming APIs.

:::

### Declarative Extractors {#extractors}

Extractors provide a high-level way to specify how to extract a flow label value
given HTTP request attributes, eliminating the need to write Rego code. Provided
extractors include:

- [Extracting values from headers][extractor]

  ```yaml
  classifiers:
    - selectors:
        - service: checkout.default.svc.cluster.local
          control_point: ingress # control point
      rules: # classification rules
        user_type: # flow label key
          extractor: # extractor
            from: request.http.headers.user-type # HTTP header
  ```

- [Parsing a field from JSON encoded request payload][json-extractor]

  ```yaml
  from: request.http.body
  pointer: /user/name
  ```

- [Parsing JWT tokens][jwt-extractor]

  ```yaml
  from: request.http.bearer
  json_pointer: /user/email
  ```

:::info See Also

[Extractor reference][extractor]

:::

> Aperture aims to expand the set of extractors to cover the most-common use
> cases.

:::caution Naming Conventions for Flow Label Keys

Keys of flow labels created by extractors must be valid [Rego][rego] identifiers
(alphanumeric characters and underscore are allowed; also, label name cannot be
a [Rego keyword][rego-kw], like `if` or `default`).

:::

:::note Benefits of Explicit Flow Label Extraction from Headers

Extracting the value from the header might not seem useful, as the value is
already available as _Flow Label_ ([as
`http.request.header.<header>`][request-labels]), but adding flow label
explicitly might still be useful, as it enables baggage propagation and
telemetry for this flow label.

:::

<!-- vale off -->

### Advanced Classification with Rego Language {#rego}

<!-- vale on -->

:::note See Also

Rego [reference][rego-rule]

:::

For more complex scenarios, [the Rego language][rego] can be used to define the
extractor. Rego allows you to define a set of labels that are extracted after
evaluating a Rego module.

Example of Rego module which also disables telemetry visibility of label:

```rego

rego:
  labels:
    user:
      telemetry: false
  module: |
    package user_from_cookie
    cookies := split(input.attributes.request.http.headers.cookie, "; ")
    user := user {
        cookie := cookies[_]
        startswith(cookie, "session=")
        session := substring(cookie, count("session="), -1)
        parts := split(session, ".")
        object := json.unmarshal(base64url.decode(parts[0]))
        user := object.user
    }
```

[ext-authz-extension]:
  https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter#config-http-filters-ext-authz
[ext-authz]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#authorization-service-proto
[attr-context]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto
[rego-playground]: https://play.openpolicyagent.org/p/mG0sXxCNdQ
[label]: /concepts/flow-label.md
[baggage]: /concepts/flow-label.md#baggage
[request-labels]: ./flow-label.md#request-labels
[reference]: /reference/configuration/spec.md#classifier
[rule]: /reference/configuration/spec.md#rule
[extractor]: /reference/configuration/spec.md#extractor
[rego-rule]: /reference/configuration/spec.md#rule-rego
[arc]: /arc/extension.md
[selector]: ./selector.md
[policies]: /concepts/advanced/policy.md
[rego]: https://www.openpolicyagent.org/docs/latest/policy-language/
[rego-kw]:
  https://www.openpolicyagent.org/docs/latest/policy-reference/#reserved-names
[control-point]: ./control-point.md
[install-istio]: /integrations/istio/istio.md
[aperturectl]: /get-started/installation/aperture-cli/aperture-cli.md
[json-extractor]: /reference/configuration/spec.md#json-extractor
[jwt-extractor]: /reference/configuration/spec.md#j-w-t-extractor
