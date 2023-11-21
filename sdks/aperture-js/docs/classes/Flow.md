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

• **new Flow**(`fcsClient`, `grpcCallOptions`, `controlPoint`, `span`,
`startDate`, `rampMode?`, `cacheKey?`, `checkResponse?`, `error?`):
[`Flow`](Flow.md)

#### Parameters

| Name              | Type                              | Default value |
| :---------------- | :-------------------------------- | :------------ |
| `fcsClient`       | `FlowControlServiceClient`        | `undefined`   |
| `grpcCallOptions` | `CallOptions`                     | `undefined`   |
| `controlPoint`    | `string`                          | `undefined`   |
| `span`            | `Span`                            | `undefined`   |
| `startDate`       | `number`                          | `undefined`   |
| `rampMode`        | `boolean`                         | `false`       |
| `cacheKey`        | `null` \| `string`                | `null`        |
| `checkResponse`   | `null` \| `CheckResponse__Output` | `null`        |
| `error`           | `null` \| `Error`                 | `null`        |

#### Returns

[`Flow`](Flow.md)

#### Defined in

[flow.ts:53](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L53)

## Properties

### cacheKey

• `Private` **cacheKey**: `null` \| `string` = `null`

#### Defined in

[flow.ts:60](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L60)

---

### checkResponse

• `Private` **checkResponse**: `null` \| `CheckResponse__Output` = `null`

#### Defined in

[flow.ts:61](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L61)

---

### controlPoint

• `Private` **controlPoint**: `string`

#### Defined in

[flow.ts:56](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L56)

---

### ended

• `Private` **ended**: `boolean` = `false`

#### Defined in

[flow.ts:50](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L50)

---

### error

• `Private` **error**: `null` \| `Error` = `null`

#### Defined in

[flow.ts:62](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L62)

---

### fcsClient

• `Private` **fcsClient**: `FlowControlServiceClient`

#### Defined in

[flow.ts:54](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L54)

---

### grpcCallOptions

• `Private` **grpcCallOptions**: `CallOptions`

#### Defined in

[flow.ts:55](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L55)

---

### rampMode

• `Private` **rampMode**: `boolean` = `false`

#### Defined in

[flow.ts:59](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L59)

---

### span

• `Private` **span**: `Span`

#### Defined in

[flow.ts:57](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L57)

---

### startDate

• `Private` **startDate**: `number`

#### Defined in

[flow.ts:58](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L58)

---

### status

• `Private` **status**: [`FlowStatus`](../README.md#flowstatus) =
`FlowStatusEnum.OK`

#### Defined in

[flow.ts:51](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L51)

## Methods

### CachedValue

▸ **CachedValue**(): [`CachedValueResponse`](CachedValueResponse.md)

Gets the cached value for the flow.

#### Returns

[`CachedValueResponse`](CachedValueResponse.md)

The cached value response.

#### Defined in

[flow.ts:174](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L174)

---

### CheckResponse

▸ **CheckResponse**(): `null` \| `CheckResponse__Output`

Gets the check response of the flow.

#### Returns

`null` \| `CheckResponse__Output`

The check response object.

#### Defined in

[flow.ts:208](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L208)

---

### DeleteCachedValue

▸ **DeleteCachedValue**(): `Promise`\<`undefined` \|
[`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)\>

Deletes the cached value for the flow.

#### Returns

`Promise`\<`undefined` \|
[`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)\>

A promise that resolves to the response of the cache delete operation.

#### Defined in

[flow.ts:135](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L135)

---

### End

▸ **End**(): `void`

Ends the flow and performs necessary cleanup.

#### Returns

`void`

#### Defined in

[flow.ts:223](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L223)

---

### Error

▸ **Error**(): `null` \| `Error`

Gets the error associated with the flow.

#### Returns

`null` \| `Error`

The error object.

#### Defined in

[flow.ts:200](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L200)

---

### SetCachedValue

▸ **SetCachedValue**(`value`, `ttl`): `Promise`\<`undefined` \|
[`SetCachedValueResponse`](SetCachedValueResponse.md)\>

Sets the cached value for the flow.

#### Parameters

| Name    | Type       | Description                            |
| :------ | :--------- | :------------------------------------- |
| `value` | `Buffer`   | The value to set.                      |
| `ttl`   | `Duration` | The time-to-live for the cached value. |

#### Returns

`Promise`\<`undefined` \|
[`SetCachedValueResponse`](SetCachedValueResponse.md)\>

A promise that resolves to the response of the cache upsert operation.

#### Defined in

[flow.ts:99](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L99)

---

### SetStatus

▸ **SetStatus**(`status`): `void`

Sets the status of the flow.

#### Parameters

| Name     | Type                                    | Description        |
| :------- | :-------------------------------------- | :----------------- |
| `status` | [`FlowStatus`](../README.md#flowstatus) | The status to set. |

#### Returns

`void`

#### Defined in

[flow.ts:89](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L89)

---

### ShouldRun

▸ **ShouldRun**(): `boolean`

Determines whether the flow should run based on the check response and ramp
mode.

#### Returns

`boolean`

A boolean value indicating whether the flow should run.

#### Defined in

[flow.ts:73](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L73)

---

### Span

▸ **Span**(): `Span`

Gets the span associated with the flow.

#### Returns

`Span`

The span object.

#### Defined in

[flow.ts:216](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L216)

---

### protoDurationToJSON

▸ **protoDurationToJSON**(`duration`): `null` \| `string`

#### Parameters

| Name       | Type                         |
| :--------- | :--------------------------- |
| `duration` | `null` \| `Duration__Output` |

#### Returns

`null` \| `string`

#### Defined in

[flow.ts:282](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L282)

---

### protoTimestampToJSON

▸ **protoTimestampToJSON**(`timestamp`): `null` \| `string`

#### Parameters

| Name        | Type                          |
| :---------- | :---------------------------- |
| `timestamp` | `null` \| `Timestamp__Output` |

#### Returns

`null` \| `string`

#### Defined in

[flow.ts:271](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/flow.ts#L271)
