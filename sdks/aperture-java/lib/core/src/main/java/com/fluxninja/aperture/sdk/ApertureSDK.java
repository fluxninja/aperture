package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.*;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckRequest;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.CheckResponse;
import com.fluxninja.generated.aperture.flowcontrol.check.v1.FlowControlServiceGrpc;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPResponse;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.FlowControlServiceHTTPGrpc;
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

public final class ApertureSDK {
    private final FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient;
    private final FlowControlServiceHTTPGrpc.FlowControlServiceHTTPBlockingStub
            httpFlowControlClient;
    private final Tracer tracer;
    private final Duration timeout;
    private final List<String> ignoredPaths;
    private final boolean ignoredPathsMatchRegex;

    ApertureSDK(
            FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient,
            FlowControlServiceHTTPGrpc.FlowControlServiceHTTPBlockingStub httpFlowControlClient,
            Tracer tracer,
            Duration timeout,
            List<String> ignoredPaths,
            boolean ignoredPathsMatchRegex) {
        this.flowControlClient = flowControlClient;
        this.tracer = tracer;
        this.timeout = timeout;
        this.httpFlowControlClient = httpFlowControlClient;
        this.ignoredPaths = ignoredPaths;
        this.ignoredPathsMatchRegex = ignoredPathsMatchRegex;
    }

    /**
     * Returns a new {@link ApertureSDKBuilder} for configuring an instance of {@linkplain
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
                value =
                        URLDecoder.decode(
                                entry.getValue().getValue(), StandardCharsets.UTF_8.name());
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

        CheckRequest req =
                CheckRequest.newBuilder()
                        .setControlPoint(controlPoint)
                        .putAllLabels(labels)
                        .build();

        Span span =
                this.tracer
                        .spanBuilder("Aperture Check")
                        .startSpan()
                        .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
                        .setAttribute(SOURCE_LABEL, "sdk");

        CheckResponse res = null;
        try {
            if (timeout.isZero()) {
                res = this.flowControlClient.check(req);
            } else {
                res =
                        this.flowControlClient
                                .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                                .check(req);
            }
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new Flow(res, span, false);
    }

    public TrafficFlow startTrafficFlow(String path, CheckHTTPRequest req) {
        if (isIgnored(path)) {
            return TrafficFlow.ignoredFlow();
        }

        Span span =
                this.tracer
                        .spanBuilder("Aperture Check")
                        .startSpan()
                        .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
                        .setAttribute(SOURCE_LABEL, "sdk");

        CheckHTTPResponse res = null;
        try {
            if (timeout.isZero()) {
                res = this.httpFlowControlClient.checkHTTP(req);
            } else {
                res =
                        this.httpFlowControlClient
                                .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                                .checkHTTP(req);
            }
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new TrafficFlow(res, span, false);
    }

    private boolean isIgnored(String path) {
        if (path == null) {
            return false;
        }
        for (String ignoredPattern : ignoredPaths) {
            if (ignoredPathsMatchRegex) {
                if (Pattern.matches(ignoredPattern, path)) {
                    return true;
                }
            } else {
                if (ignoredPattern.equals(path)) {
                    return true;
                }
            }
        }
        return false;
    }
}
