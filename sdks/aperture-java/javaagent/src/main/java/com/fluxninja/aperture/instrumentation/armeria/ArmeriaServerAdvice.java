package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.ServerBuilder;
import net.bytebuddy.asm.Advice;

public class ArmeriaServerAdvice {
    public static ApertureSDK apertureSDK = apertureSDKFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This ServerBuilder builder) {
        builder.decorator(ApertureHTTPService.newDecorator(apertureSDK));
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
            sdk = ApertureSDK.builder().setHost(host).setPort(Integer.parseInt(port)).build();
        } catch (Exception e) {
            e.printStackTrace();
            throw new RuntimeException("fail");
        }
        return sdk;
    }
}
