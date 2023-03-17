package com.fluxninja.aperture.sdk;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowControlServiceGrpc;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.AuthorizationGrpc;
import io.grpc.StatusRuntimeException;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;

import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import java.util.regex.Pattern;

import static com.fluxninja.aperture.sdk.Constants.*;

public final class ApertureSDK {
    private final FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient;
    private final AuthorizationGrpc.AuthorizationBlockingStub envoyAuthzClient;
    private final Tracer tracer;
    private final Duration timeout;
    private final List<String> blockedPaths;
    private final boolean blockedPathsMatchRegex;

    ApertureSDK(
            FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient,
            AuthorizationGrpc.AuthorizationBlockingStub envoyAuthzClient,
            Tracer tracer,
            Duration timeout,
            List<String> blockedPaths,
            boolean blockedPathsMatchRegex) {
        this.flowControlClient = flowControlClient;
        this.tracer = tracer;
        this.timeout = timeout;
        this.envoyAuthzClient = envoyAuthzClient;
        this.blockedPaths = blockedPaths;
        this.blockedPathsMatchRegex = blockedPathsMatchRegex;
    }

    /**
     * Returns a new {@link ApertureSDKBuilder} for configuring an instance of
     * {@linkplain
     * ApertureSDK the Aperture SDK}.
     */
    public static ApertureSDKBuilder builder() {
        return new ApertureSDKBuilder();
    }

    public Flow startFlow(String controlPoint, Map<String, String> explicitLabels) {
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
                .setControlPoint(controlPoint)
                .putAllLabels(labels)
                .build();

        Span span = this.tracer.spanBuilder("Aperture Check").startSpan()
                .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
                .setAttribute(SOURCE_LABEL, "sdk");

        CheckResponse res = null;
        try {
            res = this.flowControlClient
                    .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                    .check(req);
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new Flow(
                res,
                span,
                false);
    }

    public TrafficFlow startTrafficFlow(String path, AttributeContext attributes) {
        if (isBlocked(path)) {
            return TrafficFlow.ignoredFlow();
        }

        com.fluxninja.generated.envoy.service.auth.v3.CheckRequest req = com.fluxninja.generated.envoy.service.auth.v3.CheckRequest
                .newBuilder()
                .setAttributes(attributes)
                .build();

        Span span = this.tracer.spanBuilder("Aperture Check").startSpan()
                .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
                .setAttribute(SOURCE_LABEL, "sdk");

        com.fluxninja.generated.envoy.service.auth.v3.CheckResponse res = null;
        try {
            res = this.envoyAuthzClient
                    .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                    .check(req);
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new TrafficFlow(
                res,
                span,
                false);
    }

    private boolean isBlocked(String path) {
        if (path == null) {
            return false;
        }
        for (String blockedPattern : blockedPaths) {
            if (blockedPathsMatchRegex) {
                if (Pattern.matches(blockedPattern, path)) {
                    return true;
                }
            } else {
                if (blockedPattern.equals(path)) {
                    return true;
                }
            }
        }
        return false;
    }
}
