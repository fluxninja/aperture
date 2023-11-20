[@fluxninja/aperture-js](../README.md) / [Exports](../modules.md) /
ApertureClient

# Class: ApertureClient

Represents the Aperture Client used for interacting with the Aperture Agent.

## Table of contents

### Constructors

- [constructor](ApertureClient.md#constructor)

### Properties

- [exporter](ApertureClient.md#exporter)
- [fcsClient](ApertureClient.md#fcsclient)
- [tracer](ApertureClient.md#tracer)
- [tracerProvider](ApertureClient.md#tracerprovider)

### Methods

- [#newResource](ApertureClient.md##newresource)
- [GetState](ApertureClient.md#getstate)
- [Shutdown](ApertureClient.md#shutdown)
- [StartFlow](ApertureClient.md#startflow)

## Constructors

### constructor

• **new ApertureClient**(`«destructured»`):
[`ApertureClient`](ApertureClient.md)

Constructs a new instance of the ApertureClient.

#### Parameters

| Name                    | Type                 |
| :---------------------- | :------------------- |
| `«destructured»`        | `Object`             |
| › `address`             | `string`             |
| › `agentAPIKey?`        | `string`             |
| › `channelCredentials?` | `ChannelCredentials` |
| › `channelOptions?`     | `ChannelOptions`     |

#### Returns

[`ApertureClient`](ApertureClient.md)

**`Throws`**

Error if the address is not provided.

#### Defined in

[client.ts:65](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L65)

## Properties

### exporter

• `Private` `Readonly` **exporter**: `OTLPTraceExporter`

#### Defined in

[client.ts:51](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L51)

---

### fcsClient

• `Private` `Readonly` **fcsClient**: `FlowControlServiceClient`

#### Defined in

[client.ts:49](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L49)

---

### tracer

• `Private` `Readonly` **tracer**: `Tracer`

#### Defined in

[client.ts:55](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L55)

---

### tracerProvider

• `Private` `Readonly` **tracerProvider**: `NodeTracerProvider`

#### Defined in

[client.ts:53](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L53)

## Methods

### #newResource

▸ **#newResource**(): `IResource`

#### Returns

`IResource`

#### Defined in

[client.ts:198](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L198)

---

### GetState

▸ **GetState**(): `ConnectivityState`

Gets the current state of the gRPC channel.

#### Returns

`ConnectivityState`

The connectivity state of the channel.

#### Defined in

[client.ts:194](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L194)

---

### Shutdown

▸ **Shutdown**(): `void`

Shuts down the ApertureClient.

#### Returns

`void`

#### Defined in

[client.ts:183](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L183)

---

### StartFlow

▸ **StartFlow**(`controlPoint`, `params`): `Promise`\<[`Flow`](Flow.md)\>

Starts a new flow with the specified control point and parameters. StartFlow
takes a control point and labels that get passed to Aperture Agent via
flowcontrolv1.Check call. Return value is a Flow. The default semantics are
fail-to-wire. If StartFlow fails, calling Flow.ShouldRun() on returned Flow
returns as true.

#### Parameters

| Name           | Type                                        | Description                     |
| :------------- | :------------------------------------------ | :------------------------------ |
| `controlPoint` | `string`                                    | The control point for the flow. |
| `params`       | [`FlowParams`](../interfaces/FlowParams.md) | The parameters for the flow.    |

#### Returns

`Promise`\<[`Flow`](Flow.md)\>

A promise that resolves to a Flow object.

#### Defined in

[client.ts:133](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/client.ts#L133)
