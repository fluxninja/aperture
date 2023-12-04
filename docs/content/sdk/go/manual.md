---
title: Manually setting feature control points
sidebar_position: 1
slug: manually-setting-feature-control-points-using-golang-sdk
keywords:
  - go
  - sdk
  - feature
  - flow
  - control
  - points
  - manual
---

```mdx-code-block
import CodeSnippet from '../../codeSnippet.js'
```

<a href={`https://pkg.go.dev/github.com/fluxninja/aperture-go/v2`}>Aperture Go
SDK</a> can be used to manually set feature control points within a Go service.

To do so, first install the SDK:

## Install SDK

```bash
go get github.com/fluxninja/aperture-go/v2
```

Now, create an instance of ApertureClient.

## Create ApertureClient Instance

To create an instance of ApertureClient, you need to provide the address of your
Aperture Cloud instance and an API key.

Address of your Aperture Cloud is made of Organization ID. For example, if your
organization ID is `ORGANIZATION` and available at port `443` then the address
is `ORGANIZATION.app.fluxninja.com:443`.

:::info API Key

You can create an API key for your project in the Aperture Cloud UI. For
detailed instructions on locating API Keys, refer to the [API Keys][api-keys]
section.

:::

### Configuring gRPC Client Options

Aperture Go SDK uses gRPC to communicate with Aperture Cloud. You can configure
gRPC client options by passing a list of gRPC client options to the
ApertureClient constructor.

<CodeSnippet lang="go" snippetName="grpcOptions" />

### Defining Aperture Options and Creating ApertureClient

<CodeSnippet lang="go" snippetName="clientConstructor" />

## Define Control Points

Once you have created an instance of ApertureClient, you can define control
points.

### Define Labels

You can define labels for your control points. These labels can be produced by
business logic or can be static. For example, you define a label `user_id` which
is the user ID of the user making the request. You can also define a static
label `version` which is the version of your service. Depending on your use
case, you can define any number of labels.

<CodeSnippet
    lang="go"
    snippetName="defineLabels"
 />

#### Passing Labels using `FlowParams`

Labels can be passed to the `FlowParams` struct. This struct is used to pass
parameters to the `Flow` method. Additionally, this struct also defines caching
behavior for the control point.

<CodeSnippet
    lang="go"
    snippetName="defineFlowParams"
 />

### Create Control Point

Once you have defined labels, you can start a flow. To start a flow, you need to
provide the control point of the flow and the `FlowParams` struct.

<CodeSnippet
    lang="go"
    snippetName="startFlow"
 />

<details><summary>Control Point Code</summary>
<p>

<CodeSnippet
    lang="go"
    snippetName="manualFlowNoCaching"
 />

</p>
</details>

To view the working example for more context, refer to the [example
app][example] available in the repository.

## HTTP Middleware

You can also configure middleware for your HTTP server using the SDK. To do
this, after creating an instance of ApertureClient, apply the middleware to your
router as demonstrated in the example below.

For added convenience, you can specify a list of regular expression patterns.
The middleware will only be applied to request paths that match these patterns.
This feature is particularly beneficial for endpoints such as `/health`,
`/connected` which may not require Aperture intervention.

<CodeSnippet
    lang="go"
    snippetName="middleware"
 />

<!-- TODO: Fix Link -->

[example]: https://github.com/fluxninja/aperture-go/tree/main/examples
[api-keys]: /reference/cloud-ui/api-keys.md
