package com.javademoapp.javademoapp;


import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.fluxninja.aperture.sdk.Flow;
import com.fluxninja.aperture.sdk.FlowStatus;
import org.apache.catalina.LifecycleException;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import javax.servlet.FilterConfig;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.Duration;

public class SimpleService {
    private final String hostname;
    private final int port;
    private final int concurrency;
    private final Duration latency;
    private final double rejectionRate;
    private final WebClient client;
    private final Logger log;
    private final ServletInitializer tomcat;

    private ApertureSDK apertureSDK;

     public SimpleService(String hostname, int port, int concurrency, Duration latency, double rejectionRate) {
          this.hostname = hostname;
          this.port = port;
          this.concurrency = concurrency;
          this.latency = latency;
          this.rejectionRate = rejectionRate;
          this.client = WebClient.create();
          this.log = LogManager.getLogger(SimpleService.class);
          this.tomcat = new ServletInitializer();

     }

    public void run(HttpServletRequest request, HttpServletResponse response) throws LifecycleException, ServletException {
        RequestHandler requestHandler = new RequestHandler(client, hostname, concurrency, latency, rejectionRate);

        Flow flow = this.apertureSDK.startFlow("TODO", request.getTrailerFields());

        // Use WebClient to make an outbound GET request on local port 8080
        Mono<String> result = this.client.get()
                .uri("http://localhost:8081")
                .retrieve()
                .bodyToMono(String.class);

        if (flow.accepted()) {
            try {
                requestHandler.serveRequest(request, response);
                flow.end(FlowStatus.OK);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        } else {
            try {
                response.sendError(HttpServletResponse.SC_UNAUTHORIZED, "Request denied");
                flow.end(FlowStatus.Error);
            } catch (ApertureSDKException e) {
                e.printStackTrace();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }

        /*ServletContext servletContext = request.getServletContext();
        try {
            this.tomcat.onStartup(servletContext);
        } catch (ServletException e) {
            throw new RuntimeException(e);
        }*/

        // Intercept Traffic through FluxNinja Middleware, it will show up as a control point
        // TODO Try middleware or java agent later
        // wait for request to be served
        //tomcat.getServer().await();
        // Wait for the result and write it to the response
        result.subscribe(resp -> {
            try {
                response.getWriter().write(resp);
            } catch (IOException e) {
                log.atError().log("Error writing response", e);
                throw new RuntimeException(e);
            }
        });
    }

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

    public void destroy() {}
}
