[@fluxninja/aperture-js](../README.md) / [Exports](../modules.md) /
SetCachedValueResponse

# Class: SetCachedValueResponse

Represents the response of setting a cached value.

## Table of contents

### Constructors

- [constructor](SetCachedValueResponse.md#constructor)

### Properties

- [error](SetCachedValueResponse.md#error)
- [operationStatus](SetCachedValueResponse.md#operationstatus)

### Methods

- [GetError](SetCachedValueResponse.md#geterror)
- [GetOperationStatus](SetCachedValueResponse.md#getoperationstatus)

## Constructors

### constructor

• **new SetCachedValueResponse**(`error`, `operationStatus`):
[`SetCachedValueResponse`](SetCachedValueResponse.md)

Creates a new instance of SetCachedValueResponse.

#### Parameters

| Name              | Type                                             | Description                                           |
| :---------------- | :----------------------------------------------- | :---------------------------------------------------- |
| `error`           | `null` \| `Error`                                | The error that occurred during the operation, if any. |
| `operationStatus` | [`OperationStatus`](../enums/OperationStatus.md) | The status of the operation.                          |

#### Returns

[`SetCachedValueResponse`](SetCachedValueResponse.md)

#### Defined in

[cache.ts:141](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/cache.ts#L141)

## Properties

### error

• **error**: `null` \| `Error`

#### Defined in

[cache.ts:133](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/cache.ts#L133)

---

### operationStatus

• **operationStatus**: [`OperationStatus`](../enums/OperationStatus.md)

#### Defined in

[cache.ts:134](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/cache.ts#L134)

## Methods

### GetError

▸ **GetError**(): `null` \| `Error`

Gets the error that occurred during the operation.

#### Returns

`null` \| `Error`

The error that occurred during the operation, or null if no error occurred.

#### Defined in

[cache.ts:150](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/cache.ts#L150)

---

### GetOperationStatus

▸ **GetOperationStatus**(): [`OperationStatus`](../enums/OperationStatus.md)

Gets the status of the operation.

#### Returns

[`OperationStatus`](../enums/OperationStatus.md)

The status of the operation.

#### Defined in

[cache.ts:158](https://github.com/fluxninja/aperture/blob/c4fc8958b/sdks/aperture-js/sdk/cache.ts#L158)
