package com.fluxninja.aperture.servlet.javax;

import com.fluxninja.aperture.sdk.*;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import java.io.IOException;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

public class ApertureFilter implements Filter {

    private ApertureSDK apertureSDK;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws ServletException, IOException {
        CheckHTTPRequest checkRequest = ServletUtils.checkRequestFromRequest(req);

        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        String path = request.getServletPath();
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(path, checkRequest);

        if (flow.ignored()) {
            chain.doFilter(request, response);
            return;
        }

        if (flow.accepted()) {
            try {
                Map<String, String> newHeaders = new HashMap<>();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersMap();
                }
                ServletRequest newRequest = ServletUtils.updateHeaders(request, newHeaders);
                chain.doFilter(newRequest, response);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                // ending flow failed
                e.printStackTrace();
            } catch (Exception e) {
                try {
                    flow.end(FlowStatus.Error);
                } catch (ApertureSDKException ae) {
                    e.printStackTrace();
                    ae.printStackTrace();
                }
                throw e;
            }
        } else {
            int code = ServletUtils.handleRejectedFlow(flow);
            response.sendError(code, "Request denied");
        }
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        String agentHost;
        String agentPort;
        String timeoutMs;
        boolean insecureGrpc;
        String rootCertificateFile;
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
            timeoutMs = filterConfig.getInitParameter("timeout_ms");
            insecureGrpc = Boolean.parseBoolean(filterConfig.getInitParameter("insecure_grpc"));
            rootCertificateFile = filterConfig.getInitParameter("root_certificate_file");
        } catch (Exception e) {
            throw new ServletException("Could not read config parameters", e);
        }

        try {
            ApertureSDKBuilder builder = ApertureSDK.builder();
            builder.setHost(agentHost);
            builder.setPort(Integer.parseInt(agentPort));
            if (timeoutMs != null) {
                builder.setDuration(Duration.ofMillis(Integer.parseInt(timeoutMs)));
            }
            builder.useInsecureGrpc(insecureGrpc);
            if (rootCertificateFile != null && !rootCertificateFile.isEmpty()) {
                builder.setRootCertificateFile(rootCertificateFile);
            }

            this.apertureSDK = builder.build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            throw new ServletException("Couldn't create aperture SDK");
        }
    }

    @Override
    public void destroy() {}
}
