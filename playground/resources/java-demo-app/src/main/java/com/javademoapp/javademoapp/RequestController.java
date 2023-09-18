package com.javademoapp.javademoapp;

import java.io.IOException;
import java.net.URL;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Semaphore;
import java.util.concurrent.CountDownLatch;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Callable;
import java.util.concurrent.Future;
import java.util.concurrent.ExecutionException;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.core.env.Environment;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fluxninja.aperture.servlet.javax.*;

@RestController
public class RequestController {
    public static final String DEFAULT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";

    private int concurrency = Integer.parseInt(System.getenv().getOrDefault("CONCURRENCY", "10"));
    private Duration latency = Duration.ofMillis(Long.parseLong(System.getenv().getOrDefault("LATENCY", "50")));
    private double rejectRatio = Double.parseDouble(System.getenv().getOrDefault("REJECT_RATIO", "0.05"));
    private int cpuLoad = Integer.parseInt(System.getenv().getOrDefault("CPU_LOAD", "0"));
    private Logger log = LoggerFactory.getLogger(RequestController.class);
    private String hostname = System.getenv().getOrDefault("HOSTNAME", DEFAULT_HOST);
    private final AtomicInteger ongoingRequests = new AtomicInteger(0);

    // Semaphore for limiting concurrent clients
    private Semaphore limitClients = new Semaphore(concurrency);
    private ApertureFilter apertureFilter = new ApertureFilter();

    @RequestMapping(value = "/health", method = RequestMethod.GET)
    public String health() {
        String message = "Healthy";
        log.info(message);
        return message;
    }

    @RequestMapping(value = "/connected", method = RequestMethod.GET)
    public String connected() {
        String message = "Connected OK";
        log.info(message);
        return message;
    }

    @GetMapping("/")
    public String index() {
        String message = "Your request has been received!";
        log.info(message);
        return message;
    }

    @PostMapping("/request")
    public String handlePostRequest(@RequestBody String payload, HttpServletRequest request,
            HttpServletResponse response) {
        // Randomly reject requests
        if (rejectRatio > 0 && Math.random() < rejectRatio) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            return "Request rejected";
        }

