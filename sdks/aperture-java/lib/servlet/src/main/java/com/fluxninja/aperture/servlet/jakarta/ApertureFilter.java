package com.fluxninja.aperture.servlet.jakarta;

import com.fluxninja.aperture.sdk.*;
import com.fluxninja.generated.aperture.flowcontrol.checkhttp.v1.CheckHTTPRequest;
import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.FilterConfig;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

public class ApertureFilter implements Filter {

    private ApertureSDK apertureSDK;
    private String controlPointName;
    private boolean failOpen;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws ServletException, IOException {
        CheckHTTPRequest checkRequest = ServletUtils.checkRequestFromRequest(req, controlPointName);

        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        String path = request.getServletPath();
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(path, checkRequest);

        if (flow.ignored()) {
            chain.doFilter(request, response);
            return;
        }

        FlowResult flowResult = flow.result();
        boolean flowAccepted =
                (flowResult == FlowResult.Accepted
                        || (flowResult == FlowResult.Unreachable && this.failOpen));

        if (flowAccepted) {
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
        String initControlPointName;
        String timeoutMs;
        boolean insecureGrpc;
        String rootCertificateFile;
        String ignoredPaths;
        boolean ignoredPathsRegex;
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
            initControlPointName = filterConfig.getInitParameter("control_point_name");
            timeoutMs = filterConfig.getInitParameter("timeout_ms");
            insecureGrpc = Boolean.parseBoolean(filterConfig.getInitParameter("insecure_grpc"));
            rootCertificateFile = filterConfig.getInitParameter("root_certificate_file");
            ignoredPaths = filterConfig.getInitParameter("ignored_paths");
            ignoredPathsRegex =
                    Boolean.parseBoolean(
                            filterConfig.getInitParameter("ignored_paths_match_regex"));

            this.failOpen = Boolean.parseBoolean(filterConfig.getInitParameter("enable_fail_open"));

        } catch (Exception e) {
            throw new ServletException("Could not read config parameters", e);
        }

        if (initControlPointName == null || initControlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must be set");
        }
        controlPointName = initControlPointName;

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
            builder.addIgnoredPaths(ignoredPaths);
            builder.setIgnoredPathsMatchRegex(ignoredPathsRegex);

            this.apertureSDK = builder.build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            throw new ServletException("Couldn't create aperture SDK");
        }
    }

    @Override
    public void destroy() {}
}
