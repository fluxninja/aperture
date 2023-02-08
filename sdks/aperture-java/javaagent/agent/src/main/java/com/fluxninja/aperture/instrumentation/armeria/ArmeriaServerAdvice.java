package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPService;
import com.fluxninja.aperture.instrumentation.Config;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.server.ServerBuilder;
import net.bytebuddy.asm.Advice;

public class ArmeriaServerAdvice {
    public static ApertureSDK apertureSDK = Config.newSDKFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This ServerBuilder builder) {
        builder.decorator(ApertureHTTPService.newDecorator(apertureSDK));
    }
}
