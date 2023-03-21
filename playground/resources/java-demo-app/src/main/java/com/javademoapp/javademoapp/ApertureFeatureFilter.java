package com.javademoapp.javademoapp;

import java.io.IOException;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;

public class ApertureFeatureFilter implements Filter {

    private ApertureSDK apertureSDK;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain) throws ServletException, IOException {
        Map<String, String> labels = new HashMap<>();
        labels.put("app", "Demo-App");
        labels.put("instance", "instance-1");
        labels.put("ip", req.getRemoteAddr());

        Flow flow = this.apertureSDK.startFlow("awesomeFeature", labels);
        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        // See whether flow was accepted by Aperture Agent.
        try {
            if (flow.accepted()) {
                chain.doFilter(request, response);
                flow.end(FlowStatus.OK);
            } else {
                response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Request denied");
                flow.end(FlowStatus.Error);
            }
        } catch (ApertureSDKException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        String agentHost;
        String agentPort;
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
            this.apertureSDK = ApertureSDK.builder()
                    .setHost(agentHost)
                    .setPort(Integer.parseInt(agentPort))
                    .setDuration(Duration.ofMillis(1000))
                    .build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            throw new ServletException("Couldn't create aperture SDK");
        } catch (Exception e) {
            throw new ServletException("Invalid agent connection information ");
        }
    }

    @Override
    public void destroy() {}
}
