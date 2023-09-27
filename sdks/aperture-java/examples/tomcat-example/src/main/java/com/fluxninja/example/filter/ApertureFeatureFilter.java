package com.fluxninja.example.filter;

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

public class ApertureFeatureFilter implements Filter {

    private ApertureSDK apertureSDK;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws ServletException, IOException {
        Map<String, String> labels = new HashMap<>();
        // do some business logic to collect labels
        labels.put("user", "kenobi");

        Flow flow = this.apertureSDK.startFlow("awesomeFeature", labels, false);
        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        // Check whether flow was accepted by Aperture Agent
        try {
            if (flow.shouldRun()) {
                chain.doFilter(request, response);
            } else {
                flow.setStatus(FlowStatus.Unset);
                response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Request denied");
            }
        } catch (Exception e) {
            flow.setStatus(FlowStatus.Error);
            throw e;
        } finally {
            flow.end();
        }
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        String agentHost;
        String agentPort;
        boolean insecureGrpc;
        String rootCertificateFile;
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
            insecureGrpc = Boolean.parseBoolean(filterConfig.getInitParameter("insecure_grpc"));
            rootCertificateFile = filterConfig.getInitParameter("root_certificate_file");
        } catch (Exception e) {
            throw new ServletException("Could not read config parameters", e);
        }

        try {
            this.apertureSDK =
                    ApertureSDK.builder()
                            .setHost(agentHost)
                            .setPort(Integer.parseInt(agentPort))
                            .setFlowTimeout(Duration.ofMillis(1000))
                            .useInsecureGrpc(insecureGrpc)
                            .setRootCertificateFile(rootCertificateFile)
                            .build();
        } catch (IOException e) {
            e.printStackTrace();
            throw new ServletException("Couldn't create aperture SDK", e);
        }
    }

    @Override
    public void destroy() {}
}
