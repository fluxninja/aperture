package com.fluxninja.aperture.instrumentation.netty;

import com.fluxninja.aperture.netty.ApertureServerHandler;
import com.fluxninja.aperture.sdk.ApertureSDK;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpServerCodec;
import net.bytebuddy.asm.Advice;

import java.time.Duration;

public class NettyServerAdvice {
    public static ApertureSDK apertureSDK = apertureSDKFromConfig();

    @Advice.OnMethodExit
    public static void onExit(
            @Advice.This ChannelPipeline pipeline,
            @Advice.Argument(1) String handlerName,
            @Advice.Argument(2) ChannelHandler handler
            ) {
        if (handler instanceof HttpServerCodec || handler instanceof HttpRequestDecoder) {
            // only add the aperture handler after the HttpRequestDecoder or HttpServerCodec
            ApertureServerHandler apertureHandler = new ApertureServerHandler(apertureSDK);
            String hname = handlerName;
            if (hname == null) {
                ChannelHandlerContext ctx = pipeline.context(handler);
                hname = ctx.name();
            }
            pipeline.addAfter(hname, apertureHandler.getClass().getName(), apertureHandler);
        }
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
