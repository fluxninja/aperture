package com.fluxninja.example.filter;

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

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;

public class ApertureFeatureFilter implements Filter {

    private ApertureSDK apertureSDK;

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain) throws ServletException, IOException {
        Map<String, String> labels = new HashMap<>();
        // do some business logic to collect labels
        labels.put("user", "kenobi");

        Flow flow = this.apertureSDK.startFlow("awesomeFeature", labels);
        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        // See whether flow was accepted by Aperture Agent.
        if (flow.accepted()) {
            try {
                chain.doFilter(request, response);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            }
        } else {
            try {
                response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Request denied");
                flow.end(FlowStatus.Error);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            }
        }
    }

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        String agentHost;
        String agentPort;
        try {
            agentHost = filterConfig.getInitParameter("agent_host");
            agentPort = filterConfig.getInitParameter("agent_port");
        } catch (Exception e) {
            throw new ServletException("Invalid agent connection information "
                    + filterConfig.getInitParameter("agent_host")
                    + ":"
                    + filterConfig.getInitParameter("agent_port"));
        }

        try {
            this.apertureSDK = ApertureSDK.builder()
                    .setHost(agentHost)
                    .setPort(Integer.parseInt(agentPort))
                    .setDuration(Duration.ofMillis(1000))
                    .build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            throw new ServletException("Couldn't create aperture SDK");
        }
    }

    @Override
    public void destroy() {}
}
