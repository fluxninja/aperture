package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPClient;
import com.fluxninja.aperture.instrumentation.ApertureSDKWrapper;
import com.fluxninja.aperture.instrumentation.Config;
import com.linecorp.armeria.client.WebClientBuilder;
import net.bytebuddy.asm.Advice;

public class ArmeriaClientAdvice {
    public static ApertureSDKWrapper wrapper = Config.newSDKWrapperFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This WebClientBuilder builder) {
        builder.decorator(
                ApertureHTTPClient.newDecorator(
                        wrapper.apertureSDK,
                        wrapper.controlPointName,
                        wrapper.rampMode,
                        wrapper.flowTimeout));
    }
}
