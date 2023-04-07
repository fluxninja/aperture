package com.fluxninja.aperture.netty;

import com.fluxninja.aperture.sdk.Utils;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.http.HttpHeaders;
import io.netty.handler.codec.http.HttpRequest;
import io.netty.handler.codec.http.QueryStringDecoder;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import java.net.InetSocketAddress;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

public class NettyUtils {

    protected static HttpRequest updateHeaders(HttpRequest req, Map<String, String> newHeaders) {
        for (Map.Entry<String, String> newHeader : newHeaders.entrySet()) {
            String headerKey = newHeader.getKey().toLowerCase();
            String headerValue = newHeader.getValue();

            req.headers().add(headerKey, headerValue);
        }
        return req;
    }

    protected static CheckHTTPRequest checkRequestFromRequest(
            ChannelHandlerContext ctx, HttpRequest req) {
        Map<String, String> baggageLabels = new HashMap<>();

        for (Map.Entry<String, BaggageEntry> entry : Baggage.current().asMap().entrySet()) {
            String value;
            try {
                value =
                        URLDecoder.decode(
                                entry.getValue().getValue(), StandardCharsets.UTF_8.name());
            } catch (java.io.UnsupportedEncodingException e) {
                // This should never happen, as `StandardCharsets.UTF_8.name()` is a valid encoding
                throw new RuntimeException(e);
            }
            baggageLabels.put(entry.getKey(), value);
        }

        return addHttpAttributes(baggageLabels, ctx, req).build();
    }

    private static CheckHTTPRequest.Builder addHttpAttributes(
            Map<String, String> headers, ChannelHandlerContext ctx, HttpRequest req) {

        HttpHeaders originalHeaders = req.headers();
        for (Map.Entry<String, String> header : originalHeaders) {
            String headerKey = header.getKey();
            if (headerKey.startsWith(":")) {
                continue;
            }
            headers.put(headerKey, header.getValue());
        }

        String scheme = req.headers().get("scheme");
        if (scheme == null) {
            scheme = "";
        }
        String sizeStr = req.headers().get("content-length");
        int size = -1;
        if (sizeStr != null) {
            size = Integer.parseInt(sizeStr);
        }

        String sourceIp = null;
        int sourcePort = 0;
        String destinationIp = null;
        int destinationPort = 0;

        if (ctx.channel().remoteAddress() instanceof InetSocketAddress) {
            InetSocketAddress remoteAddress = (InetSocketAddress) ctx.channel().remoteAddress();
            sourceIp = remoteAddress.getAddress().getHostAddress();
            sourcePort = remoteAddress.getPort();
        }

        if (ctx.channel().localAddress() instanceof InetSocketAddress) {
            InetSocketAddress localAddress = (InetSocketAddress) ctx.channel().localAddress();
            destinationIp = localAddress.getAddress().getHostAddress();
            destinationPort = localAddress.getPort();
        }

        CheckHTTPRequest.Builder builder = CheckHTTPRequest.newBuilder();

        builder.setControlPoint("ingress")
                .setRequest(
                        CheckHTTPRequest.HttpRequest.newBuilder()
                                .setMethod(req.method().toString())
                                .setPath(new QueryStringDecoder(req.uri()).path())
                                .setHost(req.headers().get("host"))
                                .setScheme(scheme)
                                .setSize(size)
                                .setProtocol(req.protocolVersion().text())
                                .putAllHeaders(headers));

        if (sourceIp != null && sourcePort != 0) {
            builder.setSource(Utils.createSocketAddress(sourceIp, sourcePort, "TCP"));
        }
        if (destinationIp != null && destinationPort != 0) {
            builder.setDestination(
                    Utils.createSocketAddress(destinationIp, destinationPort, "TCP"));
        }
        return builder;
    }
}
