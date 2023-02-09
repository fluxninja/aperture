package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPClient;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.WebClientBuilder;
import net.bytebuddy.asm.Advice;

import java.time.Duration;

public class ArmeriaClientAdvice {
    public static ApertureSDK apertureSDK = apertureSDKFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This WebClientBuilder builder) {
        builder.decorator(ApertureHTTPClient.newDecorator(apertureSDK));
    }

    private static ApertureSDK apertureSDKFromConfig() {
        String host = System.getProperty("aperture.agent.hostname");
        String port = System.getProperty("aperture.agent.port");
        if (host == null) {
            host = "localhost";
        }
        if (port == null) {
            port = "8089";
        }
        ApertureSDK sdk;
        try {
            sdk = ApertureSDK.builder()
                    .setHost(host)
                    .setPort(Integer.parseInt(port))
                    .setDuration(Duration.ofMillis(1000))
                    .build();
        } catch (Exception e) {
            e.printStackTrace();
            throw new RuntimeException("fail");
        }
        return sdk;
    }
}
