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
import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.baggage.BaggageBuilder;
import io.opentelemetry.api.baggage.BaggageEntry;
import io.opentelemetry.api.baggage.propagation.W3CBaggagePropagator;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Context;
import io.opentelemetry.context.propagation.ContextPropagators;
import io.opentelemetry.context.propagation.TextMapGetter;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.samplers.Sampler;
import io.opentelemetry.sdk.trace.export.SimpleSpanProcessor;

import javax.annotation.Nullable;
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
    private final OtlpGrpcSpanExporter spanExporter;
    private final Duration timeout;
    private final List<String> ignoredPaths;
    private final boolean ignoredPathsMatchRegex;
    private final ContextPropagators propagators;

    ApertureSDK(
            FlowControlServiceGrpc.FlowControlServiceBlockingStub flowControlClient,
            FlowControlServiceHTTPGrpc.FlowControlServiceHTTPBlockingStub httpFlowControlClient,
            Tracer tracer,
            OtlpGrpcSpanExporter spanExporter,
            Duration timeout,
            List<String> ignoredPaths,
            boolean ignoredPathsMatchRegex) {
        this.flowControlClient = flowControlClient;
        this.tracer = tracer;
        this.spanExporter = spanExporter;
        this.timeout = timeout;
        this.httpFlowControlClient = httpFlowControlClient;
        this.ignoredPaths = ignoredPaths;
        this.ignoredPathsMatchRegex = ignoredPathsMatchRegex;


        // create a tracer provider
        SdkTracerProvider tracerProvider = SdkTracerProvider.builder()
                .addSpanProcessor(SimpleSpanProcessor.create(spanExporter))
                .setSampler(Sampler.alwaysOn())
                .build();

        // create a propagator instance
        W3CBaggagePropagator propagator = W3CBaggagePropagator.getInstance();

        // create a context propagator that includes the baggage propagator
        propagators = ContextPropagators.create(propagator);

        // create an OpenTelemetry SDK instance with the tracer provider and propagators
        OpenTelemetrySdk openTelemetrySdk = OpenTelemetrySdk.builder()
                .setTracerProvider(tracerProvider)
                .setPropagators(propagators)
                .build();

        // set the global OpenTelemetry SDK instance
        GlobalOpenTelemetry.set(openTelemetrySdk);
    }

    /**
     * Returns a new {@link ApertureSDKBuilder} for configuring an instance of {@linkplain
     * ApertureSDK the Aperture SDK}.
     */
    public static ApertureSDKBuilder builder() {
        return new ApertureSDKBuilder();
    }


    static class MyTextMapGetter implements TextMapGetter<Map<String, String>> {

        @Override
        public Iterable<String> keys(Map<String, String> carrier) {
            return carrier.keySet();
        }

        @Nullable
        @Override
        public String get(@Nullable Map<String, String> carrier, String key) {
            return carrier.get(key);
        }
    }

    public Flow startFlow(String controlPoint, Map<String, String> explicitLabels) {
        Map<String, String> labels = new HashMap<>();

        if (explicitLabels.isEmpty()) {
            // System.out.println("Exp labels empty...");
        } else {
            System.out.println("Exp labels not empty!");
        }
        // extract the baggage from the headers
        Context extractedContext = propagators.getTextMapPropagator().extract(Context.current(), explicitLabels, new MyTextMapGetter());

        // retrieve the extracted baggage
        Baggage extractedBaggage = Baggage.fromContext(extractedContext);
        if (extractedBaggage.isEmpty()) {
            // System.out.println("Baggage empty...");
        } else {
            System.out.println("Baggage not empty!");
        }

        for (Map.Entry<String, BaggageEntry> entry : extractedBaggage.asMap().entrySet()) {
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
            System.out.println("Found in baggage: " + entry.getKey() + ": " + value);
            labels.put(entry.getKey(), value);
        }
        // System.out.println("Das all");

        if (explicitLabels != null) {
            labels.putAll(explicitLabels);
        }

        BaggageBuilder bb = Baggage.builder();
        for (Map.Entry<String, String> entry : labels.entrySet()) {
            bb.put(entry.getKey(), entry.getValue());
        }
        Baggage baggage = bb.build();

        // create a context with the baggage
        Context context = Context.current().with(baggage);

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
            res =
                    this.flowControlClient
                            .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                            .check(req);
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
            res =
                    this.httpFlowControlClient
                            .withDeadlineAfter(timeout.toNanos(), TimeUnit.NANOSECONDS)
                            .checkHTTP(req);
        } catch (StatusRuntimeException e) {
            // deadline exceeded or couldn't reach agent - request should not be blocked
        }
        span.setAttribute(WORKLOAD_START_TIMESTAMP_LABEL, Utils.getCurrentEpochNanos());

        return new TrafficFlow(res, span, false);
    }

    public void addBaggage(Map<String, String> headers) {
        BaggageBuilder baggageBuilder = Baggage.builder();
        for (Map.Entry<String, String> entry : headers.entrySet()) {
            baggageBuilder.put(entry.getKey(), entry.getValue());
        }
        baggageBuilder.build().storeInContext(Context.current()).makeCurrent();
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
