package com.fluxninja.example;

import com.fluxninja.aperture.armeria.ApertureHTTPClient;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.Clients;
import com.linecorp.armeria.client.WebClient;
import com.linecorp.armeria.common.HttpResponse;
import java.io.IOException;
import java.time.Duration;

public class ArmeriaClient {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_RAMP_MODE = "false";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static void main(String[] args) {
        String agentHost = System.getenv("APERTURE_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        String agentPort = System.getenv("APERTURE_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
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

        WebClient client =
                Clients.builder("http://localhost:8080")
                        .decorator(
                                ApertureHTTPClient.newDecorator(
                                        apertureSDK,
                                        controlPointName,
                                        rampMode,
                                        Duration.ofMillis(1000)))
                        .build(WebClient.class);

        HttpResponse res = client.get("notsuper");
        System.out.println(res);
    }
}
