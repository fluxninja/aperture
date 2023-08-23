// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: aperture/flowcontrol/check/v1/check.proto

package com.fluxninja.generated.aperture.flowcontrol.check.v1;

public interface CheckResponseOrBuilder extends
    // @@protoc_insertion_point(interface_extends:aperture.flowcontrol.check.v1.CheckResponse)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * start timestamp
   * </pre>
   *
   * <code>.google.protobuf.Timestamp start = 1 [json_name = "start"];</code>
   * @return Whether the start field is set.
   */
  boolean hasStart();
  /**
   * <pre>
   * start timestamp
   * </pre>
   *
   * <code>.google.protobuf.Timestamp start = 1 [json_name = "start"];</code>
   * @return The start.
   */
  com.google.protobuf.Timestamp getStart();
  /**
   * <pre>
   * start timestamp
   * </pre>
   *
   * <code>.google.protobuf.Timestamp start = 1 [json_name = "start"];</code>
   */
  com.google.protobuf.TimestampOrBuilder getStartOrBuilder();

  /**
   * <pre>
   * end timestamp
   * </pre>
   *
   * <code>.google.protobuf.Timestamp end = 2 [json_name = "end"];</code>
   * @return Whether the end field is set.
   */
  boolean hasEnd();
  /**
   * <pre>
   * end timestamp
   * </pre>
   *
   * <code>.google.protobuf.Timestamp end = 2 [json_name = "end"];</code>
   * @return The end.
   */
  com.google.protobuf.Timestamp getEnd();
  /**
   * <pre>
   * end timestamp
   * </pre>
   *
   * <code>.google.protobuf.Timestamp end = 2 [json_name = "end"];</code>
   */
  com.google.protobuf.TimestampOrBuilder getEndOrBuilder();

  /**
   * <pre>
   * services that matched
   * </pre>
   *
   * <code>repeated string services = 4 [json_name = "services"];</code>
   * @return A list containing the services.
   */
  java.util.List<java.lang.String>
      getServicesList();
  /**
   * <pre>
   * services that matched
   * </pre>
   *
   * <code>repeated string services = 4 [json_name = "services"];</code>
   * @return The count of services.
   */
  int getServicesCount();
  /**
   * <pre>
   * services that matched
   * </pre>
   *
   * <code>repeated string services = 4 [json_name = "services"];</code>
   * @param index The index of the element to return.
   * @return The services at the given index.
   */
  java.lang.String getServices(int index);
  /**
   * <pre>
   * services that matched
   * </pre>
   *
   * <code>repeated string services = 4 [json_name = "services"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the services at the given index.
   */
  com.google.protobuf.ByteString
      getServicesBytes(int index);

  /**
   * <pre>
   * control_point of request
   * </pre>
   *
   * <code>string control_point = 5 [json_name = "controlPoint"];</code>
   * @return The controlPoint.
   */
  java.lang.String getControlPoint();
  /**
   * <pre>
   * control_point of request
   * </pre>
   *
   * <code>string control_point = 5 [json_name = "controlPoint"];</code>
   * @return The bytes for controlPoint.
   */
  com.google.protobuf.ByteString
      getControlPointBytes();

  /**
   * <pre>
   * flow label keys that were matched for this request.
   * </pre>
   *
   * <code>repeated string flow_label_keys = 6 [json_name = "flowLabelKeys"];</code>
   * @return A list containing the flowLabelKeys.
   */
  java.util.List<java.lang.String>
      getFlowLabelKeysList();
  /**
   * <pre>
   * flow label keys that were matched for this request.
   * </pre>
   *
   * <code>repeated string flow_label_keys = 6 [json_name = "flowLabelKeys"];</code>
   * @return The count of flowLabelKeys.
   */
  int getFlowLabelKeysCount();
  /**
   * <pre>
   * flow label keys that were matched for this request.
   * </pre>
   *
   * <code>repeated string flow_label_keys = 6 [json_name = "flowLabelKeys"];</code>
   * @param index The index of the element to return.
   * @return The flowLabelKeys at the given index.
   */
  java.lang.String getFlowLabelKeys(int index);
  /**
   * <pre>
   * flow label keys that were matched for this request.
   * </pre>
   *
   * <code>repeated string flow_label_keys = 6 [json_name = "flowLabelKeys"];</code>
   * @param index The index of the value to return.
   * @return The bytes of the flowLabelKeys at the given index.
   */
  com.google.protobuf.ByteString
      getFlowLabelKeysBytes(int index);

  /**
   * <pre>
   * telemetry_flow_labels are labels for telemetry purpose. The keys in telemetry_flow_labels is subset of flow_label_keys.
   * </pre>
   *
   * <code>map&lt;string, string&gt; telemetry_flow_labels = 7 [json_name = "telemetryFlowLabels"];</code>
   */
  int getTelemetryFlowLabelsCount();
  /**
   * <pre>
   * telemetry_flow_labels are labels for telemetry purpose. The keys in telemetry_flow_labels is subset of flow_label_keys.
   * </pre>
   *
   * <code>map&lt;string, string&gt; telemetry_flow_labels = 7 [json_name = "telemetryFlowLabels"];</code>
   */
  boolean containsTelemetryFlowLabels(
      java.lang.String key);
  /**
   * Use {@link #getTelemetryFlowLabelsMap()} instead.
   */
  @java.lang.Deprecated
  java.util.Map<java.lang.String, java.lang.String>
  getTelemetryFlowLabels();
  /**
   * <pre>
   * telemetry_flow_labels are labels for telemetry purpose. The keys in telemetry_flow_labels is subset of flow_label_keys.
   * </pre>
   *
   * <code>map&lt;string, string&gt; telemetry_flow_labels = 7 [json_name = "telemetryFlowLabels"];</code>
   */
  java.util.Map<java.lang.String, java.lang.String>
  getTelemetryFlowLabelsMap();
  /**
   * <pre>
   * telemetry_flow_labels are labels for telemetry purpose. The keys in telemetry_flow_labels is subset of flow_label_keys.
   * </pre>
   *
   * <code>map&lt;string, string&gt; telemetry_flow_labels = 7 [json_name = "telemetryFlowLabels"];</code>
   */
  /* nullable */
