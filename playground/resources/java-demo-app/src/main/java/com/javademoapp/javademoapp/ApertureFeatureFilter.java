package com.javademoapp.javademoapp;

import java.io.IOException;
import java.time.Duration;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;

import javax.servlet.Filter;
import javax.servlet.FilterChain;
import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpHeaders;

import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;

public class ApertureFeatureFilter implements Filter {

    private ApertureSDK apertureSDK;
    private Logger log = LoggerFactory.getLogger(ApertureFeatureFilter.class);

    /**
     * This filter uses Aperture SDK to start a flow,
     * then checks if the flow was accepted before calling next filter in the chain.
     * If the flow is denied then an error response is sent back to the client.
     */
    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws ServletException, IOException {
        Map<String, String> labels = new HashMap<>();
        labels.put("app", "demoapp");
        labels.put("instance", System.getenv().getOrDefault("HOSTNAME", "instance-1"));
        labels.put("ip", req.getRemoteAddr());

        // See whether flow was accepted by Aperture Agent.
        try {
            HttpServletRequest request = (HttpServletRequest) req;
            HttpServletResponse response = (HttpServletResponse) res;

            // Add headers sent from k6 load generator
            HttpHeaders httpHeaders = Collections.list(request.getHeaderNames())
                    .stream()
                    .collect(Collectors.toMap(
                            Function.identity(),
                            h -> Collections.list(request.getHeaders(h)),
                            (oldValue, newValue) -> newValue,
                            HttpHeaders::new));
            String userType = httpHeaders.getFirst("User-Type");
            if (userType != null) {
                labels.put("user_type", userType);
            }
            String userID = httpHeaders.getFirst("User-Id");
            if (userID != null) {
                labels.put("user_id", userID);
            }

            log.debug("Starting Aperture SDK flow");
            Flow flow = this.apertureSDK.startFlow("awesomeFeature", labels);

            if (flow.accepted()) {
                log.debug("Flow accepted by Aperture Agent");
                chain.doFilter(request, response);
                flow.end(FlowStatus.OK);
            } else {
                log.debug("Flow rejected by Aperture Agent");
                response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Request denied");
                flow.end(FlowStatus.Error);
            }
        } catch (ApertureSDKException e) {
            log.error("Aperture SDK error: " + e.getMessage());
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
            String message = "Couldn't create aperture SDK";
            log.error(message);
            throw new ServletException(message);
        } catch (Exception e) {
            String message = "Invalid agent connection information";
            log.error(message);
            throw new ServletException(message);
        }
    }

    @Override
    public void destroy() {
    }
}
