package com.fluxninja.aperture.instrumentation.springboot;

import com.fluxninja.aperture.instrumentation.Config;
import com.fluxninja.aperture.sdk.ApertureSDK;
import net.bytebuddy.description.method.MethodDescription;
import net.bytebuddy.asm.Advice;

import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletResponse;

public class SpringBootAdvice {
    public static ApertureSDK apertureSDK = Config.newSDKFromConfig();

    @Advice.OnMethodEnter
    public static void enter(@Advice.Origin MethodDescription method) {
        System.out.println("Hello, World! Intercepting request to " + method.getName());
    }

    @Advice.OnMethodExit(onThrowable = Throwable.class)
    public static void exit(@Advice.Thrown Throwable throwable) {
        if (throwable == null) {
            System.out.println("Forwarding request down the chain...");
        } else {
            System.err.println("Exception thrown: " + throwable.getMessage());
        }
    }

    @Advice.OnMethodExit
    public static void after(@Advice.Argument(value = 1, readOnly = false) ServletResponse response) {
        System.out.println("Modifying response...");
        // Modify the HTTP response code to 403
        HttpServletResponse httpServletResponse = (HttpServletResponse) response;
        httpServletResponse.setStatus(org.springframework.http.HttpStatus.FORBIDDEN.value());
    }
}
