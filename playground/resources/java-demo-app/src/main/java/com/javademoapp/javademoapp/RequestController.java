package com.javademoapp.javademoapp;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fluxninja.aperture.sdk.ApertureSDKException;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.core.env.Environment;
import org.springframework.http.*;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.List;
import java.util.ArrayList;
import java.util.Collections;
import java.util.stream.Collectors;
import java.util.function.Function;
import java.util.concurrent.Semaphore;
import java.util.concurrent.atomic.AtomicInteger;

@RestController
public class RequestController {
    public static final String DEFAULT_HOST = "localhost";
	public static final String DEFAULT_AGENT_PORT = "8089";

    private  int concurrency = Integer.parseInt(System.getenv().getOrDefault("CONCURRENCY", "10"));
    private  Duration latency = Duration.ofMillis(Long.parseLong(System.getenv().getOrDefault("LATENCY", "50")));
    private  double rejectRatio = Double.parseDouble(System.getenv().getOrDefault("REJECT_RATIO", "0.05"));
    private  Logger log = LoggerFactory.getLogger(RequestController.class);

    // Semaphore for limiting concurrent clients
    private Semaphore limitClients = new Semaphore(concurrency);
    private ApertureFeatureFilter apertureFilter = new ApertureFeatureFilter();

    @RequestMapping(value = "/super", method = RequestMethod.GET)
    public String hello() {
        String message = "Hello World";
        log.info(message);
        return message;
    }

    @RequestMapping(value = "/super2", method = RequestMethod.GET)
    public String hello2() {
        String message = "Hello World 2";
        log.info(message);
        return message;
    }

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

    @Bean
    public FilterRegistrationBean<ApertureFeatureFilter> apertureFeatureFilter(Environment env){
        FilterRegistrationBean<ApertureFeatureFilter> registrationBean = new FilterRegistrationBean<>();

        registrationBean.setFilter(apertureFilter);
        registrationBean.addUrlPatterns("/request");
        registrationBean.addInitParameter("agent_host", System.getenv().getOrDefault("FN_AGENT_HOST", DEFAULT_HOST));
        registrationBean.addInitParameter("agent_port", System.getenv().getOrDefault("FN_AGENT_PORT", DEFAULT_AGENT_PORT));

        return registrationBean;
    }

    @GetMapping("/")
    public String index() {
        String message = "Your request has been received!";
        log.info(message);
        return message;
    }

    @PostMapping("/request")
    public String handlePostRequest(@RequestBody String payload, HttpServletRequest request, HttpServletResponse response) throws ApertureSDKException {
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
                HttpHeaders::new
            ));
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
                // TODO Add check for req Dest != Hostname
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
        if (!latency.isZero() && !latency.isNegative()) {
            try {
                // Fake Overload
                Thread.sleep(latency.toMillis());
            } catch (InterruptedException e) {
                log.error("Error while sleeping: " + e.getMessage());
            }
        }

        // Release the semaphore
        if (limitClients != null) {
            limitClients.release();
        }
        return "Success";
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
