# Rollup processor

Supported pipeline types: traces, logs.

The rollup processor accepts traces or logs and rolls them up, based on attributes.

It is highly recommended to configure batch processor just before rollup processor.
Rollup processor rolls up traces or logs which are a part of one message.

## Configuration

```yaml
processors:
  batch/prerollup:
    timeout: 1s
    send_batch_size: 10000
  rollup:
    rollups:
      - from: body_size
        to: body_size_min
        type: min
      - from: body_size
        to: body_size_max
        type: max
      - from: body_size
        to: body_size_sum
        type: sum
```

| Field name     | Type       | Required | Description                                          |
| -------------- | ---------- | -------- | ---------------------------------------------------- |
| `rollups`      | `[]Rollup` | No       | Array of rollups. One per resulting rolled up field. |
| `rollups.from` | `string`   | Yes      | Field from which value to be rolled up is read.      |
| `rollups.to`   | `string`   | Yes      | Field to which rolled up value is written.           |
| `rollups.type` | `string`   | Yes      | Rollup type. See below for possible values.          |

### Rollup types

| Type name | From type | To type | Description            |
| --------- | --------- | ------- | ---------------------- |
| `min`     | `string`  | `int64` | Returns minimum value. |
| `max`     | `string`  | `int64` | Returns maximum value. |
| `sum`     | `string`  | `int64` | Returns sum of values. |

## Design

Rollup processor is basically merging maps. In order to merge two maps, they need
to have the same _key attribues_. Key attribute is an attribute which was not specified
as rollup attribute in the configuration i.e. was not specified in any of `rollups.from`
fields.

Merging of such maps results with a map with key attributes and rollup result attributes.
In addition, there is `rollup_count` attribute, which describes how many maps were
rolled up into this single result.
Please refer to the example below to fully understand.

### Example

Assume the rollup processor is configured with config given in the previous section.
On the input we have following maps (they are parts of either Traces of Logs).

```json
[
  { "path": "/foo", "user-agent": "fizz", "body_size": "5" },
  { "path": "/foo", "user-agent": "fizz", "body_size": "6" },
  { "path": "/bar", "user-agent": "fizz", "body_size": "7" }
]
```

This will result in:

```json
[
  {
    "path": "/foo",
    "user-agent": "fizz",
    "body_size_min": 5,
    "body_size_max": 6,
    "body_size_sum": 11,
    "rollup_count": 2
  },
  {
    "path": "/bar",
    "user-agent": "fizz",
    "body_size_min": 7,
    "body_size_max": 7,
    "body_size_sum": 7,
    "rollup_count": 1
  }
]
```
