[@fluxninja/aperture-js](../README.md) / FlowParams

# Interface: FlowParams

Represents the parameters for a flow.

## Table of contents

### Properties

- [grpcCallOptions](FlowParams.md#grpccalloptions)
- [labels](FlowParams.md#labels)
- [rampMode](FlowParams.md#rampmode)
- [resultCacheKey](FlowParams.md#resultcachekey)
- [stateCacheKeys](FlowParams.md#statecachekeys)
- [tryConnect](FlowParams.md#tryconnect)

## Properties

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

### stateCacheKeys

• `Optional` **stateCacheKeys**: `string`[]

Keys to state cache entries that need to be fetched at flow start.

___

### tryConnect

• `Optional` **tryConnect**: `boolean`

Specifies whether to try connecting to the flow.
