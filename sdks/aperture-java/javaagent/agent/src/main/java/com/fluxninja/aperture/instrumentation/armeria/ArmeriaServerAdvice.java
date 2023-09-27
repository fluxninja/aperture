package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.instrumentation.ApertureSDKWrapper;
import com.fluxninja.aperture.instrumentation.Config;
import com.linecorp.armeria.server.ServerBuilder;
import net.bytebuddy.asm.Advice;

public class ArmeriaServerAdvice {
    public static ApertureSDKWrapper wrapper = Config.newSDKWrapperFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This ServerBuilder builder) {
        builder.decorator(
                ApertureHTTPService.newDecorator(
                        wrapper.apertureSDK, wrapper.controlPointName, wrapper.rampMode));
    }
}
