---
title: Classifier
sidebar_position: 1
---

:::info

See also [_Classifier_ reference][reference]

:::

The _Classifier_ can be used to create additional [_Flow Labels_][label] based
on request metadata without requiring any changes to your service, if the
existing flow labels aren't sufficient.

To define a Classifier, it needs to be added as a resource in a
[policy][policies]. It specifies a set of rules to create new flow labels based
on request metadata. Envoy's [External Authorization][ext-authz] definition is
used by Aperture to describe the request metadata, specifically the
[`AttributeContext`][attr-context]. An example of how the request attributes
might look can be seen in the [INPUT section at this Rego
playground][rego-playground].

:::note

At _Feature_ [_Control Points_][control-point], developers can already provide
arbitrary flow labels by setting baggage or directly as arguments to the
`Check()` call. As flow labels can be effortlessly provided at _Feature_ control
points by the developers, _Classifiers_ are available only at _HTTP_ control
points.

:::

Any _Flow Labels_ created through the _Classifier_ are immediately available for
use in other components at the same [_Control Point_][control-point].
Additionally, the _Flow Label_ is injected as baggage, so it will be available
on every subsequent control point too (assuming you have [baggage
propagation][baggage] configured in your system). If you're a [FluxNinja ARC
extension][arc] user, such flow label will also be available for analytics.

:::note

Both these behaviors (baggage propagation and inclusion in telemetry) can be
[disabled][rule].

:::

:::caution

Although Classifier is defined as a resource in a [_Policy_][policies], _Flow
Labels_ aren't isolated in any way and are shared across policies.

:::

## Selector {#selector}

Each _Classifier_ needs to specify which control point it will be run at. For
instance, the following selector is for the "ingress" control point at a
service:

```yaml
selector:
  service_selector:
    service: service1.default.svc.cluster.local
  flow_selector:
    control_point: ingress
```

You can be more precise by adding a [_Label Matcher_][label-matcher] and, for
example, gate the Classifier to particular paths.

## Live Previewing Requests {#live-previewing-requests}

You can discover the request attributes flowing through services and control
points using [`aperturectl`][aperturectl].

For example:

```sh
aperturectl flow-control preview --kube service1-demo-app.demoapp.svc.cluster.local ingress --http
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
              ":authority": "service1-demo-app.demoapp.svc.cluster.local",
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
            "host": "service1-demo-app.demoapp.svc.cluster.local",
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

Alternatively, you can use the
[Introspection API](reference/api/agent/flow-preview-service-preview-http-requests.api.mdx)
directly on a `aperture-agent` local to the service instances (pods):

```sh
curl -X POST localhost:8080/v1/flowcontrol/preview/http_requests/service1-demo-app.demoapp.svc.cluster.local/ingress?samples=1
```

## Rules ([reference][rule]) {#rules}

In addition to the selector, a Classifier needs to specify classification rules.
Each classification rule consists of:

- _Flow Label_ key,
- A rule how to extract the flow label value based on request metadata.

There are two ways to specify a classification rule: using declarative
extractors and [Rego][rego] modules. [See examples in reference][rule].

:::caution Request body availability

The possibility of extracting values from the request body depends on how
[External Authorization in Envoy][ext-authz-extension] was configured. The
Sample [Istio Configuration][install-istio] provided by FluxNinja does not
enable request body buffering by default, as it _might_ break some streaming
APIs.

:::

### Extractors ([reference][extractor]) {#extractors}

Extractors are declarative recipes how to extract flow label value from
metadata. Provided extractors include:

- Extracting values from headers
- Parsing a field from JSON encoded request payload
- Parsing JWT tokens

Aperture aims to expand the set of extractors to cover the most-common use
cases.

:::caution

Keys of flow labels created by extractors must be valid [Rego][rego] identifiers
(alphanumeric characters and underscore are allowed; also, label name cannot be
a [Rego keyword][rego-kw], like `if` or `default`).

:::

:::note

Extracting the value from the header might not seem useful, as the value is
already available as _Flow Label_ ([as
`http.request.header.<header>`][request-labels]), but adding flow label
explicitly might still be useful, as it enables baggage propagation and
telemetry for this flow label.

:::

<!-- vale off -->

### Rego ([reference][rego-rule]) {#rego}

<!-- vale on -->

For more advanced cases, you can define the extractor in [the Rego
language][rego].

## Example

See [full example in reference][reference]

[ext-authz-extension]:
  https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter#config-http-filters-ext-authz
[ext-authz]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#authorization-service-proto
[attr-context]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto
[rego-playground]: https://play.openpolicyagent.org/p/mG0sXxCNdQ
[label]: /concepts/flow-control/flow-label.md
[baggage]: /concepts/flow-control/flow-label.md#baggage
[request-labels]: ../flow-label.md#request-labels
[reference]: /reference/policies/spec.md#classifier
[rule]: /reference/policies/spec.md#rule
[extractor]: /reference/policies/spec.md#extractor
[rego-rule]: /reference/policies/spec.md#rule-rego
[arc]: /arc/extension.md
[label-matcher]: ../flow-selector.md#label-matcher
[policies]: /concepts/policy/policy.md
[rego]: https://www.openpolicyagent.org/docs/latest/policy-language/
[rego-kw]:
  https://www.openpolicyagent.org/docs/latest/policy-reference/#reserved-names
[control-point]: ../flow-selector.md#control-point
[install-istio]: /get-started/integrations/flow-control/envoy/istio.md
[aperturectl]: /get-started/aperture-cli/aperture-cli.md
