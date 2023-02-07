package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.ServerBuilder;
import net.bytebuddy.asm.Advice;

import java.io.File;
import java.time.Duration;

public class ArmeriaServerAdvice {
    public static ApertureSDK apertureSDK = apertureSDKFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This ServerBuilder builder) {
        builder.decorator(ApertureHTTPService.newDecorator(apertureSDK));
    }

    private static ApertureSDK apertureSDKFromConfig() {
        String host = System.getProperty("aperture.agent.hostname");
        String port = System.getProperty("aperture.agent.port");
        String configFileName = System.getProperty("aperture.javaagent.config.file");
        if (configFileName != null) {
            File configFile = new File(configFileName);
            // TODO: check if file exists

        }
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
