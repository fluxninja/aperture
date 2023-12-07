package com.fluxninja.aperture.sdk;

import static com.fluxninja.aperture.sdk.Constants.FLOW_START_TIMESTAMP_LABEL;
import static com.fluxninja.aperture.sdk.Constants.SOURCE_LABEL;
import static com.fluxninja.aperture.sdk.Constants.WORKLOAD_START_TIMESTAMP_LABEL;

import com.fluxninja.generated.aperture.flowcontrol.check.v1.CacheLookupRequest;
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
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import java.util.regex.Pattern;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * The Aperture SDK provides a set of tools and functionalities for flow control by enabling the
 * initiation of flows.
 *
 * <p>To start using the Aperture SDK, create an instance of this class using the {@link
 * ApertureSDKBuilder} and utilize its methods to initiate and control flows in your application.
 */
public final class ApertureSDK {
    private final FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient;
    private final FlowControlServiceHTTPGrpc.FlowControlServiceHTTPBlockingStub
            httpFlowControlClient;
    private final Tracer tracer;
    private final List<String> ignoredPaths;
    private final boolean ignoredPathsMatchRegex;

    private static final Logger logger = LoggerFactory.getLogger(ApertureSDK.class);

    ApertureSDK(
            FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient,
            FlowControlServiceHTTPGrpc.FlowControlServiceHTTPBlockingStub httpFlowControlClient,
            Tracer tracer,
            List<String> ignoredPaths,
            boolean ignoredPathsMatchRegex) {
        this.flowControlClient = flowControlClient;
        this.tracer = tracer;
        this.httpFlowControlClient = httpFlowControlClient;
        this.ignoredPaths = ignoredPaths;
        this.ignoredPathsMatchRegex = ignoredPathsMatchRegex;
    }

    /**
     * Returns a new {@link ApertureSDKBuilder} for configuring an instance of {@linkplain
     * ApertureSDK the Aperture SDK}.
     *
     * @return A new ApertureSDKBuilder object.
     */
    public static ApertureSDKBuilder builder() {
        return new ApertureSDKBuilder();
    }

    /**
     * Starts a new flow, asking the Aperture Agent to accept or reject it based on provided labels.
     * Additional labels will be extracted from current Baggage context.
     *
     * @param parameters Flow parameters that can be built using {@link
     *     FeatureFlowParameters.Builder}
     * @return A Flow object
     */
    public Flow startFlow(FeatureFlowParameters parameters) {
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
                return new Flow(
                        null,
                        null,
                        false,
                        parameters.getRampMode(),
                        this.flowControlClient,
                        parameters.getResultCacheKey(),
                        parameters.getControlPoint(),
                        e);
            }
            labels.put(entry.getKey(), value);
        }

        if (parameters.getExplicitLabels() != null) {
            labels.putAll(parameters.getExplicitLabels());
        }

        CheckRequest req =
                CheckRequest.newBuilder()
                        .setControlPoint(parameters.getControlPoint())
                        .putAllLabels(labels)
                        .setRampMode(parameters.getRampMode())
                        .setCacheLookupRequest(
                                CacheLookupRequest.newBuilder()
                                        .addAllGlobalCacheKeys(parameters.getGlobalCacheKeys())
                                        .setResultCacheKey(parameters.getResultCacheKey())
                                        .build())
                        .build();

        Span span =
                this.tracer
                        .spanBuilder("Aperture Check")
                        .startSpan()
                        .setAttribute(FLOW_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos())
                        .setAttribute(SOURCE_LABEL, "sdk");

        CheckResponse res = null;
        try {
            if (parameters.getFlowTimeout().isZero()) {
                res = this.flowControlClient.check(req);
            } else {
                res =
                        this.flowControlClient
                                .withDeadlineAfter(
                                        parameters.getFlowTimeout().toNanos(), TimeUnit.NANOSECONDS)
                                .check(req);
            }
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
            return new Flow(
                    res,
                    span,
                    false,
                    parameters.getRampMode(),
                    this.flowControlClient,
                    parameters.getResultCacheKey(),
                    parameters.getControlPoint(),
                    e);
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new Flow(
                res,
                span,
                false,
                parameters.getRampMode(),
                this.flowControlClient,
                parameters.getResultCacheKey(),
                parameters.getControlPoint(),
                null);
    }

    /**
     * Starts a new flow, asking the Aperture Agent to accept or reject it based on provided HTTP
     * request parameters.
     *
     * @param req A {@link TrafficFlowRequest} containing configured HTTP request parameters.
     * @return A TrafficFlow object
     */
    public TrafficFlow startTrafficFlow(TrafficFlowRequest req) {
        CheckHTTPRequest checkHTTPRequest = req.getCheckHTTPRequest();
        String path = checkHTTPRequest.getRequest().getPath();

        if (isIgnored(path)) {
            logger.debug("Path " + path + " is set to be ignored.");
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
            if (req.getFlowTimeout().isZero()) {
                res = this.httpFlowControlClient.checkHTTP(checkHTTPRequest);
            } else {
                res =
                        this.httpFlowControlClient
                                .withDeadlineAfter(
                                        req.getFlowTimeout().toNanos(), TimeUnit.NANOSECONDS)
                                .checkHTTP(checkHTTPRequest);
            }
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new TrafficFlow(res, span, false, req.getCheckHTTPRequest().getRampMode());
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
