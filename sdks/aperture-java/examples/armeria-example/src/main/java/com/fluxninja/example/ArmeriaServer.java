package com.fluxninja.example;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import com.linecorp.armeria.server.*;
import java.io.IOException;
import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public class ArmeriaServer {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_RAMP_MODE = "false";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

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
        String agentHost = System.getenv("APERTURE_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        String agentPort = System.getenv("APERTURE_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
        }
        String appPort = System.getenv("FN_APP_PORT");
        if (appPort == null) {
            appPort = DEFAULT_APP_PORT;
        }
        String rampModeString = System.getenv("FN_ENABLE_RAMP_MODE");
        if (rampModeString == null) {
            rampModeString = DEFAULT_RAMP_MODE;
        }
        boolean rampMode = Boolean.parseBoolean(rampModeString);

        String controlPointName = System.getenv("FN_CONTROL_POINT_NAME");
        if (controlPointName == null) {
            controlPointName = DEFAULT_CONTROL_POINT_NAME;
        }
        String insecureGrpcString = System.getenv("FN_INSECURE_GRPC");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = System.getenv("FN_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }

        ApertureSDK apertureSDK;
        try {
            apertureSDK =
                    ApertureSDK.builder()
                            .setHost(agentHost)
                            .setPort(Integer.parseInt(agentPort))
                            .useInsecureGrpc(insecureGrpc)
                            .setRootCertificateFile(rootCertFile)
                            .build();
        } catch (IOException e) {
            e.printStackTrace();
            return;
        }
        ServerBuilder serverBuilder = Server.builder();
        serverBuilder.http(Integer.parseInt(appPort));
        serverBuilder.service("/notsuper", createHelloHTTPService());
        serverBuilder.service("/health", createHealthService());
        serverBuilder.service("/connected", createConnectedHTTPService());

        ApertureHTTPService decoratedService =
                createHelloHTTPService()
                        .decorate(
                                ApertureHTTPService.newDecorator(
                                        apertureSDK,
                                        controlPointName,
                                        rampMode,
                                        Duration.ofMillis(1000)));
        serverBuilder.service("/super", decoratedService);

        Server server = serverBuilder.build();
        CompletableFuture<Void> future = server.start();
        future.join();
    }
}
