---
title: Flow Events
sidebar_position: 1
sidebar_label: Flow Events
---

## Dimension Columns

### Common

| Name                                   | Type | Example Values                                                                                                                                       | Description                                        | Flow Control Integrations |
| -------------------------------------- | ---- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------- | ------------------------- |
| aperture.source                        |      | sdk, envoy                                                                                                                                           | Aperture Flow source                               | SDKs, Envoy               |
| workload_duration_ms                   |      | 52                                                                                                                                                   | Duration of the workload in ms                     | SDKs, Envoy               |
| flow_duration_ms                       |      | 52                                                                                                                                                   | Duration of the flow in ms                         | SDKs, Envoy               |
| aperture_processing_duration_ms        |      | 52                                                                                                                                                   | Aperture's processing duration in ms               | SDKs, Envoy               |
| aperture.decision_type                 |      | DECISION_TYPE_ACCEPTED, DECISION_TYPE_REJECTED                                                                                                       | Decision type taken by policy                      | SDKs, Envoy               |
| aperture.error                         |      | ERROR_NONE, ERROR_MISSING_TRAFFIC_DIRECTION, ERROR_INVALID_TRAFFIC_DIRECTION, ERROR_CONVERT_TO_MAP_STRUCT, ERROR_CONVERT_TO_REGO_AST, ERROR_CLASSIFY | Error reason of the decision taken by policy       | SDKs, Envoy               |
| aperture.reject_reason                 |      | REJECT_REASON_NONE, REJECT_REASON_RATE_LIMITED, REJECT_REASON_CONCURRENCY_LIMITED                                                                    | Reject reason of the decision taken by policy      | SDKs, Envoy               |
| aperture.rate_limiters                 |      | policy_name:service1-demo-app,component_index:18,policy_hash:5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo=                                            | Rate limiters matched to the traffic               | SDKs, Envoy               |
| aperture.dropping_rate_limiters        |      | policy_name:service1-demo-app,component_index:18,policy_hash:5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo=                                            | Rate limiters dropping the traffic                 | SDKs, Envoy               |
| aperture.concurrency_limiters          |      | policy_name:service1-demo-app,component_index:13,policy_hash:5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo=                                            | Concurrency limiters matched to the traffic        | SDKs, Envoy               |
| aperture.dropping_concurrency_limiters |      | policy_name:service1-demo-app,component_index:13,policy_hash:5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo=                                            | Concurrency limiters dropping the traffic          | SDKs, Envoy               |
| aperture.workloads                     |      | policy_name:service1-demo-app,component_index:13,workload_index:0,policy_hash:5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo=                           | Workloads matched to the traffic                   | SDKs, Envoy               |
| aperture.dropping_workloads            |      | policy_name:service1-demo-app,component_index:13,workload_index:0,policy_hash:5kZjjSgDAtGWmLnDT67SmQhZdHVmz0+GvKcOGTfWMVo=                           | Workloads dropping the traffic                     | SDKs, Envoy               |
| aperture.flux_meters                   |      | service1-demo-app                                                                                                                                    | Flux Meters matched to the traffic                 | SDKs, Envoy               |
| aperture.flow_label_keys               |      | http.host, http.method, http.request.header.content_length                                                                                           | Flow labels matched to the traffic                 | SDKs, Envoy               |
| aperture.classifiers                   |      | policy_name:service1-demo-app,classifier_index:0                                                                                                     | Classifiers matched to the traffic                 | SDKs, Envoy               |
| aperture.classifier_errors             |      |                                                                                                                                                      | Encountered classifier errors for specified policy | SDKs, Envoy               |
| aperture.services                      |      | service1-demo-app.demoapp.svc.cluster.local, service2-demo-app.demoapp.svc.cluster.local                                                             | Services to which metrics refer                    | SDKs, Envoy               |
| aperture.control_point                 |      | type:TYPE_INGRESS, type:TYPE_EGRESS                                                                                                                  | Control point to which metrics refer               | SDKs, Envoy               |
| aperture.response_status               |      | OK, Error                                                                                                                                            | Denotes OK or Error across all protocols           | SDKs, Envoy               |
| response_received                      |      | true, false                                                                                                                                          | Designates whether a response was received         | SDKs, envoy               |

