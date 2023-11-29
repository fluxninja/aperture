[@fluxninja/aperture-js](../README.md) / FlowParams

# Interface: FlowParams

Represents the parameters for a flow.

## Table of contents

### Properties

- [globalCacheKeys](FlowParams.md#globalcachekeys)
- [grpcCallOptions](FlowParams.md#grpccalloptions)
- [labels](FlowParams.md#labels)
- [rampMode](FlowParams.md#rampmode)
- [resultCacheKey](FlowParams.md#resultcachekey)
- [tryConnect](FlowParams.md#tryconnect)

## Properties

### globalCacheKeys

• `Optional` **globalCacheKeys**: `string`[]

Keys to global cache entries that need to be fetched at flow start.

___

### grpcCallOptions

• `Optional` **grpcCallOptions**: `CallOptions`

Additional gRPC call options for the flow.

___

### labels

• `Optional` **labels**: `Record`\<`string`, `string`\>

Optional labels for the flow.

___

### rampMode

• `Optional` **rampMode**: `boolean`

Specifies whether the flow should use ramp mode.

___

### resultCacheKey

• `Optional` **resultCacheKey**: `string`

Key to the result cache entry which needs to be fetched at flow start.

___

### tryConnect

• `Optional` **tryConnect**: `boolean`

Specifies whether to try connecting to the flow.