        try {
            HttpHeaders httpHeaders = Collections.list(request.getHeaderNames())
                    .stream()
                    .collect(Collectors.toMap(
                            Function.identity(),
                            h -> Collections.list(request.getHeaders(h)),
                            (oldValue, newValue) -> newValue,
                            HttpHeaders::new));
            Request requestObj = new ObjectMapper().readValue(payload, Request.class);
            List<List<Subrequest>> chains = requestObj.getRequest();
            for (List<Subrequest> chain : chains) {
                if (chain.size() == 0) {
                    String msg = "Empty Chain";
                    response.setStatus(HttpStatus.BAD_REQUEST.value());
                    response.getWriter().write(msg);
                    log.info(msg);
                    return msg;
                }
                String requestDestination = chain.get(0).getDestination();
                if (!requestDestination.startsWith("http")) {
                    requestDestination = String.format("http://%s", requestDestination);
                }

                URL requestDestinationURL = new URL(requestDestination);
                if (!requestDestinationURL.getHost().equals(hostname)) {
                    response.setStatus(HttpStatus.BAD_REQUEST.value());
                    String msg = "Invalid message destination";
                    log.error(String.format("%s: %s - %s", msg, hostname, requestDestinationURL.getHost()));
                    return msg;
                }
                return processChain(chain, httpHeaders);
            }

            // If all subrequests were processed successfully, return success message
            response.setStatus(HttpStatus.OK.value());
            response.getWriter().write(payload);
        } catch (Exception e) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            String msg = "Error occurred: " + e.getMessage();
            log.error(msg);
            return msg;
        }

        return "Success";
    }

    @Bean
    public FilterRegistrationBean<ApertureFilter> apertureFeatureFilter(Environment env) {
        FilterRegistrationBean<ApertureFilter> registrationBean = new FilterRegistrationBean<>();

        registrationBean.setFilter(apertureFilter);
        registrationBean.addUrlPatterns("/request");
        registrationBean.addInitParameter("agent_host", System.getenv().getOrDefault("APERTURE_AGENT_HOST", DEFAULT_HOST));
        registrationBean.addInitParameter("agent_port",
                System.getenv().getOrDefault("APERTURE_AGENT_PORT", DEFAULT_AGENT_PORT));
        registrationBean.addInitParameter("control_point_name", "awesomeFeature");
        registrationBean.addInitParameter("enable_fail_open", "true");
        registrationBean.addInitParameter("insecure_grpc", "true");

        return registrationBean;
    }

    private String processChain(List<Subrequest> chain, HttpHeaders httpHeaders) {
        if (chain.size() == 1) {
            return processRequest(chain.get(0));
        }

        List<Subrequest> trimmedChain = new ArrayList<>();
        for (int i = 1; i < chain.size(); i++) {
            trimmedChain.add(chain.get(i));
        }

        Request trimmedRequest = new Request();
        trimmedRequest.addRequest(trimmedChain);
        String requestForwardingDestination = chain.get(1).getDestination();

        return forwardRequest(trimmedRequest, requestForwardingDestination, httpHeaders);
    }

    private String processRequest(Subrequest request) {
        // Limit concurrent clients
        if (concurrency > 0 && limitClients != null) {
            try {
                limitClients.acquire();
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            }
        }

        try {
            ongoingRequests.incrementAndGet();
            if (!latency.isZero() && !latency.isNegative()) {
                if (cpuLoad > 0) {
                    // Simulate CPU load by busy waiting
                    int numCores = Runtime.getRuntime().availableProcessors();

                    // Calculate busy wait and sleep durations based on cpuLoad and ongoing requests
                    double adjustedLoad = (double) cpuLoad / 100 * (double) ongoingRequests.get();
                    if (adjustedLoad > 100) {
                        adjustedLoad = 100;
                    }
                    double totalDuration = latency.toMillis();
                    long busyWaitDuration = (long) (totalDuration * adjustedLoad / 100.0);
                    long sleepDuration = (long) (totalDuration * (100.0 - adjustedLoad) / 100.0);

                    ExecutorService executor = Executors.newFixedThreadPool(numCores);
                    List<Callable<Void>> tasks = new ArrayList<>();
                    for (int i = 0; i < numCores; i++) {
                        tasks.add(() -> {
                            long startTime = System.currentTimeMillis();
                            while (System.currentTimeMillis() - startTime < latency.toMillis()) {
                                busyWait(busyWaitDuration);
                                try {
                                    Thread.sleep(sleepDuration);
                                } catch (InterruptedException e) {
                                    log.error("Error while sleeping: " + e.getMessage());
                                }
                            }
                            return null;
                        });
                    }

                    try {
                        List<Future<Void>> futures = executor.invokeAll(tasks);
                        for (Future<Void> future : futures) {
                            future.get();
                        }
                    } catch (InterruptedException | ExecutionException e) {
                        log.error("Error while waiting for threads: " + e.getMessage());
                    }
                } else {
                    try {
                        Thread.sleep(latency.toMillis());
                    } catch (InterruptedException e) {
                        log.error("Error while sleeping: " + e.getMessage());
                    }
                }
            }

            // Release the semaphore
            if (limitClients != null) {
                limitClients.release();
            }
            return "Success";
        } finally {
            ongoingRequests.decrementAndGet();
        }
    }

    private void busyWait(long duration) {
        long startTime = System.currentTimeMillis();
        List<Object> objects = new ArrayList<>();
        while (System.currentTimeMillis() - startTime < duration) {
            // Try to make garbage collector busy
            objects.add(new Object());
        }
    }

    private String forwardRequest(Request request, String destination, HttpHeaders httpHeaders) {
        String requestJson;
        try {
            requestJson = new ObjectMapper().writeValueAsString(request);
        } catch (JsonProcessingException e) {
            String msg = "Error while parsing request: " + e.getMessage();
            log.error(msg);
            return msg;
        }

        HttpHeaders headers = httpHeaders;
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.set("X-Forwarding-Instance", getInstanceName());

        String address = String.format("http://%s", destination);
        RestTemplate restTemplate = new RestTemplate();
        HttpEntity<String> entity = new HttpEntity<>(requestJson, headers);

        ResponseEntity<String> response = restTemplate.exchange(address, HttpMethod.POST, entity, String.class);
        if (response.getStatusCode() != HttpStatus.OK) {
            String msg = "Error while forwarding request: " + response.getStatusCode();
            log.error(msg);
            return msg;
        }

        return "Success";
    }

    private String getInstanceName() {
        return System.getenv().getOrDefault("HOSTNAME", DEFAULT_HOST);
    }
}
