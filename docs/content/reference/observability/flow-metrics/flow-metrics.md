---
title: Flow Metrics
sidebar_position: 1
sidebar_label: Flow Metrics
---

Aperture Agents emit an OpenTelemetry stream for flow data, which provides a
comprehensive view of individual requests or features within services. This
stream contains high-cardinality attributes that represent key attributes of the
requests and features, allowing for a detailed analysis of system performance
and behavior. The stream can be stored and visualized in
[FluxNinja ARC](/arc/arc.md), or ingested into popular OLAP databases such as
[Apache Druid](https://druid.apache.org/).

## Dimension Columns

### Common

| Name                                   | Type        | Example Values                                                                                                                                                            | Description                                        | Flow Control Integrations |
| -------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------- | ------------------------- |
| aperture.source                        | single      | sdk, envoy                                                                                                                                                                | Aperture Flow source                               | SDKs, Envoy               |
| aperture.decision_type                 | single      | DECISION_TYPE_ACCEPTED, DECISION_TYPE_REJECTED                                                                                                                            | Decision type taken by policy                      | SDKs, Envoy               |
| aperture.reject_reason                 | single      | REJECT_REASON_NONE, REJECT_REASON_RATE_LIMITED, REJECT_REASON_CONCURRENCY_LIMITED                                                                                         | Reject reason of the decision taken by policy      | SDKs, Envoy               |
| aperture.rate_limiters                 | multi-value | "policy_name:s1, component_id:18, policy_hash:5kZjj"                                                                                                                      | Rate limiters matched to the traffic               | SDKs, Envoy               |
| aperture.dropping_rate_limiters        | multi-value | "policy_name:s1, component_id:18, policy_hash:5kZjj"                                                                                                                      | Rate limiters dropping the traffic                 | SDKs, Envoy               |
| aperture.concurrency_limiters          | multi-value | "policy_name:s1, component_id:13, policy_hash:5kZjj"                                                                                                                      | Concurrency limiters matched to the traffic        | SDKs, Envoy               |
| aperture.dropping_concurrency_limiters | multi-value | "policy_name:s1, component_id:13, policy_hash:5kZjj"                                                                                                                      | Concurrency limiters dropping the traffic          | SDKs, Envoy               |
| aperture.workloads                     | multi-value | "policy_name:s1, component_id:13, workload_index:0, policy_hash:5kZjj"                                                                                                    | Workloads matched to the traffic                   | SDKs, Envoy               |
| aperture.dropping_workloads            | multi-value | "policy_name:s1, component_id:13, workload_index:0, policy_hash:5kZjj"                                                                                                    | Workloads dropping the traffic                     | SDKs, Envoy               |
| aperture.flux_meters                   | multi-value | s1                                                                                                                                                                        | Flux Meters matched to the traffic                 | SDKs, Envoy               |
| aperture.flow_label_keys               | multi-value | http.host, http.method, http.request.header.content_length                                                                                                                | Flow labels matched to the traffic                 | SDKs, Envoy               |
| aperture.classifiers                   | multi-value | "policy_name:s1, classifier_index:0"                                                                                                                                      | Classifiers matched to the traffic                 | SDKs, Envoy               |
| aperture.classifier_errors             | multi-value | "[ERROR_NONE, ERROR_EVAL_FAILED, ERROR_EMPTY_RESULTSET, ERROR_AMBIGUOUS_RESULTSET, ERROR_MULTI_EXPRESSION, ERROR_EXPRESSION_NOT_MAP], policy_name:s1, classifier_index:0" | Encountered classifier errors for specified policy | SDKs, Envoy               |
| aperture.services                      | multi-value | s1.demoapp.svc.cluster.local, s2.demoapp.svc.cluster.local                                                                                                                | Services to which metrics refer                    | SDKs, Envoy               |
| aperture.control_point                 | single      | type:TYPE_INGRESS, type:TYPE_EGRESS                                                                                                                                       | Control point to which metrics refer               | SDKs, Envoy               |
| aperture.flow.status                   | single      | OK, Error                                                                                                                                                                 | Denotes OK or Error across all protocols           | SDKs, Envoy               |
| response_received                      | single      | true, false                                                                                                                                                               | Designates whether a response was received         | SDKs, envoy               |

### HTTP

| Name                         | Type   | Example Values                                             | Description                                            | Flow Control Integrations |
| ---------------------------- | ------ | ---------------------------------------------------------- | ------------------------------------------------------ | ------------------------- |
| http.status_code             | single | 200, 429, 503                                              | HTTP status code of the response                       | Envoy                     |
| http.request_content_length  | single | 0, 53                                                      | Length of the HTTP request content in bytes            | Envoy                     |
| http.response_content_length | single | 201, 77                                                    | Length of the HTTP response content in bytes           | Envoy                     |
| http.method                  | single | GET, POST                                                  | HTTP method of the response                            | Envoy                     |
| http.target                  | single | /request                                                   | Target endpoint of the response                        | Envoy                     |
| http.host                    | single | s1.demoapp.svc.cluster.local, s2.demoapp.svc.cluster.local | Host address of the response                           | Envoy                     |
| http.scheme                  | single | http                                                       | HTTP scheme of the response                            | Envoy                     |
| http.flavor                  | single | 1.1                                                        | HTTP protocol version                                  | Envoy                     |
| {user-defined-labels}        |        |                                                            | Configured through [Flow Classifiers][flowclassifiers] | Envoy                     |

### SDK

| Name                  | Type | Example Values | Description                                      | Flow Control Integrations |
| --------------------- | ---- | -------------- | ------------------------------------------------ | ------------------------- |
| {user-defined-labels} |      |                | Explicitly passed through FlowStart call in SDKs | SDKs                      |

## Metric Columns

| Name                                         | Type                    | Unit  | Description                                           |
| -------------------------------------------- | ----------------------- | ----- | ----------------------------------------------------- |
| workload_duration_ms_sum                     | float                   | ms    | Sum of duration of the workload                       |
| workload_duration_ms_min                     | float                   | ms    | Min of duration of the workload                       |
| workload_duration_ms_max                     | float                   | ms    | Max of duration of the workload                       |
| workload_duration_ms_sumOfSquares            | float                   | ms    | Sum of squares of duration of the workload            |
| workload_duration_ms_datasketch              | [quantilesDoubleSketch] | ms    | Datasktech of Duration of the workload                |
| flow_duration_ms_sum                         | float                   | ms    | Sum of duration of the flow                           |
| flow_duration_ms_min                         | float                   | ms    | Min of duration of the flow                           |
| flow_duration_ms_max                         | float                   | ms    | Max of duration of the flow                           |
| flow_duration_ms_sumOfSquares                | float                   | ms    | Sum of squares of duration of the flow                |
| flow_duration_ms_datasketch                  | [quantilesDoubleSketch] | ms    | Sum of Aperture's processing duration                 |
| aperture_processing_duration_ms_min          | float                   | ms    | Min of Aperture's processing duration                 |
| aperture_processing_duration_ms_max          | float                   | ms    | Max of Aperture's processing duration                 |
| aperture_processing_duration_ms_sumOfSquares | float                   | ms    | Sum of squares of Aperture's processing duration      |
| aperture_processing_duration_ms_datasketch   | [quantilesDoubleSketch] | ms    | Datasktech of Aperture's processing duration          |
| http.request_content_length_sum              | int                     | bytes | Sum of length of the HTTP request content             |
| http.request_content_length_min              | int                     | bytes | Min of length of the HTTP request content             |
| http.request_content_length_max              | int                     | bytes | Max of length of the HTTP request content             |
| http.request_content_length_sumOfSquares     | int                     | bytes | Sum of squares of length of the HTTP request content  |
| http.response_content_length_sum             | int                     | bytes | Sum of length of the HTTP response content            |
| http.response_content_length_min             | int                     | bytes | Min of length of the HTTP response content            |
| http.response_content_length_max             | int                     | bytes | Max of length of the HTTP response content            |
| http.response_content_length_sumOfSquares    | int                     | bytes | Sum of squares of length of the HTTP response content |

[quantilesdoublesketch]:
  https://druid.apache.org/docs/latest/development/extensions-core/datasketches-quantiles.html
[flowclassifiers]: /concepts/integrations/flow-control/flow-classifier.md
