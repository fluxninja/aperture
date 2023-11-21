@fluxninja/aperture-js

# @fluxninja/aperture-js

## Table of contents

### Enumerations

- [LookupStatus](enums/LookupStatus.md)
- [OperationStatus](enums/OperationStatus.md)

### Classes

- [ApertureClient](classes/ApertureClient.md)
- [CachedValueResponse](classes/CachedValueResponse.md)
- [DeleteCachedValueResponse](classes/DeleteCachedValueResponse.md)
- [Flow](classes/Flow.md)
- [SetCachedValueResponse](classes/SetCachedValueResponse.md)

### Interfaces

- [FlowParams](interfaces/FlowParams.md)

### Type Aliases

- [FlowStatus](README.md#flowstatus)

### Variables

- [FlowStatusEnum](README.md#flowstatusenum)

### Functions

- [ConvertCacheError](README.md#convertcacheerror)
- [ConvertCacheLookupStatus](README.md#convertcachelookupstatus)
- [ConvertCacheOperationStatus](README.md#convertcacheoperationstatus)

## Type Aliases

### FlowStatus

Ƭ **FlowStatus**: typeof [`FlowStatusEnum`](README.md#flowstatusenum)[keyof
typeof [`FlowStatusEnum`](README.md#flowstatusenum)]

Represents the status of a flow.

#### Defined in

[flow.ts:44](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/flow.ts#L44)

## Variables

### FlowStatusEnum

• `Const` **FlowStatusEnum**: `Object`

Enum representing the status of a flow.

#### Type declaration

| Name    | Type      |
| :------ | :-------- |
| `Error` | `"Error"` |
| `OK`    | `"OK"`    |

#### Defined in

[flow.ts:36](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/flow.ts#L36)

## Functions

### ConvertCacheError

▸ **ConvertCacheError**(`error`): `Error` \| `null`

Converts a cache error string into an Error object.

#### Parameters

| Name    | Type                    | Description             |
| :------ | :---------------------- | :---------------------- |
| `error` | `undefined` \| `string` | The cache error string. |

#### Returns

`Error` \| `null`

The Error object representing the cache error, or null if the error string is
empty.

#### Defined in

[cache.ts:61](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L61)

---

### ConvertCacheLookupStatus

▸ **ConvertCacheLookupStatus**(`status`):
[`LookupStatus`](enums/LookupStatus.md)

Converts the cache lookup status to a lookup status.

#### Parameters

| Name     | Type                                         | Description                         |
| :------- | :------------------------------------------- | :---------------------------------- |
| `status` | `undefined` \| `null` \| `CacheLookupStatus` | The cache lookup status to convert. |

#### Returns

[`LookupStatus`](enums/LookupStatus.md)

The converted lookup status.

#### Defined in

[cache.ts:17](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L17)

---

### ConvertCacheOperationStatus

▸ **ConvertCacheOperationStatus**(`status`):
[`OperationStatus`](enums/OperationStatus.md)

Converts a cache operation status to an operation status.

#### Parameters

| Name     | Type                                  | Description                            |
| :------- | :------------------------------------ | :------------------------------------- |
| `status` | `undefined` \| `CacheOperationStatus` | The cache operation status to convert. |

#### Returns

[`OperationStatus`](enums/OperationStatus.md)

The converted operation status.

#### Defined in

[cache.ts:43](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L43)
