package com.fluxninja.aperture.armeria;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.generated.envoy.service.auth.v3.Address;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
import com.fluxninja.generated.envoy.service.auth.v3.SocketAddress;
import com.linecorp.armeria.common.*;
import com.linecorp.armeria.server.ServiceRequestContext;
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
        System.out.println("NEW REQUEST HEADER EXTRACTION");
        for (Map.Entry<AsciiString, String> header : headers) {
            String headerKey = header.getKey().toString();
            System.out.println(headerKey+ ": "+ header.getValue());
            if (headerKey.startsWith(":")) {
                continue;
            }
            if (headerKey.contains("orwarded")) {
                System.out.println("Might've found source IP");
            }
            if (headerKey.contains("erver")) {
                System.out.println("Might've found destination IP");
            }
            extractedHeaders.put(headerKey, header.getValue());
        }
        System.out.println("HEADER EXTRACTION ENDED");

        java.net.SocketAddress remoteSocket = ctx.remoteAddress();
        String sourceIp = null;
        String sourceHostname = null;
        if (remoteSocket instanceof InetSocketAddress) {
            InetSocketAddress remoteAddress = (InetSocketAddress) remoteSocket;
            sourceIp = remoteAddress.getAddress().getHostAddress();
            sourceHostname = remoteAddress.getHostString();
        }

        java.net.SocketAddress localSocket = ctx.localAddress();
        String destinationIp = null;
        String destinationHostname = null;
        if (localSocket instanceof InetSocketAddress) {
            InetSocketAddress localAddress = (InetSocketAddress) localSocket;
            destinationIp = localAddress.getAddress().getHostAddress();
            destinationHostname = localAddress.getHostString();
        }

        System.out.println("SOURCE IP: " + sourceIp);
        System.out.println("SOURCE HOSTNAME: " + sourceHostname);
        System.out.println("DESTINATION IP: " + destinationIp);
        System.out.println("DESTINATION HOSTNAME: " + destinationHostname);


        // get source and dest from headers:
        // X-Forwarded-For
        // X-Armeria-Server-Address

        AttributeContext.Peer source = AttributeContext.Peer.newBuilder()
                .setAddress(
                        Address.newBuilder()
                                .setSocketAddress(
                                    SocketAddress.newBuilder()
                                        .setAddress(sourceIp)
                                        .build()
                                )
                                .build()).build();
        AttributeContext.Peer destination = AttributeContext.Peer.newBuilder()
                .setAddress(
                        Address.newBuilder()
                                .setSocketAddress(
                                    SocketAddress.newBuilder()
                                        .setAddress(destinationIp)
                                        .build()
                                )
                                .build()).build();

        return builder.putContextExtensions("control-point", "ingress")
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
                                                .putAllHeaders(extractedHeaders)))
                .setSource(source)
                .setDestination(destination);
    }
}