java.lang.String getTelemetryFlowLabelsOrDefault(
      java.lang.String key,
      /* nullable */
java.lang.String defaultValue);
  /**
   * <pre>
   * telemetry_flow_labels are labels for telemetry purpose. The keys in telemetry_flow_labels is subset of flow_label_keys.
   * </pre>
   *
   * <code>map&lt;string, string&gt; telemetry_flow_labels = 7 [json_name = "telemetryFlowLabels"];</code>
   */
  java.lang.String getTelemetryFlowLabelsOrThrow(
      java.lang.String key);

  /**
   * <pre>
   * decision_type contains what the decision was.
   * </pre>
   *
   * <code>.aperture.flowcontrol.check.v1.CheckResponse.DecisionType decision_type = 8 [json_name = "decisionType"];</code>
   * @return The enum numeric value on the wire for decisionType.
   */
  int getDecisionTypeValue();
  /**
   * <pre>
   * decision_type contains what the decision was.
   * </pre>
   *
   * <code>.aperture.flowcontrol.check.v1.CheckResponse.DecisionType decision_type = 8 [json_name = "decisionType"];</code>
   * @return The decisionType.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse.DecisionType getDecisionType();

  /**
   * <pre>
   * reject_reason contains the reason for the rejection.
   * </pre>
   *
   * <code>.aperture.flowcontrol.check.v1.CheckResponse.RejectReason reject_reason = 9 [json_name = "rejectReason"];</code>
   * @return The enum numeric value on the wire for rejectReason.
   */
  int getRejectReasonValue();
  /**
   * <pre>
   * reject_reason contains the reason for the rejection.
   * </pre>
   *
   * <code>.aperture.flowcontrol.check.v1.CheckResponse.RejectReason reject_reason = 9 [json_name = "rejectReason"];</code>
   * @return The rejectReason.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse.RejectReason getRejectReason();

  /**
   * <pre>
   * classifiers that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.ClassifierInfo classifier_infos = 10 [json_name = "classifierInfos"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.flowcontrol.check.v1.ClassifierInfo> 
      getClassifierInfosList();
  /**
   * <pre>
   * classifiers that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.ClassifierInfo classifier_infos = 10 [json_name = "classifierInfos"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.ClassifierInfo getClassifierInfos(int index);
  /**
   * <pre>
   * classifiers that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.ClassifierInfo classifier_infos = 10 [json_name = "classifierInfos"];</code>
   */
  int getClassifierInfosCount();
  /**
   * <pre>
   * classifiers that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.ClassifierInfo classifier_infos = 10 [json_name = "classifierInfos"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.flowcontrol.check.v1.ClassifierInfoOrBuilder> 
      getClassifierInfosOrBuilderList();
  /**
   * <pre>
   * classifiers that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.ClassifierInfo classifier_infos = 10 [json_name = "classifierInfos"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.ClassifierInfoOrBuilder getClassifierInfosOrBuilder(
      int index);

  /**
   * <pre>
   * flux meters that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.FluxMeterInfo flux_meter_infos = 11 [json_name = "fluxMeterInfos"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.flowcontrol.check.v1.FluxMeterInfo> 
      getFluxMeterInfosList();
  /**
   * <pre>
   * flux meters that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.FluxMeterInfo flux_meter_infos = 11 [json_name = "fluxMeterInfos"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.FluxMeterInfo getFluxMeterInfos(int index);
  /**
   * <pre>
   * flux meters that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.FluxMeterInfo flux_meter_infos = 11 [json_name = "fluxMeterInfos"];</code>
   */
  int getFluxMeterInfosCount();
  /**
   * <pre>
   * flux meters that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.FluxMeterInfo flux_meter_infos = 11 [json_name = "fluxMeterInfos"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.flowcontrol.check.v1.FluxMeterInfoOrBuilder> 
      getFluxMeterInfosOrBuilderList();
  /**
   * <pre>
   * flux meters that were matched for this request.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.FluxMeterInfo flux_meter_infos = 11 [json_name = "fluxMeterInfos"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.FluxMeterInfoOrBuilder getFluxMeterInfosOrBuilder(
      int index);

  /**
   * <pre>
   * limiter_decisions contains information about decision made by each limiter.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.LimiterDecision limiter_decisions = 12 [json_name = "limiterDecisions"];</code>
   */
  java.util.List<com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision> 
      getLimiterDecisionsList();
  /**
   * <pre>
   * limiter_decisions contains information about decision made by each limiter.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.LimiterDecision limiter_decisions = 12 [json_name = "limiterDecisions"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecision getLimiterDecisions(int index);
  /**
   * <pre>
   * limiter_decisions contains information about decision made by each limiter.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.LimiterDecision limiter_decisions = 12 [json_name = "limiterDecisions"];</code>
   */
  int getLimiterDecisionsCount();
  /**
   * <pre>
   * limiter_decisions contains information about decision made by each limiter.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.LimiterDecision limiter_decisions = 12 [json_name = "limiterDecisions"];</code>
   */
  java.util.List<? extends com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecisionOrBuilder> 
      getLimiterDecisionsOrBuilderList();
  /**
   * <pre>
   * limiter_decisions contains information about decision made by each limiter.
   * </pre>
   *
   * <code>repeated .aperture.flowcontrol.check.v1.LimiterDecision limiter_decisions = 12 [json_name = "limiterDecisions"];</code>
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.LimiterDecisionOrBuilder getLimiterDecisionsOrBuilder(
      int index);

  /**
   * <pre>
   * Recommended minimal amount of time to wait when retrying the request, if
   * decision_type is REJECTED. Optional.
   * </pre>
   *
   * <code>.google.protobuf.Duration wait_time = 13 [json_name = "waitTime"];</code>
   * @return Whether the waitTime field is set.
   */
  boolean hasWaitTime();
  /**
   * <pre>
   * Recommended minimal amount of time to wait when retrying the request, if
   * decision_type is REJECTED. Optional.
   * </pre>
   *
   * <code>.google.protobuf.Duration wait_time = 13 [json_name = "waitTime"];</code>
   * @return The waitTime.
   */
  com.google.protobuf.Duration getWaitTime();
  /**
   * <pre>
   * Recommended minimal amount of time to wait when retrying the request, if
   * decision_type is REJECTED. Optional.
   * </pre>
   *
   * <code>.google.protobuf.Duration wait_time = 13 [json_name = "waitTime"];</code>
   */
  com.google.protobuf.DurationOrBuilder getWaitTimeOrBuilder();

  /**
   * <pre>
   * http_status contains the http status code to be returned to the client, if
   * decision_type is REJECTED. Optional.
   * </pre>
   *
   * <code>.aperture.flowcontrol.check.v1.StatusCode denied_response_status_code = 14 [json_name = "deniedResponseStatusCode", (.validate.rules) = { ... }</code>
   * @return The enum numeric value on the wire for deniedResponseStatusCode.
   */
  int getDeniedResponseStatusCodeValue();
  /**
   * <pre>
   * http_status contains the http status code to be returned to the client, if
   * decision_type is REJECTED. Optional.
   * </pre>
   *
   * <code>.aperture.flowcontrol.check.v1.StatusCode denied_response_status_code = 14 [json_name = "deniedResponseStatusCode", (.validate.rules) = { ... }</code>
   * @return The deniedResponseStatusCode.
   */
  com.fluxninja.generated.aperture.flowcontrol.check.v1.StatusCode getDeniedResponseStatusCode();
}
