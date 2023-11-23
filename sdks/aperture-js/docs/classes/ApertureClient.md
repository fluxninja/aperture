[@fluxninja/aperture-js](../README.md) / ApertureClient

# Class: ApertureClient

Represents the Aperture Client used for interacting with the Aperture Agent.

**`Example`**

```ts
const apertureClient = new ApertureClient({
 address:
   process.env.APERTURE_AGENT_ADDRESS !== undefined
     ? process.env.APERTURE_AGENT_ADDRESS
     : "localhost:8089",
 apiKey: process.env.APERTURE_API_KEY || undefined,
 // if process.env.APERTURE_AGENT_INSECURE set channelCredentials to insecure
 channelCredentials:
   process.env.APERTURE_AGENT_INSECURE !== undefined
     ? grpc.credentials.createInsecure()
     : grpc.credentials.createSsl(),
});
```

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

• **new ApertureClient**(`«destructured»`): [`ApertureClient`](ApertureClient.md)

Constructs a new instance of the ApertureClient.

#### Parameters

| Name | Type |
| :------ | :------ |
| `«destructured»` | `Object` |
| › `address` | `string` |
| › `apiKey?` | `string` |
| › `channelCredentials?` | `ChannelCredentials` |
| › `channelOptions?` | `ChannelOptions` |

#### Returns

[`ApertureClient`](ApertureClient.md)

**`Throws`**

Error if the address is not provided.

## Properties

### exporter

• `Private` `Readonly` **exporter**: `OTLPTraceExporter`

___

### fcsClient

• `Private` `Readonly` **fcsClient**: `FlowControlServiceClient`

___

### tracer

• `Private` `Readonly` **tracer**: `Tracer`

___

### tracerProvider

• `Private` `Readonly` **tracerProvider**: `NodeTracerProvider`

## Methods

### #newResource

▸ **#newResource**(): `IResource`

#### Returns

`IResource`

___

### GetState

▸ **GetState**(): `ConnectivityState`

Gets the current state of the gRPC channel.

#### Returns

`ConnectivityState`

The connectivity state of the channel.

___

### Shutdown

▸ **Shutdown**(): `void`

Shuts down the ApertureClient.

#### Returns

`void`

___

### StartFlow

▸ **StartFlow**(`controlPoint`, `params`): `Promise`\<[`Flow`](Flow.md)\>

Starts a new flow with the specified control point and parameters.
StartFlow takes a control point and labels that get passed to Aperture Agent via flowcontrolv1.Check call.
Return value is a Flow.
The default semantics are fail-to-wire. If StartFlow fails, calling Flow.ShouldRun() on returned Flow returns as true.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `controlPoint` | `string` | The control point for the flow. |
| `params` | [`FlowParams`](../interfaces/FlowParams.md) | The parameters for the flow. |

#### Returns

`Promise`\<[`Flow`](Flow.md)\>

A promise that resolves to a Flow object.

**`Example`**

```ts
apertureClient.StartFlow("awesomeFeature", {
 labels: labels,
 grpcCallOptions: {
   deadline: Date.now() + 30000,
 },
 rampMode: false,
 cacheKey: "cache",
});
