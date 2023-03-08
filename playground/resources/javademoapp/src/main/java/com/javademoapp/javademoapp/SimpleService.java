package com.javademoapp.javademoapp;


import org.apache.catalina.Context;
import org.apache.catalina.LifecycleException;
import org.apache.catalina.startup.Tomcat;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import javax.servlet.Servlet;
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

     public SimpleService(String hostname, int port, int concurrency, Duration latency, double rejectionRate) {
          this.hostname = hostname;
          this.port = port;
          this.concurrency = concurrency;
          this.latency = latency;
          this.rejectionRate = rejectionRate;
          this.client = WebClient.create();
          this.log = LogManager.getLogger(SimpleService.class);
     }

    public void run(HttpServletRequest request, HttpServletResponse response) throws LifecycleException {
        RequestHandler requestHandler = new RequestHandler(client, hostname, concurrency, latency, rejectionRate);

        // Use WebClient to make an outbound GET request on local port 8080
        Mono<String> result = this.client.get()
                .uri("http://localhost:8081")
                .retrieve()
                .bodyToMono(String.class);

        // handle the incoming request
        try {
            requestHandler.serveRequest(request, response);
        } catch (ServletException e) {
            log.atError().log("Error serving request", e);
            throw new RuntimeException(e);
        } catch (IOException e) {
            log.atError().log("Error serving request", e);
            throw new RuntimeException(e);
        }

        // NOTE - Use of same port may make them clash
        // Start a Tomcat server on the specified port
        Tomcat tomcat = new Tomcat();
        tomcat.setPort(port);
        Context context = tomcat.addContext("", null);
        Tomcat.addServlet(context, "Simple Service", (Servlet) requestHandler).setAsyncSupported(true);
        context.addServletMappingDecoded("/", "Simple Service");
        tomcat.start();

        // Intercept Traffic through FluxNinja Middleware, it will show up as a control point
        // TODO Try middleware or java agent later
        // wait for request to be served
        tomcat.getServer().await();
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
}
