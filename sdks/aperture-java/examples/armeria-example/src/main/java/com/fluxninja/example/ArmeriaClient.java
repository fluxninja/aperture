package com.fluxninja.example;

import com.fluxninja.aperture.armeria.ApertureHTTPClient;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.fluxninja.aperture.sdk.ApertureSDKException;
import com.linecorp.armeria.client.Clients;
import com.linecorp.armeria.client.WebClient;
import com.linecorp.armeria.common.HttpResponse;
import java.time.Duration;

public class ArmeriaClient {

    public static final String DEFAULT_APP_PORT = "8080";
    public static final String DEFAULT_AGENT_HOST = "localhost";
    public static final String DEFAULT_AGENT_PORT = "8089";
    public static final String DEFAULT_INSECURE_GRPC = "true";
    public static final String DEFAULT_SSL_CERT = "";

    public static void main(String[] args) {
        String agentHost = System.getenv("FN_AGENT_HOST");
        if (agentHost == null) {
            agentHost = DEFAULT_AGENT_HOST;
        }
        String agentPort = System.getenv("FN_AGENT_PORT");
        if (agentPort == null) {
            agentPort = DEFAULT_AGENT_PORT;
        }
        String insecureGrpcString = System.getenv("FN_INSECURE_GRPC");
        if (insecureGrpcString == null) {
            insecureGrpcString = DEFAULT_INSECURE_GRPC;
        }
        boolean insecureGrpc = Boolean.parseBoolean(insecureGrpcString);

        String sslCertFile = System.getenv("FN_SSL_CERTIFICATE_FILE");
        if (sslCertFile == null) {
            sslCertFile = DEFAULT_SSL_CERT;
        }

        ApertureSDK apertureSDK;
        try {
            apertureSDK =
                    ApertureSDK.builder()
                            .setHost(agentHost)
                            .setPort(Integer.parseInt(agentPort))
                            .setDuration(Duration.ofMillis(1000))
                            .useInsecureGrpc(insecureGrpc)
                            .setCACertificateFile(sslCertFile)
                            .build();
        } catch (ApertureSDKException e) {
            e.printStackTrace();
            return;
        }

        WebClient client =
                Clients.builder("http://localhost:8080")
                        .decorator(ApertureHTTPClient.newDecorator(apertureSDK))
                        .build(WebClient.class);

        HttpResponse res = client.get("notsuper");
        System.out.println(res);
    }
}
