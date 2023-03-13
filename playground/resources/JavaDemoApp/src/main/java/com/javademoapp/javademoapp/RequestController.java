package com.javademoapp.javademoapp;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
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
import java.util.concurrent.atomic.AtomicInteger;

@RestController
public class RequestController {

    private static final String FORWARDING_INSTANCE_HEADER = "X-Forwarding-Instance";
    private static final AtomicInteger INSTANCE_COUNTER = new AtomicInteger(0);
    private static final String APP_NAME = "Demo-App";

    Logger logger = LoggerFactory.getLogger(getClass());
   // private int concurrency = Integer.parseInt(System.getenv("CONCURRENCY"));
   // private Duration latency = Duration.ofMillis(Integer.parseInt(System.getenv("LATENCY")));
   // private double rejectRatio = Double.parseDouble(System.getenv("REJECT_RATIO"));
    private  int concurrency = 10;
    private  Duration latency = Duration.ofMillis(50);
    private  double rejectRatio = 0.05;

    @RequestMapping(value = "/super", method = RequestMethod.GET)
    // /super endpoint is protected by a Filter created using Aperture SDK feature flow
    public String hello() {
        return "Hello World";
    }

    @RequestMapping(value = "/super2", method = RequestMethod.GET)
    // /super2 endpoint is protected by imported, ready-to-use Aperture Filter
    public String hello2() {
        return "Hello World 2";
    }

    @RequestMapping(value = "/health", method = RequestMethod.GET)
    public String health() {
        return "Healthy";
    }

    @RequestMapping(value = "/connected", method = RequestMethod.GET)
    public String connected() {
        return "";
    }

    @Bean
    public FilterRegistrationBean<ApertureFeatureFilter> apertureFeatureFilter(Environment env){
        FilterRegistrationBean<ApertureFeatureFilter> registrationBean = new FilterRegistrationBean<>();

        registrationBean.setFilter(new ApertureFeatureFilter());
        registrationBean.addUrlPatterns("/super2");

        String agentHost = env.getProperty("FN_AGENT_HOST");
        String agentPort = env.getProperty("FN_AGENT_PORT");
        registrationBean.addInitParameter("agent_host", agentHost);
        registrationBean.addInitParameter("agent_port", agentPort);

        return registrationBean;
    }

    @GetMapping("/")
    public String index() {
        return "Your request has been received!";
    }

    @PostMapping("/process")
    public String handlePostRequest(@RequestBody String payload, HttpServletRequest request, HttpServletResponse response) throws ApertureSDKException {
        // Create a hashmap of labels and values from the payload
        Map<String,String> labels = new HashMap<String,String>();
        labels.put("app", APP_NAME);
        labels.put("instance", getInstanceName());
        labels.put("ip", request.getRemoteAddr());
        // TODO Add other relevant information to the labels


        // Randomly reject requests
        if (rejectRatio > 0 && Math.random() < rejectRatio) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            logger.warn("Request rejected.");
            return "Request rejected.";
        }

        ObjectMapper objectMapper = new ObjectMapper();
        try {
            Request requestObj = objectMapper.readValue(payload, Request.class);
            for (SubrequestChain chain : requestObj.getRequest()) {
                if (chain.getSubrequest().size() == 0) {
                    response.setStatus(HttpStatus.BAD_REQUEST.value());
                    response.getWriter().write("Empty Chain");
                    logger.warn("Empty chain.");
                    return "Empty Chain";
                }
                String requestDestination = chain.getSubrequest().get(0).getDestination();
                // TODO: Add Check for Hostname != request Destination
                return processChain(chain);
            }
            response.setStatus(HttpStatus.OK.value());
            response.getWriter().write(payload);
            logger.info("Request processed successfully.");
        } catch (JsonMappingException e) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            logger.error("Error while mapping JSON.", e);
            throw new RuntimeException(e);
        } catch (JsonProcessingException e) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            logger.error("Error while processing JSON.", e);
            throw new RuntimeException(e);
        } catch (IOException e) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            logger.error("IO exception occurred.", e);
            throw new RuntimeException(e);
        }
        return "Success";
    }

    private String processChain(SubrequestChain chain){
        if (chain.getSubrequest().size() == 1) {
            return processRequest(chain.getSubrequest().get(0));
        }
        SubrequestChain trimmedSubrequestChain = new SubrequestChain();
        for (int i = 1; i < chain.getSubrequest().size(); i++) {
            trimmedSubrequestChain.addSubrequest(chain.getSubrequest().get(i));
        }
        Request trimmedRequest = new Request();
        trimmedRequest.addRequest(trimmedSubrequestChain);
        String requestForwardingDestination = chain.getSubrequest().get(1).getDestination();
        return forwardRequest(trimmedRequest, requestForwardingDestination);
    }
    private String processRequest(Subrequest request){
        if (concurrency > 0) {
            // TODO Add logic to handle concurrency and limit clients
        }
        if (!latency.isNegative() && !latency.isZero()) {
            try {
                // Fake Overload
                Thread.sleep(latency.toMillis());
            } catch (InterruptedException e) {
                System.out.println("Error while sleeping: " + e.getMessage());
            }
        }
        return "Success";
    }
    private String forwardRequest(Request request, String destination) {
        String address = String.format("http://%s/request", destination);

        RestTemplate restTemplate = new RestTemplate();
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.set(FORWARDING_INSTANCE_HEADER, getInstanceName());

        // Get Json of request, if error then return error
        String requestJson;
        try {
            requestJson = new ObjectMapper().writeValueAsString(request);
        } catch (JsonProcessingException e) {
            return "Error while parsing request: " + e.getMessage();
        }

        HttpEntity<String> entity = new HttpEntity<>(requestJson, headers);
        ResponseEntity<String> response = restTemplate.exchange(address, HttpMethod.POST, entity, String.class);
        if (response.getStatusCode() != HttpStatus.OK) {
            return "Error while forwarding request: " + response.getStatusCode();
        }

        return "Success";
    }
    private String getInstanceName() {
        return APP_NAME + "-" + INSTANCE_COUNTER.incrementAndGet();
    }

}
