package com.fluxninja.aperture.servlet.jakarta;

import com.fluxninja.aperture.sdk.*;
import com.fluxninja.generated.envoy.service.auth.v3.AttributeContext;
import com.fluxninja.generated.envoy.service.auth.v3.HeaderValueOption;
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
import java.util.ArrayList;
import java.util.List;

public class ApertureFilter implements Filter {

    private ApertureSDK apertureSDK;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws ServletException, IOException {
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
                List<HeaderValueOption> newHeaders = new ArrayList<>();
                if (flow.checkResponse() != null) {
                    newHeaders = flow.checkResponse().getOkResponse().getHeadersList();
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
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
            timeoutMs = filterConfig.getInitParameter("timeout_ms");
        } catch (Exception e) {
            throw new ServletException(
                    "Invalid agent connection information "
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
