package com.fluxninja.aperture.tomcat7;

import java.io.IOException;
import java.time.Duration;

import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.fluxninja.aperture.sdk.*;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;

public class ApertureFilter implements Filter {

    private ApertureSDK apertureSDK;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain) throws ServletException, IOException {
        AttributeContext attributes = ServletUtils.attributesFromRequest(req);

        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        String path = request.getServletPath();
        TrafficFlow flow = this.apertureSDK.startTrafficFlow(path, attributes);

        if (flow.ignored()) {
            chain.doFilter(request, response);
            return;
        }

        if (flow.accepted()) {
            try {
                chain.doFilter(request, response);
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
            try {
                flow.end(FlowStatus.Error);
            } catch (Exception e) {
                throw new ServletException(e);
            }
            response.sendError(code, "Request denied");
        }
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        String agentHost;
        String agentPort;
        String timeoutMs;
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
            timeoutMs = filterConfig.getInitParameter("timeout_ms");
        } catch (Exception e) {
            throw new ServletException("Invalid agent connection information "
                    + filterConfig.getInitParameter("agent_host")
                    + ":"
                    + filterConfig.getInitParameter("agent_port"));
        }

        try {
            ApertureSDKBuilder builder = ApertureSDK.builder();
            builder.setHost(agentHost);
            builder.setPort(Integer.parseInt(agentPort));
            if (timeoutMs != null) {
                builder.setDuration(Duration.ofMillis(Integer.parseInt(timeoutMs)));
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
