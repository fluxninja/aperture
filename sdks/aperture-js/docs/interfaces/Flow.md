[@fluxninja/aperture-js](../README.md) / Flow

# Interface: Flow

## Table of contents

### Methods

- [checkResponse](Flow.md#checkresponse)
- [deleteGlobalCache](Flow.md#deleteglobalcache)
- [deleteResultCache](Flow.md#deleteresultcache)
- [end](Flow.md#end)
- [error](Flow.md#error)
- [globalCache](Flow.md#globalcache)
- [resultCache](Flow.md#resultcache)
- [setGlobalCache](Flow.md#setglobalcache)
- [setResultCache](Flow.md#setresultcache)
- [setStatus](Flow.md#setstatus)
- [shouldRun](Flow.md#shouldrun)
- [span](Flow.md#span)

## Methods

### checkResponse

▸ **checkResponse**(): ``null`` \| `CheckResponse__Output`

#### Returns

``null`` \| `CheckResponse__Output`

___

### deleteGlobalCache

▸ **deleteGlobalCache**(`key`, `grpcOptions?`): `Promise`\<[`KeyDeleteResponse`](KeyDeleteResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `key` | `string` |
| `grpcOptions?` | `CallOptions` |

#### Returns

`Promise`\<[`KeyDeleteResponse`](KeyDeleteResponse.md)\>

___

### deleteResultCache

▸ **deleteResultCache**(`grpcOptions?`): `Promise`\<`undefined` \| [`KeyDeleteResponse`](KeyDeleteResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `grpcOptions?` | `CallOptions` |

#### Returns

`Promise`\<`undefined` \| [`KeyDeleteResponse`](KeyDeleteResponse.md)\>

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

### globalCache

▸ **globalCache**(`key`): [`KeyLookupResponse`](KeyLookupResponse.md)

#### Parameters

| Name | Type |
| :------ | :------ |
| `key` | `string` |

#### Returns

[`KeyLookupResponse`](KeyLookupResponse.md)

___

### resultCache

▸ **resultCache**(): [`KeyLookupResponse`](KeyLookupResponse.md)

#### Returns

[`KeyLookupResponse`](KeyLookupResponse.md)

___

### setGlobalCache

▸ **setGlobalCache**(`key`, `cacheEntry`, `grpcOptions?`): `Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `key` | `string` |
| `cacheEntry` | [`CacheEntry`](CacheEntry.md) |
| `grpcOptions?` | `CallOptions` |

#### Returns

`Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

___

### setResultCache

▸ **setResultCache**(`cacheEntry`, `grpcOptions?`): `Promise`\<[`KeyUpsertResponse`](KeyUpsertResponse.md)\>

#### Parameters

| Name | Type |
| :------ | :------ |
| `cacheEntry` | [`CacheEntry`](CacheEntry.md) |
| `grpcOptions?` | `CallOptions` |

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
