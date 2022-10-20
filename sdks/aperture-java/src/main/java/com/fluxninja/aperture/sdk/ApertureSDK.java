package com.fluxninja.aperture.sdk;

import com.fluxninja.aperture.flowcontrol.v1.CheckRequest;
import com.fluxninja.aperture.flowcontrol.v1.CheckResponse;
import com.fluxninja.aperture.flowcontrol.v1.FlowControlServiceGrpc;
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
  private final Tracer tracer;
  private final Duration timeout;

  ApertureSDK(
      FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient,
      Tracer tracer,
      Duration timeout) {
    this.flowControlClient = flowControlClient;
    this.tracer = tracer;
    this.timeout = timeout;
  }

  /**
   * Returns a new {@link ApertureSDKBuilder} for configuring an instance of {@linkplain
   * ApertureSDK the Aperture SDK}.
   */
  public static ApertureSDKBuilder builder() {
    return new ApertureSDKBuilder();
  }

  public Flow startFlow(String feature, Map<String, String> explicitLabels) {
    Map<String, String> labels = new HashMap<>();

    for (Map.Entry<String, BaggageEntry> entry: Baggage.current().asMap().entrySet()) {
      String value;
      try {
        value = URLDecoder.decode(entry.getValue().getValue(), StandardCharsets.UTF_8.name());
      } catch (java.io.UnsupportedEncodingException e) {
        // This should never happen, as `StandardCharsets.UTF_8.name()` is a valid encoding
        throw new RuntimeException(e);
      }
      labels.put(entry.getKey(), value);
    }

    labels.putAll(explicitLabels);


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

    return new Flow(
            res,
            span,
            false
    );
  }
}
