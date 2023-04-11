package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.aperture.sdk.Utils;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import com.linecorp.armeria.common.*;
import io.netty.util.AsciiString;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import java.net.InetSocketAddress;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
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
                && flow.checkResponse().getDeniedResponse().getStatus() != 0) {
            int httpStatusCode = flow.checkResponse().getDeniedResponse().getStatus();
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

    protected static CheckHTTPRequest checkRequestFromRequest(
            RequestContext ctx, HttpRequest req, String controlPointName) {
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

        return addHttpAttributes(baggageLabels, ctx, req, controlPointName).build();
    }

    protected static HttpRequest updateHeaders(HttpRequest req, Map<String, String> newHeaders) {
        RequestHeadersBuilder newHeadersBuilder = req.headers().toBuilder();
        for (Map.Entry<String, String> newHeader : newHeaders.entrySet()) {
            String headerKey = newHeader.getKey().toLowerCase();
            String headerValue = newHeader.getValue();
            newHeadersBuilder = newHeadersBuilder.add(headerKey, headerValue);
        }
        return req.withHeaders(newHeadersBuilder.build());
    }

    private static CheckHTTPRequest.Builder addHttpAttributes(
            Map<String, String> headers,
            RequestContext ctx,
            HttpRequest req,
            String controlPointName) {
        RequestHeaders originalHeaders = req.headers();
        for (Map.Entry<AsciiString, String> header : originalHeaders) {
            String headerKey = header.getKey().toString();
            if (headerKey.startsWith(":")) {
                continue;
            }
            headers.put(headerKey, header.getValue());
        }

        java.net.SocketAddress remoteSocket = ctx.remoteAddress();
        String sourceIp = null;
        int sourcePort = 0;
        if (remoteSocket instanceof InetSocketAddress) {
            InetSocketAddress remoteAddress = (InetSocketAddress) remoteSocket;
            sourceIp = remoteAddress.getAddress().getHostAddress();
            sourcePort = remoteAddress.getPort();
        }

        java.net.SocketAddress localSocket = ctx.localAddress();
        String destinationIp = null;
        int destinationPort = 0;
        if (localSocket instanceof InetSocketAddress) {
            InetSocketAddress localAddress = (InetSocketAddress) localSocket;
            destinationIp = localAddress.getAddress().getHostAddress();
            destinationPort = localAddress.getPort();
        }

        CheckHTTPRequest.Builder builder = CheckHTTPRequest.newBuilder();

        builder.setControlPoint(controlPointName)
                .setRequest(
                        CheckHTTPRequest.HttpRequest.newBuilder()
                                .setMethod(req.method().toString())
                                .setPath(req.path())
                                .setHost(req.authority())
                                .setScheme(req.scheme())
                                .setSize(originalHeaders.contentLength())
                                .setProtocol("HTTP/2")
                                .putAllHeaders(headers));

        if (sourceIp != null) {
            builder.setSource(Utils.createSocketAddress(sourceIp, sourcePort, "TCP"));
        }
        if (destinationIp != null) {
            builder.setDestination(
                    Utils.createSocketAddress(destinationIp, destinationPort, "TCP"));
        }
        return builder;
    }
}
