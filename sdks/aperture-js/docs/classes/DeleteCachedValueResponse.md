[@fluxninja/aperture-js](../README.md) / DeleteCachedValueResponse

# Class: DeleteCachedValueResponse

Represents the response of deleting a cached value.

## Table of contents

### Constructors

- [constructor](DeleteCachedValueResponse.md#constructor)

### Properties

- [error](DeleteCachedValueResponse.md#error)
- [operationStatus](DeleteCachedValueResponse.md#operationstatus)

### Methods

- [GetError](DeleteCachedValueResponse.md#geterror)
- [GetOperationStatus](DeleteCachedValueResponse.md#getoperationstatus)

## Constructors

### constructor

• **new DeleteCachedValueResponse**(`error`, `operationStatus`):
[`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)

Creates a new instance of DeleteCachedValueResponse.

#### Parameters

| Name              | Type                                             | Description                                                  |
| :---------------- | :----------------------------------------------- | :----------------------------------------------------------- |
| `error`           | `null` \| `Error`                                | The error that occurred during the delete operation, if any. |
| `operationStatus` | [`OperationStatus`](../enums/OperationStatus.md) | The status of the delete operation.                          |

#### Returns

[`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)

#### Defined in

[cache.ts:175](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L175)

## Properties

### error

• **error**: `null` \| `Error`

#### Defined in

[cache.ts:167](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L167)

---

### operationStatus

• **operationStatus**: [`OperationStatus`](../enums/OperationStatus.md)

#### Defined in

[cache.ts:168](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L168)

## Methods

### GetError

▸ **GetError**(): `null` \| `Error`

Gets the error that occurred during the delete operation, if any.

#### Returns

`null` \| `Error`

The error object or null if no error occurred.

#### Defined in

[cache.ts:184](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L184)

---

### GetOperationStatus

▸ **GetOperationStatus**(): [`OperationStatus`](../enums/OperationStatus.md)

Gets the status of the delete operation.

#### Returns

[`OperationStatus`](../enums/OperationStatus.md)

The operation status.

#### Defined in

[cache.ts:192](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/cache.ts#L192)
