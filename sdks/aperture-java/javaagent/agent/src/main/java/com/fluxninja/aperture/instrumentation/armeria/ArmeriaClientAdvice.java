package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.armeria.ApertureHTTPClient;
import com.fluxninja.aperture.instrumentation.Config;
import com.fluxninja.aperture.sdk.ApertureSDK;
import com.linecorp.armeria.client.WebClientBuilder;
import net.bytebuddy.asm.Advice;

public class ArmeriaClientAdvice {
    public static ApertureSDK apertureSDK = Config.newSDKFromConfig();

    @Advice.OnMethodEnter
    public static void onEnter(@Advice.This WebClientBuilder builder) {
        builder.decorator(ApertureHTTPClient.newDecorator(apertureSDK));
    }
}
