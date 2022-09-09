---
title: Flow Classifier
---

:::info

See also [Classifier reference][reference]

:::

If existing [flow labels](label.md) are not sufficent, Flow Classifier can be
used to inject additional ones without any changes to your service.

A classifier is defined as a resource in a [policy][policies] and describes a
set of _rules_ on how to create new flow labels based on request metadata.
Aperture uses Envoy's [External Authorization][ext-authz] definition to describe
the request metadata (more specifically, the [AttributeContext][attr-context]).
The [INPUT section at this Rego playground][rego-playground] is an example how
the request attributes may look like.

:::note

At _feature_ [control points][control-point] developers already can provide
arbitrary flow labels – either by setting baggage, or directly as arguments to
the `Check()` call. Since at feature control points flow labels can be easily
controlled, classifiers are available only on at _traffic_ control points.

:::

Classifier-created flow label is immediately available for usage in other
components at the same [control point][control-point]. The flow-label is also
injected as baggage, so it will be available on every subsequent control point
too (assuming you have [baggage propagation][baggage] configured in your
system). If you're a [FluxNinja Cloud plugin][plugin] user, such flow label will
also be available in the Cloud for analytics.

:::note

Both these behaviours (baggage propagation and inclusion in telemetry) are
[disablable][rule].

:::

:::caution

Although the classifier is defined as a resource in a [policy][policies], flow
labels are not namespaced in any way.

:::

## Selector

Each classifier needs to specify which control point it will be run at. Most
likely, you'll want to classify ingress traffic to your service:

```yaml
selector:
  service: service1.default.svc.cluster.local
  control_point:
    traffic: ingress
```

You can be more precise by adding a [label matcher][label-matcher] and eg. gate
the classifier to particular paths.

## Rules ([reference][rule]) {#rules}

In addition to the selector, a classifier needs to specify classification rules.
Each classification rule consists of:

- flow label key,
- a rule how to extract the flow label value based on request metadata.

There are two ways to specify a classification rule: using declarative
extractors and [rego][rego] modules. [See examples in reference][rule].

### Extractors ([reference][extractor]) {#extractors}

Extractors are declarative recipes how to extract flow label value fom metadata.
Provided extractors include:

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
available as flow label
([as `http.request.header.<header>`](label.md#request-labels)), but adding flow
label explicitly may still be userful, as it enables baggage propagation and
telemetry for this flow label.

:::

### Rego ([reference][rego-rule]) {#rego}

For more advanced cases, you can define the extractor in [the Rego
language][rego].

## Example

See [full example in reference][reference]

[ext-authz]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto#authorization-service-proto
[attr-context]:
  https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/attribute_context.proto
[rego-playground]: https://play.openpolicyagent.org/p/mG0sXxCNdQ
[baggage]: label.md#baggage
[reference]: /reference/configuration/policies.md#v1-classifier
[rule]: /reference/configuration/policies.md#v1-rule
[extractor]: /reference/configuration/policies.md#v1-extractor
[rego-rule]: /reference/configuration/policies.md#rule-rego
[plugin]: /cloud/plugin.md
[label-matcher]: ../selector.md#label-matcher
[policies]: /concepts/policies/policies.md
[rego]: https://www.openpolicyagent.org/docs/latest/policy-language/
[rego-kw]:
  https://www.openpolicyagent.org/docs/latest/policy-reference/#reserved-names
[control-point]: ../flow-control.md#control-point
