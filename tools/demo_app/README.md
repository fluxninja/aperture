# Demo App

Simple service is a test service that can be composed to create a whole mesh of
micro services in order to create a test cluster and mimic real world scenarios
in a fully controlled test environment.

The service has one endpoint and depending on the incoming payload, it will
either forward the request onwards to one of more upstream services or just
respond to the request. The request payload has a list of requests where each
request is a chain of service destinations. The first element in each list
should be the service itself, otherwise the request was routed incorrectly and
an error should be returned.

## How it works

Service receives the "request chain" description as JSON on `POST /request`
endpoint and forwards the traffic to all the downstream services. It then waits
for responses and if all services responded, returns 200 OK.

### Request chain spec

The request chain description should satisfy the folllowing
[json schema](https://json-schema.org/):

```json
{
  "$schema": "http://json-schema.org/schema#",
  "type": "object",
  "properties": {
    "request": {
      "type": "array",
      "items": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "destination": { "type": "string" }
          },
          "required": ["destination"],
          "additionalProperties": false
        }
      }
    }
  },
  "required": ["request"]
}
```

#### Destination

Destination should be a valid HTTP hostname (if the port is different than 80,
the destination should also include port), examples:

- `simple-srv-a.fn-demo-app-cluster.svc.cluster.local`
- `127.0.0.1:8090`

> Note: The destination of the first service in a chain is not used as an actual
> destination, but should **match exactly** the HOSTNAME env. This is a sanity
> check to verify that the correct service is initiating the chain.

#### Examples

Let's assume three services with hostnames A, B and C.

1. Single request (no chain)

   ```json
   { "request": [[{ "destination": "A" }]] }
   ```

2. A→B→C chain:

   ```json
   {
     "request": [
       [{ "destination": "A" }, { "destination": "B" }, { "destination": "C" }]
     ]
   }
   ```

3. Multiple chains:

   ```text
   A → B
   │   ↓
   └─→ C
   ```

   ```json
   {
     "request": [
       [{ "destination": "A" }, { "destination": "B" }, { "destination": "C" }],
       [{ "destination": "A" }, { "destination": "C" }]
     ]
   }
   ```

## Configuration

The app reads the following env variables:

- `SIMPLE_SERVICE_CONCURRENCY` - number of concurrent requests to process.
- `SIMPLE_SERVICE_PORT` – port to listen on.
- `ENVOY_EGRESS_PORT` – port to be used as http proxy for egress traffic.
- `HOSTNAME` – the services own hostname. This is used as a sanity check to
  verify validity of request chain.
