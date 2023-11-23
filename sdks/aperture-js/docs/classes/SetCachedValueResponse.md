[@fluxninja/aperture-js](../README.md) / SetCachedValueResponse

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

• **new SetCachedValueResponse**(`error`, `operationStatus`): [`SetCachedValueResponse`](SetCachedValueResponse.md)

Creates a new instance of SetCachedValueResponse.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `error` | ``null`` \| `Error` | The error that occurred during the operation, if any. |
| `operationStatus` | [`OperationStatus`](../enums/OperationStatus.md) | The status of the operation. |

#### Returns

[`SetCachedValueResponse`](SetCachedValueResponse.md)

## Properties

### error

• **error**: ``null`` \| `Error`

___

### operationStatus

• **operationStatus**: [`OperationStatus`](../enums/OperationStatus.md)

## Methods

### GetError

▸ **GetError**(): ``null`` \| `Error`

Gets the error that occurred during the operation.

#### Returns

``null`` \| `Error`

The error that occurred during the operation, or null if no error occurred.

___

### GetOperationStatus

▸ **GetOperationStatus**(): [`OperationStatus`](../enums/OperationStatus.md)

Gets the status of the operation.

#### Returns

[`OperationStatus`](../enums/OperationStatus.md)

The status of the operation.
