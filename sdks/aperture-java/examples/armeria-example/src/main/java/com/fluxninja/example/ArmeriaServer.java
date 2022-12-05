package com.fluxninja.example;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import com.linecorp.armeria.server.*;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public class ArmeriaServer {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";

    public static HttpService createHelloHTTPService() {
        return new AbstractHttpService() {
            @Override
            protected HttpResponse doGet(ServiceRequestContext ctx, HttpRequest req) {
                return HttpResponse.of("Hello, world!");
            }
        };
    }
    public static HttpService createHealthService() {
        return new AbstractHttpService() {
            @Override
            protected HttpResponse doGet(ServiceRequestContext ctx, HttpRequest req) {
                return HttpResponse.of("Healthy");
            }
        };
    }
    public static HttpService createConnectedHTTPService() {
        return new AbstractHttpService() {
            @Override
            protected HttpResponse doGet(ServiceRequestContext ctx, HttpRequest req) {
                return HttpResponse.of("");
            }
        };
    }

    public static void main(String[] args) {
        String agentHost = System.getenv("FN_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        String agentPort = System.getenv("FN_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
        }

        ApertureSDK apertureSDK;
        try {
            apertureSDK = ApertureSDK.builder()
                    .setHost(agentHost)
                    .setPort(Integer.parseInt(agentPort))
                    .setDuration(Duration.ofMillis(1000))
                    .build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            return;
        }
        ServerBuilder serverBuilder = Server.builder();
        serverBuilder.http(8080);
        serverBuilder.service("/notsuper", createHelloHTTPService());
        serverBuilder.service("/health", createHealthService());
        serverBuilder.service("/connected", createConnectedHTTPService());

        ApertureHTTPService decoratedService = createHelloHTTPService()
            .decorate(ApertureHTTPService.newDecorator(apertureSDK));
        serverBuilder.service("/super", decoratedService);

        Server server = serverBuilder.build();
        CompletableFuture<Void> future = server.start();
        future.join();
    }
}
