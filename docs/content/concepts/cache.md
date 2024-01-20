---
title: Cache
sidebar_position: 9
---

Aperture's _Cache_ can be used to reduce the load on a service by caching the
results of expensive computations. It is built on top of a reliable distributed
cache so that each Aperture Agent can quickly and efficiently operate on the
cached data.

The _Cache_ functionality can be used via [Aperture SDKs][skds] as part of the
[Flow][flow-label]. While some SDKs might differ in their implementation and
function names, the general idea is as follows:

1. Create an instance of Aperture `Client`.
2. Instantiate a `Flow` by calling the `StartFlow` method with the
   `resultCacheKey` parameter set to your desired value. The first call will let
   Aperture initialize a cache entry for the flow, uniquely identified by the
   `ControlPoint` and `ResultCacheKey` values. Subsequent calls will return the
   cached value as part of the response object.
3. The value stored in the cache can be retrieved by calling the `ResultCache`
   method on the `Flow` object. It returns an object that can be used to perform
   the following operations:

   - `GetLookupStatus` - retrieve the status of the lookup operation, whether it
     was a `HIT` or a `MISS`.
   - `GetError` - retrieve the error that occurred during the lookup operation.
   - `GetValue` - retrieve the cached value.

4. While the flow is active, the following additional _Cache_ related operations
   can be performed on the `Flow` object:

   - `SetResultCache` - set the cached value.
   - `DeleteResultCache` - delete the cached value.

5. Call the `EndFlow` method to complete the flow.

:::info

Refer to the [Caching Guide][guide] for more information on how to use the
_Cache_ via [`aperture-js`][aperture-js] SDK.

:::

Additionally, Aperture also provides a _Global Cache_ that can be used to access
cached data across multiple flows. This is useful when the same data is required
by multiple flows. The _Global Cache_ can be accessed via the `GlobalCache`
method with `key` parameter on the `Flow` object. It returns the same object as
the `ResultCache` method.

Similar to `ResultCache`, `Set` and `Delete` methods can be used to set and
delete entries in the _Global Cache_.

[skds]: /sdk/sdk.md
[flow-label]: /concepts/flow-label.md
[guide]: /guides/caching.md
[aperture-js]: https://github.com/fluxninja/aperture-js
