package com.fluxninja.aperture.servlet.javax;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.*;
import javax.servlet.ServletRequest;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletRequestWrapper;

public class ServletUtils {
    protected static int handleRejectedFlow(TrafficFlow flow) {
        try {
            flow.end(FlowStatus.Unset);
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }
        if (flow.checkResponse() != null
                && flow.checkResponse().hasDeniedResponse()
                && flow.checkResponse().getDeniedResponse().hasStatus()) {
            return flow.checkResponse().getDeniedResponse().getStatus().getCodeValue();
        }
        return 403;
    }

    protected static AttributeContext attributesFromRequest(ServletRequest req) {
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

        return addHttpAttributes(builder, req).build();
    }

    protected static ServletRequest updateHeaders(
            ServletRequest req, List<HeaderValueOption> newHeaders) {
        HttpServletRequest httpReq = (HttpServletRequest) req;
        Map<String, String> headerMap = new HashMap<>();
        for (HeaderValueOption option : newHeaders) {
            headerMap.put(option.getHeader().getKey(), option.getHeader().getValue());
        }
        return new HttpServletRequestWrapper(httpReq) {
            @Override
            public Enumeration<String> getHeaderNames() {
                Set<String> headerNames = new HashSet<>(Collections.list(super.getHeaderNames()));
                headerNames.addAll(headerMap.keySet());
                return Collections.enumeration(headerNames);
            }

            @Override
            public String getHeader(String name) {
                String header = headerMap.get(name);
                return header != null ? header : super.getHeader(name);
            }

            @Override
            public Enumeration<String> getHeaders(String name) {
                String header = headerMap.get(name);
                if (header != null) {
                    List<String> values = Arrays.asList(header.split(","));
                    return Collections.enumeration(values);
                } else {
                    return super.getHeaders(name);
                }
            }
        };
    }

    private static AttributeContext.Builder addHttpAttributes(
            AttributeContext.Builder builder, ServletRequest req) {
        HttpServletRequest request = (HttpServletRequest) req;
        Map<String, String> extractedHeaders = new HashMap<>();
        Enumeration<String> headers = request.getHeaderNames();
        while (headers.hasMoreElements()) {
            String headerKey = headers.nextElement();
            extractedHeaders.put(headerKey, request.getHeader(headerKey));
        }

        return builder.putContextExtensions("control-point", "ingress")
                .setRequest(
                        AttributeContext.Request.newBuilder()
                                .setHttp(
                                        AttributeContext.HttpRequest.newBuilder()
                                                .setMethod(request.getMethod())
                                                .setPath(request.getServletPath())
                                                .setHost(req.getRemoteHost())
                                                .setScheme(req.getScheme())
                                                .setSize(req.getContentLength())
                                                .setProtocol(req.getProtocol())
                                                .putAllHeaders(extractedHeaders)));
    }
}
