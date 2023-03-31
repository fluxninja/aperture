package com.fluxninja.aperture.netty;

import com.fluxninja.aperture.sdk.Utils;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
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
import java.util.List;
import java.util.Map;

public class NettyUtils {

    protected static HttpRequest updateHeaders(
            HttpRequest req, List<HeaderValueOption> newHeaders) {
        for (HeaderValueOption newHeader : newHeaders) {
            String headerKey = newHeader.getHeader().getKey().toLowerCase();
            String headerValue = newHeader.getHeader().getValue();
            if (!newHeader.getKeepEmptyValue() && headerValue.isEmpty()) {
                req.headers().remove(headerKey);
                continue;
            }
            if (newHeader.getAppend().getValue()) {
                req.headers().add(headerKey, headerValue);
            } else {
                req.headers().set(headerKey, headerValue);
            }
        }
        return req;
    }

    protected static AttributeContext attributesFromRequest(
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

        AttributeContext.Builder builder = AttributeContext.newBuilder();
        builder.putAllContextExtensions(baggageLabels);

        return addHttpAttributes(builder, ctx, req).build();
    }

    private static AttributeContext.Builder addHttpAttributes(
            AttributeContext.Builder builder, ChannelHandlerContext ctx, HttpRequest req) {

        Map<String, String> extractedHeaders = new HashMap<>();
        HttpHeaders headers = req.headers();
        for (Map.Entry<String, String> header : headers) {
            String headerKey = header.getKey();
            if (headerKey.startsWith(":")) {
                continue;
            }
            extractedHeaders.put(headerKey, header.getValue());
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
        String destinationIp = null;

        if (ctx.channel().remoteAddress() instanceof InetSocketAddress) {
            InetSocketAddress remoteAddress = (InetSocketAddress) ctx.channel().remoteAddress();
            sourceIp = remoteAddress.getAddress().getHostAddress();
        }

        if (ctx.channel().localAddress() instanceof InetSocketAddress) {
            InetSocketAddress localAddress = (InetSocketAddress) ctx.channel().localAddress();
            destinationIp = localAddress.getAddress().getHostAddress();
        }

        builder.putContextExtensions("control-point", "ingress")
                .setRequest(
                        AttributeContext.Request.newBuilder()
                                .setHttp(
                                        AttributeContext.HttpRequest.newBuilder()
                                                .setMethod(req.method().toString())
                                                .setPath(new QueryStringDecoder(req.uri()).path())
                                                .setHost(req.headers().get("host"))
                                                .setScheme(scheme)
                                                .setSize(size)
                                                .setProtocol(req.protocolVersion().text())
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
