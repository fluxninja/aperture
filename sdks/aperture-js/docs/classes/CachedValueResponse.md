[@fluxninja/aperture-js](../README.md) / CachedValueResponse

# Class: CachedValueResponse

Represents a response from a cached value lookup.

## Table of contents

### Constructors

- [constructor](CachedValueResponse.md#constructor)

### Properties

- [error](CachedValueResponse.md#error)
- [lookupStatus](CachedValueResponse.md#lookupstatus)
- [operationStatus](CachedValueResponse.md#operationstatus)
- [value](CachedValueResponse.md#value)

### Methods

- [GetError](CachedValueResponse.md#geterror)
- [GetLookupStatus](CachedValueResponse.md#getlookupstatus)
- [GetOperationStatus](CachedValueResponse.md#getoperationstatus)
- [GetValue](CachedValueResponse.md#getvalue)

## Constructors

### constructor

• **new CachedValueResponse**(`lookupStatus`, `operationStatus`, `error`, `value`): [`CachedValueResponse`](CachedValueResponse.md)

Creates a new CachedValueResponse instance.

#### Parameters

| Name | Type | Description |
| :------ | :------ | :------ |
| `lookupStatus` | [`LookupStatus`](../enums/LookupStatus.md) | The lookup status. |
| `operationStatus` | [`OperationStatus`](../enums/OperationStatus.md) | The operation status. |
| `error` | ``null`` \| `Error` | The error, if any. |
| `value` | ``null`` \| `Buffer` | The cached value, if any. |

#### Returns

[`CachedValueResponse`](CachedValueResponse.md)

## Properties

### error

• **error**: ``null`` \| `Error`

___

### lookupStatus

• **lookupStatus**: [`LookupStatus`](../enums/LookupStatus.md)

___

### operationStatus

• **operationStatus**: [`OperationStatus`](../enums/OperationStatus.md)

___

### value

• **value**: ``null`` \| `Buffer`

## Methods

### GetError

▸ **GetError**(): ``null`` \| `Error`

Gets the error, if any.

#### Returns

``null`` \| `Error`

The error, or null if no error occurred.

___

### GetLookupStatus

▸ **GetLookupStatus**(): [`LookupStatus`](../enums/LookupStatus.md)

Gets the lookup status.

#### Returns

[`LookupStatus`](../enums/LookupStatus.md)

The lookup status.

___

### GetOperationStatus

▸ **GetOperationStatus**(): [`OperationStatus`](../enums/OperationStatus.md)

Gets the operation status.

#### Returns

[`OperationStatus`](../enums/OperationStatus.md)

The operation status.

___

### GetValue

▸ **GetValue**(): ``null`` \| `Buffer`

Gets the cached value, if any.

#### Returns

``null`` \| `Buffer`

The cached value, or null if no value is available.
