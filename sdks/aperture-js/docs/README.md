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

Ƭ **FlowStatus**: typeof [`FlowStatusEnum`](README.md#flowstatusenum)[keyof typeof [`FlowStatusEnum`](README.md#flowstatusenum)]

Represents the status of a flow.

## Variables

### FlowStatusEnum

• `Const` **FlowStatusEnum**: `Object`

Enum representing the status of a flow.

#### Type declaration

| Name | Type |
| :------ | :------ |
| `Error` | ``"Error"`` |
| `OK` | ``"OK"`` |

## Functions

### ConvertCacheError

▸ **ConvertCacheError**(`error`): `Error` \| ``null``

Converts a cache error string into an Error object.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `error` | `undefined` \| `string` | The cache error string. |

#### Returns

`Error` \| ``null``

The Error object representing the cache error, or null if the error string is empty.

___

### ConvertCacheLookupStatus

▸ **ConvertCacheLookupStatus**(`status`): [`LookupStatus`](enums/LookupStatus.md)

Converts the cache lookup status to a lookup status.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `status` | `undefined` \| ``null`` \| `CacheLookupStatus` | The cache lookup status to convert. |

#### Returns

[`LookupStatus`](enums/LookupStatus.md)

The converted lookup status.

___

### ConvertCacheOperationStatus

▸ **ConvertCacheOperationStatus**(`status`): [`OperationStatus`](enums/OperationStatus.md)

Converts a cache operation status to an operation status.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `status` | `undefined` \| `CacheOperationStatus` | The cache operation status to convert. |

#### Returns

[`OperationStatus`](enums/OperationStatus.md)

The converted operation status.
