## Aperture Requirements

Aperture provides Flow Control as a Service for reliably running modern web-scale cloud applications with little or no impact on the application code itself. Mostly all modern web-scale cloud applications can be deployed in an Aperture-enabled cluster without any changes at all. This document describes some application considerations and specific requirements of Aperture enablement.

### Hardware Resource Utilization Metrics

We support two modes of deployment for Aperture Agent:

- Sidecar
- DaemonSet

| Resource | Sidecar  | DaemonSet |
| -------- | -------- | --------- |
| CPU      | 0.5 vCPU | 1 vCPU    |
| Memory   | 0.5 Gi   | 1 Gi      |

**Note**: The above metrics may vary based on the load received on the Agent.

### Ports used by Aperture

The following ports and protocols are used by the Aperture Agent.

| Port | Protocol | Description                                      | Pod-internal only | Configurable |
| ---- | -------- | ------------------------------------------------ | ----------------- | ------------ |
| 80   | TCP      | Socket listener                                  | No                | Yes          |
| 4317 | TCP      | OpenTelemetry collector port                     | No                | No           |
| 3320 | TCP      | Distributed Cache                                | Yes               | Yes          |
| 3322 | TCP      | Memberlist to advertise to other cluster members | Yes               | Yes          |

The following ports and protocols are used by the Aperture Controller.

| Port | Protocol | Description        | Pod-internal only | Configurable |
| ---- | -------- | ------------------ | ----------------- | ------------ |
| 80   | TCP      | Socket listener    | No                | Yes          |
| 8086 | TCP      | Validating webhook | No                | Yes          |

### Integration

Aperture leverages recent advancements and standards for tracing (e.g., [Open Tracing](https://opentracing.io)) for observability and control. Aperture provides plugins for popular traffic proxies and gateways; and libraries that work with popular languages & frameworks. Please contact [Support](mailto:support@fluxninja.com) for getting started with custom integration.
