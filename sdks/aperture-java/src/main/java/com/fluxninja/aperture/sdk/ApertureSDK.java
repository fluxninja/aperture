package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowControlServiceGrpc;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.AuthorizationGrpc;
import com.fluxninja.generated.google.rpc.Status;
import com.google.rpc.Code;
import io.grpc.StatusRuntimeException;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;

import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import static com.fluxninja.aperture.sdk.Constants.*;

public final class ApertureSDK {
  private final FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient;
  private final AuthorizationGrpc.AuthorizationBlockingStub envoyAuthzClient;
  private final Tracer tracer;
  private final Duration timeout;

  ApertureSDK(
      FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient,
      AuthorizationGrpc.AuthorizationBlockingStub envoyAuthzClient,
      Tracer tracer,
      Duration timeout) {
    this.flowControlClient = flowControlClient;
    this.tracer = tracer;
    this.timeout = timeout;
    this.envoyAuthzClient = envoyAuthzClient;
  }

  /**
   * Returns a new {@link ApertureSDKBuilder} for configuring an instance of
   * {@linkplain
   * ApertureSDK the Aperture SDK}.
   */
  public static ApertureSDKBuilder builder() {
    return new ApertureSDKBuilder();
  }

  public FeatureFlow startFlow(String feature, Map<String, String> explicitLabels) {
    Map<String, String> labels = new HashMap<>();

    for (Map.Entry<String, BaggageEntry> entry : Baggage.current().asMap().entrySet()) {
      String value;
      try {
        value = URLDecoder.decode(entry.getValue().getValue(), StandardCharsets.UTF_8.name());
      } catch (java.io.UnsupportedEncodingException e) {
        // This should never happen, as `StandardCharsets.UTF_8.name()` is a valid
        // encoding
        throw new RuntimeException(e);
      }
      labels.put(entry.getKey(), value);
    }

    if (explicitLabels != null) {
      labels.putAll(explicitLabels);
    }

    CheckRequest req = CheckRequest.newBuilder()
        .setFeature(feature)
        .putAllLabels(labels)
        .build();

    Span span = this.tracer.spanBuilder("Aperture Check").startSpan()
        .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
        .setAttribute(SOURCE_LABEL, "sdk");

    CheckResponse res;
    try {
      res = this.flowControlClient
          .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
          .check(req);
    } catch (StatusRuntimeException e) {
      // deadline exceeded or couldn't reach agent - request should not be blocked
      res = CheckResponse.newBuilder()
          .setDecisionType(CheckResponse.DecisionType.DECISION_TYPE_ACCEPTED)
          .build();
    }
    span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

    return new FeatureFlow(
        res,
        span,
        false);
  }

  public TrafficFlow startTrafficFlow(AttributeContext attributes) {
    com.fluxninja.generated.envoy.service.auth.v3.CheckRequest req = com.fluxninja.generated.envoy.service.auth.v3.CheckRequest
        .newBuilder()
        .setAttributes(attributes)
        .build();

    Span span = this.tracer.spanBuilder("Aperture Check").startSpan()
        .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
        .setAttribute(SOURCE_LABEL, "sdk");

    com.fluxninja.generated.envoy.service.auth.v3.CheckResponse res;
    try {
      res = this.envoyAuthzClient
          .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
          .check(req);
    } catch (StatusRuntimeException e) {
      // deadline exceeded or couldn't reach agent - request should not be blocked
      res = com.fluxninja.generated.envoy.service.auth.v3.CheckResponse.newBuilder()
          .setStatus(Status.newBuilder().setCode(Code.OK_VALUE).build())
          .build();
    }
    span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

    return new TrafficFlow(
        res,
        span,
        false);
  }
}
