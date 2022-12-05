package com.fluxninja.aperture.tomcat7;

import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.FlowStatus;
import com.fluxninja.aperture.sdk.TrafficFlow;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;

import io.opentelemetry.api.baggage.Baggage;
import io.opentelemetry.api.baggage.BaggageEntry;

import javax.servlet.ServletRequest;
import javax.servlet.http.HttpServletRequest;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ServletUtils {
    protected static int handleRejectedFlow(TrafficFlow flow) {
        try {
            flow.end(FlowStatus.Unset);
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }
        if (flow.checkResponse().hasDeniedResponse() && flow.checkResponse().getDeniedResponse().hasStatus()) {
            return flow.checkResponse().getDeniedResponse().getStatus().getCodeValue();
        }
        return 403;
    }
    protected static AttributeContext attributesFromRequest(ServletRequest req) {
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
        builder.putAllContextExtensions(baggageLabels);

        return addHttpAttributes(builder, req).build();
    }

    protected static ServletRequest updateHeaders(ServletRequest req, List<HeaderValueOption> newHeaders) {
        // TODO: Update headers of ServletRequest (probably need to create requestWrapper)
        return req;
    }

    private static AttributeContext.Builder addHttpAttributes(AttributeContext.Builder builder, ServletRequest req) {
        HttpServletRequest request = (HttpServletRequest) req;
        Map<String, String> extractedHeaders = new HashMap<>();
        Enumeration<String> headers = request.getHeaderNames();
        while (headers.hasMoreElements()) {
            String headerKey = headers.nextElement();
            extractedHeaders.put(headerKey, request.getHeader(headerKey));
        }

        return builder
                .putContextExtensions("control-point", "ingress")
                .setRequest(AttributeContext.Request.newBuilder()
                        .setHttp(AttributeContext.HttpRequest.newBuilder()
                                .setMethod(request.getMethod())
                                .setPath(request.getServletPath())
                                .setHost(req.getRemoteHost())
                                .setScheme(req.getScheme())
                                .setSize(req.getContentLength())
                                .setProtocol(req.getProtocol())
                                .putAllHeaders(extractedHeaders)));
    }
}
