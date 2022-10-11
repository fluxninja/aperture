---
title: Flow Classifier
sidebar_position: 4
---

:::info

See also [Classifier reference][reference]

:::

If existing [Flow Labels][label] are not sufficient, Flow Classifier can be used
to inject additional ones without any changes to your service.

A classifier is defined as a resource in a [policy][policies] and describes a
set of _rules_ on how to create new flow labels based on request metadata.
Aperture uses Envoy's [External Authorization][ext-authz] definition to describe
the request metadata (more specifically, the [AttributeContext][attr-context]).
The [INPUT section at this Rego playground][rego-playground] is an example how
the request attributes may look like.

:::note

At _feature_ [Control Points][control-point] developers already can provide
arbitrary flow labels – either by setting baggage, or directly as arguments to
the `Check()` call. Since at feature control points flow labels can be easily
controlled, classifiers are available only on at _traffic_ control points.

:::

Flow Labels created via Classifier are immediately available for use in other
components at the same [Control Point][control-point]. The Flow Label is also
injected as baggage, so it will be available on every subsequent control point
too (assuming you have [baggage propagation][baggage] configured in your
system). If you're a [FluxNinja Cloud plugin][plugin] user, such flow label will
also be available in the Cloud for analytics.

:::note

Both these behaviours (baggage propagation and inclusion in telemetry) can be
[disabled][rule].

:::

:::caution

Although Classifier is defined as a resource in a [policy][policies], Flow
Labels are not namespaced in any way.

:::

## Selector

Each classifier needs to specify which control point it will be run at. Most
likely, you'll want to classify ingress traffic to your service:

```yaml
selector:
  service_selector:
    service: service1.default.svc.cluster.local
  flow_selector:
    control_point:
      traffic: ingress
```

You can be more precise by adding a [label matcher][label-matcher] and e.g. gate
the classifier to particular paths.

## Rules ([reference][rule]) {#rules}

In addition to the selector, a classifier needs to specify classification rules.
Each classification rule consists of:

- flow label key,
- a rule how to extract the flow label value based on request metadata.

There are two ways to specify a classification rule: using declarative
extractors and [rego][rego] modules. [See examples in reference][rule].

:::caution Request body availability

Possibility of extracting values from request body depends on how [External
Authorization in Envoy][ext-authz-extension] was configured. Sample [Istio
Configuration][install-istio] provided by FluxNinja doesn't enable request body
buffering by default, as it _might_ break some streaming APIs.

:::

### Extractors ([reference][extractor]) {#extractors}

Extractors are declarative recipes how to extract flow label value from
metadata. Provided extractors include:

- extracting values from headers,
- grabbing a field from json-encoded request payload,
- parsing JWT tokens,
- and others.

Aperture aims to expand the set of extractors to cover most-common usecases.

:::caution

Keys of flow labels created by extractors must be valid [Rego][rego] identifiers
(allowed are alphanumeric characters and underscore; also label name cannot be a
[Rego keyword][rego-kw], like `if` or `default`). This limitation may be lifted
in future.

:::

:::note

Extracting value from header may seem not useful, as the value is already
available as flow label ([as `http.request.header.<header>`][request-labels]),
but adding flow label explicitly may still be useful, as it enables baggage
propagation and telemetry for this flow label.

:::

### Rego ([reference][rego-rule]) {#rego}

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
[request-labels]: /concepts/flow-control/flow-label.md#request-labels
[reference]: /references/configuration/policy.md#v1-classifier
[rule]: /references/configuration/policy.md#v1-rule
[extractor]: /references/configuration/policy.md#v1-extractor
[rego-rule]: /references/configuration/policy.md#rule-rego
[plugin]: /cloud/plugin.md
[label-matcher]: /concepts/flow-control/selector.md#label-matcher
[policies]: /concepts/policy/policy.md
[rego]: https://www.openpolicyagent.org/docs/latest/policy-language/
[rego-kw]:
  https://www.openpolicyagent.org/docs/latest/policy-reference/#reserved-names
[control-point]: /concepts/flow-control/flow-control.md#control-point
[install-istio]: /get-started/installation/agent/envoy/istio.md
