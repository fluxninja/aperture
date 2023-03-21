package com.javademoapp.javademoapp;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
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
import java.util.concurrent.Semaphore;
import java.util.concurrent.atomic.AtomicInteger;

@RestController
public class RequestController {
    public static final String DEFAULT_HOST = "localhost";
	public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_CONCURRENCY = "10";
	public static final String DEFAULT_LATENCY = "50";
	public static final String DEFAULT_REJECT_RATIO = "0.05";

    private  int concurrency = Integer.parseInt(System.getenv().getOrDefault("CONCURRENCY", DEFAULT_CONCURRENCY));
    private  Duration latency = Duration.ofMillis(Long.parseLong(System.getenv().getOrDefault("LATENCY", DEFAULT_LATENCY)));
    private  double rejectRatio = Double.parseDouble(System.getenv().getOrDefault("REJECT_RATIO", DEFAULT_REJECT_RATIO));

    // Semaphore for limiting concurrent clients
    private Semaphore limitClients = new Semaphore(concurrency);

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
        registrationBean.addInitParameter("agent_host", System.getenv().getOrDefault("FN_AGENT_HOST", DEFAULT_HOST));
        registrationBean.addInitParameter("agent_port", System.getenv().getOrDefault("FN_AGENT_PORT", DEFAULT_AGENT_PORT));

        return registrationBean;
    }

    @GetMapping("/")
    public String index() {
        return "Your request has been received!";
    }

    @PostMapping("/process")
    public String handlePostRequest(@RequestBody String payload, HttpServletRequest request, HttpServletResponse response) throws ApertureSDKException {
        Map<String,String> labels = new HashMap<String,String>();

        labels.put("app", "demoapp");
        labels.put("instance", getInstanceName());
        labels.put("ip", request.getRemoteAddr());

        // Randomly reject requests
        if (rejectRatio > 0 && Math.random() < rejectRatio) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            return "Request rejected.";
        }

        ObjectMapper objectMapper = new ObjectMapper();
        try {
            // Processing the request object's subrequests
            Request requestObj = objectMapper.readValue(payload, Request.class);
            for (SubrequestChain chain : requestObj.getRequest()) {
                if (chain.getSubrequest().size() == 0) {
                    response.setStatus(HttpStatus.BAD_REQUEST.value());
                    response.getWriter().write("Empty Chain");
                    return "Empty Chain";
                }
                String requestDestination = chain.getSubrequest().get(0).getDestination();
                // TODO Add check for req Dest != Hostname
                return processChain(chain);
            }

            // If all subrequests were processed successfully, return success message
            response.setStatus(HttpStatus.OK.value());
            response.setStatus(HttpStatus.OK.value());
            response.getWriter().write(payload);

        } catch (Exception e) {
            response.setStatus(HttpStatus.BAD_REQUEST.value());
            return "Error occurred" + e.getMessage();
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
        // Limit concurrent clients
        if (limitClients != null && concurrency > 0) {
            try {
                limitClients.acquire();
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            }
        }
        if (!latency.isNegative() && !latency.isZero()) {
            try {
                // Fake Overload
                Thread.sleep(latency.toMillis());
            } catch (InterruptedException e) {
                System.out.println("Error while sleeping: " + e.getMessage());
            }
        }

        // Release the semaphore
        if (limitClients != null) {
            limitClients.release();
        }
        return "Success";
    }
    private String forwardRequest(Request request, String destination) {
        String requestJson;
        try {
            requestJson = new ObjectMapper().writeValueAsString(request);
        } catch (JsonProcessingException e) {
            return "Error while parsing request: " + e.getMessage();
        }

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.set("X-Forwarding-Instance", getInstanceName());

        String address = String.format("http://%s/request", destination);
        RestTemplate restTemplate = new RestTemplate();
        HttpEntity<String> entity = new HttpEntity<>(requestJson, headers);

        ResponseEntity<String> response = restTemplate.exchange(address, HttpMethod.POST, entity, String.class);
        if (response.getStatusCode() != HttpStatus.OK) {
            return "Error while forwarding request: " + response.getStatusCode();
        }

        return "Success";
    }
    private String getInstanceName() {
        return System.getenv().getOrDefault("HOSTNAME", DEFAULT_HOST);
    }
    public void setRejectRatio(double rejectRatio) {
        this.rejectRatio = rejectRatio;
    }
}
