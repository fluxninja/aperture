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
    public static final String DEFAULT_AGENT_ADDRESS = "localhost:8089";
    public static final String DEFAULT_RAMP_MODE = "false";
    public static final String DEFAULT_CONTROL_POINT_NAME = "awesome_feature";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_ROOT_CERT = "";

    public static void main(String[] args) {
        String agentAddress = System.getenv("APERTURE_AGENT_ADDRESS");
        if (agentAddress == null) {
            agentAddress = DEFAULT_AGENT_ADDRESS;
        }
        String apiKey = System.getenv("APERTURE_API_KEY");
        if (apiKey == null) {
            apiKey = "";
        }
        String rampModeString = System.getenv("APERTURE_ENABLE_RAMP_MODE");
        if (rampModeString == null) {
            rampModeString = DEFAULT_RAMP_MODE;
        }
        boolean rampMode = Boolean.parseBoolean(rampModeString);

        String controlPointName = System.getenv("APERTURE_CONTROL_POINT_NAME");
        if (controlPointName == null) {
            controlPointName = DEFAULT_CONTROL_POINT_NAME;
        }
        String insecureGrpcString = System.getenv("APERTURE_AGENT_INSECURE");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String rootCertFile = System.getenv("APERTURE_ROOT_CERTIFICATE_FILE");
        if (rootCertFile == null) {
            rootCertFile = DEFAULT_ROOT_CERT;
        }

        ApertureSDK apertureSDK;
        try {
            apertureSDK =
                    ApertureSDK.builder()
                            .setAddress(agentAddress)
                            .setAPIKey(apiKey)
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
