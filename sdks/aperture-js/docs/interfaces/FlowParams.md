[@fluxninja/aperture-js](../README.md) / FlowParams

# Interface: FlowParams

Represents the parameters for a flow.

## Table of contents

### Properties

- [cacheKey](FlowParams.md#cachekey)
- [grpcCallOptions](FlowParams.md#grpccalloptions)
- [labels](FlowParams.md#labels)
- [rampMode](FlowParams.md#rampmode)
- [tryConnect](FlowParams.md#tryconnect)

## Properties

### cacheKey

• `Optional` **cacheKey**: `string`

The cache key for the flow.

#### Defined in

[client.ts:42](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/client.ts#L42)

---

### grpcCallOptions

• `Optional` **grpcCallOptions**: `CallOptions`

Additional gRPC call options for the flow.

#### Defined in

[client.ts:34](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/client.ts#L34)

---

### labels

• `Optional` **labels**: `Record`\<`string`, `string`\>

Optional labels for the flow.

#### Defined in

[client.ts:26](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/client.ts#L26)

---

### rampMode

• `Optional` **rampMode**: `boolean`

Specifies whether the flow should use ramp mode.

#### Defined in

[client.ts:30](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/client.ts#L30)

---

### tryConnect

• `Optional` **tryConnect**: `boolean`

Specifies whether to try connecting to the flow.

#### Defined in

[client.ts:38](https://github.com/fluxninja/aperture/blob/5ab1329aa/sdks/aperture-js/sdk/client.ts#L38)
