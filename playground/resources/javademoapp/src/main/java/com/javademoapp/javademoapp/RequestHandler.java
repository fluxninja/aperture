package com.javademoapp.javademoapp;

import com.google.gson.Gson;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.http.*;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.reactive.function.client.WebClient;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.LinkedBlockingQueue;


public class RequestHandler {
    private WebClient webClient;
    private LinkedBlockingQueue<WebClient> limitClients;
    private String hostname;
    private int concurrency;
    private Duration latency;
    private double rejectRatio;
    private final Logger log;

    public RequestHandler(WebClient webClient, String hostname, int concurrency, Duration latency, double rejectRatio) {
        this.webClient = webClient;
        this.limitClients = new LinkedBlockingQueue<>(concurrency);
        this.hostname = hostname;
        this.concurrency = concurrency;
        this.latency = latency;
        this.rejectRatio = rejectRatio;
        this.log = LogManager.getLogger(RequestHandler.class);
    }

    public void serveRequest(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        Map<String, String[]> labels = request.getParameterMap();
        int status = HttpServletResponse.SC_OK;

        try {
            if (labels.isEmpty()) {
                log.atWarn().log("Request body is empty");
                throw new IllegalArgumentException("Request body is empty");
            }
        } catch (IllegalArgumentException ex) {
            log.atWarn().log("Bad request", ex);
            writeErrorResponse(response, ex.getMessage(), HttpServletResponse.SC_BAD_REQUEST);
            return;
        }

        // NOTE - Harjot mentioned to try without span at first because metrics should be collected in javasdk

        // Randomly reject requests based on rejectRatio
        if (getRejectRatio() > 0 && Math.random() < getRejectRatio()) {
            status = HttpServletResponse.SC_SERVICE_UNAVAILABLE;
            return;
        }

        Request req = new Request();
        req.addChain(createChain(labels));
        for (SubrequestChain chain : req.getChains()) {
            if (chain.getSubrequests().isEmpty()) {
                log.atWarn().log("Request body is empty");
                throw new IllegalArgumentException("Request body is empty");
            }
            try {
                status = processChain(chain);
            } catch (IllegalArgumentException ex) {
                log.atWarn().log("Bad request", ex);
                writeErrorResponse(response, ex.getMessage(), HttpServletResponse.SC_BAD_REQUEST);
                return;
            } catch (Exception ex) {
                log.atError().log("Internal server error", ex);
                writeErrorResponse(response, ex.getMessage(), HttpServletResponse.SC_INTERNAL_SERVER_ERROR);
                return;
            }
            if (status != HttpServletResponse.SC_OK) {
                writeErrorResponse(response, "Bad request", status);
                return;
            }
        }

        // Write success response
        response.setHeader("Content-Type", "application/json");
        response.setStatus(status);
        writeSuccessResponse(response);
    }
    private void writeSuccessResponse(HttpServletResponse response) throws IOException {
        response.getWriter().write("{\"message\": \"success\"}");
    }

    private void writeErrorResponse(HttpServletResponse response, String message, int status) throws IOException {
        response.setHeader("Content-Type", "application/json");
        response.setStatus(status);
        response.getWriter().write("{\"error\": \"" + message + "\"}");
    }
    public double getRejectRatio() {
        return rejectRatio;
    }
    public int getConcurrency() {
        return concurrency;
    }
    public Duration getLatency() {
        return latency;
    }
    public SubrequestChain createChain(Map<String, String[]> labels) {

        SubrequestChain chain = new SubrequestChain();
        for (Map.Entry<String, String[]> entry : labels.entrySet()) {
            String key = entry.getKey();
            String[] values = entry.getValue();
            for (String value : values) {
                Subrequest subrequest = new Subrequest(key, value);
                chain.addSubrequest(subrequest);
            }
        }
        return chain;
    }

    public int processChain(SubrequestChain chain) throws InterruptedException {
        String requestForwardDestination, trimmedSubrequestChain;
        Request trimmedRequest;
        if (chain.getSubrequests().size() == 1) {
            return processRequest(chain.getSubrequests().get(0));
        }
        requestForwardDestination = chain.getSubrequests().get(1).getDestination();
        trimmedSubrequestChain = chain.getSubrequests().subList(1, chain.getSubrequests().size()).toString();
        trimmedRequest = new Request();
        trimmedRequest.addChain(createChain(Map.of(requestForwardDestination, new String[]{trimmedSubrequestChain})));
        return forwardRequest(requestForwardDestination, trimmedRequest);
    }

    public int processRequest(Subrequest subReq) throws InterruptedException {
        if (getConcurrency() > 0) {
            try {
                limitClients.put(this.webClient);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                log.atWarn().log("Interrupted while waiting to process request");
                return HttpServletResponse.SC_INTERNAL_SERVER_ERROR;
            } finally {
                try {
                    limitClients.take();
                } catch (InterruptedException e) {
                    log.atWarn().log("Interrupted while waiting to process request");
                    throw new RuntimeException(e);
                }
            }
        }

        if (!getLatency().isNegative() && !getLatency().isZero()) {
            try {
                // Fake workload
                Thread.sleep(Long.parseLong(getLatency().toString()));
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                log.atWarn().log("Interrupted while waiting to process request");
                return HttpServletResponse.SC_INTERNAL_SERVER_ERROR;
            }
        }
        return HttpServletResponse.SC_OK;
    }
    public int forwardRequest(String destinationHostname, Request requestBody) {
        String address = "http://" + destinationHostname + "/request";

        Gson gson = new Gson();
        String jsonRequest = gson.toJson(requestBody);

        RestTemplate restTemplate = new RestTemplate();
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<String> requestEntity = new HttpEntity<>(jsonRequest, headers);
        ResponseEntity<String> responseEntity = restTemplate.exchange(address, HttpMethod.POST, requestEntity, String.class);

        if (responseEntity.getStatusCode() == HttpStatus.OK) {
            return HttpServletResponse.SC_OK;
        } else {
            return HttpServletResponse.SC_INTERNAL_SERVER_ERROR;
        }
    }
}
