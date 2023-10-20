package com.fluxninja.aperture.servlet.javax;

import com.fluxninja.aperture.sdk.*;
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
    private String controlPointName;
    private boolean rampMode;
    private Duration flowTimeout;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws ServletException, IOException {
        TrafficFlowRequest trafficFlowRequest =
                ServletUtils.trafficFlowRequestFromRequest(req, controlPointName, flowTimeout);

        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        String path = request.getServletPath();
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(trafficFlowRequest);

        if (flow.ignored()) {
            chain.doFilter(request, response);
            return;
        }

        FlowDecision flowDecision = flow.getDecision();
        boolean flowAccepted =
                (flowDecision == FlowDecision.Accepted
                        || (flowDecision == FlowDecision.Unreachable && this.rampMode));

        if (flowAccepted) {
            try {
                Map<String, String> newHeaders = new HashMap<>();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersMap();
                }
                ServletRequest newRequest = ServletUtils.updateHeaders(request, newHeaders);
                chain.doFilter(newRequest, response);
            } catch (Exception e) {
                flow.setStatus(FlowStatus.Error);
                throw e;
            } finally {
                flow.end();
            }
        } else {
            ServletUtils.handleRejectedFlow(flow, response);
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
            insecureGrpc = Boolean.parseBoolean(filterConfig.getInitParameter("insecure_grpc"));
            rootCertificateFile = filterConfig.getInitParameter("root_certificate_file");
            ignoredPaths = filterConfig.getInitParameter("ignored_paths");
            ignoredPathsRegex =
                    Boolean.parseBoolean(
                            filterConfig.getInitParameter("ignored_paths_match_regex"));

            timeoutMs = filterConfig.getInitParameter("timeout_ms");
            if (timeoutMs != null) {
                this.flowTimeout = Duration.ofMillis(Integer.parseInt(timeoutMs));
            } else {
                this.flowTimeout = Constants.DEFAULT_RPC_TIMEOUT;
            }
            this.rampMode = Boolean.parseBoolean(filterConfig.getInitParameter("enable_ramp_mode"));

        } catch (Exception e) {
            throw new ServletException("Could not read config parameters", e);
        }

        if (initControlPointName == null || initControlPointName.trim().isEmpty()) {
            throw new IllegalArgumentException("Control Point name must be set");
        }
        controlPointName = initControlPointName;

        ApertureSDKBuilder builder = ApertureSDK.builder();
        builder.setHost(agentHost);
        builder.setPort(Integer.parseInt(agentPort));
        builder.useInsecureGrpc(insecureGrpc);
        if (rootCertificateFile != null && !rootCertificateFile.isEmpty()) {
            try {
                builder.setRootCertificateFile(rootCertificateFile);
            } catch (IOException e) {
                e.printStackTrace();
                throw new ServletException("Couldn't create aperture SDK", e);
            }
        }
        builder.addIgnoredPaths(ignoredPaths);
        builder.setIgnoredPathsMatchRegex(ignoredPathsRegex);

        this.apertureSDK = builder.build();
    }

    @Override
    public void destroy() {}
}
