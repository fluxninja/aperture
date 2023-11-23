[@fluxninja/aperture-js](../README.md) / Flow

# Class: Flow

Represents a flow in the SDK.

## Table of contents

### Constructors

- [constructor](Flow.md#constructor)

### Properties

- [cacheKey](Flow.md#cachekey)
- [checkResponse](Flow.md#checkresponse)
- [controlPoint](Flow.md#controlpoint)
- [ended](Flow.md#ended)
- [error](Flow.md#error)
- [fcsClient](Flow.md#fcsclient)
- [grpcCallOptions](Flow.md#grpccalloptions)
- [rampMode](Flow.md#rampmode)
- [span](Flow.md#span)
- [startDate](Flow.md#startdate)
- [status](Flow.md#status)

### Methods

- [CachedValue](Flow.md#cachedvalue)
- [CheckResponse](Flow.md#checkresponse-1)
- [DeleteCachedValue](Flow.md#deletecachedvalue)
- [End](Flow.md#end)
- [Error](Flow.md#error-1)
- [SetCachedValue](Flow.md#setcachedvalue)
- [SetStatus](Flow.md#setstatus)
- [ShouldRun](Flow.md#shouldrun)
- [Span](Flow.md#span-1)
- [protoDurationToJSON](Flow.md#protodurationtojson)
- [protoTimestampToJSON](Flow.md#prototimestamptojson)

## Constructors

### constructor

• **new Flow**(`fcsClient`, `grpcCallOptions`, `controlPoint`, `span`, `startDate`, `rampMode?`, `cacheKey?`, `checkResponse?`, `error?`): [`Flow`](Flow.md)

#### Parameters

| Name | Type | Default value |
| :------ | :------ | :------ |
| `fcsClient` | `FlowControlServiceClient` | `undefined` |
| `grpcCallOptions` | `CallOptions` | `undefined` |
| `controlPoint` | `string` | `undefined` |
| `span` | `Span` | `undefined` |
| `startDate` | `number` | `undefined` |
| `rampMode` | `boolean` | `false` |
| `cacheKey` | ``null`` \| `string` | `null` |
| `checkResponse` | ``null`` \| `CheckResponse__Output` | `null` |
| `error` | ``null`` \| `Error` | `null` |

#### Returns

[`Flow`](Flow.md)

## Properties

### cacheKey

• `Private` **cacheKey**: ``null`` \| `string` = `null`

___

### checkResponse

• `Private` **checkResponse**: ``null`` \| `CheckResponse__Output` = `null`

___

### controlPoint

• `Private` **controlPoint**: `string`

___

### ended

• `Private` **ended**: `boolean` = `false`

___

### error

• `Private` **error**: ``null`` \| `Error` = `null`

___

### fcsClient

• `Private` **fcsClient**: `FlowControlServiceClient`

___

### grpcCallOptions

• `Private` **grpcCallOptions**: `CallOptions`

___

### rampMode

• `Private` **rampMode**: `boolean` = `false`

___

### span

• `Private` **span**: `Span`

___

### startDate

• `Private` **startDate**: `number`

___

### status

• `Private` **status**: [`FlowStatus`](../README.md#flowstatus) = `FlowStatusEnum.OK`

## Methods

### CachedValue

▸ **CachedValue**(): [`CachedValueResponse`](CachedValueResponse.md)

Gets the cached value for the flow.

#### Returns

[`CachedValueResponse`](CachedValueResponse.md)

The cached value response.

___

### CheckResponse

▸ **CheckResponse**(): ``null`` \| `CheckResponse__Output`

Gets the check response of the flow.

#### Returns

``null`` \| `CheckResponse__Output`

The check response object.

___

### DeleteCachedValue

▸ **DeleteCachedValue**(): `Promise`\<`undefined` \| [`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)\>

Deletes the cached value for the flow.

#### Returns

`Promise`\<`undefined` \| [`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)\>

A promise that resolves to the response of the cache delete operation.

___

### End

▸ **End**(): `void`

Ends the flow and performs necessary cleanup.

#### Returns

`void`

___

### Error

▸ **Error**(): ``null`` \| `Error`

Gets the error associated with the flow.

#### Returns

``null`` \| `Error`

The error object.

___

### SetCachedValue

▸ **SetCachedValue**(`value`, `ttl`): `Promise`\<`undefined` \| [`SetCachedValueResponse`](SetCachedValueResponse.md)\>

Sets the cached value for the flow.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `value` | `Buffer` | The value to set. |
| `ttl` | `Duration` | The time-to-live for the cached value. |

#### Returns

`Promise`\<`undefined` \| [`SetCachedValueResponse`](SetCachedValueResponse.md)\>

A promise that resolves to the response of the cache upsert operation.

___

### SetStatus

▸ **SetStatus**(`status`): `void`

Sets the status of the flow.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `status` | [`FlowStatus`](../README.md#flowstatus) | The status to set. |

#### Returns

`void`

___

### ShouldRun

▸ **ShouldRun**(): `boolean`

Determines whether the flow should run based on the check response and ramp mode.

#### Returns

`boolean`

A boolean value indicating whether the flow should run.

___

### Span

▸ **Span**(): `Span`

Gets the span associated with the flow.

#### Returns

`Span`

The span object.

___

### protoDurationToJSON

▸ **protoDurationToJSON**(`duration`): ``null`` \| `string`

#### Parameters

| Name | Type |
| :------ | :------ |
| `duration` | ``null`` \| `Duration__Output` |

#### Returns

``null`` \| `string`

___

### protoTimestampToJSON

▸ **protoTimestampToJSON**(`timestamp`): ``null`` \| `string`

#### Parameters

| Name | Type |
| :------ | :------ |
| `timestamp` | ``null`` \| `Timestamp__Output` |

#### Returns

``null`` \| `string`
