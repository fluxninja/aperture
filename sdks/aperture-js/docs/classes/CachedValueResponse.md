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

• **new CachedValueResponse**(`lookupStatus`, `operationStatus`, `error`,
`value`): [`CachedValueResponse`](CachedValueResponse.md)

Creates a new CachedValueResponse instance.

#### Parameters

| Name              | Type                                             | Description               |
| :---------------- | :----------------------------------------------- | :------------------------ |
| `lookupStatus`    | [`LookupStatus`](../enums/LookupStatus.md)       | The lookup status.        |
| `operationStatus` | [`OperationStatus`](../enums/OperationStatus.md) | The operation status.     |
| `error`           | `null` \| `Error`                                | The error, if any.        |
| `value`           | `null` \| `Buffer`                               | The cached value, if any. |

#### Returns

[`CachedValueResponse`](CachedValueResponse.md)

#### Defined in

[cache.ts:84](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L84)

## Properties

### error

• **error**: `null` \| `Error`

#### Defined in

[cache.ts:74](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L74)

---

### lookupStatus

• **lookupStatus**: [`LookupStatus`](../enums/LookupStatus.md)

#### Defined in

[cache.ts:72](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L72)

---

### operationStatus

• **operationStatus**: [`OperationStatus`](../enums/OperationStatus.md)

#### Defined in

[cache.ts:73](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L73)

---

### value

• **value**: `null` \| `Buffer`

#### Defined in

[cache.ts:75](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L75)

## Methods

### GetError

▸ **GetError**(): `null` \| `Error`

Gets the error, if any.

#### Returns

`null` \| `Error`

The error, or null if no error occurred.

#### Defined in

[cache.ts:116](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L116)

---

### GetLookupStatus

▸ **GetLookupStatus**(): [`LookupStatus`](../enums/LookupStatus.md)

Gets the lookup status.

#### Returns

[`LookupStatus`](../enums/LookupStatus.md)

The lookup status.

#### Defined in

[cache.ts:100](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L100)

---

### GetOperationStatus

▸ **GetOperationStatus**(): [`OperationStatus`](../enums/OperationStatus.md)

Gets the operation status.

#### Returns

[`OperationStatus`](../enums/OperationStatus.md)

The operation status.

#### Defined in

[cache.ts:108](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L108)

---

### GetValue

▸ **GetValue**(): `null` \| `Buffer`

Gets the cached value, if any.

#### Returns

`null` \| `Buffer`

The cached value, or null if no value is available.

#### Defined in

[cache.ts:124](https://github.com/fluxninja/aperture/blob/a92f6b393/sdks/aperture-js/sdk/cache.ts#L124)
