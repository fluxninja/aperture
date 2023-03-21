package com.fluxninja.aperture.instrumentation.netty;

import com.fluxninja.aperture.instrumentation.Config;
import com.fluxninja.aperture.netty.ApertureServerHandler;
import com.fluxninja.aperture.sdk.ApertureSDK;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpServerCodec;
import net.bytebuddy.asm.Advice;

public class NettyServerAdvice {
    public static ApertureSDK apertureSDK = Config.newSDKFromConfig();

    @Advice.OnMethodExit
    public static void onExit(
            @Advice.This ChannelPipeline pipeline,
            @Advice.Argument(1) String handlerName,
            @Advice.Argument(2) ChannelHandler handler) {
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
}