### HTTP

| Name                         | Type | Example Values                                                                           | Description                                  | Flow Control Integrations |
| ---------------------------- | ---- | ---------------------------------------------------------------------------------------- | -------------------------------------------- | ------------------------- |
| http.status_code             |      | 200, 429, 503                                                                            | HTTP status code of the response             | Envoy                     |
| http.request_content_length  |      | 0, 53                                                                                    | Length of the HTTP request content in bytes  | Envoy                     |
| http.response_content_length |      | 201, 77                                                                                  | Length of the HTTP response content in bytes | Envoy                     |
| http.method                  |      | GET, POST                                                                                | HTTP method of the response                  | Envoy                     |
| http.target                  |      | /request                                                                                 | Target endpoint of the response              | Envoy                     |
| http.host                    |      | service1-demo-app.demoapp.svc.cluster.local, service2-demo-app.demoapp.svc.cluster.local | Host address of the response                 | Envoy                     |
| http.scheme                  |      | http                                                                                     | HTTP scheme of the response                  | Envoy                     |
| http.flavor                  |      | 1.1                                                                                      | HTTP protocol version                        | Envoy                     |

### SDK

| Name                    | Type | Example Values | Description           | Flow Control Integrations |
| ----------------------- | ---- | -------------- | --------------------- | ------------------------- |
| aperture.feature.status |      | OK, Error      | Status of the feature | SDKs                      |

## Metric Columns

| Name                                         | Type  | Unit  | Description                                           |
| -------------------------------------------- | ----- | ----- | ----------------------------------------------------- |
| workload_duration_ms_sum                     | float | ms    | Sum of duration of the workload                       |
| workload_duration_ms_min                     | float | ms    | Min of duration of the workload                       |
| workload_duration_ms_max                     | float | ms    | Max of duration of the workload                       |
| workload_duration_ms_sumOfSquares            | float | ms    | Sum of squares of duration of the workload            |
| workload_duration_ms_datasketch              | float | ms    | Datasktech of Duration of the workload                |
| flow_duration_ms_sum                         | float | ms    | Sum of duration of the flow                           |
| flow_duration_ms_min                         | float | ms    | Min of duration of the flow                           |
| flow_duration_ms_max                         | float | ms    | Max of duration of the flow                           |
| flow_duration_ms_sumOfSquares                | float | ms    | Sum of squares of duration of the flow                |
| flow_duration_ms_datasketch                  | float | ms    | Datasktech of duration of the flow                    |
| aperture_processing_duration_ms_sum          | float | ms    | Sum of Aperture's processing duration                 |
| aperture_processing_duration_ms_min          | float | ms    | Min of Aperture's processing duration                 |
| aperture_processing_duration_ms_max          | float | ms    | Max of Aperture's processing duration                 |
| aperture_processing_duration_ms_sumOfSquares | float | ms    | Sum of squares of Aperture's processing duration      |
| aperture_processing_duration_ms_datasketch   | float | ms    | Datasktech of Aperture's processing duration          |
| http.request_content_length_sum              | int   | bytes | Sum of length of the HTTP request content             |
| http.request_content_length_min              | int   | bytes | Min of length of the HTTP request content             |
| http.request_content_length_max              | int   | bytes | Max of length of the HTTP request content             |
| http.request_content_length_sumOfSquares     | int   | bytes | Sum of squares of length of the HTTP request content  |
| http.response_content_length_sum             | int   | bytes | Sum of length of the HTTP response content            |
| http.response_content_length_min             | int   | bytes | Min of length of the HTTP response content            |
| http.response_content_length_max             | int   | bytes | Max of length of the HTTP response content            |
| http.response_content_length_sumOfSquares    | int   | bytes | Sum of squares of length of the HTTP response content |
