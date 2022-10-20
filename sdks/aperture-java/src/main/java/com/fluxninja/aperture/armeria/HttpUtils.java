package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.EnvoyAuthzFlow;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpStatus;
import com.linecorp.armeria.common.RequestHeaders;
import io.netty.util.AsciiString;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;

import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

class HttpUtils {
    protected static HttpStatus handleRejectedFlow(EnvoyAuthzFlow flow) {
        try {
            flow.end(FlowStatus.Unset);
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }
        if (flow.checkResponse().hasDeniedResponse() && flow.checkResponse().getDeniedResponse().hasStatus()) {
            int httpStatusCode = flow.checkResponse().getDeniedResponse().getStatus().getCodeValue();
            return HttpStatus.valueOf(httpStatusCode);
        }
        return HttpStatus.BAD_REQUEST;
    }

    protected static Map<String, String> labelsFromRequest(HttpRequest req) {
        Map<String, String> labels = new HashMap<>();
        RequestHeaders headers = req.headers();
        for (Map.Entry<AsciiString, String> header: headers) {
            String headerKey = header.getKey().toString();
            if (headerKey.startsWith(":")) {
                continue;
            }
            String labelName = String.format("http.request.header.%s", headerKey);
            labels.put(labelName, header.getValue());
        }
        labels.put("http.method", req.method().toString());
        labels.put("http.uri", req.uri().toString());
        return labels;
    }

    protected static AttributeContext attributesFromRequest(HttpRequest req) {
        Map<String, String> baggageLabels = new HashMap<>();

        for (Map.Entry<String, BaggageEntry> entry : Baggage.current().asMap().entrySet()) {
            String value;
            try {
                value = URLDecoder.decode(entry.getValue().getValue(), StandardCharsets.UTF_8.name());
            } catch (java.io.UnsupportedEncodingException e) {
                // This should never happen, as `StandardCharsets.UTF_8.name()` is a valid encoding
                throw new RuntimeException(e);
            }
            baggageLabels.put(entry.getKey(), value);
        }

        AttributeContext.Builder builder = AttributeContext.newBuilder();
        builder.putAllContextExtensions(baggageLabels);  // TODO: is this even needed here?

        return addHttpAttributes(builder, req).build();
    }

    private static AttributeContext.Builder addHttpAttributes(AttributeContext.Builder builder, HttpRequest req) {
        Map<String, String> extractedHeaders = new HashMap<>();
        RequestHeaders headers = req.headers();
        for (Map.Entry<AsciiString, String> header: headers) {
            String headerKey = header.getKey().toString();
            if (headerKey.startsWith(":")) {
                continue;
            }
            extractedHeaders.put(headerKey, header.getValue());
        }
        return builder.setRequest(AttributeContext.Request.newBuilder()
                .setHttp(AttributeContext.HttpRequest.newBuilder()
                        .setMethod(req.method().toString())
                        .setPath(req.path())
                        .setHost(req.authority())
                        .setScheme(req.scheme())
                        .setSize(-1) // TODO: should we set something here?
                        .setProtocol("HTTP/2")
                        .putAllHeaders(extractedHeaders)));
    }
}
