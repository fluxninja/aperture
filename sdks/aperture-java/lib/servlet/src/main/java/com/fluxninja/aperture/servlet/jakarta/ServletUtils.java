package com.fluxninja.aperture.servlet.jakarta;

import com.fluxninja.aperture.sdk.*;
import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletRequestWrapper;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.*;

public class ServletUtils {
    protected static void handleRejectedFlow(TrafficFlow flow, HttpServletResponse response)
            throws IOException {
        try {
            flow.end(FlowStatus.Unset);
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }

        int code = 403;
        if (flow.checkResponse() != null && flow.checkResponse().hasDeniedResponse()) {
            code = flow.getRejectionHttpStatusCode();
            Map<String, String> headers = flow.checkResponse().getDeniedResponse().getHeadersMap();
            for (Map.Entry<String, String> entry : headers.entrySet()) {
                response.setHeader(entry.getKey(), entry.getValue());
            }
        }

        response.sendError(code, "Request denied");
    }

    protected static TrafficFlowRequest trafficFlowRequestFromRequest(
            ServletRequest req, String controlPointName) {
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

        TrafficFlowRequestBuilder builder = addHttpAttributes(baggageLabels, req);
        builder.setControlPoint(controlPointName);
        return builder.build();
    }

    protected static ServletRequest updateHeaders(
            ServletRequest req, Map<String, String> newHeaders) {
        HttpServletRequest httpReq = (HttpServletRequest) req;
        Map<String, String> headerMap = new HashMap<>(newHeaders);
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

    private static TrafficFlowRequestBuilder addHttpAttributes(
            Map<String, String> headers, ServletRequest req) {
        HttpServletRequest request = (HttpServletRequest) req;
        Enumeration<String> originalHeaders = request.getHeaderNames();
        while (originalHeaders.hasMoreElements()) {
            String headerKey = originalHeaders.nextElement();
            headers.put(headerKey, request.getHeader(headerKey));
        }

        String sourceIp = req.getRemoteAddr();
        int sourcePort = req.getRemotePort();
        String destinationIp = req.getLocalAddr();
        int destinationPort = req.getLocalPort();

        TrafficFlowRequestBuilder builder = TrafficFlowRequest.newBuilder();

        builder.setHttpMethod(request.getMethod())
                .setHttpPath(request.getServletPath())
                .setHttpHost(req.getRemoteHost())
                .setHttpScheme(req.getScheme())
                .setHttpSize(req.getContentLength())
                .setHttpProtocol(req.getProtocol())
                .setHttpHeaders(headers);

        if (sourceIp != null) {
            builder.setSource(sourceIp, sourcePort, "TCP");
        }
        if (destinationIp != null) {
            builder.setDestination(destinationIp, destinationPort, "TCP");
        }
        return builder;
    }
}
