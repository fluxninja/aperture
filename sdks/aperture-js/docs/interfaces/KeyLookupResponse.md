[@fluxninja/aperture-js](../README.md) / KeyLookupResponse

# Interface: KeyLookupResponse

Represents a cache value lookup.

## Table of contents

### Methods

- [getError](KeyLookupResponse.md#geterror)
- [getLookupStatus](KeyLookupResponse.md#getlookupstatus)
- [getValue](KeyLookupResponse.md#getvalue)

## Methods

### getError

▸ **getError**(): ``null`` \| `Error`

Gets the error, if any.

#### Returns

``null`` \| `Error`

The error, or null if no error occurred.

___

### getLookupStatus

▸ **getLookupStatus**(): [`LookupStatus`](../enums/LookupStatus.md)

Gets the lookup status.

#### Returns

[`LookupStatus`](../enums/LookupStatus.md)

The lookup status.

___

### getValue

▸ **getValue**(): ``null`` \| `Buffer`

Gets the cached value, if any.

#### Returns

``null`` \| `Buffer`

The cached value, or null if no value is available.
