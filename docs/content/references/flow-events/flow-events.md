---
title: Flow Events
sidebar_position: 1
sidebar_label: Flow Events
---

## Dimension Columns

### Common

| Name                                   | Type   | Example Values | Description                                | Flow Control Integrations |
| -------------------------------------- | ------ | -------------- | ------------------------------------------ | ------------------------- |
| aperture.source                        | single | sdk, envoy     | Aperture Flow source                       | SDKs, Envoy               |
| workload_duration_ms                   |        |                | Duration of the workload in ms             | SDKs, Envoy               |
| flow_duration_ms                       |        |                | Duration of the flow in ms                 | SDKs, Envoy               |
| aperture_processing_duration_ms        | single |                |                                            | SDKs, Envoy               |
| aperture.decision_type                 |        |                |                                            | SDKs, Envoy               |
| aperture.error                         |        |                |                                            | SDKs, Envoy               |
| aperture.reject_reason                 |        |                |                                            | SDKs, Envoy               |
| aperture.rate_limiters                 |        |                |                                            | SDKs, Envoy               |
| aperture.dropping_rate_limiters        |        |                |                                            | SDKs, Envoy               |
| aperture.concurrency_limiters          |        |                |                                            | SDKs, Envoy               |
| aperture.dropping_concurrency_limiters |        |                |                                            | SDKs, Envoy               |
| aperture.workloads                     |        |                |                                            | SDKs, Envoy               |
| aperture.dropping_workloads            |        |                |                                            | SDKs, Envoy               |
| aperture.flux_meters                   |        |                |                                            | SDKs, Envoy               |
| aperture.flow_label_keys               |        |                |                                            | SDKs, Envoy               |
| aperture.classifiers                   |        |                |                                            | SDKs, Envoy               |
| aperture.classifier_errors             |        |                |                                            | SDKs, Envoy               |
| aperture.services                      |        |                |                                            | SDKs, Envoy               |
| aperture.control_point                 |        |                |                                            | SDKs, Envoy               |
| aperture.response_status               | single | ok, error      | Denotes OK or Error across all protocols   | SDKs, Envoy               |
| response_received                      |        |                | Designates whether a response was received | SDKs, envoy               |

### HTTP

| Name                         | Type | Example Values | Description | Flow Control Integrations |
| ---------------------------- | ---- | -------------- | ----------- | ------------------------- |
| http.status_code             |      |                |             | Envoy                     |
| http.request_content_length  |      |                |             | Envoy                     |
| http.response_content_length |      |                |             | Envoy                     |
| http.method                  |      |                |             | Envoy                     |
| http.target                  |      |                |             | Envoy                     |
| http.host                    |      |                |             | Envoy                     |
| http.scheme                  |      |                |             | Envoy                     |
| http.request_content_length  |      |                |             | Envoy                     |
| http.flavor                  |      |                |             | Envoy                     |

### SDK

| Name                    | Type | Example Values | Description | Flow Control Integrations |
| ----------------------- | ---- | -------------- | ----------- | ------------------------- |
| aperture.feature.status |      |                |             | SDKs                      |

## Metric Columns

| Name                            | Type | Rollup Type | Unit  | Description                         |
| ------------------------------- | ---- | ----------- | ----- | ----------------------------------- |
| workload_duration_ms            |      |             | ms    | Duration of the workload            |
| flow_duration_ms                |      |             | ms    | Duration of the flow                |
| aperture_processing_duration_ms |      |             | ms    | Aperture's processing duration      |
| http.request_content_length     |      |             | bytes | Length of the HTTP request content  |
| http.response_content_length    |      |             | bytes | Length of the HTTP response content |
