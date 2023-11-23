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

• **new DeleteCachedValueResponse**(`error`, `operationStatus`): [`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)

Creates a new instance of DeleteCachedValueResponse.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `error` | ``null`` \| `Error` | The error that occurred during the delete operation, if any. |
| `operationStatus` | [`OperationStatus`](../enums/OperationStatus.md) | The status of the delete operation. |

#### Returns

[`DeleteCachedValueResponse`](DeleteCachedValueResponse.md)

## Properties

### error

• **error**: ``null`` \| `Error`

___

### operationStatus

• **operationStatus**: [`OperationStatus`](../enums/OperationStatus.md)

## Methods

### GetError

▸ **GetError**(): ``null`` \| `Error`

Gets the error that occurred during the delete operation, if any.

#### Returns

``null`` \| `Error`

The error object or null if no error occurred.

___

### GetOperationStatus

▸ **GetOperationStatus**(): [`OperationStatus`](../enums/OperationStatus.md)

Gets the status of the delete operation.

#### Returns

[`OperationStatus`](../enums/OperationStatus.md)

The operation status.
