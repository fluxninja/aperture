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

        WebClient client = Clients.builder("http://localhost:8080")
                .decorator(ApertureHTTPClient.newDecorator(apertureSDK))
                .build(WebClient.class);

        HttpResponse res = client.get("notsuper");
        System.out.println(res);
    }
}
