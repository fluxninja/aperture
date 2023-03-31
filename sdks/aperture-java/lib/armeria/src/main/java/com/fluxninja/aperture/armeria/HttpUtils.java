package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.aperture.sdk.Utils;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
import com.linecorp.armeria.common.*;
import io.netty.util.AsciiString;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import java.net.InetSocketAddress;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class HttpUtils {
    protected static HttpStatus handleRejectedFlow(TrafficFlow flow) {
        try {
            flow.end(FlowStatus.Unset);
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }
        if (flow.checkResponse() != null
                && flow.checkResponse().hasDeniedResponse()
                && flow.checkResponse().getDeniedResponse().hasStatus()) {
            int httpStatusCode =
                    flow.checkResponse().getDeniedResponse().getStatus().getCodeValue();
            return HttpStatus.valueOf(httpStatusCode);
        }
        return HttpStatus.FORBIDDEN;
    }

    protected static Map<String, String> labelsFromRequest(HttpRequest req) {
        Map<String, String> labels = new HashMap<>();
        RequestHeaders headers = req.headers();
        for (Map.Entry<AsciiString, String> header : headers) {
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

    protected static AttributeContext attributesFromRequest(RequestContext ctx, HttpRequest req) {
        Map<String, String> baggageLabels = new HashMap<>();

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
            baggageLabels.put(entry.getKey(), value);
        }

        AttributeContext.Builder builder = AttributeContext.newBuilder();
        builder.putAllContextExtensions(baggageLabels);

        return addHttpAttributes(builder, ctx, req).build();
    }

    // getAppend is deprecated but agent does set it, so we should use it
    @SuppressWarnings("deprecation")
    protected static HttpRequest updateHeaders(
            HttpRequest req, List<HeaderValueOption> newHeaders) {
        RequestHeadersBuilder newHeadersBuilder = req.headers().toBuilder();
        for (HeaderValueOption newHeader : newHeaders) {
            String headerKey = newHeader.getHeader().getKey().toLowerCase();
            String headerValue = newHeader.getHeader().getValue();
            if (!newHeader.getKeepEmptyValue() && headerValue.isEmpty()) {
                newHeadersBuilder = newHeadersBuilder.removeAndThen(headerKey);
                continue;
            }
            if (newHeader.getAppend().getValue()) {
                newHeadersBuilder = newHeadersBuilder.add(headerKey, headerValue);
            } else {
                newHeadersBuilder = newHeadersBuilder.set(headerKey, headerValue);
            }
        }
        return req.withHeaders(newHeadersBuilder.build());
    }

    private static AttributeContext.Builder addHttpAttributes(
            AttributeContext.Builder builder, RequestContext ctx, HttpRequest req) {
        Map<String, String> extractedHeaders = new HashMap<>();
        RequestHeaders headers = req.headers();
        for (Map.Entry<AsciiString, String> header : headers) {
            String headerKey = header.getKey().toString();
            if (headerKey.startsWith(":")) {
                continue;
            }
            extractedHeaders.put(headerKey, header.getValue());
        }

        java.net.SocketAddress remoteSocket = ctx.remoteAddress();
        String sourceIp = null;
        if (remoteSocket instanceof InetSocketAddress) {
            InetSocketAddress remoteAddress = (InetSocketAddress) remoteSocket;
            sourceIp = remoteAddress.getAddress().getHostAddress();
        }

        java.net.SocketAddress localSocket = ctx.localAddress();
        String destinationIp = null;
        if (localSocket instanceof InetSocketAddress) {
            InetSocketAddress localAddress = (InetSocketAddress) localSocket;
            destinationIp = localAddress.getAddress().getHostAddress();
        }

        builder.putContextExtensions("control-point", "ingress")
                .setRequest(
                        AttributeContext.Request.newBuilder()
                                .setHttp(
                                        AttributeContext.HttpRequest.newBuilder()
                                                .setMethod(req.method().toString())
                                                .setPath(req.path())
                                                .setHost(req.authority())
                                                .setScheme(req.scheme())
                                                .setSize(headers.contentLength())
                                                .setProtocol("HTTP/2")
                                                .putAllHeaders(extractedHeaders)));
        if (sourceIp != null) {
            builder.setSource(Utils.peerFromAddress(sourceIp));
        }
        if (destinationIp != null) {
            builder.setDestination(Utils.peerFromAddress(destinationIp));
        }
        return builder;
    }
}
