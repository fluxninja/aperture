package com.fluxninja.example;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.common.HttpRequest;
import com.linecorp.armeria.common.HttpResponse;
import com.linecorp.armeria.server.AbstractHttpService;
import com.linecorp.armeria.server.HttpService;
import com.linecorp.armeria.server.Server;
import com.linecorp.armeria.server.ServerBuilder;
import com.linecorp.armeria.server.ServiceRequestContext;
import java.io.IOException;
import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public class ArmeriaServer {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_ADDRESS = "localhost:8089";
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

    public static String getEnv(String key, String defaultValue) {
        String value = System.getenv(key);
        return value != null ? value : defaultValue;
    }

    public static void main(String[] args) {
        String agentHost = getEnv("APERTURE_AGENT_ADDRESS", DEFAULT_AGENT_ADDRESS);
        String agentAPIKey = getEnv("APERTURE_AGENT_API_KEY", "");

        String appPort = getEnv("APERTURE_APP_PORT", DEFAULT_APP_PORT);

        String rampModeString = getEnv("APERTURE_ENABLE_RAMP_MODE", DEFAULT_RAMP_MODE);
        boolean rampMode = Boolean.parseBoolean(rampModeString);

        String controlPointName = getEnv("APERTURE_CONTROL_POINT_NAME", DEFAULT_CONTROL_POINT_NAME);

        String insecureGrpcString = getEnv("APERTURE_AGENT_INSECURE", DEFAULT_INSECURE_GRPC);

        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = getEnv("APERTURE_ROOT_CERTIFICATE_FILE", DEFAULT_ROOT_CERT);

        ApertureSDK apertureSDK;
        try {
            apertureSDK =
                    ApertureSDK.builder()
                            .setAddress(agentHost)
                            .setAgentAPIKey(agentAPIKey)
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
