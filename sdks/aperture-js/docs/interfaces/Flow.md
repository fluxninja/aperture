[@fluxninja/aperture-js](../README.md) / Flow

# Interface: Flow

## Table of contents

### Methods

- [checkResponse](Flow.md#checkresponse)
- [deleteResultCache](Flow.md#deleteresultcache)
- [deleteStateCache](Flow.md#deletestatecache)
- [end](Flow.md#end)
- [error](Flow.md#error)
- [resultCache](Flow.md#resultcache)
- [setResultCache](Flow.md#setresultcache)
- [setStateCache](Flow.md#setstatecache)
- [setStatus](Flow.md#setstatus)
- [shouldRun](Flow.md#shouldrun)
- [span](Flow.md#span)
- [stateCache](Flow.md#statecache)

## Methods

### checkResponse

▸ **checkResponse**(): ``null`` \| `CheckResponse__Output`

#### Returns

``null`` \| `CheckResponse__Output`

___

### deleteResultCache

▸ **deleteResultCache**(): `Promise`\<`undefined` \| [`KeyDeleteResponse`](KeyDeleteResponse.md)\>

#### Returns

`Promise`\<`undefined` \| [`KeyDeleteResponse`](KeyDeleteResponse.md)\>

___

### deleteStateCache

▸ **deleteStateCache**(`key`): `Promise`\<[`KeyDeleteResponse`](KeyDeleteResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `key` | `string` |

#### Returns

`Promise`\<[`KeyDeleteResponse`](KeyDeleteResponse.md)\>

___

### end

▸ **end**(): `void`

#### Returns

`void`

___

### error

▸ **error**(): ``null`` \| `Error`

#### Returns

``null`` \| `Error`

___

### resultCache

▸ **resultCache**(): [`KeyLookupResponse`](KeyLookupResponse.md)

#### Returns

[`KeyLookupResponse`](KeyLookupResponse.md)

___

### setResultCache

▸ **setResultCache**(`cacheEntry`): `Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `cacheEntry` | [`CacheEntry`](CacheEntry.md) |

#### Returns

`Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

___

### setStateCache

▸ **setStateCache**(`key`, `cacheEntry`): `Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `key` | `string` |
| `cacheEntry` | [`CacheEntry`](CacheEntry.md) |

#### Returns

`Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

___

### setStatus

▸ **setStatus**(`status`): `void`

#### Parameters

| Name | Type |
| :------ | :------ |
| `status` | [`FlowStatus`](../enums/FlowStatus.md) |

#### Returns

`void`

___

### shouldRun

▸ **shouldRun**(): `boolean`

#### Returns

`boolean`

___

### span

▸ **span**(): `Span`

#### Returns

`Span`

___

### stateCache

▸ **stateCache**(`key`): [`KeyLookupResponse`](KeyLookupResponse.md)

#### Parameters

| Name | Type |
| :------ | :------ |
| `key` | `string` |

#### Returns

[`KeyLookupResponse`](KeyLookupResponse.md)
